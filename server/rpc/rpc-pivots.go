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

	"github.com/gsmith257-cyber/better-sliver-package/protobuf/clientpb"
	"github.com/gsmith257-cyber/better-sliver-package/protobuf/commonpb"
	"github.com/gsmith257-cyber/better-sliver-package/protobuf/sliverpb"
	"github.com/gsmith257-cyber/better-sliver-package/server/core"
)

// PivotGraph - Return the server's pivot graph
func (rpc *Server) PivotGraph(ctx context.Context, req *commonpb.Empty) (*clientpb.PivotGraph, error) {
	pivotGraph := &clientpb.PivotGraph{
		Children: []*clientpb.PivotGraphEntry{},
	}
	for _, topLevel := range core.PivotGraph() {
		pivotGraph.Children = append(pivotGraph.Children, topLevel.ToProtobuf())
	}
	return pivotGraph, nil
}

// PivotSessionListeners - Get a list of all pivot listeners from an implant
func (rpc *Server) PivotSessionListeners(ctx context.Context, req *sliverpb.PivotListenersReq) (*sliverpb.PivotListeners, error) {
	resp := &sliverpb.PivotListeners{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PivotStartListener - Instruct the implant to start a pivot listener
func (rpc *Server) PivotStartListener(ctx context.Context, req *sliverpb.PivotStartListenerReq) (*sliverpb.PivotListener, error) {
	resp := &sliverpb.PivotListener{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PivotStopListener - Instruct the implant to stop a pivot listener
func (rpc *Server) PivotStopListener(ctx context.Context, req *sliverpb.PivotStopListenerReq) (*commonpb.Empty, error) {
	resp := &sliverpb.PivotListener{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return &commonpb.Empty{}, nil
}

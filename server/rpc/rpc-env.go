package rpc

/*
	Sliver Implant Framework
	Copyright (C) 2021  Bishop Fox

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

	"github.com/gsmith257-cyber/better-sliver-package/protobuf/commonpb"
	"github.com/gsmith257-cyber/better-sliver-package/protobuf/sliverpb"
)

// GetEnv - Retrieve the environment variables list from the current session
func (rpc *Server) GetEnv(ctx context.Context, req *sliverpb.EnvReq) (*sliverpb.EnvInfo, error) {
	resp := &sliverpb.EnvInfo{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SetEnv - Set an environment variable
func (rpc *Server) SetEnv(ctx context.Context, req *sliverpb.SetEnvReq) (*sliverpb.SetEnv, error) {
	resp := &sliverpb.SetEnv{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UnsetEnv - Set an environment variable
func (rpc *Server) UnsetEnv(ctx context.Context, req *sliverpb.UnsetEnvReq) (*sliverpb.UnsetEnv, error) {
	resp := &sliverpb.UnsetEnv{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

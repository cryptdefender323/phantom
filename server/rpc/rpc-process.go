package rpc

/*
	Phantom Implant Framework
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

	"github.com/cryptdefender3232/phantom/protobuf/commonpb"
	"github.com/cryptdefender3232/phantom/protobuf/phantompb"
)

// Ps - List the processes on the remote machine
func (rpc *Server) Ps(ctx context.Context, req *phantompb.PsReq) (*phantompb.Ps, error) {
	resp := &phantompb.Ps{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ProcessDump - Dump the memory of a remote process
func (rpc *Server) ProcessDump(ctx context.Context, req *phantompb.ProcessDumpReq) (*phantompb.ProcessDump, error) {
	resp := &phantompb.ProcessDump{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Terminate - Terminate a remote process
func (rpc *Server) Terminate(ctx context.Context, req *phantompb.TerminateReq) (*phantompb.Terminate, error) {
	resp := &phantompb.Terminate{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

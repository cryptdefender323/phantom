package rpc

/*
	Phantom Implant Framework
	Copyright (C) 2023  Bishop Fox

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
	"github.com/cryptdefender3232/phantom/server/log"
)

var (
	rpcWasmLog = log.NamedLogger("rpc", "wasm")
)

// RegisterWasmExtension - Register a new wasm extension with the implant
func (rpc *Server) RegisterWasmExtension(ctx context.Context, req *phantompb.RegisterWasmExtensionReq) (*phantompb.RegisterWasmExtension, error) {
	resp := &phantompb.RegisterWasmExtension{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListWasmExtensions - List registered wasm extensions
func (rpc *Server) ListWasmExtensions(ctx context.Context, req *phantompb.ListWasmExtensionsReq) (*phantompb.ListWasmExtensions, error) {
	resp := &phantompb.ListWasmExtensions{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ExecWasmExtension - Execute a wasm extension
func (rpc *Server) ExecWasmExtension(ctx context.Context, req *phantompb.ExecWasmExtensionReq) (*phantompb.ExecWasmExtension, error) {
	resp := &phantompb.ExecWasmExtension{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

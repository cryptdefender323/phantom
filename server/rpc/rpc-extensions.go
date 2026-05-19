package rpc

import (
	"context"

	"github.com/cryptdefender3232/phantom/protobuf/commonpb"
	"github.com/cryptdefender3232/phantom/protobuf/phantompb"
)

// RegisterExtension registers a new extension in the implant
func (rpc *Server) RegisterExtension(ctx context.Context, req *phantompb.RegisterExtensionReq) (*phantompb.RegisterExtension, error) {
	resp := &phantompb.RegisterExtension{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListExtensions lists the registered extensions
func (rpc *Server) ListExtensions(ctx context.Context, req *phantompb.ListExtensionsReq) (*phantompb.ListExtensions, error) {
	resp := &phantompb.ListExtensions{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CallExtension calls a specific export of the loaded extension
func (rpc *Server) CallExtension(ctx context.Context, req *phantompb.CallExtensionReq) (*phantompb.CallExtension, error) {
	resp := &phantompb.CallExtension{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

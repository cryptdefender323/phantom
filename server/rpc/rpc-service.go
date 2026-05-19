package rpc

import (
	"context"

	"github.com/cryptdefender3232/phantom/protobuf/commonpb"
	"github.com/cryptdefender3232/phantom/protobuf/phantompb"
)

// Services - List and control services
func (rpc *Server) Services(ctx context.Context, req *phantompb.ServicesReq) (*phantompb.Services, error) {
	resp := &phantompb.Services{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rpc *Server) ServiceDetail(ctx context.Context, req *phantompb.ServiceDetailReq) (*phantompb.ServiceDetail, error) {
	resp := &phantompb.ServiceDetail{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// StartService creates and starts a Windows service on a remote host
func (rpc *Server) StartService(ctx context.Context, req *phantompb.StartServiceReq) (*phantompb.ServiceInfo, error) {
	resp := &phantompb.ServiceInfo{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rpc *Server) StartServiceByName(ctx context.Context, req *phantompb.StartServiceByNameReq) (*phantompb.ServiceInfo, error) {
	resp := &phantompb.ServiceInfo{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// StopService stops a remote service
func (rpc *Server) StopService(ctx context.Context, req *phantompb.StopServiceReq) (*phantompb.ServiceInfo, error) {
	resp := &phantompb.ServiceInfo{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RemoveService deletes a service from the remote system
func (rpc *Server) RemoveService(ctx context.Context, req *phantompb.RemoveServiceReq) (*phantompb.ServiceInfo, error) {
	resp := &phantompb.ServiceInfo{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

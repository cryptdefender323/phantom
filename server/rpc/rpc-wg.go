package rpc

import (
	"context"

	"github.com/cryptdefender3232/phantom/protobuf/clientpb"
	"github.com/cryptdefender3232/phantom/protobuf/commonpb"
	"github.com/cryptdefender3232/phantom/protobuf/phantompb"
	"github.com/cryptdefender3232/phantom/server/certs"
	"github.com/cryptdefender3232/phantom/server/generate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GenerateWGClientConfig - Generate a client config for a WG interface
func (rpc *Server) GenerateWGClientConfig(ctx context.Context, _ *commonpb.Empty) (*clientpb.WGClientConfig, error) {
	clientIP, privkey, pubkey, err := generate.GenerateUniqueWGPeerKeys()
	if err != nil {
		rpcLog.Errorf("Could not generate WG keys: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	_, serverPubKey, err := certs.GetWGServerKeys()
	if err != nil {
		rpcLog.Errorf("Could not get WG server keys: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &clientpb.WGClientConfig{
		ClientPrivateKey: privkey,
		ClientIP:         clientIP,
		ClientPubKey:     pubkey,
		ServerPubKey:     serverPubKey,
	}

	return resp, nil
}

// WGStartPortForward - Start a port forward
func (rpc *Server) WGStartPortForward(ctx context.Context, req *phantompb.WGPortForwardStartReq) (*phantompb.WGPortForward, error) {
	resp := &phantompb.WGPortForward{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WGStopPortForward - Stop a port forward
func (rpc *Server) WGStopPortForward(ctx context.Context, req *phantompb.WGPortForwardStopReq) (*phantompb.WGPortForward, error) {
	resp := &phantompb.WGPortForward{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WGAddForwarder - Add a TCP forwarder
func (rpc *Server) WGStartSocks(ctx context.Context, req *phantompb.WGSocksStartReq) (*phantompb.WGSocks, error) {
	resp := &phantompb.WGSocks{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WGStopForwarder - Stop a TCP forwarder
func (rpc *Server) WGStopSocks(ctx context.Context, req *phantompb.WGSocksStopReq) (*phantompb.WGSocks, error) {
	resp := &phantompb.WGSocks{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rpc *Server) WGListSocksServers(ctx context.Context, req *phantompb.WGSocksServersReq) (*phantompb.WGSocksServers, error) {
	resp := &phantompb.WGSocksServers{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WGAddForwarder - List wireguard forwarders
func (rpc *Server) WGListForwarders(ctx context.Context, req *phantompb.WGTCPForwardersReq) (*phantompb.WGTCPForwarders, error) {
	resp := &phantompb.WGTCPForwarders{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

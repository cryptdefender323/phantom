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
	"crypto/sha256"
	"fmt"

	"github.com/cryptdefender3232/phantom/protobuf/commonpb"
	"github.com/cryptdefender3232/phantom/protobuf/phantompb"
	"github.com/cryptdefender3232/phantom/server/core"
	"github.com/cryptdefender3232/phantom/server/db"
	"github.com/cryptdefender3232/phantom/server/db/models"
	"github.com/cryptdefender3232/phantom/server/log"
	"github.com/cryptdefender3232/phantom/util/encoders"
)

var (
	fsLog = log.NamedLogger("rcp", "fs")
)

// Ls - List a directory
func (rpc *Server) Ls(ctx context.Context, req *phantompb.LsReq) (*phantompb.Ls, error) {
	resp := &phantompb.Ls{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Mv - Move or rename a file
func (rpc *Server) Mv(ctx context.Context, req *phantompb.MvReq) (*phantompb.Mv, error) {
	resp := &phantompb.Mv{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Cp - Copy a file to another location
func (rpc *Server) Cp(ctx context.Context, req *phantompb.CpReq) (*phantompb.Cp, error) {
	resp := &phantompb.Cp{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Rm - Remove file or directory
func (rpc *Server) Rm(ctx context.Context, req *phantompb.RmReq) (*phantompb.Rm, error) {
	resp := &phantompb.Rm{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Mkdir - Make a directory
func (rpc *Server) Mkdir(ctx context.Context, req *phantompb.MkdirReq) (*phantompb.Mkdir, error) {
	resp := &phantompb.Mkdir{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Cd - Change directory
func (rpc *Server) Cd(ctx context.Context, req *phantompb.CdReq) (*phantompb.Pwd, error) {
	resp := &phantompb.Pwd{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Pwd - Print working directory
func (rpc *Server) Pwd(ctx context.Context, req *phantompb.PwdReq) (*phantompb.Pwd, error) {
	resp := &phantompb.Pwd{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Download - Download a file from the remote file system
func (rpc *Server) Download(ctx context.Context, req *phantompb.DownloadReq) (*phantompb.Download, error) {
	resp := &phantompb.Download{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Upload - Upload a file from the remote file system
func (rpc *Server) Upload(ctx context.Context, req *phantompb.UploadReq) (*phantompb.Upload, error) {
	resp := &phantompb.Upload{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	if req.IsIOC {
		go trackIOC(req, resp)
	}
	return resp, nil
}

// Chmod - Change permission on a file or directory
func (rpc *Server) Chmod(ctx context.Context, req *phantompb.ChmodReq) (*phantompb.Chmod, error) {
	resp := &phantompb.Chmod{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Chown - Change owner on a file or directory
func (rpc *Server) Chown(ctx context.Context, req *phantompb.ChownReq) (*phantompb.Chown, error) {
	resp := &phantompb.Chown{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Chtimes - Change file access and modification times on a file or directory
func (rpc *Server) Chtimes(ctx context.Context, req *phantompb.ChtimesReq) (*phantompb.Chtimes, error) {
	resp := &phantompb.Chtimes{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MemfilesList - List memfiles
func (rpc *Server) MemfilesList(ctx context.Context, req *phantompb.MemfilesListReq) (*phantompb.Ls, error) {
	resp := &phantompb.Ls{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MemfilesAdd - Add memfile
func (rpc *Server) MemfilesAdd(ctx context.Context, req *phantompb.MemfilesAddReq) (*phantompb.MemfilesAdd, error) {
	resp := &phantompb.MemfilesAdd{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MemfilesRm - Close memfile
func (rpc *Server) MemfilesRm(ctx context.Context, req *phantompb.MemfilesRmReq) (*phantompb.MemfilesRm, error) {
	resp := &phantompb.MemfilesRm{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func hashUploadData(encoder string, data []byte) [32]byte {
	if encoder == "gzip" {
		decodedData, err := new(encoders.Gzip).Decode(data)
		if err != nil {
			return sha256.Sum256(nil)
		}
		return sha256.Sum256(decodedData)
	} else {
		return sha256.Sum256(data)
	}
}

func trackIOC(req *phantompb.UploadReq, resp *phantompb.Upload) {
	fsLog.Debugf("Adding IOC to database ...")
	request := req.GetRequest()
	if request == nil {
		fsLog.Error("No request for upload")
		return
	}
	session := core.Sessions.Get(request.SessionID)
	if session == nil {
		fsLog.Error("No session for upload request")
		return
	}
	host, err := db.HostByHostUUID(session.UUID)
	if err != nil {
		fsLog.Errorf("No host for session uuid %v", session.UUID)
		return
	}

	sum := hashUploadData(req.Encoder, req.Data)
	ioc := &models.IOC{
		HostID:   host.HostUUID,
		Path:     resp.Path,
		FileHash: fmt.Sprintf("%x", sum),
	}
	if db.Session().Create(ioc).Error != nil {
		fsLog.Error("Failed to create IOC")
	}
}

// Grep - Search a file or directory for text matching a regex
func (rpc *Server) Grep(ctx context.Context, req *phantompb.GrepReq) (*phantompb.Grep, error) {
	resp := &phantompb.Grep{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Mount - Get information on mounted filesystems
func (rpc *Server) Mount(ctx context.Context, req *phantompb.MountReq) (*phantompb.Mount, error) {
	resp := &phantompb.Mount{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

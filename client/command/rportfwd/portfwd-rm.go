package rportfwd

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

	"github.com/cryptdefender3232/phantom/client/console"
	"github.com/cryptdefender3232/phantom/protobuf/phantompb"
	"github.com/spf13/cobra"
)

// StartRportFwdListenerCmd - Start listener for reverse port forwarding on implant.
func StopRportFwdListenerCmd(cmd *cobra.Command, con *console.PhantomClient, args []string) {
	session := con.ActiveTarget.GetSessionInteractive()
	if session == nil {
		return
	}

	listenerID, _ := cmd.Flags().GetUint32("id")
	rportfwdListener, err := con.Rpc.StopRportFwdListener(context.Background(), &phantompb.RportFwdStopListenerReq{
		Request: con.ActiveTarget.Request(cmd),
		ID:      listenerID,
	})
	if err != nil {
		con.PrintWarnf("%s\n", err)
		return
	}
	printStoppedRportFwdListener(rportfwdListener, con)
}

func printStoppedRportFwdListener(rportfwdListener *phantompb.RportFwdListener, con *console.PhantomClient) {
	if rportfwdListener.Response != nil && rportfwdListener.Response.Err != "" {
		con.PrintErrorf("%s", rportfwdListener.Response.Err)
		return
	}
	con.PrintInfof("Stopped reverse port forwarding %s <- %s\n", rportfwdListener.ForwardAddress, rportfwdListener.BindAddress)
}

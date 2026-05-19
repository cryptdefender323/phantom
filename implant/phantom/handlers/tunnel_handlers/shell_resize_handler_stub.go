//go:build !darwin && !linux && !freebsd && !openbsd && !dragonfly

package tunnel_handlers

import (
	"github.com/cryptdefender3232/phantom/implant/phantom/transports"
	"github.com/cryptdefender3232/phantom/protobuf/commonpb"
	"github.com/cryptdefender3232/phantom/protobuf/phantompb"
	"google.golang.org/protobuf/proto"
)

func ShellResizeReqHandler(envelope *phantompb.Envelope, connection *transports.Connection) {
	resp, _ := proto.Marshal(&commonpb.Empty{})
	connection.Send <- &phantompb.Envelope{
		ID:   envelope.ID,
		Data: resp,
	}
}

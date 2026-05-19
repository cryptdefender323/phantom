//go:build !client

package serverctx

import (
	"github.com/cryptdefender3232/phantom/client/console"
	"github.com/spf13/cobra"
)

// Commands is a no-op when building without the `client` build tag (e.g. phantom-server).
func Commands(_ *console.PhantomClient) []*cobra.Command {
	return nil
}

//go:build server && go_sqlite && phantom_e2e

package c2

import "net"

// HandleWGPhantomConnectionForTest exposes the raw connection handler to external
// (package `c2_test`) end-to-end tests without shipping test-only symbols in
// production builds.
func HandleWGPhantomConnectionForTest(conn net.Conn) {
	handleWGPhantomConnection(conn)
}

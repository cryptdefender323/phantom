//go:build server && go_sqlite && phantom_e2e

package c2

import "net"

// HandlePhantomConnectionForTest exposes the raw connection handler to external
// (package `c2_test`) end-to-end tests without shipping test-only symbols in
// production builds.
func HandlePhantomConnectionForTest(conn net.Conn) {
	handlePhantomConnection(conn)
}

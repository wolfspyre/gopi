/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register rpc/server
	gopi.RegisterModule(gopi.Module{
		Name: "rpc/server",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("rpc.port", 0, "Server Port")
			config.AppFlags.FlagString("rpc.sslcert", "", "SSL Certificate Path")
			config.AppFlags.FlagString("rpc.sslkey", "", "SSL Key Path")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			port, _ := app.AppFlags.GetUint("rpc.port")
			key, _ := app.AppFlags.GetString("rpc.sslkey")
			cert, _ := app.AppFlags.GetString("rpc.sslcert")
			return gopi.Open(Server{Port: port, SSLCertificate: cert, SSLKey: key}, app.Logger)
		},
	})

	// Register rpc/client
	gopi.RegisterModule(gopi.Module{
		Name: "rpc/client",
		Type: gopi.MODULE_TYPE_OTHER,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("rpc.addr", "localhost:8001", "Address")
			config.AppFlags.FlagBool("rpc.ssl", false, "SSL Enabled")
			config.AppFlags.FlagBool("rpc.skipverify", true, "Skip SSL Verification")
			config.AppFlags.FlagDuration("rpc.timeout", 0, "Connection Timeout")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			addr, _ := app.AppFlags.GetString("rpc.addr")
			ssl, _ := app.AppFlags.GetBool("rpc.ssl")
			skipverify, _ := app.AppFlags.GetBool("rpc.skipverify")
			timeout, _ := app.AppFlags.GetDuration("rpc.timeout")
			return gopi.Open(Client{Host: addr, SSL: ssl, SkipVerify: skipverify, Timeout: timeout}, app.Logger)
		},
	})

	// Register rpc/discovery
	gopi.RegisterModule(gopi.Module{
		Name: "rpc/discovery",
		Type: gopi.MODULE_TYPE_MDNS,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagString("mdns.domain", MDNS_DEFAULT_DOMAIN, "Domain")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			domain, _ := app.AppFlags.GetString("mdns.domain")
			return gopi.Open(Config{Domain: domain}, app.Logger)
		},
	})
}

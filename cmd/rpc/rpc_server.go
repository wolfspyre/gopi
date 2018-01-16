/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// The server serves the GRPC reflection package and the
// helloworld package, which is described in helloworld/helloworld.proto
// In order to install this package, you will need to run go generate with
// both the protoc compiler and the GRPC GO plugin available:
//
// mac# brew install protobuf
// rpi# sudo apt install protobuf-compiler
// go get -u github.com/golang/protobuf/protoc-gen-go
//
// Then:
//
// cd "${GOPATH}/src/github.com/djthorpe/gopi"
// go generate cmd/rpc/rpc_server.go
// go install cmd/rpc/rpc_server.go
//
package main

//go:generate protoc helloworld/helloworld.proto --go_out=plugins=grpc:.

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	context "golang.org/x/net/context"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"

	// Helloworld Protocol Buffer
	hw "github.com/djthorpe/gopi/cmd/rpc/helloworld"
)

////////////////////////////////////////////////////////////////////////////////
// HelloworldService implementation

type HelloworldService struct{}

func (this *HelloworldService) Register(server gopi.RPCServer) error {
	// Check to make sure we satisfy the interface
	var _ hw.GreeterServer = (*HelloworldService)(nil)

	return server.Fudge(reflect.ValueOf(hw.RegisterGreeterServer), this)
}

func (this *HelloworldService) ServiceType() string {
	return "helloworld"
}

func (this *HelloworldService) SayHello(ctx context.Context, req *hw.HelloRequest) (*hw.HelloReply, error) {
	if req.Name == "" {
		req.Name = "World"
	}
	return &hw.HelloReply{
		Message: "Hello, " + req.Name,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////

func EventLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	server, ok := app.ModuleInstance("rpc/server").(gopi.RPCServer)
	if server == nil || ok == false {
		return errors.New("rpc/server missing")
	}

	// Listen for events
	c := server.Subscribe()
FOR_LOOP:
	for {
		select {
		case evt := <-c:
			if evt != nil {
				fmt.Println(evt)
			}
		case <-done:
			break FOR_LOOP
		}
	}

	// Stop listening for events
	server.Unsubscribe(c)

	return nil
}

func ServerLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	server, ok := app.ModuleInstance("rpc/server").(gopi.RPCServer)
	if server == nil || ok == false {
		return errors.New("rpc/server missing")
	}

	// Create the helloworld module
	if hw_service := new(HelloworldService); hw_service == nil {
		return errors.New("HelloworldService missing")
	} else {
		// Start server - will end when Stop is called
		server.Start(hw_service)
	}

	// wait for done
	<-done

	// Bomb out
	return nil
}

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	server, ok := app.ModuleInstance("rpc/server").(gopi.RPCServer)
	if server == nil || ok == false {
		return errors.New("rpc/server missing")
	}

	app.WaitForSignal()
	server.Close()

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the lirc instance
	config := gopi.NewAppConfig("rpc/server")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop, ServerLoop, EventLoop))
}
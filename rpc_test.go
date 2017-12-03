package gopi_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi"
	logger "github.com/djthorpe/gopi/sys/logger"
	mdns "github.com/djthorpe/gopi/sys/mdns"
)

func TestRPCDiscovery_000(t *testing.T) {
	if logger, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Fatal(err)
	} else if driver, err := gopi.Open(mdns.Config{}, logger.(gopi.Logger)); err != nil {
		t.Fatal(err)
	} else {
		defer driver.Close()

		mdns := driver.(gopi.RPCServiceDiscovery)
		serviceType := "_smb._tcp"

		// Wait for service records
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err := mdns.Browse(ctx, serviceType, func(service *gopi.RPCService) {
			if service != nil {
				fmt.Println(service)
			} else {
				fmt.Println("No more records")
			}
		}); err != nil {
			t.Error(err)
		}

		// Register a service
		if err := mdns.Register(&gopi.RPCService{Name: "My Service", Type: serviceType, Port: 8000}); err != nil {
			t.Error(err)
		}

		time.Sleep(5 * time.Second)
	}
}
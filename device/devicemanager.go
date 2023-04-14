package devicemanager

import (
	"fmt"
	"github.com/google/gousb"
	"image/color"
	"log"
	"nzxt-driver-dev/driver"
)

type RgbFanHub struct {
	hidEndpoint *gousb.OutEndpoint
}

func (hub RgbFanHub) SetColor(channel int, c color.Color) {
	packet := driver.MakeColorPacket(channel, driver.CommandRgbRing, []color.Color{c})
	hub.hidEndpoint.Write(packet)
}

func (hub RgbFanHub) SetColors(channel int, colors []color.Color) {
	packet := driver.MakeColorPacket(channel, driver.CommandRgbPixels, colors)
	hub.hidEndpoint.Write(packet)
}

func GetManagedDriver(id int, handler func(hub RgbFanHub)) {
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Find the USB device by vendor ID and product ID
	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Product == gousb.ID(id) // Replace with your product ID
	})

	if err != nil {
		fmt.Printf("Error finding device: %v\n", err)
		return
	}
	defer func() {
		for _, dev := range devs {
			dev.Close()
		}
	}()

	// Detach the kernel driver from the USB device

	for _, dev := range devs {
		if err := dev.SetAutoDetach(true); err != nil {
			fmt.Printf("Error setting auto detach: %v\n", err)
			return
		}
		if err := dev.Reset(); err != nil {
			fmt.Printf("Error resetting device: %v\n", err)
			return
		}

		cfg, err := dev.Config(1)
		if err != nil {
			log.Fatalf("Failed to get config: %v", err)
		}
		defer cfg.Close()

		intf, err := cfg.Interface(0, 0)
		if err != nil {
			log.Fatalf("Failed to get interface: %v", err)
		}
		defer intf.Close()

		endpoints := intf.Setting.Endpoints
		if err != nil {
			log.Fatalf("Failed to get endpoints: %v", err)
		}

		var outEndpoint *gousb.OutEndpoint
		for _, endpoint := range endpoints {
			if endpoint.Address&0x80 == 0x00 && endpoint.TransferType == gousb.TransferTypeInterrupt {
				outEndpoint, err = intf.OutEndpoint(int(endpoint.Address))
				if err != nil {
					log.Fatalf("Failed to get OutEndpoint: %v", err)
				}
				break
			}
		}

		if outEndpoint == nil {
			log.Fatalf("Failed to find OutEndpoint")
		}

		// send stuff

		var hub = RgbFanHub{hidEndpoint: outEndpoint}
		handler(hub)
	}
	fmt.Println("Driver detached from USB device")
}

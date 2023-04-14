package main

import (
	"image/color"
	devicemanager "nzxt-driver-dev/device"
)

func main() {

	devicemanager.GetManagedDriver(0x2011, func(hub devicemanager.RgbFanHub) {
		var c = color.RGBA{R: 50, G: 0, B: 25, A: 255}
		// loop 4 times
		for i := 0; i < 4; i++ {
			hub.SetColor(i+1, c)
		}
	})

}

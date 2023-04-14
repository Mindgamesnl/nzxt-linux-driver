package main

import (
	"image/color"
	"math"
	devicemanager "nzxt-driver-dev/device"
	"time"
)

func main() {

	devicemanager.GetManagedDriver(0x2011, func(hub devicemanager.RgbFanHub) {
		hub.SetColor(1, color.RGBA{R: 0, B: 255, G: 0, A: 255})
		var r, g, b float64 = 0, 0, 0
		for {
			r = math.Sin(0.5*float64(time.Now().UnixNano())/1000000000)*127 + 128
			g = math.Sin(0.6*float64(time.Now().UnixNano())/1000000000)*127 + 128
			b = math.Sin(0.7*float64(time.Now().UnixNano())/1000000000)*127 + 128

			c := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}

			for i := 0; i < 4; i++ {
				hub.SetColor(i+1, c)
			}

			time.Sleep(time.Millisecond * 10)
		}
	})

}

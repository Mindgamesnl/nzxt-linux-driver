package main

import (
	"fmt"
	"image/color"
	devicemanager "nzxt-driver-dev/device"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]

	if len(args) < 4 {
		fmt.Println("Usage: nzxt-led <int|all> <int> <int> <int>")
		return
	}

	commandOrTarget := args[0]
	colors := args[1:]

	r, err := strconv.Atoi(colors[0])
	if err != nil && colors[0] != "all" {
		fmt.Println("Invalid argument:", colors[0])
		return
	}

	g, err := strconv.Atoi(colors[1])
	if err != nil {
		fmt.Println("Invalid argument:", colors[1])
		return
	}

	b, err := strconv.Atoi(colors[2])
	if err != nil {
		fmt.Println("Invalid argument:", colors[2])
		return
	}

	var isCommandNumber = false
	if commandOrTarget != "all" {
		_, err := strconv.Atoi(commandOrTarget)
		if err == nil {
			isCommandNumber = true
		}
	}

	var cmdStr string
	if commandOrTarget == "all" {
		cmdStr = fmt.Sprintf("Setting color of all LEDs to (%d,%d,%d)", r, g, b)
		set(-1, r, g, b)
	} else if isCommandNumber {
		device, err := strconv.Atoi(commandOrTarget)
		if err != nil {
			fmt.Println("Invalid argument:", commandOrTarget)
			return
		}
		cmdStr = fmt.Sprintf("Setting color of LED %d to (%d,%d,%d)", device, r, g, b)
		set(device, r, g, b)
	}

	// send the command to the NZXT device here...

	fmt.Println(cmdStr)
}

func set(id int, r int, g int, b int) {
	devicemanager.GetManagedDriver(0x2011, func(hub devicemanager.RgbFanHub) {
		var c = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}

		if id == -1 {
			for i := 0; i < 6; i++ {
				hub.SetColor(i+1, c)
			}
		} else {
			hub.SetColor(id, c)
		}
	})
}

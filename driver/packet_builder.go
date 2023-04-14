package driver

import (
	"image/color"
)

// CommandType packet types,
type CommandType byte
type DeviceId byte

// commands
const (
	CommandRgbRing   = CommandType(0b00101010)
	CommandRgbPixels = CommandType(0b00100010)
)

func MakeDeviceId(readableId int) DeviceId {
	var v = 0b00000001
	if readableId > 1 {
		v = v << (readableId - 1)
	}
	return DeviceId(byte(v))
}

func MakePacketHeader(deviceId DeviceId, packetType CommandType) []byte {
	if packetType == CommandRgbRing {
		return []byte{
			byte(packetType),
			byte(0b00000100),
			byte(deviceId),
			byte(deviceId),
			byte(0b00000000),
			byte(0b00110010),
			byte(0b00000000),
		}
	}

	if packetType == CommandRgbPixels {
		return []byte{
			byte(packetType),
			byte(0x10),
			byte(deviceId),
			byte(0b00000000),
		}
	}

	panic("invalid packet type")
}

func MakeColorSection(colors []color.Color) []byte {
	// make an array of a byte for each channel, 16 pixels * rgb
	var colorSection = make([]byte, 18*3)

	// fil the array with 0 bytes
	for i := range colorSection {
		colorSection[i] = 0b00000000
	}

	// for each color, add the rgb values to the color section
	for i, c := range colors {
		var r, g, b, _ = c.RGBA()
		colorSection[i*3] = byte(g)
		colorSection[i*3+1] = byte(r)
		colorSection[i*3+2] = byte(b)
	}

	return colorSection
}

func MakeColorPacket(readableDeviceId int, mode CommandType, colors []color.Color) []byte {
	var packet = make([]byte, 64)

	// add the header
	var header = MakePacketHeader(MakeDeviceId(readableDeviceId), mode)
	copy(packet, header)

	// add the color section
	var colorSection = MakeColorSection(colors)
	copy(packet[len(header):], colorSection)

	// add 6 bytes of padding (from byte 60)
	for i := 60; i < 64; i++ {
		packet[i] = 0b00000000
	}

	//var out = ""
	//for _, b := range packet {
	//	out += fmt.Sprintf("%02x", b)
	//	out += ""
	//}
	//fmt.Println(out)

	return packet
}

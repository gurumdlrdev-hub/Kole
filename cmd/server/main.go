package main

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/go-vgo/robotgo"
)

func main() {

	screenWidth, screenHeight := robotgo.GetScreenSize()
	fmt.Printf("Laptop Screen Resolution: %dx%d\n", screenWidth, screenHeight)

	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:9000")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening on UDP:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Kole Server (UDP) Port 9000 is Ready...")
	fmt.Println("Waiting for Data...")

	buffer := make([]byte, 5)
	isMouseDown := false

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil || n < 5 {
			continue
		}

		fmt.Printf("Received raw: %v\n", buffer)

		eventType := buffer[0]
		rawX := binary.BigEndian.Uint16(buffer[1:3])
		rawY := binary.BigEndian.Uint16(buffer[3:5])

		targetX := int((float64(rawX) / 65535.0) * float64(screenWidth))
		targetY := int((float64(rawY) / 65535.0) * float64(screenHeight))

		if eventType == 1 {

			robotgo.Move(targetX, targetY)

			if !isMouseDown {
				robotgo.Move(targetX, targetY)
				robotgo.MouseDown("left")
				isMouseDown = true

			} else {
				robotgo.DragSmooth(targetX, targetY)

			}
		} else if eventType == 0 {
			if isMouseDown {
				robotgo.Toggle("left")
				isMouseDown = false
			}

		}

	}
}

package main

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"net"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type TouchPad struct {
	widget.BaseWidget
	conn *net.UDPConn
}

func NewTouchPad(conn *net.UDPConn) *TouchPad {
	tp := &TouchPad{conn: conn}
	tp.ExtendBaseWidget(tp)
	return tp
}

func (tp *TouchPad) send(eventType byte, pos fyne.Position) {
	size := tp.Size()
	if size.Width <= 0 || size.Height <= 0 {
		return
	}

	normX := float64(pos.X) / float64(size.Width)
	normY := float64(pos.Y) / float64(size.Height)

	if normX < 0 {
		normX = 0
	} else if normX > 1 {
		normX = 1
	}
	if normY < 0 {
		normY = 0
	} else if normY > 1 {
		normY = 1
	}

	uintX := uint16(normX * 65535.0)
	uintY := uint16(normY * 65535.0)

	packet := make([]byte, 5)
	packet[0] = eventType
	binary.BigEndian.PutUint16(packet[1:3], uintX)
	binary.BigEndian.PutUint16(packet[3:5], uintY)

	tp.conn.Write(packet)
}

func (tp *TouchPad) Dragged(e *fyne.DragEvent) {
	tp.send(1, e.Position)
}

func (tp *TouchPad) DragEnd() {
	packet := make([]byte, 5)
	packet[0] = 0
	tp.conn.Write(packet)
}

func (tp *TouchPad) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.NRGBA{R: 15, G: 15, B: 15, A: 255})
	return &touchPadRenderer{bg: bg}
}

type touchPadRenderer struct {
	bg *canvas.Rectangle
}

func (r *touchPadRenderer) Layout(s fyne.Size)           { r.bg.Resize(s) }
func (r *touchPadRenderer) MinSize() fyne.Size           { return fyne.NewSize(300, 200) }
func (r *touchPadRenderer) Refresh()                     { r.bg.Refresh() }
func (r *touchPadRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.bg} }
func (r *touchPadRenderer) Destroy()                     {}

func main() {

	serverAddr, err := net.ResolveUDPAddr("udp", "192.168.1.34:9000")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Address error: %v\n", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "UDP error: %v\n", err)
		return
	}
	defer conn.Close()

	myApp := app.New()
	myWindow := myApp.NewWindow("Kole")

	pad := NewTouchPad(conn)
	myWindow.SetContent(pad)
	myWindow.SetFullScreen(true)
	myWindow.ShowAndRun()
}

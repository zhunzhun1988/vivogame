package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type DrawerInfo struct {
	seatX, seatY   int
	seatW, seatH   int
	titleX, titleY int
	titleSize      float64
}

type Window struct {
	width      int
	height     int
	background color.Color
	drawer     *Drawer
	xu         *xgbutil.XUtil
	canvas     *image.RGBA
	xwin       *xwindow.Window
	drawInfos  []DrawerInfo
	size       int
}

func NewWindow(w, h int, size int, bg color.Color) *Window {
	xu, err := xgbutil.NewConn()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// just create a id for the window
	xwin, err := xwindow.Generate(xu)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// now, create the window
	err = xwin.CreateChecked(
		xu.RootWin(), // parent window
		0, 0, w, h,   // window size
		0) // related to event, not considered here
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// now we can see the window on the screen
	xwin.Map()

	// 'rstr' calculates the data needed to draw
	// 'painter' draw with the data on 'canvas'
	canvas := image.NewRGBA(image.Rect(0, 0, w, h))

	d := NewDrawer(xu, xwin, canvas, color.RGBA{0x00, 0x00, 0x00, 0xff})
	win := &Window{
		width:      w,
		height:     h,
		background: bg,
		drawer:     d,
		xu:         xu,
		canvas:     canvas,
		xwin:       xwin,
		size:       size,
		drawInfos:  make([]DrawerInfo, size*2+1),
	}
	win.initDrawInfos()
	return win
}
func (w *Window) getStrColor(str string) color.RGBA {
	if str >= "A" && str <= "Z" {
		return color.RGBA{0xff, 0xff, 0x00, 0xff}
	} else {
		return color.RGBA{0x00, 0xff, 0x00, 0xff}
	}
}
func (w *Window) Move(g *Game, str string) {
	var moveStartX, moveStartY int
	var moveEndX, moveEndY int
	seatW := (w.width/(g.size+1))*2/3 - 10
	seatH := w.height / 4
	l := seatW
	if l > seatH {
		l = seatH
	}
	for i := 1; i <= g.size; i++ {
		if g.emptyIndex != i-1 && g.items[i-1].title == str {
			moveStartX = (seatW*3/2)*(i-1) + seatW/2 + l/2 - 25
			moveStartY = seatH/2 + l/2 - 25
		} else if g.emptyIndex == i-1 {
			moveEndX = (seatW*3/2)*(i-1) + seatW/2 + l/2 - 25
			moveEndY = seatH/2 + l/2 - 25
		}
		if g.emptyIndex != g.size*2-i+1 && g.items[g.size*2+1-i].title == str {
			moveStartX = (seatW*3/2)*(i-1) + seatW/2 + l/2 - 25
			moveStartY = w.height - (seatH / 2) - l + l/2 - 25
		} else if g.emptyIndex == g.size*2-i+1 {
			moveEndX = (seatW*3/2)*(i-1) + seatW/2 + l/2 - 25
			moveEndY = w.height - (seatH / 2) - l + l/2 - 25
		}
	}
	if g.emptyIndex != g.size && g.items[g.size].title == str {
		moveStartX = (seatW*3/2)*g.size + seatW/2 + l/2 - 25
		moveStartY = w.height/2 - 25
	} else if g.emptyIndex == g.size {
		moveEndX = (seatW*3/2)*g.size + seatW/2 + l/2 - 25
		moveEndY = w.height/2 - 25
	}
	drawMove := func(fromX, fromY, toX, toY int) {
		w.drawer.FillRect(DrawPoint{fromX, fromY}, 50, 50, w.background)
		w.DrawStatus(g, str)
		w.drawer.DrawText(DrawPoint{toX, toY}, str, 50, w.getStrColor(str))
		w.drawer.Show()
	}
	stepX := (moveEndX - moveStartX) / 4
	stepY := (moveEndY - moveStartY) / 4
	for i := 0; i < 3; i++ {
		drawMove(moveStartX, moveStartY, moveStartX+stepX, moveStartY+stepY)
		moveStartX = moveStartX + stepX
		moveStartY = moveStartY + stepY
		time.Sleep(10 * time.Millisecond)
	}
	drawMove(moveStartX, moveStartY, moveEndX, moveEndY)
	time.Sleep(50 * time.Millisecond)
}

func (w *Window) initDrawInfos() {
	seatW := (w.width/(w.size+1))*2/3 - 10
	seatH := w.height / 4
	l := seatW
	if l > seatH {
		l = seatH
	}
	for i := 1; i <= w.size; i++ {
		w.drawInfos[i-1].seatX = (seatW*3/2)*(i-1) + seatW/2
		w.drawInfos[i-1].seatY = seatH / 2
		w.drawInfos[i-1].seatW = l
		w.drawInfos[i-1].seatH = l

		w.drawInfos[w.size*2+1-i].seatX = (seatW*3/2)*(i-1) + seatW/2
		w.drawInfos[w.size*2+1-i].seatY = w.height - (seatH / 2) - l
		w.drawInfos[w.size*2+1-i].seatW = l
		w.drawInfos[w.size*2+1-i].seatH = l

		w.drawInfos[i-1].titleX = (seatW*3/2)*(i-1) + seatW/2 + l/2 - 25
		w.drawInfos[i-1].titleY = seatH/2 + l/2 - 25
		w.drawInfos[i-1].titleSize = 50

		w.drawInfos[w.size*2+1-i].titleX = (seatW*3/2)*(i-1) + seatW/2 + l/2 - 25
		w.drawInfos[w.size*2+1-i].titleY = w.height - (seatH / 2) - l + l/2 - 25
		w.drawInfos[w.size*2+1-i].titleSize = 50
	}

	w.drawInfos[w.size].seatX = (seatW*3/2)*w.size + seatW/2
	w.drawInfos[w.size].seatY = w.height/2 - l/2
	w.drawInfos[w.size].seatW = l
	w.drawInfos[w.size].seatH = l

	w.drawInfos[w.size].titleX = (seatW*3/2)*w.size + seatW/2 + l/2 - 25
	w.drawInfos[w.size].titleY = w.height/2 - 25
	w.drawInfos[w.size].titleSize = 50
}
func (w *Window) draw(di DrawerInfo, str string, seatColor color.RGBA) {
	w.drawer.DrawRect(DrawPoint{di.seatX, di.seatY}, di.seatW, di.seatH, seatColor)
	w.drawer.FillRect(DrawPoint{di.titleX, di.titleY}, int(di.titleSize), int(di.titleSize), w.background)
	if str != "" {
		w.drawer.DrawText(DrawPoint{di.titleX, di.titleY}, str, di.titleSize, w.getStrColor(str))
	}
}
func (w *Window) DrawStatus(g *Game, expect string) {
	for i, di := range w.drawInfos {
		title := g.items[i].title
		if expect != "" && expect == g.items[i].title {
			title = ""
		}
		if i < g.size {
			w.draw(di, title, w.getStrColor("A"))
		} else if i > g.size {
			w.draw(di, title, w.getStrColor("1"))
		} else {
			w.draw(di, title, color.RGBA{0xff, 0x00, 0x00, 0xff})
		}

	}
	w.drawer.Show()
}

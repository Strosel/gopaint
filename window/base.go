package window

import (
	"image"
	"log"
	"os"
	"time"

	"github.com/strosel/gopaint"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
)

//Window Defines a basic Window
type Window struct {
	*gopaint.Painter

	Fps          int
	Draw         func()
	OnResize     func(size.Event)
	OnKey        func(key.Event)
	OnMouseMove  func(mouse.Event)
	OnMouseClick func(mouse.Event)

	window     screen.Window
	mainBuffer screen.Buffer
	loop       bool
}

//CreateWindow Create a New Window
func (win *Window) CreateWindow(Title string, Width, Height int) {
	var err error

	if Height < 250 {
		Height = 250
	}
	if Width < 250 {
		Width = 250
	}

	if win.Fps <= 0 {
		win.loop = false
	} else {
		win.loop = true
	}

	win.Painter = gopaint.NewPainter(image.NewNRGBA(image.Rect(0, 0, Width, Height)))

	driver.Main(func(s screen.Screen) {
		win.window, err = s.NewWindow(&screen.NewWindowOptions{
			Height: Height,
			Width:  Width,
			Title:  Title,
		})
		if err != nil {
			log.Fatal("Something went wrong creating the window; ", err)
		}
		defer win.window.Release()

		win.mainBuffer, err = s.NewBuffer(image.Point{Width, Height})
		if err != nil {
			log.Fatal("Something went wrong creating the buffer; ", err)
		}
		defer win.mainBuffer.Release()

		if win.Draw != nil {
			go func(win *Window) {
				var start time.Time
				t := time.Duration((1000 * 1000 * 1000) / win.Fps) // one Fps:th of a second
				cnt := 0
				for win.loop || cnt == 0 {
					start = time.Now()
					win.Draw()
					win.doDraw()
					cnt++

					if time.Since(start) < t {
						time.Sleep(time.Since(start) - t)
					}
				}
			}(win)
		}

		for {
			switch e := win.window.NextEvent().(type) {

			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					log.Println("Window Closed")
					os.Exit(0)
				}

			case size.Event:
				if win.OnResize != nil {
					win.OnResize(e)
				}

			case key.Event:
				if win.OnKey != nil {
					win.OnKey(e)
				}

			case mouse.Event:
				if e.Direction == mouse.DirNone {
					if win.OnMouseMove != nil {
						win.OnMouseMove(e)
					}
				} else {
					if win.OnMouseClick != nil {
						win.OnMouseClick(e)
					}
				}

			case paint.Event:
				win.window.Upload(image.Point{0, 0}, win.mainBuffer, image.Rect(0, 0, Width, Height))
				win.window.Publish()
			}
		}
	})
}

func (win *Window) doDraw() {
	drawBuffer := win.mainBuffer.RGBA()
	bounds := win.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			drawBuffer.Set(x, y, win.At(x, y))
		}
	}
	win.window.Send(paint.Event{External: true})
}

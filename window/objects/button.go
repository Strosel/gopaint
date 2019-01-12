package objects

import (
	"image"
	"image/color"

	"github.com/strosel/gopaint/window"
)

type Button struct {
	rect         image.Rectangle
	border, fill color.Color
	weight       int
	onClick      func()
}

//NewButton Create a new Button
func NewButton(rect image.Rectangle, onClick func()) Button {
	return Button{
		rect:    rect,
		border:  color.Black,
		fill:    color.Gray{127},
		weight:  3,
		onClick: onClick,
	}
}

//NewCustomButton Create a new Button with cutom styling
func NewCustomButton(rect image.Rectangle, borderColor, fillColor color.Color, borderWeight int, onClick func()) Button {
	return Button{
		rect:    rect,
		border:  borderColor,
		fill:    fillColor,
		weight:  borderWeight,
		onClick: onClick,
	}
}

//Draw Draw the button
func (b Button) Draw(win *window.Window) {
	win.Push()
	win.Rotate(0, 0, 0)
	win.Color(b.border)
	win.Fill(b.fill)
	win.Weight(b.weight)

	win.Rect(b.rect.Min.X, b.rect.Min.Y, b.rect.Dx(), b.rect.Dy())

	win.Pop()
}

//Click handle a click
func (b Button) Click(x, y int) {
	p := image.Pt(x, y)
	if p.In(b.rect) {
		b.onClick()
	}
}

package gopaint

import (
	"image/color"
	"image/draw"
	"math"
)

type Painter struct {
	draw.Image

	color, fill        color.Color
	weight, rotX, rotY int
	deg                float64
}

func NewPainter(buffer draw.Image) *Painter {
	return &Painter{
		Image:  buffer,
		color:  color.Black,
		fill:   color.White,
		weight: 1,
		rotX:   0,
		rotY:   0,
		deg:    0,
	}
}

// INTERNAL FUNCTIONS

func colorF(clr color.Color) (float64, float64, float64, float64) {
	colorR, colorG, colorB, colorA := clr.RGBA()
	alpha := float64(colorA%256) / 255
	return float64(colorR%256) * alpha, float64(colorG%256) * alpha, float64(colorB%256) * alpha, alpha
}

// INTERNAL METHODS

func (p Painter) set(x, y int, fill bool) {
	//Alpha handling
	bottomColorR, bottomColorG, bottomColorB, _ := colorF(p.At(x, y))
	var topColorR, topColorG, topColorB, topColorA float64
	if !fill {
		topColorR, topColorG, topColorB, topColorA = colorF(p.color)
	} else {
		topColorR, topColorG, topColorB, topColorA = colorF(p.fill)
	}

	finalColor := color.RGBA{
		R: uint8(topColorR + bottomColorR*(1-topColorA)),
		G: uint8(topColorG + bottomColorG*(1-topColorA)),
		B: uint8(topColorB + bottomColorB*(1-topColorA)),
	}

	p.Set(x, y, finalColor)
}

func (p Painter) doFill(x, y, w, h int, check func(int, int) bool) {
	for fx := x; fx <= x+w; fx++ {
		for fy := y; fy <= y+h; fy++ {
			if check(fx, fy) {
				p.set(fx, fy, true)
			}
		}
	}
}

func (p Painter) doRotate(x, y int) (int, int) {
	// rotation handling
	s := math.Sin(p.deg)
	c := math.Cos(p.deg)
	x -= p.rotX
	y -= p.rotY
	xnew := float64(x)*c - float64(y)*s
	ynew := float64(x)*s + float64(y)*c
	x = int(xnew) + p.rotX
	y = int(ynew) + p.rotY
	return x, y
}

// EXPORTED METHODS

//Color Set the line color
func (p *Painter) Color(c color.Color) {
	p.color = c
}

//Fill Set the fill color
func (p *Painter) Fill(c color.Color) {
	p.fill = c
}

//Weight Set the line weight
func (p *Painter) Weight(w int) {
	p.weight = w
}

//Rotate Set the rotation transform values, rotate deg degrees around (x,y)
func (p *Painter) Rotate(x, y int, deg float64) {
	p.rotX = x
	p.rotY = y
	p.deg = deg
}

//Background Set the background color. DOES NOT INTERPRET ALPHA
func (p Painter) Background(clr color.Color) {
	bounds := p.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			p.Set(x, y, clr)
		}
	}
}

//Line Draw a line from (x0,y0) to (x1,y1)
func (p Painter) Line(x0, y0, x1, y1 int) {
	x0, y0 = p.doRotate(x0, y0)
	x1, y1 = p.doRotate(x1, y1)

	l := makeLine(x0, y0, x1, y1)

	mod := 0
	for w := 1; w <= p.weight; w++ {
		if w%2 == 0 {
			mod++
		}
		mod *= -1

		if !l.v {
			sign := 1
			if l.k <= 1 && l.k >= -1 {
				if x0 > x1 {
					sign = -1
				}
				for x := x0; x != x1; x += sign {
					y := int(l.k*float64(x)+l.m) + mod
					p.set(x, y, false)
				}
			} else {
				if y0 > y1 {
					sign = -1
				}
				for y := y0; y != y1; y += sign {
					x := float64(y-int(l.m)) / l.k
					p.set(int(x)-mod, y, false)
				}
			}
		} else {
			sign := 1
			if y0 > y1 {
				sign = -1
			}
			for y := y0; y != y1; y += sign {
				p.set(x0+mod, y, false)
			}
		}
	}
}

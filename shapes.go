package gopaint

import (
	"math"
)

//Ellipse Draws an ellipse with center (x,y)
func (p Painter) Ellipse(x, y, w, h int) {
	wr := float64(w) / 2
	hr := float64(h) / 2
	px := x + int(wr*math.Sin(0))
	py := y + int(hr*math.Cos(0))

	for a := 0.1; a <= 360; a += 0.1 {
		cx := x + int(wr*math.Sin(a))
		cy := y + int(hr*math.Cos(a))
		p.Line(px, py, cx, cy)
		px, py = cx, cy
	}

	buff := float64(p.weight) / 2.0
	wr -= buff
	hr -= buff
	p.doFill(x-int(wr), y-int(hr), w, h, func(ix, iy int) bool {
		return (math.Pow(float64(ix-x), 2)/math.Pow(wr, 2))+(math.Pow(float64(iy-y), 2)/math.Pow(hr, 2)) < 1
	})
}

//Rect Draw a Rectangle with opposing corners (x0,y0) to (x1,y1)
func (p Painter) Rect(x, y, w, h int) {
	x1 := x + w
	y1 := y + h
	buff := int(math.Round(float64(p.weight) / 2.0))
	p.Line(x, y+buff, x, y1-buff+1)
	p.Line(x1, y+buff, x1, y1-buff+1)
	p.Line(x-buff+1, y, x1+buff, y)
	p.Line(x-buff+1, y1, x1+buff, y1)

	p.doFill(x+buff, y+buff, w-p.weight-1, h-p.weight-1, func(ix, iy int) bool {
		return true
	})
}

//Polyline Draw a Polyline through points (x[n],y[n])
func (p Painter) Polyline(x, y []int) {
	if len(x) != len(y) {
		panic("Coordinate-list length error")
	}

	for i := 0; i < len(x)-1; i++ {
		p.Line(x[i], y[i], x[i+1], y[i+1])
	}
}

//Polygon Draw a Polygon with corners (x[n],y[n])
func (p Painter) Polygon(x, y []int) {
	if len(x) != len(y) {
		panic("Coordinate-list length error")
	}

	if x[0] != x[len(x)-1] && y[0] != y[len(y)-1] {
		x = append(x, x[0])
		y = append(y, y[0])
	}

	lines := []line{}
	for i := 0; i < len(x)-1; i++ {
		p.Line(x[i], y[i], x[i+1], y[i+1])
		lines = append(lines, makeLine(x[i], y[i], x[i+1], y[i+1]))
	}

	bounds := p.Bounds()
	xmin, ymin, xmax, ymax := bounds.Max.X, bounds.Max.Y, bounds.Min.X, bounds.Min.Y
	for i := range x {
		if x[i] < xmin {
			xmin = x[i]
		} else if x[i] > xmax {
			xmax = x[i]
		}

		if y[i] < ymin {
			ymin = y[i]
		} else if y[i] > ymax {
			ymax = y[i]
		}
	}

	p.doFill(xmin, ymin, xmax-xmin, ymax-ymin, func(ix, iy int) bool {
		flat := makeLine(-1, iy, ix, iy)
		intersects := 0
		for _, l := range lines {
			if flat.intersects(l) {
				intersects++
			}
		}

		return intersects%2 != 0
	})
}

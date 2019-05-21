package gopaint

import (
	"image"
	"image/color"
)

type line struct {
	b, e image.Point //beginning, end
	k, m float64     //k, m from y = kx + m
	v    bool        //vertical?
}

// func makeLine(x0, y0, x1, y1 int) line {
// 	l := line{
// 		b: image.Point{x0, y0},
// 		e: image.Point{x1, y1},
// 	}
// 	return l
// }

func (l line) intersects(l2 line) bool {
	tNum := (l.b.X-l2.b.X)*(l2.b.Y-l2.e.Y) - (l.b.Y-l2.b.Y)*(l2.b.X-l2.e.X)
	uNum := (l.b.X-l.e.X)*(l.b.Y-l2.b.Y) - (l.b.Y-l.e.Y)*(l.b.X-l2.b.X)
	denom := (l.b.X-l.e.X)*(l2.b.Y-l2.e.Y) - (l.b.Y-l.e.Y)*(l2.b.X-l2.e.X)

	if denom == 0 {
		return false
	}

	t := tNum / denom
	u := -uNum / denom

	return 0.0 <= t && t <= 1.0 && 0.0 <= u && u <= 1.0
}

func makeLine(x0, y0, x1, y1 int) line {
	l := line{
		b: image.Point{x0, y0},
		e: image.Point{x1, y1},
	}
	if x0 != x1 {
		l.k = float64(y1-y0) / float64(x1-x0)
		l.m = float64(y0) - l.k*float64(x0)
		l.v = false
	} else {
		l.k = 0
		l.m = 0
		l.v = true
	}
	return l
}

// func (l line) on(x int) bool {
// 	if l.b.X < l.e.X {
// 		return x > l.b.X && x < l.e.X
// 	}
// 	return x < l.b.X && x > l.e.X
// }

// func (l line) intersects(l2 line) bool {
// 	var x int
// 	if l.v {
// 		x = l.b.X
// 	} else if l2.v {
// 		x = l2.b.X
// 	} else if l.k == l2.k {
// 		return false
// 	} else {
// 		dm := l.m - l2.m
// 		dk := l2.k - l.k
// 		x = int(dm / dk)
// 	}

// 	return l.on(x) && l2.on(x)
// }

type style struct {
	color, fill        color.Color
	weight, rotX, rotY int
	deg                float64
}

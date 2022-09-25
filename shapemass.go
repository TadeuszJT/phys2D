package phys2D

import (
	geom "github.com/tadeuszjt/geom/32"
)

/* Masses of common shapes */
func MassRectangle(r geom.Rect) geom.Ori2 {
	w := r.Width()
	h := r.Height()
	mass := w * h
	angM := geom.Angle((mass * (w*w + h*h)) / 12)
	return geom.Ori2{mass, mass, angM}
}

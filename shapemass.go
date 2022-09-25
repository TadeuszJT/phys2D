package phys2D

import (
	"github.com/tadeuszjt/geom/64"
)

/* Masses of common shapes */
func MassRectangle(r geom.Rect) geom.Ori2 {
	w := r.Width()
	h := r.Height()
	mass := w*h
	angM := (mass * (w*w + h*h)) / 12
	return geom.Ori2{mass, mass, angM}
}

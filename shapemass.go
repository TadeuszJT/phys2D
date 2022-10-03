package phys2D

import (
	"github.com/tadeuszjt/geom/generic"
)

/* Masses of common shapes */
func MassRectangle(r geom.Rect[float64]) geom.Ori2[float64] {
	w := r.Width()
	h := r.Height()
	mass := w * h
	angM := (mass * (w*w + h*h)) / 12
	return geom.Ori2[float64]{mass, mass, angM}
}

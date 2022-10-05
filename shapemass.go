package phys2D

import (
	"github.com/tadeuszjt/geom/generic"
)

// from wikipedia: drag coefficient
const (
	DragCoefficientSphere   = 0.47
	DragCoefficientCube     = 1.05
	DragCoefficientTeardrop = 0.04
)

/* Masses of common shapes */

/*
moment of inertia about the x axis: Ix = bh^3 / 12
moment of inertia about the y axis: Iy = b^3h / 12
Perpedicular axis theorem states that Iz = Ix + Iy
*/
func MassRectangle(r geom.Rect[float64]) geom.Ori2[float64] {
	w := r.Width()
	h := r.Height()

	mass := w * h
	angM := (mass * (w*w + h*h)) / 12
	return geom.Ori2[float64]{mass, mass, angM}
}

//func DragRectangle(r geom.Rect[float64]) geom.Ori2[float64] {
//    ACdX := r.Height() * DragCoefficientCube
//    ACdY := r.Width() * DragCoefficientCube
//}

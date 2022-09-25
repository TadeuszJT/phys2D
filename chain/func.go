package main

import (
	geom "github.com/tadeuszjt/geom/64"
)

func addBox(pos geom.Vec2, w, h float64) {
	orientation := geom.Ori2{pos, 0}
	rectangle := geom.RectCentre(w, h)

	// body
	m := w * h
	I := (m * (w*w + h*h)) / 12
	mass := geom.Ori2{geom.Vec2{m, m}, I}
	bodies.AddBody(orientation, mass)

	// sprite
	sprites.AddSprite(orientation, rectangle)
}

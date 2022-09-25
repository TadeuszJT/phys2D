package main

import (
	geom "github.com/tadeuszjt/geom/64"
)

var view struct {
	zoom        float64
	orientation geom.Ori2
	window, gl  geom.Rect
}

func init() {
	view.zoom = 1.0
	view.orientation = geom.Ori2{geom.Vec2{100, 0}, 0}
	view.gl = geom.RectCentre(2, -2)
}

func windowToWorld() geom.Mat33 {
	cam := geom.RectCentre(
		view.window.Width()*view.zoom,
		view.window.Height()*view.zoom,
	)
	return transforms.Cam2DToCamera(cam, view.window, view.orientation)
}

func worldToGl() geom.Mat33 {
	cam := geom.RectCentre(
		view.window.Width()*view.zoom,
		view.window.Height()*view.zoom,
	)
	return transforms.Cam2DToDisplay(cam, view.gl, view.orientation)
}

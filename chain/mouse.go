package main

import (
	geom "github.com/tadeuszjt/geom/64"
)

var mouse struct {
	window, world geom.Vec2
	held          bool
	rectHeld      int
	rectOffset    geom.Vec2
}

// -1 no box
func boxAt(p geom.Vec2) int {
	for i := range sprites.Rec {
		pLocal := sprites.Ori[i].ToLocal(p)
		if sprites.Rec[i].Contains(pLocal) {
			return i
		}
	}
	return -1
}


func mouseCB(x, y, scroll float64, event uint8) {
	switch event {
	case smeg.MousePress:
		mouse.held = true
		mouse.rectHeld = boxAt(mouse.world)
		if mouse.rectHeld > -1 {
			mouse.rectOffset = sprites.Ori[mouse.rectHeld].ToLocal(mouse.world)
		}

	case smeg.MouseRelease:
		mouse.held = false

	case smeg.MouseMove:
		prev := mouse.window
		mouse.window = geom.Vec2{x, y}
		if mouse.held {
			if mouse.rectHeld < 0 { // world held
				view.orientation.PlusEquals(
					geom.Ori2{
						windowToWorld().TimesVec2(prev.Minus(mouse.window), 0),
						0,
					},
				)
			} else { // box held

			}
		}
		mouse.world = windowToWorld().TimesVec2(mouse.window, 1)

	case smeg.MouseScroll:
		view.zoom *= (1 - scroll*.04)
		mNew := windowToWorld().TimesVec2(mouse.window, 1)
		view.orientation.PlusEquals(geom.Ori2{mouse.world.Minus(mNew), 0})
	}
}

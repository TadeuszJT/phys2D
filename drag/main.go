package main

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/tadeuszjt/data"
	"github.com/tadeuszjt/geom/generic"
	"github.com/tadeuszjt/gfx"
	"github.com/tadeuszjt/phys2D"
)

var (
	camRect geom.Rect[float32]

	poly = geom.Poly[float32]{ // centroid is origin after init
		{0, 0},
		{2, -.5},
		{2.5, 0},
		{2, .5},
	}
	polyDrawData = []float32{}
	polyColour   = gfx.Red
	polyKey      data.Key

	timeStep = 1. / 60.

	world               = phys2D.NewWorld()
	mousePos            geom.Vec2[float32]
	mousePolyHeld       bool
	mousePolyHeldOffset geom.Vec2[float32]
)

func init() {
	centroid := geom.PolyCentroid(poly)
	for i := range poly {
		poly[i] = poly[i].Minus(centroid)
	}

	polyDrawData = []float32{
		poly[0].X, poly[0].Y, 0, 0, polyColour.R, polyColour.G, polyColour.B, polyColour.A,
		poly[1].X, poly[1].Y, 0, 0, polyColour.R, polyColour.G, polyColour.B, polyColour.A,
		poly[2].X, poly[2].Y, 0, 0, polyColour.R, polyColour.G, polyColour.B, polyColour.A,
		poly[0].X, poly[0].Y, 0, 0, polyColour.R, polyColour.G, polyColour.B, polyColour.A,
		poly[2].X, poly[2].Y, 0, 0, polyColour.R, polyColour.G, polyColour.B, polyColour.A,
		poly[3].X, poly[3].Y, 0, 0, polyColour.R, polyColour.G, polyColour.B, polyColour.A,
	}

	world.AirDensity = 0.2
	world.Gravity = geom.Ori2[float64]{0, 1, 0}
	massXY := float64(geom.PolyArea(poly))
	massTheta := float64(geom.PolyMomentOfInertia(poly))
	fmt.Println(massXY, massTheta)
	polyOri := geom.Ori2[float32]{4, 4, 1.1}
	ori64 := geom.Ori2Convert[float32, float64](polyOri)
	polyKey = world.AddBody(ori64, geom.Ori2[float64]{massXY, massXY, massTheta})

	for i := range poly {
		if i == (len(poly) - 1) {
			world.AddDragPlate(
				polyKey,
				geom.Vec2Convert[float32, float64](poly[i]),
				geom.Vec2Convert[float32, float64](poly[0]),
			)
		} else {
			world.AddDragPlate(
				polyKey,
				geom.Vec2Convert[float32, float64](poly[i]),
				geom.Vec2Convert[float32, float64](poly[i+1]),
			)
		}
	}
}

func setup(w *gfx.Win) error {
	winSize := w.Size().ScaledBy(0.02)
	camRect = geom.RectOrigin[float32](winSize.X, winSize.Y)

	return nil
}

func draw(w *gfx.Win, c gfx.Canvas) {
	polyOri := geom.Ori2Convert[float64, float32](world.GetOrientation(polyKey))
	modelMat := polyOri.Mat3Transform()
	winSize := w.Size()
	winRect := geom.RectOrigin(winSize.X, winSize.Y)
	camMat := geom.Mat3Camera2D(camRect, winRect)

	if mousePolyHeld {
		displayToCam := geom.Mat3Camera2D(winRect, camRect)
		mouseWorld := displayToCam.TimesVec2(mousePos, 1).Vec2()
		start := modelMat.TimesVec2(mousePolyHeldOffset, 1).Vec2()
		end := mouseWorld
		gfx.Draw2DArrow(c, start, end, gfx.Green, 0.1, camMat)

		f := geom.Vec2Convert[float32, float64](end.Minus(start))
		offset := geom.Vec2Convert[float32, float64](mousePolyHeldOffset.RotatedBy(polyOri.Theta))

		world.ApplyImpulse(polyKey, f, offset, timeStep)
	}

	world.Update(timeStep)

	mat := camMat.Product(modelMat)

	c.Draw2DVertexData(polyDrawData, nil, &mat)
}

func mouse(w *gfx.Win, ev gfx.MouseEvent) {
	switch e := ev.(type) {
	case gfx.MouseMove:
		mousePos = e.Position
	case gfx.MouseButton:
		if e.Action == glfw.Press {
			if e.Button == glfw.MouseButtonLeft {
				winSize := w.Size()
				winRect := geom.RectOrigin(winSize.X, winSize.Y)
				displayToCam := geom.Mat3Camera2D(winRect, camRect)
				polyOri := geom.Ori2Convert[float64, float32](world.GetOrientation(polyKey))

				trans := geom.Mat3Translation(polyOri.Vec2().ScaledBy(-1))
				rot := geom.Mat3Rotation(-polyOri.Theta)
				displayToModel := rot.Product(trans).Product(displayToCam)
				mouseModel := displayToModel.TimesVec2(mousePos, 1).Vec2()

				if poly.Contains(mouseModel) {
					mousePolyHeld = true
					mousePolyHeldOffset = mouseModel
				}

			}
		} else if e.Action == glfw.Release {
			if e.Button == glfw.MouseButtonLeft {
				mousePolyHeld = false
			}
		}

	default:
	}
}

func main() {
	gfx.RunWindow(gfx.WinConfig{
		DrawFunc:  draw,
		SetupFunc: setup,
		MouseFunc: mouse,
		Width:     1024,
		Height:    800,
		Title:     "Chain",
	})

	fmt.Println("benis")
}

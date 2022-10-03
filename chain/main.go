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
	world = phys2D.NewWorld()

	rectSize = geom.RectCentred[float32](20, 20)
	numRects = 10
	timeStep = 4. / 60.

	rects struct {
		data.KeyMap
		physKeys data.RowT[data.Key]
		colours  data.RowT[gfx.Colour]
	}

	jointKeys []data.Key
	rectKeys  []data.Key

	mousePos = geom.Vec2[float32]{}
	rectHeld = data.Key(-1)
)

func init() {
	rects.KeyMap.Row = data.Table{&rects.physKeys, &rects.colours}
}

func setup(w *gfx.Win) error {
	for i := 0; i < numRects; i++ {
		rectSize64 := geom.RectConvert[float32, float64](rectSize)

		mass := phys2D.MassRectangle(rectSize64)
		if i == 0 || i == (numRects-1) {
			mass = geom.Ori2[float64]{0, 0, 0}
		}
		ori := geom.Ori2[float64]{300 + float64(i)*20, 60 + float64(i)*20, 0}

		rectKeys = append(
			rectKeys,
			rects.Append(world.AddBody(ori, mass), gfx.ColourRand()),
		)

		if i > 0 {
			jointKeys = append(
				jointKeys,
				world.AddJoint(
					rects.physKeys[i-1],
					rects.physKeys[i],
					rectSize64.Max,
					rectSize64.Min,
				),
			)
		}
	}

	world.DeleteJoint(jointKeys[len(jointKeys)/2])

	return nil
}

func draw(w *gfx.Win, c gfx.Canvas) {
	for i := range rects.physKeys {
		ori := world.GetOrientations(rects.physKeys[i])
		ori32 := geom.Ori2Convert[float64, float32](ori[0])
		colour := rects.colours[i]
		gfx.DrawSprite(c, ori32, rectSize, colour, nil, nil)
	}

	if rectHeld != (-1) {
		index := rects.KeyMap.KeyToIndex[rectHeld]
		ori1 := world.GetOrientations(rects.physKeys[index])

		mousePos64 := geom.Vec2Convert[float32, float64](mousePos)
		delta := mousePos64.Minus(ori1[0].Vec2())
		deltaOri2 := geom.Ori2[float64]{delta.X, delta.Y, 0}

		world.ApplyImpulse(rects.physKeys[index], deltaOri2.ScaledBy(10))

	}

	world.Update(timeStep)
}

func mouse(w *gfx.Win, ev gfx.MouseEvent) {
	switch e := ev.(type) {
	case gfx.MouseMove:
		mousePos = e.Position
	case gfx.MouseButton:
		if e.Button == glfw.MouseButtonLeft {
			if e.Action == glfw.Press {

				for _, key := range rectKeys {
					i := rects.KeyMap.KeyToIndex[key]

					ori := world.GetOrientations(rects.physKeys[i])
					ori32 := geom.Ori2Convert[float64, float32](ori[0])

					trans := geom.Mat3Translation(ori32.Vec2().ScaledBy(-1))
					rot := geom.Mat3Rotation(ori32.Theta * (-1))
					mat := rot.Product(trans)

					v := mat.TimesVec2(mousePos, 1).Vec2()

					if rectSize.Contains(v) {
						rectHeld = key
					}
				}

			} else if e.Action == glfw.Release {
				rectHeld = data.Key(-1)
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
	})

	fmt.Println("benis")
}

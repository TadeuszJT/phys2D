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
	rectSize = geom.RectCentred[float32](20, 20)
	numRects = 16
	timeStep = 1. / 60.
	timeMul  = 10

	world    = phys2D.NewWorld()
	mousePos = geom.Vec2[float32]{}

	rects struct {
		data.KeyMap
		physKeys data.RowT[data.Key]
		colours  data.RowT[gfx.Colour]
	}

	rectKeys data.RowT[data.Key]
	rectHeld = data.KeyInvalid
)

func init() {
	rects.KeyMap = data.MakeKeyMap(data.Table{&rects.physKeys, &rects.colours})
}

func deleteRect(keyIndex int) {
	if rectHeld == rectKeys[keyIndex] {
		rectHeld = data.KeyInvalid
	}

	rectIndex := rects.GetIndex(rectKeys[keyIndex])
	world.DeleteBody(rects.physKeys[rectIndex])
	rects.Delete(rectKeys[keyIndex])
	rectKeys.Delete(keyIndex)
}

func setup(w *gfx.Win) error {
	for i := 0; i < numRects; i++ {
		rectSize64 := geom.RectConvert[float32, float64](rectSize)

		mass := phys2D.MassRectangle(rectSize64)
		if i == 0 || i == (numRects-1) {
			mass = geom.Ori2[float64]{0, 0, 0}
		}
		ori := geom.Ori2[float64]{300 + float64(i)*20, 60 + float64(i)*20, 0}

		rectKeys.Append(rects.Append(world.AddBody(ori, mass), gfx.ColourRand()))

		if i > 0 {
			world.AddJoint(
				rects.physKeys[i-1],
				rects.physKeys[i],
				rectSize64.Max,
				rectSize64.Min,
			)
		}
	}

	deleteRect(len(rectKeys) / 2)

	return nil
}

func draw(w *gfx.Win, c gfx.Canvas) {
	for i := 0; i < timeMul; i++ {
		if rectHeld != data.KeyInvalid {
			index := rects.GetIndex(rectHeld)
			ori1 := world.GetOrientation(rects.physKeys[index])

			mousePos64 := geom.Vec2Convert[float32, float64](mousePos)
			delta := mousePos64.Minus(ori1.Vec2())

			world.ApplyImpulse(
				rects.physKeys[index],
				delta.ScaledBy(100),
				geom.Vec2[float64]{},
				timeStep,
			)

		}

		world.Update(timeStep)
	}

	for i := range rects.physKeys {
		ori := world.GetOrientation(rects.physKeys[i])
		ori32 := geom.Ori2Convert[float64, float32](ori)
		colour := rects.colours[i]
		gfx.DrawSprite(c, ori32, rectSize, colour, nil, nil)
	}
}

func mouse(w *gfx.Win, ev gfx.MouseEvent) {
	switch e := ev.(type) {
	case gfx.MouseMove:
		mousePos = e.Position
	case gfx.MouseButton:
		if e.Action == glfw.Press {
			for _, key := range rectKeys {
				i := rects.GetIndex(key)

				ori := world.GetOrientation(rects.physKeys[i])
				ori32 := geom.Ori2Convert[float64, float32](ori)
				trans := geom.Mat3Translation(ori32.Vec2().ScaledBy(-1))
				rot := geom.Mat3Rotation(ori32.Theta * (-1))
				mat := rot.Product(trans)
				v := mat.TimesVec2(mousePos, 1)

				if rectSize.Contains(v) {
					if e.Button == glfw.MouseButtonLeft {
						rectHeld = key
					} else if e.Button == glfw.MouseButtonRight {
						deleteRect(i)
					}

					break
				}
			}

		} else if e.Action == glfw.Release {
			if e.Button == glfw.MouseButtonLeft {
				rectHeld = data.KeyInvalid
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

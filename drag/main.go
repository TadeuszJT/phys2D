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
	rectSize = geom.RectCentred[float32](200, 50)
	numRects = 1
	timeStep = 1. / 60.
	timeMul  = 4

	world    = phys2D.NewWorld()
	mousePos geom.Vec2[float32]

	rects struct {
		data.KeyMap
		physKeys data.RowT[data.Key]
		colours  data.RowT[gfx.Colour]
	}

	rectKeys       data.RowT[data.Key]
	rectHeld       = data.KeyInvalid
	rectHeldOffset geom.Vec2[float64]
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
	world.AirDensity = 0.2

	for i := 0; i < numRects; i++ {
		rectSize64 := geom.RectConvert[float32, float64](rectSize)

		mass := phys2D.MassRectangle(rectSize64)
		ori := geom.Ori2[float64]{200, 200, 0}

		bodyKey := world.AddBody(ori, mass)

		rectKeys.Append(rects.Append(bodyKey, gfx.ColourRand()))

		verts := [4]geom.Vec2[float64]{
			{rectSize64.Min.X, rectSize64.Min.Y}, // top left
			{rectSize64.Max.X, rectSize64.Min.Y}, // top right
			{rectSize64.Min.X, rectSize64.Max.Y}, // bottom left
			{rectSize64.Max.X, rectSize64.Max.Y}, // bottom right
		}

		world.AddDragPlate(bodyKey, verts[0], verts[1])
		world.AddDragPlate(bodyKey, verts[1], verts[3])
		world.AddDragPlate(bodyKey, verts[3], verts[2])
		world.AddDragPlate(bodyKey, verts[2], verts[0])

	}

	return nil
}

func draw(w *gfx.Win, c gfx.Canvas) {
	for i := 0; i < timeMul; i++ {
		if rectHeld != data.KeyInvalid {
			index := rects.GetIndex(rectHeld)
			ori1 := world.GetOrientations(rects.physKeys[index])
			ori := ori1[0]

			mousePos64 := geom.Vec2Convert[float32, float64](mousePos)
			offset := rectHeldOffset.RotatedBy(ori.Theta)
			force := mousePos64.Minus(ori.Vec2().Plus(offset))

			world.ApplyImpulse(rects.physKeys[index], force.ScaledBy(1000), offset, timeStep)

		}

		world.Update(timeStep)


//        for _, imp := range phys2D.Impulses {
//            ori := world.GetOrientations(imp.Key)
//
//            start := ori[0].Vec2().Plus(imp.Offset)
//            start32 := geom.Vec2Convert[float64, float32](start)
//            force32 := geom.Vec2Convert[float64, float32](imp.F)
//            mat := geom.Mat3Identity[float32]()
//
//            gfx.Draw2DArrow(c, start32, start32.Plus(force32.ScaledBy(0.02)), gfx.Blue, 4, mat)
//        }
//        phys2D.Impulses = []phys2D.Impulse{}
	}

	for i := range rects.physKeys {
		ori := world.GetOrientations(rects.physKeys[i])
		ori32 := geom.Ori2Convert[float64, float32](ori[0])
		colour := rects.colours[i]
		gfx.DrawSprite(c, ori32, rectSize, colour, nil, nil)
	}

	if rectHeld != data.KeyInvalid {
		index := rects.GetIndex(rectHeld)
		ori1 := world.GetOrientations(rects.physKeys[index])
		ori := ori1[0]

		offset := rectHeldOffset.RotatedBy(ori.Theta)
		start := geom.Vec2Convert[float64, float32](ori.Vec2().Plus(offset))
		end := mousePos
		mat := geom.Mat3Identity[float32]()

		gfx.Draw2DArrow(c, start, end, gfx.Red, 10, mat)
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

				oris := world.GetOrientations(rects.physKeys[i])
				ori := oris[0]
				rectSize64 := geom.RectConvert[float32, float64](rectSize)
				mousePos64 := geom.Vec2Convert[float32, float64](mousePos)

				trans := geom.Mat3Translation(ori.Vec2().ScaledBy(-1))
				rot := geom.Mat3Rotation(ori.Theta * (-1))
				mat := rot.Product(trans)
				v := mat.TimesVec2(mousePos64, 1).Vec2()

				if rectSize64.Contains(v) {
					if e.Button == glfw.MouseButtonLeft {
						rectHeld = key
						rectHeldOffset = v
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

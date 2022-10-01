package main

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/tadeuszjt/data"
	"github.com/tadeuszjt/geom/64"
	"github.com/tadeuszjt/gfx"
	phys2D "github.com/tadeuszjt/phys2D"

	geom32 "github.com/tadeuszjt/geom/32"
)

var (
	world = phys2D.NewWorld()

	rectSize = geom.Rect{geom.Vec2{-10, -10}, geom.Vec2{10, 10}}
	numRects = 10

    rects struct {
        data.Table
        physKeys data.RowT[data.Key]
        colours  data.RowT[gfx.Colour]
    }


	jointKeys []data.Key

	mousePos  = geom.Vec2{}
	mouseHeld = false
)

func init() {
    rects.Table = data.Table { &rects.physKeys, &rects.colours }
}

func setup(w *gfx.Win) error {
	for i := 0; i < numRects; i++ {
		mass := phys2D.MassRectangle(rectSize)
		if i == 0 || i == (numRects-1) {
			mass = geom.Ori2{0, 0, 0}
		}

		ori := geom.Ori2{ 300 + float64(i)*20, 60 + float64(i)*20, 0 }

        rectPhysKey := world.AddBody(ori, mass)

        rects.Append(rectPhysKey, gfx.ColourRand())

		if i > 0 {
			jointKeys = append(
				jointKeys,
				world.AddJoint(
                    rects.physKeys[i-1],
                    rects.physKeys[i],
                    rectSize.Max,
                    rectSize.Min,
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
        ori32 := geom32.Ori2{float32(ori[0].X), float32(ori[0].Y), geom32.Angle(ori[0].Theta) }
        colour := rects.colours[i]
		rect32 := geom32.Rect{
			geom32.Vec2{float32(rectSize.Min.X), float32(rectSize.Min.Y)},
			geom32.Vec2{float32(rectSize.Max.X), float32(rectSize.Max.Y)},
		}
		gfx.DrawSprite(c, ori32, rect32, colour, nil, nil)
    }


	world.Update(4 / 60.) // 60 fps
	world.Update(4 / 60.) // 60 fps
	world.Update(4 / 60.) // 60 fps
	world.Update(4 / 60.) // 60 fps
}

func mouse(w *gfx.Win, ev gfx.MouseEvent) {
	switch e := ev.(type) {
	case gfx.MouseMove:
		mousePos.X = float64(e.Position.X)
		mousePos.Y = float64(e.Position.Y)
	case gfx.MouseButton:
		if e.Button == glfw.MouseButtonLeft {
			if e.Action == glfw.Press {
				mouseHeld = true

				//                for i := range rectKeys {
				//                    ori := phys2D.GetOrientations(rectKeys[i])
				//                }

			} else if e.Action == glfw.Release {
				mouseHeld = false
			}
		}

		fmt.Println(mouseHeld)

	default:
	}
}

//type MouseButton struct {
//	Button glfw.MouseButton
//	Action glfw.Action
//	Mods   glfw.ModifierKey
//}

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

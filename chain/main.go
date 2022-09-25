package main

import (
    "fmt"
	"github.com/tadeuszjt/geom/64"
	"github.com/tadeuszjt/phys2D"
    "github.com/tadeuszjt/gfx"

    geom32 "github.com/tadeuszjt/geom/32"
)

var (
	world phys2D.World
    rectSize = geom.Rect{geom.Vec2{-10, -10}, geom.Vec2{10, 10}}
    rectKeys []phys2D.Key
    numRects = 7
)

func setup(w *gfx.Win) error {
    for i := 0; i < numRects; i++ {
        mass := phys2D.MassRectangle(rectSize)
        if i == 0 {
            mass = geom.Ori2{0, 0, 0}
        }

        ori := geom.Ori2 {
            300 + float64(i) * 20,
            60 + float64(i)*20,
            0,
        }

        rectKeys = append(rectKeys, world.AddBody(ori, mass))
        if i > 0 {
            world.AddJoint(rectKeys[i-1], rectKeys[i], geom.Vec2{10, 10}, geom.Vec2{-10, -10})
        }
    }

    world.Gravity = geom.Ori2{0, 10, 0}

    return nil
}


func draw(w* gfx.Win, c gfx.Canvas) {
    oris := world.GetOrientations(rectKeys)

    orisConverted := []geom32.Ori2{}
    for _, o := range oris {
        orisConverted = append(orisConverted, geom32.Ori2{float32(o.X), float32(o.Y), geom32.Angle(o.Theta)})
    }

    for i := range orisConverted {
        gfx.DrawSprite(c, orisConverted[i], geom32.RectCentred(20, 20), gfx.Red, nil, nil)
    }

    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
    world.Update(1 / 60.) // 60 fps
}

func main() {
    gfx.RunWindow(gfx.WinConfig{
        DrawFunc: draw,
        SetupFunc: setup,
        Width: 1024,
        Height: 800,
    })


    fmt.Println("benis")
}

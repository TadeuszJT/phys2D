package main

import (
    "fmt"
	geom "github.com/tadeuszjt/geom/64"
	"github.com/tadeuszjt/phys2D"
    "github.com/tadeuszjt/gfx"

    newGeom "github.com/tadeuszjt/geom/32"
)

var (
	world phys2D.World
    rectKeys []phys2D.Key

    rectSize = geom.Rect{geom.Vec2{-10, -10}, geom.Vec2{10, 10}}
    rectStartPositions = []geom.Ori2 {
        {100, 60, 0},
        {120, 80, 0},
        {140, 100, 0},
        {160, 120, 0},
        {180, 140, 0},
        {200, 160, 0},
        {220, 180, 0},
    }
)

func setup(w *gfx.Win) error {

    for i, p := range rectStartPositions {
        mass := phys2D.MassRectangle(rectSize)

        if i == 0 {
            mass = geom.Ori2{0, 0, 0}
        }

        rectKeys = append(rectKeys, world.AddBody(p, mass))
    }

    world.Gravity = geom.Ori2{0, 10, 0}

    world.AddJoint(rectKeys[0], rectKeys[1], geom.Vec2{10, 10}, geom.Vec2{-10, -10})
    world.AddJoint(rectKeys[1], rectKeys[2], geom.Vec2{10, 10}, geom.Vec2{-10, -10})
    world.AddJoint(rectKeys[2], rectKeys[3], geom.Vec2{10, 10}, geom.Vec2{-10, -10})
    world.AddJoint(rectKeys[3], rectKeys[4], geom.Vec2{10, 10}, geom.Vec2{-10, -10})
    world.AddJoint(rectKeys[4], rectKeys[5], geom.Vec2{10, 10}, geom.Vec2{-10, -10})
    world.AddJoint(rectKeys[5], rectKeys[6], geom.Vec2{10, 10}, geom.Vec2{-10, -10})

    return nil
}


func draw(w* gfx.Win, c gfx.Canvas) {
    oris := world.GetOrientations(rectKeys)

    orisConverted := []newGeom.Ori2{}

    for _, o := range oris {
        orisConverted = append(orisConverted, newGeom.Ori2{float32(o.X), float32(o.Y), newGeom.Angle(o.Theta)})
    }

    for i := range orisConverted {
        gfx.DrawSprite(c, orisConverted[i], newGeom.RectCentred(20, 20), gfx.Red, nil, nil)
    }

    world.Update(0.1)
    fmt.Println("benis")
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

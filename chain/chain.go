package main

import (
	geom "github.com/tadeuszjt/geom/64"
	"github.com/tadeuszjt/phys2D"
)

var (
	// systems
	sprites sprite.System
	bodies  phys2D.System
)

func main() {
	// make boxes
	for i := 0; i < 40; i++ {
		pos := geom.Vec2{float64(i*10 + 100), 0}
		addBox(pos, 10, 10)
		if i > 0 {
			bodies.AddJoint(i-1, i, geom.Vec2{5, 0}, geom.Vec2{-5, 0})
		}
	}

	bodies.InvMass[0] = geom.Ori2{}
	bodies.Gravity.Y = 100.

	smeg.Init()
	smeg.SetMouseCB(mouseCB)
	boxTex, _ := smeg.AddTexture("box.png")

	for smeg.IsOpen() {
		smeg.PollEvents()

		// set view
		w, h := smeg.GetSize()
		view.window = geom.RectOrigin(float64(w), float64(h))
		m := worldToGl()
		smeg.SetViewMat([16]float32{
			float32(m[0][0]), float32(m[0][1]), 0, float32(m[0][2]),
			float32(m[1][0]), float32(m[1][1]), 0, float32(m[1][2]),
			0, 0, 1, 0,
			0, 0, 0, 1,
		})

		/* update systems */
		if mouse.held && mouse.rectHeld > -1 {
			i := mouse.rectHeld
			p := sprites.Ori[i].ToWorld(mouse.rectOffset)
			bodies.ApplyForce(i, p, mouse.world.Minus(p).Scaled(100.0))
		}
		n := 10
		for i := 0; i < n; i++ {
			bodies.Update(1.0 / (60 * float64(n)))
		}
		copy(sprites.Ori, bodies.COM)
		sprites.Update()

		/* draw */
		smeg.Clear()
		smeg.Draw(sprites.Verts, sprites.Indices, boxTex)
		mVerts, mInd := sprite.Single(geom.Ori2{mouse.world, 0}, geom.RectCentre(10, 10), smeg.Colour{1, 1, 1, 1})
		smeg.Draw(mVerts, mInd, boxTex)
		smeg.Display()
	}

	smeg.Terminate()
}

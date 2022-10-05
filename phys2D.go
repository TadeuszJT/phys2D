package phys2D

import (
	"fmt"
	"github.com/tadeuszjt/data"
	"github.com/tadeuszjt/geom/generic"
)

type joint struct {
	bodyKey [2]data.Key
	offset  [2]geom.Vec2[float64]

	// precompute
	index     [2]int
	jacobian  [2]geom.Ori2[float64]
	bias, jmj float64
}

type World struct {
	Gravity geom.Ori2[float64]

	bodies struct {
		data.KeyMap
		orientation data.RowT[geom.Ori2[float64]]
		velocity    data.RowT[geom.Ori2[float64]]
		invMass     data.RowT[geom.Ori2[float64]]
	}

	joints data.RowT[joint]
}

func NewWorld() *World {
	world := World{
		Gravity: geom.Ori2[float64]{0, 10, 0},
	}

	world.bodies.KeyMap = data.MakeKeyMap(data.Table{
		&world.bodies.orientation,
		&world.bodies.velocity,
		&world.bodies.invMass,
	})

	return &world
}

func (w *World) AddBody(orientation, mass geom.Ori2[float64]) data.Key {
	inv := mass
	if mass.X != 0 {
		inv.X = 1 / mass.X
	}
	if mass.Y != 0 {
		inv.Y = 1 / mass.Y
	}
	if mass.Theta != 0 {
		inv.Theta = 1 / mass.Theta
	}

	fmt.Println("body added, mass:", mass)

	return w.bodies.Append(orientation, geom.Ori2[float64]{}, inv)
}

func (w *World) DeleteBody(key data.Key) {
	for i := 0; i < len(w.joints); i++ {
		if w.joints[i].bodyKey[0] == key || w.joints[i].bodyKey[1] == key {
			w.joints.Delete(i)
			i--
		}
	}
	w.bodies.Delete(key)
}

func (w *World) AddJoint(bodyA, bodyB data.Key, offsetA, offsetB geom.Vec2[float64]) {
	w.joints.Append(joint{
		bodyKey: [2]data.Key{bodyA, bodyB},
		offset:  [2]geom.Vec2[float64]{offsetA, offsetB},
	})
}

func (w *World) ApplyImpulse(key data.Key, mag geom.Ori2[float64]) {
	index := w.bodies.GetIndex(key)
	w.applyImpulse(index, mag)
}

/*
Apply a force with an offset to a body.

	 offset
	 ^----->Fmag
	/

O-V Ftorque
mass
*/
func (w *World) ApplyForce(key data.Key, fMag, offset geom.Vec2[float64], dt float64) {
	// Ftorque = offset X F
	// a = f / m
	// dv = a * dt
	index := w.bodies.GetIndex(key)
	fOri := geom.Ori2[float64]{fMag.X, fMag.Y, offset.Cross(fMag)}
	accel := fOri.Times(w.bodies.invMass[index])
	w.bodies.velocity[index].PlusEquals(accel.ScaledBy(dt))
}

func (w *World) SetOrientations(keys []data.Key, orientations []geom.Ori2[float64]) {
	for i := range keys {
		index := w.bodies.GetIndex(keys[i])
		w.bodies.orientation[index] = orientations[i]
	}
}

func (w *World) GetOrientations(keys ...data.Key) []geom.Ori2[float64] {
	orientations := make([]geom.Ori2[float64], len(keys))
	for i := range keys {
		index := w.bodies.GetIndex(keys[i])
		orientations[i] = w.bodies.orientation[index]
	}
	return orientations
}

func (w *World) SetVelocities(keys []data.Key, velocities []geom.Ori2[float64]) {
	for i := range keys {
		index := w.bodies.GetIndex(keys[i])
		w.bodies.velocity[index] = velocities[i]
	}
}

package phys2D

import (
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

	joints struct {
		data.KeyMap
		row data.RowT[joint]
	}
}

func NewWorld() *World {
	world := World{
		Gravity: geom.Ori2[float64]{0, 10, 0},
	}

	world.bodies.KeyMap = data.KeyMap{
		Row: &data.Table{
			&world.bodies.orientation,
			&world.bodies.velocity,
			&world.bodies.invMass,
		},
	}

	world.joints.KeyMap = data.KeyMap{
		Row: &world.joints.row,
	}

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

	return w.bodies.Append(orientation, geom.Ori2[float64]{}, inv)
}

func (w *World) DeleteBody(key data.Key) {
	w.bodies.Delete(key)
}

func (w *World) AddJoint(bodyA, bodyB data.Key, offsetA, offsetB geom.Vec2[float64]) data.Key {
	return w.joints.Append(joint{
		bodyKey: [2]data.Key{bodyA, bodyB},
		offset:  [2]geom.Vec2[float64]{offsetA, offsetB},
	})
}

func (w *World) DeleteJoint(key data.Key) {
	w.joints.Delete(key)
}

func (w *World) ApplyImpulse(key data.Key, mag geom.Ori2[float64]) {
    index := w.bodies.KeyMap.KeyToIndex[key]
    w.applyImpulse(index, mag)
}

func (w *World) SetOrientations(keys []data.Key, orientations []geom.Ori2[float64]) {
	for i := range keys {
		index := w.bodies.KeyToIndex[keys[i]]
		w.bodies.orientation[index] = orientations[i]
	}
}

func (w *World) GetOrientations(keys ...data.Key) []geom.Ori2[float64] {
	orientations := make([]geom.Ori2[float64], len(keys))
	for i := range keys {
		index := w.bodies.KeyToIndex[keys[i]]
		orientations[i] = w.bodies.orientation[index]
	}
	return orientations
}

func (w *World) SetVelocities(keys []data.Key, velocities []geom.Ori2[float64]) {
	for i := range keys {
		index := w.bodies.KeyToIndex[keys[i]]
		w.bodies.velocity[index] = velocities[i]
	}
}

package phys2D

import (
	"github.com/tadeuszjt/geom/64"
)

type World struct {
	Gravity geom.Ori2
	bodySystem
	jointSystem
}

func (w *World) AddBody(orientation, mass geom.Ori2) Key {
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

	return w.bodySystem.Append(orientation, geom.Ori2{}, inv)
}

func (w *World) DeleteBody(key Key) {
	w.bodySystem.Delete(key)
}

func (w *World) AddJoint(bodyA, bodyB Key, offsetA, offsetB geom.Vec2) Key {
	return w.jointSystem.Append(joint{
		bodyKey: [2]Key{bodyA, bodyB},
		offset:  [2]geom.Vec2{offsetA, offsetB},
	})
}

func (w *World) DeleteJoint(key Key) {
	w.jointSystem.Delete(key)
}

func (w *World) SetOrientations(keys []Key, orientations []geom.Ori2) {
	b := &w.bodySystem
	for i := range keys {
		index := b.keyMap.keyToIndex[keys[i]]
		b.orientation[index] = orientations[i]
	}
}

func (w *World) GetOrientations(keys []Key) []geom.Ori2 {
	orientations := make([]geom.Ori2, len(keys))
	b := &w.bodySystem
	for i := range keys {
		index := b.keyMap.keyToIndex[keys[i]]
		orientations[i] = b.orientation[index]
	}
	return orientations
}

func (w *World) SetVelocities(keys []Key, velocities []geom.Ori2) {
	b := &w.bodySystem
	for i := range keys {
		index := b.keyMap.keyToIndex[keys[i]]
		b.velocity[index] = velocities[i]
	}
}

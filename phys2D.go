package phys2D

import (
	"github.com/tadeuszjt/geom/64"
	"github.com/tadeuszjt/data"
)

type joint struct {
	bodyKey [2]Key
	offset  [2]geom.Vec2

	// precompute
	index     [2]int
	jacobian  [2]geom.Ori2
	bias, jmj float64
}

type World struct {
	Gravity geom.Ori2

    bodies struct {
        orientation data.RowT[geom.Ori2]
        velocity    data.RowT[geom.Ori2]
        invMass     data.RowT[geom.Ori2]
        keyMap
    }

    joints struct {
        row data.RowT[joint]
        keyMap
    }
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

    w.bodies.orientation.Append(orientation)
    w.bodies.velocity.Append(geom.Ori2{})
    w.bodies.invMass.Append(inv)
    return w.bodies.keyMap.Append()
}

func (w *World) DeleteBody(key Key) {
    index := w.bodies.keyMap.keyToIndex[key]
    w.bodies.orientation.Delete(index)
    w.bodies.velocity.Delete(index)
    w.bodies.invMass.Delete(index)
	w.bodies.keyMap.Delete(key)
}

func (w *World) AddJoint(bodyA, bodyB Key, offsetA, offsetB geom.Vec2) Key {
	w.joints.row.Append(joint{
		bodyKey: [2]Key{bodyA, bodyB},
		offset:  [2]geom.Vec2{offsetA, offsetB},
	})
    return w.joints.keyMap.Append()
}

func (w *World) DeleteJoint(key Key) {
    index := w.joints.keyMap.keyToIndex[key]
	w.joints.row.Delete(index)
    w.joints.keyMap.Delete(key)
}

func (w *World) SetOrientations(keys []Key, orientations []geom.Ori2) {
	for i := range keys {
		index := w.bodies.keyMap.keyToIndex[keys[i]]
		w.bodies.orientation[index] = orientations[i]
	}
}

func (w *World) GetOrientations(keys []Key) []geom.Ori2 {
	orientations := make([]geom.Ori2, len(keys))
	for i := range keys {
		index := w.bodies.keyMap.keyToIndex[keys[i]]
		orientations[i] = w.bodies.orientation[index]
	}
	return orientations
}

func (w *World) SetVelocities(keys []Key, velocities []geom.Ori2) {
	for i := range keys {
		index := w.bodies.keyMap.keyToIndex[keys[i]]
		w.bodies.velocity[index] = velocities[i]
	}
}

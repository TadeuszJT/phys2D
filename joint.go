package phys2D

import (
	"github.com/tadeuszjt/geom/64"
)

type joint struct {
	bodyKey [2]Key
	offset  [2]geom.Vec2

	// precompute
	index     [2]int
	jacobian  [2]geom.Ori2
	bias, jmj float64
}

type jointSystem struct {
	keyMap
	joints []joint
}

func (j *jointSystem) Append(jnt joint) Key {
	j.joints = append(j.joints, jnt)
	return j.keyMap.Append()
}

func (j *jointSystem) Delete(key Key) {
	index := j.keyMap.keyToIndex[key]
	end := len(j.joints) - 1

	if index != end { // swap
		j.joints[index] = j.joints[end]
	}

	j.keyMap.Delete(key)
}

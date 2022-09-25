package phys2D

import (
	geom "github.com/tadeuszjt/geom/32"
)

type bodySystem struct {
	keyMap
	orientation, velocity, invMass []geom.Ori2
}

func (b *bodySystem) Append(orientation, velocity, invMass geom.Ori2) Key {
	b.orientation = append(b.orientation, orientation)
	b.velocity = append(b.velocity, velocity)
	b.invMass = append(b.invMass, invMass)
	return b.keyMap.Append()
}

func (b *bodySystem) Delete(key Key) {
	index := b.keyMap.keyToIndex[key]
	end := len(b.orientation) - 1

	if index != end { // swap
		b.orientation[index] = b.orientation[end]
		b.velocity[index] = b.velocity[end]
		b.invMass[index] = b.invMass[end]
	}

	/* Shrink arrays */
	b.orientation = b.orientation[:end]
	b.velocity = b.velocity[:end]
	b.invMass = b.invMass[:end]

	b.keyMap.Delete(key)
}

func (b *bodySystem) applyImpulse(index int, mag geom.Ori2) {
	t := geom.Ori2{
		mag.X * b.invMass[index].X,
		mag.Y * b.invMass[index].Y,
		b.invMass[index].Theta,
	}
	b.velocity[index].PlusEquals(t)
}

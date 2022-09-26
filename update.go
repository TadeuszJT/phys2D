package phys2D

import (
	"github.com/tadeuszjt/geom/64"
)

func (w *World) applyImpulse(index int, mag geom.Ori2) {
	w.bodies.velocity[index].PlusEquals(w.bodies.invMass[index].Times(mag))
}


func (w *World) Update(dt float64) {
	/* Apply forces */
	for i := range w.bodies.velocity {
		if w.bodies.invMass[i].X != 0 {
			w.bodies.velocity[i].X += w.Gravity.X * dt
		}
		if w.bodies.invMass[i].Y != 0 {
			w.bodies.velocity[i].Y += w.Gravity.Y * dt
		}
		if w.bodies.invMass[i].Theta != 0 {
			w.bodies.velocity[i].Theta += w.Gravity.Theta * dt
		}
	}

	/* Precompute constraints */
	for i := range w.joints.row {
		joint := &w.joints.row[i]
		index0 := w.bodies.keyMap.keyToIndex[joint.bodyKey[0]]
		index1 := w.bodies.keyMap.keyToIndex[joint.bodyKey[1]]

		joint.index[0] = index0
		joint.index[1] = index1

		o0 := w.bodies.orientation[index0]
		o1 := w.bodies.orientation[index1]

		// joint world positions
		p0 := o0.Mat3Transform().TimesVec2(joint.offset[0], 1).Vec2()
		p1 := o1.Mat3Transform().TimesVec2(joint.offset[1], 1).Vec2()

		// joint separation
		d := p0.Minus(p1)

		// jacobian
		joint.jacobian[0].X = 2 * d.X
		joint.jacobian[0].Y = 2 * d.Y
		joint.jacobian[0].Theta = -d.Cross(p0.Minus(o0.Vec2()))

		joint.jacobian[1].X = 2 * -d.X
		joint.jacobian[1].Y = 2 * -d.Y
		joint.jacobian[1].Theta = 2 * d.Cross(p1.Minus(o1.Vec2()))

		// bias
		joint.bias = (2 / dt) * (d.X*d.X + d.Y*d.Y)

		// J^T * M^-1 * J
		joint.jmj = joint.jacobian[0].Dot(w.bodies.invMass[index0].Times(joint.jacobian[0])) +
			joint.jacobian[1].Dot(w.bodies.invMass[index1].Times(joint.jacobian[1]))
	}

	/* Correct velocities */
	for num := 0; num < 1; num++ {
		for i := range w.joints.row {
			joint := &w.joints.row[i]

			vel0 := w.bodies.velocity[joint.index[0]]
			vel1 := w.bodies.velocity[joint.index[1]]

			Jv := joint.jacobian[0].Dot(vel0) + joint.jacobian[1].Dot(vel1)
			lambda := 0.0
			if joint.jmj != 0.0 {
				lambda = -(Jv + joint.bias) / joint.jmj
			}

			w.applyImpulse(joint.index[0], joint.jacobian[0].ScaledBy(lambda))
			w.applyImpulse(joint.index[1], joint.jacobian[1].ScaledBy(lambda))
		}
	}

	/* Set new positions */
	for i := range w.bodies.orientation {
		w.bodies.orientation[i].PlusEquals(w.bodies.velocity[i].ScaledBy(dt))
	}
}

package phys2D

import (
	//"fmt"
	"github.com/tadeuszjt/geom/generic"
)

func (w *World) applyImpulse(index int, mag geom.Ori2[float64]) {
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

	/* Apply Drag */
	for i := range w.dragPlates {
		plate := &w.dragPlates[i]
		bodyIndex := w.bodies.GetIndex(plate.bodyKey)
		bodyOrientation := w.bodies.orientation[bodyIndex]
		bodyVelocity := w.bodies.velocity[bodyIndex]

		delta := plate.point[1].Minus(plate.point[0])
		surface := delta.Perpendicular().RotatedBy(bodyOrientation.Theta)

		// simulate two plates to include an approximation to rotational drag
		points := [2]geom.Vec2[float64]{
			plate.point[0].Plus(delta.ScaledBy(0.25)),
			plate.point[1].Minus(delta.ScaledBy(0.25)),
		}

		for _, p := range points {
			offset := p.RotatedBy(bodyOrientation.Theta)
			plateVelFromRot := offset.Perpendicular().ScaledBy(bodyVelocity.Theta)
			plateVel := bodyVelocity.Vec2().Plus(plateVelFromRot)

			// Fd = 1/2 * v^2 * p * Cd * A
			scalar := plateVel.Dot(surface.Normal()) * 0.5 * w.AirDensity * plateVel.Len() * 0.5
			force := surface.ScaledBy(scalar)
			if force.Dot(plateVel) > 0 {
				force = force.ScaledBy(-1)
			}

			w.ApplyImpulse(plate.bodyKey, force, offset, dt)
		}
	}

	/* Precompute constraints */
	for i := range w.joints {
		joint := &w.joints[i]
		index0 := w.bodies.GetIndex(joint.bodyKey[0])
		index1 := w.bodies.GetIndex(joint.bodyKey[1])

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
	for num := 0; num < 10; num++ {
		for i := range w.joints {
			joint := &w.joints[i]

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

package phys2D

import (
	geom "github.com/tadeuszjt/geom/32"
)

func (w *World) Update(dt float32) {
	b := &w.bodySystem
	j := &w.jointSystem

	/* Apply forces */
	for i := range b.velocity {
		if b.invMass[i].X != 0 {
			b.velocity[i].X += w.Gravity.X * dt
		}
		if b.invMass[i].Y != 0 {
			b.velocity[i].Y += w.Gravity.Y * dt
		}
		if b.invMass[i].Theta != 0 {
			b.velocity[i].Theta += w.Gravity.Theta * geom.Angle(dt)
		}
	}

	/* Precompute constraints */
	for i := range j.joints {
		joint := &j.joints[i]
		index0 := b.keyMap.keyToIndex[joint.bodyKey[0]]
		index1 := b.keyMap.keyToIndex[joint.bodyKey[1]]

		joint.index[0] = index0
		joint.index[1] = index1

		o0 := b.orientation[index0]
		o1 := b.orientation[index1]

		// joint world positions
		p0 := o0.Mat3Transform().TimesVec2(joint.offset[0], 1).Vec2()
		p1 := o1.Mat3Transform().TimesVec2(joint.offset[1], 1).Vec2()

		// joint separation
		d := p0.Minus(p1)

		// jacobian
		joint.jacobian[0].X = 2 * d.X
		joint.jacobian[0].Y = 2 * d.Y
		joint.jacobian[0].Theta = geom.Angle(-d.Cross(p0.Minus(o0.Vec2())))

		joint.jacobian[1].X = 2 * -d.X
		joint.jacobian[1].Y = 2 * -d.Y
		joint.jacobian[1].Theta = geom.Angle(2 * d.Cross(p1.Minus(o1.Vec2())))

		// bias
		joint.bias = (2 / dt) * (d.X*d.X + d.Y*d.Y)

		// J^T * M^-1 * J
		joint.jmj = joint.jacobian[0].Vec3().Dot(b.invMass[index0].Vec3().Times(joint.jacobian[0].Vec3())) +
			joint.jacobian[1].Vec3().Dot(b.invMass[index1].Vec3().Times(joint.jacobian[1].Vec3()))
	}

	/* Correct velocities */
	for num := 0; num < 10; num++ {
		for i := range j.joints {
			joint := &j.joints[i]

			vel0 := b.velocity[joint.index[0]]
			vel1 := b.velocity[joint.index[1]]

			Jv := joint.jacobian[0].Vec3().Dot(vel0.Vec3()) + joint.jacobian[1].Vec3().Dot(vel1.Vec3())
			lambda := float32(0.0)
			if joint.jmj != 0.0 {
				lambda = -(Jv + joint.bias) / joint.jmj
			}

			j0Scaled := geom.Ori2{
				joint.jacobian[0].X * lambda,
				joint.jacobian[0].Y * lambda,
				joint.jacobian[0].Theta * geom.Angle(lambda),
			}
			j1Scaled := geom.Ori2{
				joint.jacobian[1].X * lambda,
				joint.jacobian[1].Y * lambda,
				joint.jacobian[1].Theta * geom.Angle(lambda),
			}

			b.applyImpulse(joint.index[0], j0Scaled)
			b.applyImpulse(joint.index[1], j1Scaled)
		}
	}

	/* Set new positions */
	for i := range b.orientation {
		velScaled := geom.Ori2{
			b.velocity[i].X * dt,
			b.velocity[i].Y * dt,
			b.velocity[i].Theta * geom.Angle(dt),
		}

		b.orientation[i].PlusEquals(velScaled)
	}
}

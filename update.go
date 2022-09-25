package phys2D

func (w *World) Update(dt float64) {
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
			b.velocity[i].Theta += w.Gravity.Theta * dt
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
		joint.jacobian[0].Theta = -d.Cross(p0.Minus(o0.Vec2()))

		joint.jacobian[1].X = 2 * -d.X
		joint.jacobian[1].Y = 2 * -d.Y
		joint.jacobian[1].Theta = 2 * d.Cross(p1.Minus(o1.Vec2()))

		// bias
		joint.bias = (2 / dt) * (d.X*d.X + d.Y*d.Y)

		// J^T * M^-1 * J
		joint.jmj = joint.jacobian[0].Dot(b.invMass[index0].Times(joint.jacobian[0])) +
			joint.jacobian[1].Dot(b.invMass[index1].Times(joint.jacobian[1]))
	}

	/* Correct velocities */
	for num := 0; num < 10; num++ {
		for i := range j.joints {
			joint := &j.joints[i]

			vel0 := b.velocity[joint.index[0]]
			vel1 := b.velocity[joint.index[1]]

			Jv := joint.jacobian[0].Dot(vel0) + joint.jacobian[1].Dot(vel1)
			lambda := 0.0
			if joint.jmj != 0.0 {
				lambda = -(Jv + joint.bias) / joint.jmj
			}

			b.applyImpulse(joint.index[0], joint.jacobian[0].ScaledBy(lambda))
			b.applyImpulse(joint.index[1], joint.jacobian[1].ScaledBy(lambda))
		}
	}

	/* Set new positions */
	for i := range b.orientation {
		b.orientation[i].PlusEquals(b.velocity[i].ScaledBy(dt))
	}
}

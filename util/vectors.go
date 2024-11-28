package util

import (
	"github.com/chewxy/math32"
	"github.com/et-nik/metamod-go/vector"
)

func WrapPitch(angle float32) float32 {
	for angle > 89 {
		angle = 89
	}

	for angle < -89 {
		angle = -89
	}

	return angle
}

func WrapAngle90(angle float32) float32 {
	for angle > 90 {
		angle -= 180
	}

	for angle < -90 {
		angle += 180
	}

	return angle
}

func WrapAngle180(angle float32) float32 {
	for angle > 180 {
		angle -= 360
	}

	for angle < -180 {
		angle += 360
	}

	return angle
}

func WrapAngle360(angle float32) float32 {
	for angle > 360 {
		angle -= 360
	}

	for angle < 0 {
		angle += 360
	}

	return angle
}

func SmoothAngle(current float32, target float32, speed float32) float32 {
	delta := target - current

	if delta > 180 {
		delta -= 360
	}

	if delta < -180 {
		delta += 360
	}

	if delta > speed {
		return current + speed
	}

	if delta < -speed {
		return current - speed
	}

	return target
}

func CurveValue(target float32, curvatureFactor float32) float32 {
	return target + (math32.Sin(target) * curvatureFactor)
}

// AnglesToRight converts angles to right vector
func AnglesToRight(angles vector.Vector) vector.Vector {
	pitch := angles[0] * (2 * math32.Pi / 360)
	sp := math32.Sin(pitch)
	cp := math32.Cos(pitch)

	yaw := angles[1] * (2 * math32.Pi / 360)
	sy := math32.Sin(yaw)
	cy := math32.Cos(yaw)

	roll := angles[2] * (2 * math32.Pi / 360)
	sr := math32.Sin(roll)
	cr := math32.Cos(roll)

	return vector.Vector{
		-1*sr*sp*cy + -1*cr*-sy,
		-1*sr*sp*sy + -1*cr*cy,
		-1 * sr * cp,
	}
}

// AnglesToForward converts angles to forward vector
func AnglesToForward(angles vector.Vector) vector.Vector {
	pitch := angles[0] * (2 * math32.Pi / 360)
	sp := math32.Sin(pitch)
	cp := math32.Cos(pitch)

	yaw := angles[1] * (2 * math32.Pi / 360)
	sy := math32.Sin(yaw)
	cy := math32.Cos(yaw)

	return vector.Vector{
		cp * cy,
		cp * sy,
		-1 * sp,
	}
}

func MiddleOfDistance(v1, v2 vector.Vector) vector.Vector {
	return vector.Vector{
		(v1[0] + v2[0]) / 2,
		(v1[1] + v2[1]) / 2,
		(v1[2] + v2[2]) / 2,
	}
}

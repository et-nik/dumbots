package main

func wrapAngle90(angle float32) float32 {
	for angle > 90 {
		angle -= 180
	}

	for angle < -90 {
		angle += 180
	}

	return angle
}

func wrapAngle180(angle float32) float32 {
	for angle > 180 {
		angle -= 360
	}

	for angle < -180 {
		angle += 360
	}

	return angle
}

func wrapAngle360(angle float32) float32 {
	for angle > 360 {
		angle -= 360
	}

	for angle < 0 {
		angle += 360
	}

	return angle
}

//func distance(a, b [3]float32) float32 {
//	return vek32.Distance(a[:], b[:])
//}
//
//func distance2D(a, b [3]float32) float32 {
//	return float32(math.Sqrt(float64((a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1]))))
//}
//
//func normalize(v [3]float32) [3]float32 {
//	ln := vek32.Norm(v[:])
//
//	if math.Abs(float64(ln)) < 1e-6 {
//		return [3]float32{0.0, 0.0, 1e-6}
//	}
//
//	ln = 1.0 / ln
//
//	return [3]float32{
//		v[0] * ln,
//		v[1] * ln,
//		v[2] * ln,
//	}
//}
//
//func vector32To64(v [3]float32) [3]float64 {
//	return [3]float64{float64(v[0]), float64(v[1]), float64(v[2])}
//}
//
//func vectorToSlice(v [3]float32) []float32 {
//	return []float32{v[0], v[1], v[2]}
//}
//
//func sliceFloat32To64(v []float32) []float64 {
//	result := make([]float64, len(v))
//
//	for i, val := range v {
//		result[i] = float64(val)
//	}
//
//	return result
//}

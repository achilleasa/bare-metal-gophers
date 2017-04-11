package main

import "math"

var zeroVec, hitThreshold, noHitThreshold, shadowSoftness, palette, scene, lights = vec3{}, 0.1, 20.0, 4.0,
	[]uint16{
		' ', '.', '`', '-', ':', ',', '_', '\'', ';', '=', '^', '/', '+', '"', '>', '<', '(', ')', 'i', '\\', '%', 'l', 'r', 'c', 'I', ']', '!', '[', '*', '?', 's', '1', '7', 'a', 'f', 'n', 'e', 'L', 'C', 'J', 'o', 'j', '3', '2', 'Y', 'T', '4', '5', 'F', 'S', '9', 'P', 'k', 'b', 'd', '6', 'X', 'G', 'h', 'V', 'q', 'p', 'Z', 'A', 'E', 'g', 'U', '&', 'O', 'D', 'K', '#', '@', 'm', 'H', '$', '8', '0', 'R', 'B', 'W', 'N', 'Q', 'M',
	}, []struct {
		fn                    func(time float64, point vec3, origin vec3, params vec3) float64
		origin, params, color vec3
	}{
		{dRepBouncingSpheres, vec3{5, 3, 5}, vec3{0.8, 0, 0}, vec3{8, 7, 4, 12}}, {dEllipsoid, vec3{0, 0, 0}, vec3{1.5, 3.0, 1.5}, vec3{8, 7, 3, 11}}, {dEllipsoid, vec3{-0.6, 2.8, 0}, vec3{0.1, 0.3, 0.1}, vec3{8, 7, 3, 11}},
		{dEllipsoid, vec3{0.6, 2.8, 0}, vec3{0.1, 0.3, 0.1}, vec3{8, 7, 3, 11}}, {dEllipsoid, vec3{-0.5, 2.0, 1.0}, vec3{0.5, 0.5, 0.5}, vec3{8, 7, 8, 15}}, {dEllipsoid, vec3{0.5, 2.0, 1.0}, vec3{0.5, 0.5, 0.5}, vec3{8, 7, 8, 15}},
		{dEllipsoid, vec3{0, 1.2, 1.0}, vec3{0.5, 0.3, 0.5}, vec3{8, 7, 6, 15}}, {dEllipsoid, vec3{-1.5, 1.0, 0.0}, vec3{0.7, 0.3, 0.3}, vec3{8, 7, 6, 15}}, {dEllipsoid, vec3{1.5, 1.0, 0.0}, vec3{0.7, 0.3, 0.3}, vec3{8, 7, 6, 15}},
		{dEllipsoid, vec3{-0.5, -1.5, 1.0}, vec3{0.3, 0.3, 0.8}, vec3{8, 7, 6, 15}}, {dEllipsoid, vec3{0.5, -1.5, 1.0}, vec3{0.3, 0.3, 0.8}, vec3{8, 7, 6, 15}}, {dPlane, vec3{0, 1, 0}, vec3{1.5, 0, 0}, vec3{8, 7, 1, 9}},
	}, []vec3{{0, 3.2, -3}, {0, 3.2, 3}}

type vec3 [4]float64 // dirty hack: the 4th float is ignored by vector methods and only used for the primitives' 4-color palette

func (v *vec3) Normalize() *vec3                { return v.MulS(1 / v.Len()) }
func (v *vec3) Len() float64                    { return math.Sqrt(v.Dot(*v) + 0.00001) }
func (v *vec3) Dot(o vec3) float64              { return v[0]*o[0] + v[1]*o[1] + v[2]*o[2] }
func (v *vec3) Sub(o vec3) *vec3                { v[0], v[1], v[2] = v[0]-o[0], v[1]-o[1], v[2]-o[2]; return v }
func (v *vec3) MulS(s float64) *vec3            { v[0], v[1], v[2] = v[0]*s, v[1]*s, v[2]*s; return v }
func (v *vec3) AddS(index int, s float64) *vec3 { v[index] += s; return v }
func (v *vec3) Cross(a, b vec3) *vec3 {
	v[0], v[1], v[2] = a[1]*b[2]-a[2]*b[1], a[2]*b[0]-a[0]*b[2], a[0]*b[1]-a[1]*b[0]
	return v
}
func (v *vec3) MulMat(m [16]float64, w float64) *vec3 {
	v[0], v[1], v[2] = m[0]*v[0]+m[1]*v[1]+m[2]*v[2]+m[3]*w, m[4]*v[0]+m[5]*v[1]+m[6]*v[2]+m[7]*w, m[8]*v[0]+m[9]*v[1]+m[10]*v[2]+m[11]*w
	return v
}
func (v *vec3) MulAdd(a, b vec3, s float64) *vec3 {
	v[0], v[1], v[2] = a[0]+s*b[0], a[1]+s*b[1], a[2]+s*b[2]
	return v
}
func dRepBouncingSpheres(t float64, p vec3, rep, params vec3) float64 {
	p[0], p[1], p[2] = (p[0]-rep[0]*math.Floor(p[0]/rep[0]))-0.5*rep[0], p[1]-rep[1]*math.Abs(math.Sin(2*t)), (p[2]-rep[2]*math.Floor(p[2]/rep[2]))-0.5*rep[2]
	return p.Len() - params[0]
}
func dEllipsoid(_ float64, p, o, params vec3) float64 {
	p[0], p[1], p[2] = (p[0]-o[0])/params[0], (p[1]-o[1])/params[1], (p[2]-o[2])/params[2]
	return (p.Len() - 1) * math.Min(math.Min(params[0], params[1]), params[2])
}
func dPlane(_ float64, p, origin, params vec3) float64 { return p.Dot(origin) + params[0] }
func distAt(t float64, p vec3, nearestPrim *int) (minDist float64) {
	for index, prim := range scene {
		if d := prim.fn(t, p, prim.origin, prim.params); index == 0 || d < minDist {
			minDist, *nearestPrim = d, index
		}
	}
	return
}
func normalAt(t float64, p vec3) (normal vec3) {
	for sample, axis, nearestPrim := p, 0, 0; axis < 3; sample, axis = p, axis+1 {
		normal[axis] = distAt(t, *sample.AddS(axis, 0.01), &nearestPrim) - distAt(t, *sample.AddS(axis, -0.02), &nearestPrim)
	}
	return *normal.Normalize()
}
func sampleAt(t float64, origin, p vec3, nearestPrim int, hitDist float64) uint16 {
	var lDiff float64
	for nAtP, lDir, occlusion, lIndex := normalAt(t, p), p, p, 0; lIndex < len(lights); lIndex++ {
		dToL := lDir.Sub(lights[lIndex]).MulS(-1).Len() - 2*hitThreshold
		_, _, visFactor := trace(t, p, *lDir.Normalize(), &occlusion, 2*hitThreshold, dToL)
		lDiff += visFactor * math.Max(0.0, lDir.Dot(nAtP))
	}
	return uint16(math.Max(0, math.Copysign(1.0, noHitThreshold-hitDist))) * ((uint16(scene[nearestPrim].color[int(math.Min(3, lDiff*4))]) << 8) | palette[int(math.Min(0.99, lDiff)*float64(len(palette)))])
}
func trace(t float64, origin, dir vec3, hit *vec3, minTraceDist, maxTraceDist float64) (nearestPrim int, hitDist, visFactor float64) {
	visFactor = 1.0
	for stepDist := maxTraceDist; minTraceDist < maxTraceDist && stepDist >= hitThreshold; minTraceDist, hit = minTraceDist+stepDist, hit.MulAdd(origin, dir, minTraceDist+stepDist) {
		stepDist = distAt(t, *hit.MulAdd(origin, dir, minTraceDist), &nearestPrim)
		visFactor = math.Copysign(1, stepDist-hitThreshold) * math.Min(visFactor, shadowSoftness*stepDist/(minTraceDist+0.001))
	}
	return nearestPrim, minTraceDist, math.Max(0, math.Min(1, visFactor))
}
func bareMetalGophers(fb []uint16, fbWidth, fbHeight float64, eye, look, up vec3, tickStep float64) {
	for t, startSpin, dx, dy, yaw, rotEye, dir, p, xa, ya, za := 0.0, -150.0, 1/fbWidth, 1/fbHeight, 0.0, zeroVec, zeroVec, zeroVec, zeroVec, zeroVec, zeroVec; ; t, yaw, startSpin = t+tickStep, yaw+tickStep*math.Min(1, math.Max(0, startSpin)), startSpin+1 {
		rotEye[0], rotEye[1], rotEye[2] = eye[2]*math.Sin(yaw)+eye[0]*math.Cos(yaw), eye[1], eye[2]*math.Cos(yaw)-eye[0]*math.Sin(yaw)
		za.MulAdd(rotEye, look, -1).Normalize()
		xa.Cross(up, za).Normalize()
		ya.Cross(za, xa)
		viewMat := [16]float64{xa[0], ya[0], za[0], rotEye[0], xa[1], ya[1], za[1], rotEye[1], xa[2], ya[2], za[2], rotEye[2], 0, 0, 0, 1}
		for y, ry, offset, origin := 0, -0.5+0.5*dy, 0, zeroVec.MulS(0).MulMat(viewMat, 1.0); y < int(fbHeight); y, ry = y+1, ry+dy {
			for x, rx := 0, -0.5+0.5*dx; x < int(fbWidth); x, rx, offset = x+1, rx+dx, offset+1 {
				dir[0], dir[1], dir[2] = rx, -ry, -1
				nearest, hitDist, _ := trace(t, *origin, *dir.MulMat(viewMat, 0).Normalize(), &p, 0, noHitThreshold)
				fb[offset] = sampleAt(t, *origin, p, nearest, hitDist)
			}
		}
	}
}

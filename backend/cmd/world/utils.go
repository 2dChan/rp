// Copyright (c) 2026 Andrey Kriulin
// Licensed under the MIT License.
// See the LICENSE file in the project root for full license text.

package main

import (
	"math"

	"github.com/golang/geo/r2"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
)

// ========================================================================
//    This file created by AI and not checked. Don't using in production
// ========================================================================

func PointToScreenFloat(p s2.Point) r2.Point {
	r2p := proj.Project(p)

	x := (r2p.X + xScale) / (2 * xScale)
	y := (-r2p.Y + xScale/2) / xScale

	return r2.Point{X: x * float64(width), Y: y * float64(height)}
}

func SplitPolygonAtAntimeridian(region interface {
	NumVertices() int
	Vertex(int) s2.Point
	Center() s2.Point
}) (polygons [2][]r2.Point) {
	n := region.NumVertices()
	if n == 0 {
		return [2][]r2.Point{nil, nil}
	}

	type vertex struct {
		lat, lng float64
		screen   r2.Point
	}

	vertices := make([]vertex, n)
	for i := range n {
		p := region.Vertex(i)
		ll := s2.LatLngFromPoint(p)
		vertices[i] = vertex{
			lat:    ll.Lat.Radians(),
			lng:    ll.Lng.Radians(),
			screen: PointToScreenFloat(p),
		}
	}

	hasCrossing := false
	for i := range n {
		j := (i + 1) % n
		if crossesAntimeridian(vertices[i].lng, vertices[j].lng) {
			hasCrossing = true
			break
		}
	}

	if !hasCrossing {
		points := make([]r2.Point, n)
		for i, v := range vertices {
			points[i] = v.screen
		}
		return [2][]r2.Point{points, nil}
	}

	for i := range n {
		j := (i + 1) % n
		v1, v2 := vertices[i], vertices[j]

		if crossesAntimeridian(v1.lng, v2.lng) {
			crossLat := interpolateLatAtAntimeridian(v1.lat, v1.lng, v2.lat, v2.lng)
			crossPoint := s2.PointFromLatLng(s2.LatLng{Lat: s1.Angle(crossLat), Lng: s1.Angle(math.Pi)})
			crossY := PointToScreenFloat(crossPoint).Y

			if v1.lng > 0 {
				polygons[0] = append(polygons[0], v1.screen)
				polygons[0] = append(polygons[0], r2.Point{X: float64(width), Y: crossY})
				polygons[1] = append(polygons[1], r2.Point{X: 0, Y: crossY})
			} else {
				polygons[1] = append(polygons[1], v1.screen)
				polygons[1] = append(polygons[1], r2.Point{X: 0, Y: crossY})
				polygons[0] = append(polygons[0], r2.Point{X: float64(width), Y: crossY})
			}
		} else {
			if v1.lng > 0 {
				polygons[0] = append(polygons[0], v1.screen)
			} else {
				polygons[1] = append(polygons[1], v1.screen)
			}
		}
	}

	return
}

func interpolateLatAtAntimeridian(lat1, lng1, lat2, lng2 float64) float64 {
	var dLng float64
	if lng1 > 0 {
		dLng = (math.Pi - lng1) + (math.Pi + lng2)
	} else {
		dLng = (math.Pi + lng1) + (math.Pi - lng2)
	}

	var t float64
	if lng1 > 0 {
		t = (math.Pi - lng1) / dLng
	} else {
		t = (math.Pi + lng1) / dLng
	}

	return lat1 + t*(lat2-lat1)
}

func crossesAntimeridian(lng1, lng2 float64) bool {
	return math.Abs(lng2-lng1) > math.Pi
}

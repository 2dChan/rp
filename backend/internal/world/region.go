// Copyright (c) 2026 Andrey Kriulin
// Licensed under the MIT License.
// See the LICENSE file in the project root for full license text.

package world

import (
	"fmt"

	"github.com/golang/geo/s2"
)

type Region struct {
	r   *Regions
	idx int
}

func (r Region) Center() s2.Point {
	return r.r.centers[r.idx]
}

func (r Region) Height() uint8 {
	return r.r.heights[r.idx]
}

func (r Region) NumVertices() int {
	return r.r.regionOffsets[r.idx+1] - r.r.regionOffsets[r.idx]
}

func (r Region) Vertex(i int) s2.Point {
	start := r.r.regionOffsets[r.idx]
	end := r.r.regionOffsets[r.idx+1]
	if i < 0 || i >= end-start {
		panic(fmt.Sprintf("Vertex: index %d out of range [0 %d)", i, end-start))
	}
	return r.r.vertices[r.r.borderIndices[start+i]]
}

func (r Region) NumNeighbors() int {
	return r.r.regionOffsets[r.idx+1] - r.r.regionOffsets[r.idx]
}

func (r Region) Neighbor(i int) Region {
	start := r.r.regionOffsets[r.idx]
	end := r.r.regionOffsets[r.idx+1]
	if i < 0 || i >= end-start {
		panic(fmt.Sprintf("Neighbor: index %d out of range [0 %d)", i, end-start))
	}
	return r.r.At(r.r.neighborIndices[start+i])
}

func (r Region) IsGoingToWater() bool {
	// TODO: Add logic.
	return false
}

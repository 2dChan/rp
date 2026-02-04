// Copyright (c) 2026 Andrey Kriulin
// Licensed under the MIT License.
// See the LICENSE file in the project root for full license text.

package world

import (
	"fmt"
	"math"

	"github.com/2dChan/s2voronoi"
	"github.com/2dChan/s2voronoi/utils"
	"github.com/golang/geo/s2"
	"github.com/ojrac/opensimplex-go"
)

type Regions struct {
	centers []s2.Point
	heights []uint8

	vertices []s2.Point
	// NOTE: CCW sort(look out of sphere)
	borderIndices []int
	// NOTE: CCW sort(look out of sphere)
	neighborIndices []int
	regionOffsets   []int
}

func (r *Regions) Len() int {
	return len(r.regionOffsets) - 1
}

func (r *Regions) At(i int) Region {
	if i < 0 || i+1 >= len(r.regionOffsets) {
		right := len(r.regionOffsets) - 1
		panic(fmt.Sprintf("Region: index %d out of range [0 %d)", i, right))
	}
	return Region{r: r, idx: i}
}

type World struct {
	regions Regions
}

type Options struct {
	Scale float64
	Seed  int64
}

type Option func(*Options)

func WithScale(scale float64) Option {
	return func(o *Options) {
		o.Scale = scale
	}
}

func WithSeed(seed int64) Option {
	return func(o *Options) {
		o.Seed = seed
	}
}

func NewWorld(numRegions int, setters ...Option) (*World, error) {
	if numRegions < 4 {
		return nil, fmt.Errorf("NewWorld: insufficient regions for world, minimum 4 required")
	}

	opts := &Options{
		Scale: 2,
		Seed:  0,
	}
	for _, set := range setters {
		set(opts)
	}

	sites := utils.GenerateRandomPoints(numRegions, opts.Seed)
	vd, err := s2voronoi.NewDiagram(sites)
	if err != nil {
		return nil, err
	}
	// TODO: Add to Options.
	if err := vd.Relax(3); err != nil {
		return nil, err
	}

	heights := make([]uint8, numRegions)
	noise := opensimplex.NewNormalized(opts.Seed)
	for i := range vd.NumCells() {
		c := vd.Cell(i)
		x, y, z := c.Site().X*opts.Scale, c.Site().Y*opts.Scale, c.Site().Z*opts.Scale
		heights[i] = uint8(noise.Eval3(x, y, z) * math.MaxUint8)
	}

	world := &World{
		regions: Regions{
			centers: vd.Sites,
			heights: heights,

			vertices:        vd.Vertices,
			borderIndices:   vd.CellVertices,
			neighborIndices: vd.CellNeighbors,
			regionOffsets:   vd.CellOffsets,
		},
	}

	return world, nil
}

func (w *World) NumRegions() int {
	return w.regions.Len()
}

func (w *World) Region(i int) Region {
	return w.regions.At(i)
}

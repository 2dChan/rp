// Copyright (c) 2026 Andrey Kriulin
// Licensed under the MIT License.
// See the LICENSE file in the project root for full license text.

package main

import (
	"log"
	"os"

	"github.com/2dChan/rp/backend/internal/world"
	svg "github.com/ajstarks/svgo"
	"github.com/golang/geo/s2"
)

const (
	filename = "world.svg"

	width  = 1500
	height = width / 2

	regionStyle = "fill:rgb(230,230,230);stroke:rgb(170,170,170);stroke-width:1;stroke-opacity:1.0"
	waterStyle  = "fill:rgb(170,210,230)"
)

var (
	xScale = float64(width)
	proj   = s2.NewMercatorProjection(xScale)
)

func renderWorld(world *world.World) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	canvas := svg.New(file)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, waterStyle)
	xPoints := make([]int, 0)
	yPoints := make([]int, 0)

	for i := range world.NumRegions() {
		r := world.Region(i)
		if r.Height() < 140 {
			continue
		}

		for _, poly := range SplitPolygonAtAntimeridian(r) {
			if len(poly) >= 3 {
				xPoints = xPoints[:0]
				yPoints = yPoints[:0]
				for _, p := range poly {
					xPoints = append(xPoints, int(p.X))
					yPoints = append(yPoints, int(p.Y))
				}
				canvas.Polygon(xPoints, yPoints, regionStyle)
			}
		}
	}

	canvas.End()
}

func main() {
	world, err := world.NewWorld(5000, 1000)
	if err != nil {
		log.Fatal(err)
	}
	renderWorld(world)
}

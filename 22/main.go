package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
)

type cube struct {
	points  map[grid.XYZVec]int
	indexes map[int][]grid.XYZVec
}

func (c cube) String() string {
	pieces := len(c.indexes)
	cubes := make([][]grid.XYZVec, pieces)
	for i := 0; i < pieces; i++ {
		cubes[i] = c.indexes[i]
	}
	var b strings.Builder
	for i, z := range cubes {
		b.WriteString(fmt.Sprintf("%d: %v\n", i, z))
	}
	return b.String()
}

func (c *cube) Drop(index int, l line3D) {
	// fmt.Printf("drop: %d, %v\n", index, l)
	z := l.from.Z
	var depth int
	for {
		if z == 0 {
			for _, p := range l.points() {
				p = grid.XYZVec{
					X: p.X,
					Y: p.Y,
					Z: p.Z - depth + 1,
				}
				c.indexes[index] = append(c.indexes[index], p)
				c.points[p] = index
			}
			break
		}
		z--
		depth++
		var stop bool
		for _, p := range l.points() {
			p = grid.XYZVec{
				X: p.X,
				Y: p.Y,
				Z: p.Z - depth,
			}
			if _, ok := c.points[p]; ok {
				stop = true
				break
			}
		}
		if stop {
			for _, p := range l.points() {
				p := grid.XYZVec{
					X: p.X,
					Y: p.Y,
					Z: p.Z - depth + 1,
				}
				c.indexes[index] = append(c.indexes[index], p)
				c.points[p] = index
			}
			break
		}
	}
}

func (c cube) Count() int32 {
	var count int32
	for i, ps := range c.indexes {
		var reserveFlag bool
		for _, p := range ps {
			pp := grid.XYZVec{
				X: p.X,
				Y: p.Y,
				Z: p.Z + 1,
			}
			// 被其他 cube 压着
			// 被压着的不能是自己
			// parent 只能有一个 child
			if parent, ok := c.points[pp]; ok && c.points[pp] != i {
				var moreChild bool
				for _, ppa := range c.indexes[parent] {
					ppap := grid.XYZVec{
						X: ppa.X,
						Y: ppa.Y,
						Z: ppa.Z - 1,
					}
					if c.points[ppap] != i && c.points[ppap] != parent {
						moreChild = true
						break
					}
				}
				if !moreChild {
					reserveFlag = true
					break
				}
			}
		}
		if !reserveFlag {
			fmt.Printf("destroy: %d\n", i)
			count++
		}
	}
	return count
}

type line3D struct {
	from grid.XYZVec
	to   grid.XYZVec
}

func (l line3D) points() []grid.XYZVec {
	var points []grid.XYZVec
	switch {
	case l.from.X == l.to.X:
		for y := l.from.Y; y <= l.to.Y; y++ {
			for z := l.from.Z; z <= l.to.Z; z++ {
				points = append(points, grid.XYZVec{
					X: l.from.X,
					Y: y,
					Z: z,
				})
			}
		}
	case l.from.Y == l.to.Y:
		for x := l.from.X; x <= l.to.X; x++ {
			for z := l.from.Z; z <= l.to.Z; z++ {
				points = append(points, grid.XYZVec{
					X: x,
					Y: l.from.Y,
					Z: z,
				})
			}
		}
	case l.from.Z == l.to.Z:
		for x := l.from.X; x <= l.to.X; x++ {
			for y := l.from.Y; y <= l.to.Y; y++ {
				points = append(points, grid.XYZVec{
					X: x,
					Y: y,
					Z: l.from.Z,
				})
			}
		}
	default:
		panic("invalid line")
	}
	return points
}

func main() {
	txt := input.NewTXTFile("input.txt")
	var lines []line3D
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		parts := strings.Split(line, "~")
		xyz := strings.Split(parts[0], ",")
		xyz2 := strings.Split(parts[1], ",")
		lines = append(lines, line3D{
			from: grid.XYZVec{
				X: input.Atoi(xyz[0]),
				Y: input.Atoi(xyz[1]),
				Z: input.Atoi(xyz[2]),
			},
			to: grid.XYZVec{
				X: input.Atoi(xyz2[0]),
				Y: input.Atoi(xyz2[1]),
				Z: input.Atoi(xyz2[2]),
			},
		})
		return nil
	})
	c := cube{
		points:  make(map[grid.XYZVec]int),
		indexes: make(map[int][]grid.XYZVec),
	}
	for i, l := range lines {
		c.Drop(i, l)
	}
	fmt.Println(c)
	fmt.Printf("p1: %d\n", c.Count())
}

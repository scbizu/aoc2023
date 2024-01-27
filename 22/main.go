package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/set"
)

type cube struct {
	points  map[grid.XYZVec]int
	indexes map[int][]grid.XYZVec
}

type supportGraph struct {
	G       map[int]*set.Set[int]
	reverse map[int]*set.Set[int]
}

func (s supportGraph) Reverse() supportGraph {
	s.reverse = make(map[int]*set.Set[int])
	for f, ts := range s.G {
		for _, t := range ts.List() {
			if _, ok := s.reverse[t]; !ok {
				s.reverse[t] = set.New[int]()
			}
			s.reverse[t].Add(f)
		}
	}
	return s
}

func (s supportGraph) ReverseString() string {
	b := strings.Builder{}
	for f, ts := range s.reverse {
		b.WriteString(fmt.Sprintf("%d -> %s\n", f, ts.String()))
	}
	return b.String()
}

func (s supportGraph) String() string {
	b := strings.Builder{}
	for f, ts := range s.G {
		for _, t := range ts.List() {
			b.WriteString(fmt.Sprintf("%d-->%d\n", f, t))
		}
	}
	return b.String()
}

func (s supportGraph) FindHasSupportNode() int {
	var count int
	for _, ts := range s.G {
		var reserve bool
		for _, t := range ts.List() {
			if s.reverse[t].Size() == 1 {
				reserve = true
				break
			}
		}
		if !reserve {
			count++
		}
	}
	for i := range s.reverse {
		if _, ok := s.G[i]; !ok {
			count++
		}
	}
	return count
}

func (c cube) ToGraph() supportGraph {
	sg := supportGraph{
		G: make(map[int]*set.Set[int]),
	}
	for i := 0; i < len(c.indexes); i++ {
		for _, p := range c.indexes[i] {
			pp := grid.XYZVec{
				X: p.X,
				Y: p.Y,
				Z: p.Z + 1,
			}
			if _, ok := c.points[pp]; ok && c.points[pp] != i {
				if _, ok := sg.G[i]; !ok {
					sg.G[i] = set.New[int]()
				}
				sg.G[i].Add(c.points[pp])
			}
		}
	}
	return sg
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
	fmt.Printf("%s\n", c.ToGraph())
	fmt.Printf("p1:%d\n", c.ToGraph().Reverse().FindHasSupportNode())
}

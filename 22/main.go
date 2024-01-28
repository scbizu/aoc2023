package main

import (
	"context"
	"fmt"
	"sort"
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
	G        map[int]*set.Set[int]
	reverse  map[int]*set.Set[int]
	destroys []int
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

func (s supportGraph) String() string {
	b := strings.Builder{}
	for f, ts := range s.G {
		b.WriteString(fmt.Sprintf("%d -> %v\n", f, ts.String()))
	}
	return b.String()
}

func (s supportGraph) chainDestroy() int {
	var count int
	for _, i := range s.destroys {
		c := s.chain(i, set.New[int](i))
		// fmt.Printf("chain: %d, %d\n", i, c)
		count += c
	}
	return count
}

func (s supportGraph) chain(i int, bks *set.Set[int]) int {
	var count int
	if _, ok := s.G[i]; !ok {
		return 0
	}
	for _, t := range s.G[i].List() {
		if bks.Has(t) {
			continue
		}
		re, ok := s.reverse[t]
		if ok {
			// 存在其它支撑当前砖块的砖块不在「已经确定将要坠落的砖块列表」里面
			// 那么这个砖块就不会坠落
			// re:  所有支撑当前砖块的砖块列表
			// bks: 已经确定将要坠落的砖块列表
			if set.Intersection[int](bks, re).Size() < re.Size() {
				// fmt.Printf("size: %d, %d\n", set.Intersection[int](bks, re).Size(), re.Size())
				continue
			}
		}
		count++
		bks.Add(t)
		count += s.chain(t, bks)
	}
	return count
}

func (s supportGraph) CountDisintegrate() (int, *set.Set[int]) {
	var count int
	ds := set.New[int]()
	for i, ts := range s.G {
		var reserve bool
		for _, t := range ts.List() {
			if s.reverse[t].Size() == 1 {
				reserve = true
				break
			}
		}
		if !reserve {
			ds.Add(i)
			count++
		}
	}
	for i := range s.reverse {
		if _, ok := s.G[i]; !ok {
			ds.Add(i)
			count++
		}
	}
	return count, ds
}

func (c cube) ToGraph() supportGraph {
	sg := supportGraph{
		G: make(map[int]*set.Set[int]),
	}
	for i := 0; i < len(c.indexes); i++ {
		sg.G[i] = set.New[int]()
		for _, p := range c.indexes[i] {
			pp := grid.XYZVec{
				X: p.X,
				Y: p.Y,
				Z: p.Z + 1,
			}
			if _, ok := c.points[pp]; ok && c.points[pp] != i {
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
	// order by z asc
	sort.Slice(lines, func(i, j int) bool {
		return lines[i].from.Z < lines[j].from.Z
	})
	for i, l := range lines {
		c.Drop(i, l)
	}
	g := c.ToGraph()
	// fmt.Printf("%v\n", g)
	g = g.Reverse()
	count, d := g.CountDisintegrate()
	fmt.Printf("p1:%d\n", count)
	ds := set.New[int]()
	for i := 0; i < len(lines); i++ {
		ds.Add(i)
	}
	diff := set.Difference[int](ds, d)
	g.destroys = diff.List()
	// fmt.Printf("diff: %v\n", diff)

	chain := g.chainDestroy()
	fmt.Printf("p2:%d\n", chain)
}

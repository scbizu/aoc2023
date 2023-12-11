package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/magejiCoder/set"
	"github.com/scbizu/aoc2022/helper/input"
	"github.com/scbizu/aoc2022/helper/matrix"
)

const (
	line, col = 140, 140
)

// gImage is the galaxies image
type gImage struct {
	m               matrix.Matrix[byte]
	galaxies        []matrix.Point[byte]
	doubleLineIndex []int
	doubleColIndex  []int
	gap             int
}

func getLine(m matrix.Matrix[byte], i int) []byte {
	return m[i]
}

func getCol(m matrix.Matrix[byte], j int) []byte {
	col := make([]byte, len(m))
	for i := range m {
		col[i] = m[i][j]
	}
	return col
}

func (gi gImage) Lines() int {
	return len(gi.m)
}

func (gi gImage) Cols() int {
	return len(gi.m[0])
}

func isAllEmpty(raw string) bool {
	return strings.Count(raw, ".") == len(raw)
}

func (gi *gImage) ReIndex() {
	gi.expand()
	gi.reIndexGalaxies(gi.gap)
}

func (gi *gImage) Print() {
	gi.m.PrintEx("%c")
	fmt.Printf("gals: %v\n", gi.galaxies)
	fmt.Printf("lIndex: %v\n", gi.doubleLineIndex)
	fmt.Printf("cIndex: %v\n", gi.doubleColIndex)
}

type fromTo struct {
	from matrix.Point[byte]
	to   matrix.Point[byte]
}

func (gi gImage) ShortestPath() int {
	if len(gi.galaxies) == 0 {
		return 0
	}
	if len(gi.galaxies) == 1 {
		return 0
	}
	visit := set.New[fromTo]()
	var path int
	for _, g := range gi.galaxies {
		for _, g2 := range gi.galaxies {
			if g.X == g2.X && g.Y == g2.Y {
				continue
			}
			if visit.Has(fromTo{from: g, to: g2}) {
				continue
			}
			p := abs(g.X, g2.X) + abs(g.Y, g2.Y)
			fmt.Printf("(%d,%d) -> (%d,%d) : %d\n", g.X, g.Y, g2.X, g2.Y, p)
			path += p
			visit.Add(fromTo{from: g, to: g2})
			visit.Add(fromTo{from: g2, to: g})
		}
	}
	return path
}

func abs(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func (gi *gImage) expand() {
	for i := 0; i < gi.Lines(); i++ {
		if isAllEmpty(string(getLine(gi.m, i))) {
			gi.doubleLineIndex = append(gi.doubleLineIndex, i)
		}
	}
	for j := 0; j < gi.Cols(); j++ {
		if isAllEmpty(string(getCol(gi.m, j))) {
			gi.doubleColIndex = append(gi.doubleColIndex, j)
		}
	}
}

func (gi *gImage) reIndexGalaxies(gap int) {
	// fmt.Printf("index gals: %+v\n", gi.galaxies)
	newGals := make([]matrix.Point[byte], len(gi.galaxies))
	copy(newGals, gi.galaxies)
	for i := range gi.galaxies {
		x, y := gi.galaxies[i].X, gi.galaxies[i].Y
		for _, lIndex := range gi.doubleLineIndex {
			if x > lIndex {
				newGals[i].X += gap
			}
		}
		for _, cIndex := range gi.doubleColIndex {
			if y > cIndex {
				newGals[i].Y += gap
			}
		}
	}
	gi.galaxies = newGals
	// fmt.Printf("indexed gals: %+v\n", gi.galaxies)
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	image := matrix.NewMatrix[byte](line, col)
	var points []matrix.Point[byte]
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		for j, c := range line {
			image.Add(i, j, byte(c))
			if c != '.' {
				points = append(points, matrix.Point[byte]{X: i, Y: j})
			}
		}
		return nil
	})
	img := &gImage{
		m:        image,
		galaxies: points,
		gap:      1,
	}
	// img.Print()
	img.ReIndex()
	// img.Print()
	p := img.ShortestPath()
	fmt.Fprintf(os.Stdout, "p1: %d\n", p)
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	image := matrix.NewMatrix[byte](line, col)
	var points []matrix.Point[byte]
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		for j, c := range line {
			image.Add(i, j, byte(c))
			if c != '.' {
				points = append(points, matrix.Point[byte]{X: i, Y: j})
			}
		}
		return nil
	})
	img := &gImage{
		m:        image,
		galaxies: points,
		gap:      1000000 - 1,
	}
	// img.Print()
	img.ReIndex()
	// img.Print()
	p := img.ShortestPath()
	fmt.Fprintf(os.Stdout, "p2: %d\n", p)
}

func main() {
	p1()
	p2()
}

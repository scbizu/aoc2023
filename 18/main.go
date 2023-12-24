package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/matrix"
)

type planItem struct {
	direction byte
	distance  int
	color     string
}

func (p planItem) String() string {
	return fmt.Sprintf("%c %d %s", p.direction, p.distance, p.color)
}

type trench struct {
	m       matrix.Matrix[byte]
	current matrix.Point[byte]
}

func (t *trench) init() {
	t.m.ForEach(func(x, y int, _ byte) {
		t.m[x][y] = '.'
	})
}

func (t *trench) draw(p planItem) {
	// fmt.Printf("p : %s\n", p.color)
	for i := 0; i < p.distance; i++ {
		c := move(t.current, p.direction, 1)
		// fmt.Printf("current: %v, next: %v\n", t.current, c)
		t.m[c.X][c.Y] = '#'
		t.current = c
	}
}

func (t *trench) apply(p planItem) matrix.Point[byte] {
	c := move(t.current, p.direction, p.distance)
	t.current = c
	return c
}

// area return the area of the polygon
// uses [shoelace formula](https://en.wikipedia.org/wiki/Shoelace_formula)
func area(p []matrix.Point[byte]) int {
	var sum int
	for i := 0; i < len(p)-1; i++ {
		sum += p[i].X*p[i+1].Y - p[i].Y*p[i+1].X
	}
	sum += p[len(p)-1].X*p[0].Y - p[len(p)-1].Y*p[0].X
	return sum / 2
}

func perimeter(p []matrix.Point[byte]) int {
	var sum int
	for i := 0; i < len(p)-1; i++ {
		sum += abs(p[i].X-p[i+1].X) + abs(p[i].Y-p[i+1].Y)
	}
	sum += abs(p[len(p)-1].X-p[0].X) + abs(p[len(p)-1].Y-p[0].Y)
	return sum
}

// count return the number of points inside the polygon
// uses [pick's theorem](https://en.wikipedia.org/wiki/Pick%27s_theorem)
func count(area int, perimeter int) int {
	return abs(area - perimeter/2 - 1)
}

func (t *trench) Col() int {
	return len(t.m[0])
}

func (t *trench) passCorners() map[grid.Vec]struct{} {
	cn := make(map[grid.Vec]struct{})
	t.m.ForEach(func(x, y int, v byte) {
		// left-bottom corner and right-top corner
		// left-bottom corner
		if x-1 > 0 && y+1 < t.Col() {
			if v == '#' && t.m.Get(x-1, y) == '#' && t.m.Get(x, y+1) == '#' {
				cn[grid.Vec{X: x, Y: y}] = struct{}{}
			}
		}
		// right-top corner
		if x+1 < t.m.Rows() && y-1 > 0 {
			if v == '#' && t.m.Get(x+1, y) == '#' && t.m.Get(x, y-1) == '#' {
				cn[grid.Vec{X: x, Y: y}] = struct{}{}
			}
		}
	})
	return cn
}

func (t *trench) fill() {
	var insiders []grid.Vec
	pcs := t.passCorners()
	t.m.ForEach(func(x, y int, v byte) {
		if t.m.IsPointInside(
			matrix.Point[byte]{
				X:     x,
				Y:     y,
				Value: v,
			},
			matrix.WithCornerPass[byte](pcs),
			matrix.WithMatch[byte](func(p matrix.Point[byte]) bool {
				return p.Value == '#'
			}),
		) {
			insiders = append(insiders, grid.Vec{X: x, Y: y})
		}
	})
	for _, p := range insiders {
		t.m[p.X][p.Y] = '#'
	}
}

func (t *trench) countTrench() int {
	var count int
	t.m.ForEach(func(_, _ int, v byte) {
		if v == '#' {
			count++
		}
	})
	return count
}

func move(p matrix.Point[byte], d byte, distance int) matrix.Point[byte] {
	switch d {
	case 'R':
		return matrix.Point[byte]{
			X: p.X,
			Y: p.Y + distance,
		}
	case 'L':
		return matrix.Point[byte]{
			X: p.X,
			Y: p.Y - distance,
		}
	case 'U':
		return matrix.Point[byte]{
			X: p.X - distance,
			Y: p.Y,
		}
	case 'D':
		return matrix.Point[byte]{
			X: p.X + distance,
			Y: p.Y,
		}
	default:
		panic("invalid direction")
	}
}

func board(pis []planItem) (int, int) {
	var maxY, maxX int
	var minY, minX int
	var currentY, currentX int
	for _, p := range pis {
		switch p.direction {
		case 'U':
			currentX -= p.distance
			if currentX < minX {
				minX = currentX
			}
		case 'D':
			currentX += p.distance
			if currentX > maxX {
				maxX = currentX
			}
		case 'L':
			currentY -= p.distance
			if currentY < minY {
				minY = currentY
			}
		case 'R':
			currentY += p.distance
			if currentY > maxY {
				maxY = currentY
			}
		}
	}
	return maxX - minX + 1, maxY - minY + 1
}

func findLeftTop(pis []planItem) matrix.Point[byte] {
	var currentY, currentX int
	var minX, minY int
	for _, p := range pis {
		switch p.direction {
		case 'U':
			currentX -= p.distance
			if currentX < minX {
				minX = currentX
			}
		case 'D':
			currentX += p.distance
		case 'L':
			currentY -= p.distance
			if currentY < minY {
				minY = currentY
			}
		case 'R':
			currentY += p.distance
		}
	}
	return matrix.Point[byte]{
		X: minX,
		Y: minY,
	}
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	var planItems []planItem
	planItems = append(planItems, planItem{
		direction: 'R',
		distance:  1,
	})
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			panic("invalid input")
		}
		d := parts[0][0]
		planItems = append(planItems, planItem{
			direction: d,
			distance:  atoi(parts[1]),
			color:     parts[2],
		})
		return nil
	})
	row, col := board(planItems[1:])
	original := findLeftTop(planItems[1:])
	fmt.Printf("original: %+v\n", original)
	t := &trench{
		m: matrix.NewMatrix[byte](row, col),
		current: matrix.Point[byte]{
			X: -original.X,
			Y: -original.Y - 1,
		},
	}
	t.init()
	// fmt.Printf("row: %d, col: %d\n", row, col)
	// t.m.PrintEx("%c")
	for _, p := range planItems {
		t.draw(p)
	}
	// t.m.PrintEx("%c")
	// println()
	t.fill()
	// t.m.PrintEx("%c")
	fmt.Fprintf(os.Stdout, "p1: %d\n", t.countTrench())
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	var planItems []planItem
	planItems = append(planItems, planItem{
		direction: 'R',
		distance:  1,
	})
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			panic("invalid input")
		}
		var hexD string
		hexD = strings.TrimSuffix(parts[2], ")")
		hexD = strings.TrimPrefix(hexD, "(#")
		d := parseDirection(hexD[len(hexD)-1:])
		dst := parseHex(hexD[:len(hexD)-1])
		planItems = append(planItems, planItem{
			direction: d,
			distance:  int(dst),
			color:     hexD,
		})
		return nil
	})
	original := findLeftTop(planItems[1:])
	t := &trench{
		current: matrix.Point[byte]{
			X: -original.X,
			Y: -original.Y - 1,
		},
	}
	var ps []matrix.Point[byte]
	for _, p := range planItems {
		ps = append(ps, t.apply(p))
	}
	area, perimeter := area(ps), perimeter(ps)
	c := count(area, perimeter)
	fmt.Fprintf(os.Stdout, "p2: %d\n", c)
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

func parseHex(s string) int64 {
	var n int64
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*16 + int64(c-'0')
		} else if c >= 'a' && c <= 'f' {
			n = n*16 + int64(c-'a'+10)
		} else if c >= 'A' && c <= 'F' {
			n = n*16 + int64(c-'A'+10)
		}
	}
	return n
}

func parseDirection(d string) byte {
	switch d {
	case "0":
		return 'R'
	case "1":
		return 'D'
	case "2":
		return 'L'
	case "3":
		return 'U'
	default:
		panic("invalid direction")
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func main() {
	p1()
	p2()
}

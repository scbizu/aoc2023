package main

import (
	"context"
	"fmt"
	"os"

	"github.com/magejiCoder/set"
	"github.com/scbizu/aoc2022/helper/grid"
	"github.com/scbizu/aoc2022/helper/input"
	"github.com/scbizu/aoc2022/helper/matrix"
)

const (
	line, col = 140, 140
)

type pipeMaze struct {
	maze     matrix.Matrix[byte]
	startAt  matrix.Point[byte]
	distance int
	visited  map[grid.Vec]struct{}
}

func (p pipeMaze) traverse() pipeMaze {
	nbs := getNeighbor(p.maze, p.startAt.X, p.startAt.Y)
	for _, nb := range nbs {
		if nb.Value == 'S' || nb.Value == '.' {
			continue
		}
		if _, ok := p.visited[grid.Vec{X: nb.X, Y: nb.Y}]; ok {
			continue
		}
		p.visited[grid.Vec{X: nb.X, Y: nb.Y}] = struct{}{}
		newMaze := pipeMaze{
			maze:     p.maze,
			startAt:  nb,
			distance: p.distance + 1,
			visited:  p.visited,
		}
		// fmt.Println("new maze")
		// newMaze.maze.Print()
		newMaze = newMaze.traverse()
		if p.distance < newMaze.distance {
			p.distance = newMaze.distance
		}
		// fmt.Println("maze")
		// p.maze.Print()
	}
	return p
}

func availablePath(p grid.Vec, c byte) []grid.Vec {
	switch c {
	// north / south
	case '|':
		return []grid.Vec{
			{
				X: p.X - 1,
				Y: p.Y,
			},
			{
				X: p.X + 1,
				Y: p.Y,
			},
		}
	// east / west
	case '-':
		return []grid.Vec{
			{
				X: p.X,
				Y: p.Y - 1,
			},
			{
				X: p.X,
				Y: p.Y + 1,
			},
		}
	// north / east
	case 'L':
		return []grid.Vec{
			{
				X: p.X - 1,
				Y: p.Y,
			},
			{
				X: p.X,
				Y: p.Y + 1,
			},
		}
	// north / west
	case 'J':
		return []grid.Vec{
			{
				X: p.X - 1,
				Y: p.Y,
			},
			{
				X: p.X,
				Y: p.Y - 1,
			},
		}
	// south / west
	case '7':
		return []grid.Vec{
			{
				X: p.X + 1,
				Y: p.Y,
			},
			{
				X: p.X,
				Y: p.Y - 1,
			},
		}
	// south / east
	case 'F':
		return []grid.Vec{
			{
				X: p.X + 1,
				Y: p.Y,
			},
			{
				X: p.X,
				Y: p.Y + 1,
			},
		}
	default:
		return []grid.Vec{
			{
				X: p.X - 1,
				Y: p.Y,
			},
			{
				X: p.X,
				Y: p.Y + 1,
			},
			{
				X: p.X,
				Y: p.Y + 1,
			},
			{
				X: p.X,
				Y: p.Y - 1,
			},
		}
	}
}

func getNeighbor(maze matrix.Matrix[byte], x, y int) []matrix.Point[byte] {
	nbSet := set.New[grid.Vec]()
	for _, p := range maze.GetNeighbor(x, y) {
		nbSet.Add(grid.Vec{
			X: p.X,
			Y: p.Y,
		})
	}
	avaSet := set.New[grid.Vec]()
	for _, p := range availablePath(grid.Vec{X: x, Y: y}, maze.Get(x, y)) {
		avaSet.Add(p)
	}
	inSet := set.Intersection[grid.Vec](nbSet, avaSet)
	points := make([]matrix.Point[byte], 0, inSet.Size())
	inSet.Each(func(item grid.Vec) bool {
		points = append(points, matrix.Point[byte]{
			X:     item.X,
			Y:     item.Y,
			Value: maze.Get(item.X, item.Y),
		})
		return true
	})
	return points
}

func p1() {
	maze := matrix.NewMatrix[byte](line, col)
	txt := input.NewTXTFile("input.txt")
	startAt := matrix.Point[byte]{}
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		for idx, c := range line {
			if c == 'S' {
				startAt = matrix.Point[byte]{
					X: i, Y: idx,
					Value: byte(c),
				}
			}
			maze.Add(i, idx, byte(c))
		}
		return nil
	})
	pm := pipeMaze{
		maze:     maze,
		startAt:  startAt,
		distance: 0,
		visited: map[grid.Vec]struct{}{
			{X: startAt.X, Y: startAt.Y}: {},
		},
	}
	pm.maze.PrintEx("%c")
	pm = pm.traverse()
	fmt.Fprintf(os.Stdout, "p1: %d\n", (pm.distance+1)/2)
}

type squareMaze struct {
	maze    matrix.Matrix[byte]
	startAt matrix.Point[byte]
	visited map[grid.Vec]struct{}
	traces  map[grid.Vec]struct{}
}

func (m squareMaze) traverse() squareMaze {
	nbs := getNeighbor(m.maze, m.startAt.X, m.startAt.Y)
	for _, nb := range nbs {
		if nb.Value == 'S' || nb.Value == '.' {
			continue
		}
		if _, ok := m.visited[grid.Vec{X: nb.X, Y: nb.Y}]; ok {
			m.traces[grid.Vec{X: nb.X, Y: nb.Y}] = struct{}{}
			continue
		}
		m.visited[grid.Vec{X: nb.X, Y: nb.Y}] = struct{}{}
		newMaze := squareMaze{
			maze:    m.maze,
			startAt: nb,
			visited: m.visited,
			traces:  m.traces,
		}
		_ = newMaze.traverse()
	}
	return m
}

func (m squareMaze) mostRight() int {
	var right int
	for p := range m.visited {
		if p.Y > right {
			right = p.Y
		}
	}
	return right
}

// isPointInside 判断点是否在平面内(射线法实现)
func (m squareMaze) isPointInside(x, y int) bool {
	var intersectCount int
	i, j := x, y
	for {
		if i < 0 || j < 0 {
			break
		}
		c := m.maze.Get(i, j)
		// 从当前点到左上边界的射线上的交点
		// 1. 不同横竖的射线: 需要计算太多重合的时候方向的case，在这个场景下不如直接用斜线
		// 2. 把转角点看成曲线，则有: 7,L 对于 k = -1 的斜线来说，都是凸点，属于外侧;
		//    F,J 对于 k = -1 的斜线来说，都是凹点，属于内侧;
		//    S 既是凸点也是凹点
		if c == '7' || c == 'L' || c == 'S' {
			i = i - 1
			j = j - 1
			continue
		}
		if _, ok := m.traces[grid.Vec{X: i, Y: j}]; ok {
			intersectCount++
		}
		i--
		j--
	}
	// fmt.Printf("x,y = (%d,%d), intersectCount: %d\n", x, y, intersectCount)
	// 奇数: 在平面中
	// 偶数: 不在平面中
	return intersectCount%2 > 0 && intersectCount > 0
}

func p2() {
	maze := matrix.NewMatrix[byte](line, col)
	txt := input.NewTXTFile("input.txt")
	startAt := matrix.Point[byte]{}
	// var points []matrix.Point[byte]
	points := set.New[matrix.Point[byte]]()
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		if line == "" {
			return nil
		}
		for idx, c := range line {
			if c == 'S' {
				startAt = matrix.Point[byte]{
					X: i, Y: idx,
					Value: byte(c),
				}
			}
			if c == '.' {
				points.Add(matrix.Point[byte]{
					X: i,
					Y: idx,
				})
			}
			maze.Add(i, idx, byte(c))
		}
		return nil
	})
	pm := squareMaze{
		maze:    maze,
		startAt: startAt,
		visited: map[grid.Vec]struct{}{
			{X: startAt.X, Y: startAt.Y}: {},
		},
		traces: map[grid.Vec]struct{}{
			{X: startAt.X, Y: startAt.Y}: {},
		},
	}
	// pm.maze.PrintEx("%c")
	pm = pm.traverse()

	fmt.Printf("traces: %v\n", pm.traces)

	maze.ForEach(func(x, y int, _ byte) {
		if _, ok := pm.visited[grid.Vec{X: x, Y: y}]; !ok {
			points.Add(matrix.Point[byte]{
				X: x,
				Y: y,
			})
		}
	})

	var insidePoint int

	for _, p := range points.List() {
		if p.Y > pm.mostRight() || p.Y == 0 {
			continue
		}
		if pm.isPointInside(p.X, p.Y) {
			fmt.Printf("inside: %d %d\n", p.X, p.Y)
			insidePoint++
		}
	}
	fmt.Fprintf(os.Stdout, "p2: %d\n", insidePoint)
}

func main() {
	p1()
	p2()
}

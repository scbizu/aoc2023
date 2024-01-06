package main

import (
	"container/list"
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/math/sequence"
	"github.com/magejiCoder/magejiAoc/matrix"
	"github.com/magejiCoder/magejiAoc/queue"
)

func main() {
	p1()
	p2()
}

const (
	row = 131
	col = 131
)

const (
	maxStep = 327
)

type gardens struct {
	gq         map[grid.Vec]*garden
	gqCount    map[grid.Vec]int
	gardenSeq  map[grid.Vec][]int
	g1, g2     int
	max1, max2 int
	count      map[gardenPoint]struct{}
	remains    int64
}

type gardenPoint struct {
	gardenIndex grid.Vec
	point       matrix.Point[byte]
}

func (gs *gardens) walk2() {
	for {
		if gs.remains == 0 {
			break
		}
		gs.count = make(map[gardenPoint]struct{})
		gs.remains--
		gs.g1, gs.g2 = gs.g2, gs.g1
		// fmt.Printf("gardens: %d\n", len(gs.gq))
		next := map[grid.Vec][]matrix.Point[byte]{}
		for v, g := range gs.gq {
			if g.isReachMax {
				continue
			}
			// fmt.Printf("at garden: %+v\n", v)
			g.walk()
			// g.print()
			// println()
			// fmt.Printf("count: %d\n", g.count())
			// self
			// for p := range g.all {
			// 	gs.count[gardenPoint{
			// 		gardenIndex: v,
			// 		point:       p,
			// 	}] = struct{}{}
			// }
			// others
			for _, p := range g.borders {
				if gg, ok := gs.gq[v.Add(p.offset)]; ok {
					if gg.isReachMax {
						continue
					}
					if _, ok := gs.gq[v.Add(p.offset)].all[p.start]; !ok && gs.gq[v.Add(p.offset)].m.Get(p.start.X, p.start.Y) != '#' {
						// fmt.Printf("offset: %+v,start: %+v\n", v.Add(p.offset), p.start)
						next[v.Add(p.offset)] = append(next[v.Add(p.offset)], p.start)
						// gs.count[gardenPoint{
						// 	gardenIndex: v.Add(p.offset),
						// 	point:       p.start,
						// }] = struct{}{}
					}
					continue
				}
				// fmt.Printf("new garden: %+v\n", p)
				newG := &garden{
					m: g.m,
					paths: queue.Queue[matrix.Point[byte]]{
						List: list.New(),
					},
					all: make(map[matrix.Point[byte]]struct{}),
				}
				newG.all[p.start] = struct{}{}
				// newG.print()
				// println()
				// fmt.Printf("count: %d\n", newG.count())
				// gs.gq[v.Add(p.offset)] = newG
				next[v.Add(p.offset)] = append(next[v.Add(p.offset)], p.start)
				// gs.count[gardenPoint{
				// 	gardenIndex: v.Add(p.offset),
				// 	point:       p.start,
				// }] = struct{}{}
			}
			g.borders = []border{}
		}
		for v, ps := range next {
			for _, p := range ps {
				if _, ok := gs.gq[v]; !ok {
					gs.gq[v] = &garden{
						m:     gs.gq[grid.Vec{X: 0, Y: 0}].m,
						paths: queue.Queue[matrix.Point[byte]]{List: list.New()},
						all:   make(map[matrix.Point[byte]]struct{}),
					}
				}
				gs.gq[v].all[p] = struct{}{}
			}
		}
		// if gs.gq[grid.Vec{X: 1, Y: 0}] != nil {
		// 	fmt.Printf("garden : %d\n", len(gs.gq[grid.Vec{X: 0, Y: -1}].all))
		// }
		// gc := make(map[grid.Vec]int)
		// for g := range gs.count {
		// 	gc[g.gardenIndex]++
		// }
		// gs.gqCount = gc
		for v, g := range gs.gq {
			if g.isReachMax {
				continue
			}
			c := g.count()
			// gs.gardenSeq[v] = append(gs.gardenSeq[v], c)
			if ok := sequence.IsCircleStable[int](gs.gardenSeq[v]); ok {
				seq := sequence.NewCircleSeq[int](gs.gardenSeq[v])
				start, _, _ := seq.FindCircle(true)
				gs.max1 = gs.gardenSeq[v][start]
				gs.max2 = gs.gardenSeq[v][start+1]
				// fmt.Printf("remains: %d\n", gs.remains)
				if gs.max1 > gs.max2 {
					gs.g1++
				} else {
					gs.g2++
				}
				g.isReachMax = true
				delete(gs.gardenSeq, v)
			} else {
				gs.gardenSeq[v] = append(gs.gardenSeq[v], c)
			}
		}
	}
}

type border struct {
	start  matrix.Point[byte]
	offset grid.Vec
}

type garden struct {
	m          matrix.Matrix[byte]
	paths      queue.Queue[matrix.Point[byte]]
	all        map[matrix.Point[byte]]struct{}
	borders    []border
	isReachMax bool
}

func (g *garden) load() {
	for p := range g.all {
		g.paths.Push(p)
	}
	g.all = make(map[matrix.Point[byte]]struct{})
}

func (g *garden) print() {
	pMaze := matrix.NewMatrix[byte](row, col)
	g.m.ForEach(func(i, j int, v byte) {
		pMaze.Add(i, j, v)
	})
	pMaze.ForEach(func(i, j int, _ byte) {
		if _, ok := g.all[matrix.Point[byte]{
			X: i,
			Y: j,
		}]; ok {
			pMaze.Set(i, j, 'O')
		}
	})
	pMaze.PrintEx("%c")
}

func (g *garden) count() int {
	return len(g.all)
}

func (g *garden) walk() {
	g.load()
	for {
		if g.paths.Len() == 0 {
			break
		}
		p := g.paths.Pop()
		nbs := g.m.GetNeighbors(p.X, p.Y, true)
		for _, nb := range nbs {
			if nb.X < 0 {
				b := border{
					start: matrix.Point[byte]{
						X: row - 1,
						Y: nb.Y,
					},
					offset: grid.Vec{
						X: -1,
						Y: 0,
					},
				}
				g.borders = append(g.borders, b)
			}
			if nb.X > row-1 {
				b := border{
					start: matrix.Point[byte]{
						X: 0,
						Y: nb.Y,
					},
					offset: grid.Vec{
						X: 1,
						Y: 0,
					},
				}
				g.borders = append(g.borders, b)
			}
			if nb.Y < 0 {
				b := border{
					start: matrix.Point[byte]{
						X: nb.X,
						Y: col - 1,
					},
					offset: grid.Vec{
						X: 0,
						Y: -1,
					},
				}
				g.borders = append(g.borders, b)
			}
			if nb.Y > col-1 {
				b := border{
					start: matrix.Point[byte]{
						X: nb.X,
						Y: 0,
					},
					offset: grid.Vec{
						X: 0,
						Y: 1,
					},
				}
				g.borders = append(g.borders, b)
			}
			if nb.X < 0 || nb.X > row-1 || nb.Y < 0 || nb.Y > col-1 {
				continue
			}
			if g.m.Get(nb.X, nb.Y) == '#' {
				continue
			}
			if _, ok := g.all[matrix.Point[byte]{
				X: nb.X,
				Y: nb.Y,
			}]; ok {
				continue
			}
			g.all[matrix.Point[byte]{
				X: nb.X,
				Y: nb.Y,
			}] = struct{}{}
		}
	}
}

func (g *garden) TopCalibration() {
	newPoints := make(map[matrix.Point[byte]]struct{})
	for p := range g.all {
		newPoints[matrix.Point[byte]{
			X: p.X + row,
			Y: p.Y,
		}] = struct{}{}
	}
}

func (g *garden) LeftCalibration() {
	newPoints := make(map[matrix.Point[byte]]struct{})
	for p := range g.all {
		newPoints[matrix.Point[byte]{
			X: p.X,
			Y: p.Y + col,
		}] = struct{}{}
	}
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[byte](row, col)
	start := matrix.Point[byte]{}
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		for j, c := range line {
			if c == 'S' {
				start.X = i
				start.Y = j
				start.Value = byte(c)
			}
			m.Add(i, j, byte(c))
		}
		return nil
	})
	g := &garden{
		m: m,
		paths: queue.Queue[matrix.Point[byte]]{
			List: list.New(),
		},
		all: make(map[matrix.Point[byte]]struct{}),
	}
	g.all[start] = struct{}{}
	for i := 0; i < maxStep; i++ {
		g.walk()
		// g.print()
		// println()
	}
	fmt.Printf("p1: %d\n", g.count())
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[byte](row, col)
	start := matrix.Point[byte]{}
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		for j, c := range line {
			if c == 'S' {
				start.X = i
				start.Y = j
				start.Value = byte(c)
			}
			m.Add(i, j, byte(c))
		}
		return nil
	})
	steps := []int{65, 131 + 65, 2*131 + 65}
	totals := make([]int, 0, len(steps))
	for _, s := range steps {
		g := &garden{
			m: m,
			paths: queue.Queue[matrix.Point[byte]]{
				List: list.New(),
			},
			all: make(map[matrix.Point[byte]]struct{}),
		}
		g.all[start] = struct{}{}
		gardens := &gardens{
			gq:        map[grid.Vec]*garden{{X: 0, Y: 0}: g},
			remains:   int64(s),
			gqCount:   map[grid.Vec]int{},
			gardenSeq: map[grid.Vec][]int{},
		}
		gardens.walk2()
		var total int
		for _, c := range gardens.gardenSeq {
			// fmt.Printf("garden: %+v, count: %d\n", v, c)
			total += c[len(c)-1]
		}
		total += gardens.g2 * gardens.max1
		total += gardens.g1 * gardens.max2
		totals = append(totals, total)
	}
	// NOTE: 纯几何解
	// 1. 观察输入和最终输出可以知道，在面积为 131 x 131 的矩形中，从中间往外扩张 ，可以每次都得到一个"完美"的菱形
	// 第一次扩张到边界需要 (131 / 2) = 65 步，扩张完第二个矩形(边界也可以自动扩展)需要 131 + (131/2) = 196 步，以此类推
	// 2. 根据[皮克定理](https://github.com/magejiCoder/magejiAoc/blob/master/math/geometry/polygon.go#L30)
	//    可以知道，菱形中包含的点为 T = `area(面积) - perimeter(周长)/2 + 1 - count("#")` , 假设菱形边长为 n(或者说是第n个产生的菱形), 根据菱形面积公式，可以得到:
	//    最后的 T 与 n 可以构成一个二次函数: `T(n) = A * n^2 + B*n + C`
	// 根据拉格朗日插值法，和 T(0) (第0个菱形，也就是当前没有扩张的矩形), T(1) (左侧第1个生成的菱形), T(2) (左侧第二个生成的菱形) 可以得到 A, B, C 的值，进而可以得到 T(n) 的值
	// 对于最终的输出，需要计算第`26501365`步之后包含的点，而 `26501365 = 202300 * 131 + 65`，也就是说要计算扩张至第 202300 个菱形之后的点数
	// 根据上面推出饿 T(n) = A * n^2 + B*n + C, 可以得到 T(202300) = 最后答案
	y2, y1, y0 := totals[2], totals[1], totals[0]
	a, b, c := (y2-2*y1+y0)/2, (y2-y0)/2, y1
	n := (26501365 - 65) / 131
	fmt.Printf("p2: %d\n", a*n*n+b*n+c)
}

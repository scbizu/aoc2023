package main

import (
	"container/heap"
	"context"
	"fmt"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/matrix"
)

const (
	defaultLimit = 2
)

type direction uint8

const (
	up direction = iota + 1
	down
	left
	right
)

func (d direction) String() string {
	switch d {
	case up:
		return "⬆"
	case down:
		return "⬇"
	case right:
		return "->"
	case left:
		return "<-"
	}
	return ""
}

func (d direction) getOpposite() direction {
	switch d {
	case up:
		return down
	case down:
		return up
	case left:
		return right
	case right:
		return left
	default:
		panic("unknown direction")
	}
}

var caches = make(map[orderVec]int)

type vec struct {
	matrix.Point[int] // point
	d                 directionLimit
	di                direction
}

func (v vec) String() string {
	return fmt.Sprintf("(%d,%d)[%d](%v)", v.X, v.Y, v.Value, v.d)
}

type orderVec struct {
	p          matrix.Point[int]
	l, r, u, d int
	di         direction
}

var _ heap.Interface = (*vecQueue)(nil)

type vecQueue []*vec

func (vs vecQueue) Len() int {
	return len(vs)
}

func (vs vecQueue) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs vecQueue) Less(i, j int) bool {
	return vs[i].Value < vs[j].Value
}

func (vs *vecQueue) Push(v any) {
	item := v.(*vec)
	*vs = append(*vs, item)
}

func (vs *vecQueue) Pop() any {
	old := *vs
	n := len(old)
	item := old[n-1]
	*vs = old[0 : n-1]
	return item
}

type heatLossMap struct {
	m  matrix.Matrix[int]
	vq *vecQueue
}

func (h heatLossMap) neighbor(x, y int) []matrix.Point[int] {
	var ns []matrix.Point[int]
	if y < h.m.Cols()-1 {
		ns = append(ns, matrix.Point[int]{X: x, Y: y + 1, Value: h.m.Get(x, y+1)})
	}
	if x < h.m.Rows()-1 {
		ns = append(ns, matrix.Point[int]{X: x + 1, Y: y, Value: h.m.Get(x+1, y)})
	}
	if x > 0 {
		ns = append(ns, matrix.Point[int]{X: x - 1, Y: y, Value: h.m.Get(x-1, y)})
	}
	if y > 0 {
		ns = append(ns, matrix.Point[int]{X: x, Y: y - 1, Value: h.m.Get(x, y-1)})
	}
	return ns
}

func getDirection(p1, p2 matrix.Point[int]) direction {
	if p1.X == p2.X {
		if p1.Y > p2.Y {
			return left
		} else {
			return right
		}
	} else {
		if p1.X > p2.X {
			return up
		} else {
			return down
		}
	}
}

type directionLimit map[direction]int

func (h *heatLossMap) goDirection(p matrix.Point[int],
	d direction,
	distance int,
) matrix.Point[int] {
	switch d {
	case up:
		return matrix.Point[int]{
			X: p.X - distance,
			Y: p.Y,
		}
	case down:
		return matrix.Point[int]{
			X: p.X + distance,
			Y: p.Y,
		}
	case right:
		return matrix.Point[int]{
			X: p.X,
			Y: p.Y + distance,
		}
	case left:
		return matrix.Point[int]{
			X: p.X,
			Y: p.Y - distance,
		}
	}
	panic("invalid direction")
}

func (h *heatLossMap) getState(p matrix.Point[int], di direction, dl directionLimit) int {
	return caches[orderVec{
		p: matrix.Point[int]{
			X: p.X,
			Y: p.Y,
		},
		l:  dl[left],
		r:  dl[right],
		u:  dl[up],
		d:  dl[down],
		di: di,
	}]
}

func (h *heatLossMap) storeState(p matrix.Point[int], di direction, dl directionLimit, v int) {
	if _, ok := caches[orderVec{
		p: matrix.Point[int]{
			X: p.X,
			Y: p.Y,
		},
		l:  dl[left],
		r:  dl[right],
		u:  dl[up],
		d:  dl[down],
		di: di,
	}]; !ok {
		caches[orderVec{
			p: matrix.Point[int]{
				X: p.X,
				Y: p.Y,
			},
			l:  dl[left],
			r:  dl[right],
			u:  dl[up],
			d:  dl[down],
			di: di,
		}] = v
		heap.Push(h.vq, &vec{
			Point: matrix.Point[int]{
				X:     p.X,
				Y:     p.Y,
				Value: v,
			},
			d:  dl,
			di: di,
		})
	}
}

func (h *heatLossMap) fastPathWithRange() int {
	if h.vq.Len() == 0 {
		panic("no path")
	}
	p := heap.Pop(h.vq).(*vec)
	dl := p.d
	value := h.getState(p.Point, p.di, dl)
	fmt.Printf("pop: (%d,%d,%d)\n", p.X, p.Y, value)
	if p.X == h.m.Rows()-1 && p.Y == h.m.Cols()-1 && dl[p.di] <= 6 {
		return value
	}
	vecs := h.neighbor(p.X, p.Y)

	for _, v := range vecs {
		newdl := make(directionLimit)
		di := getDirection(p.Point, matrix.Point[int]{X: v.X, Y: v.Y})
		if di == p.di.getOpposite() {
			continue
		}
		for d := range dl {
			if d == di {
				newdl[direction(d)] = dl[direction(d)] - 1
			} else {
				newdl[direction(d)] = 10
			}
		}
		if dl[direction(p.di)] <= 6 && di != p.di {
			h.storeState(v, di, newdl, value+v.Value)
		}
		limit, ok := dl[di]
		if ok && limit == 0 {
			continue
		}
		if di == p.di {
			h.storeState(v, di, newdl, value+v.Value)
		}
	}
	return h.fastPathWithRange()
}

func (h *heatLossMap) fastPath() int {
	if h.vq.Len() == 0 {
		panic("no path")
	}
	p := heap.Pop(h.vq).(*vec)
	dl := p.d
	value := caches[orderVec{
		p: matrix.Point[int]{
			X: p.X,
			Y: p.Y,
		},
		l:  dl[left],
		r:  dl[right],
		u:  dl[up],
		d:  dl[down],
		di: p.di,
	}]
	if p.X == h.m.Rows()-1 && p.Y == h.m.Cols()-1 {
		return value
	}
	vecs := h.neighbor(p.X, p.Y)
	for _, v := range vecs {
		newdl := make(directionLimit)
		di := getDirection(p.Point, matrix.Point[int]{X: v.X, Y: v.Y})
		if di == p.di.getOpposite() {
			continue
		}
		if limit, ok := dl[di]; ok && limit == 0 {
			continue
		}
		for d := range dl {
			if d == di {
				newdl[direction(d)] = dl[direction(d)] - 1
			} else {
				newdl[direction(d)] = defaultLimit + 1
			}
		}
		_, ok := caches[orderVec{
			p: matrix.Point[int]{
				X: v.X,
				Y: v.Y,
			},
			l:  newdl[left],
			r:  newdl[right],
			u:  newdl[up],
			d:  newdl[down],
			di: di,
		}]
		if ok {
			continue
		} else {
			caches[orderVec{
				p: matrix.Point[int]{
					X: v.X,
					Y: v.Y,
				},
				l:  newdl[left],
				r:  newdl[right],
				u:  newdl[up],
				d:  newdl[down],
				di: di,
			}] = value + v.Value
			heap.Push(h.vq, &vec{
				Point: matrix.Point[int]{
					X:     v.X,
					Y:     v.Y,
					Value: value + v.Value,
				},
				d:  newdl,
				di: di,
			})
		}
	}
	return h.fastPath()
}

const (
	row = 141
	col = 141
)

func p2() {
	caches = make(map[orderVec]int)
	vq := &vecQueue{}
	heap.Init(vq)
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[int](row, col)
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		for j, v := range line {
			m.Add(i, j, int(v-'0'))
		}
		return nil
	})
	hl := &heatLossMap{
		m:  m,
		vq: vq,
	}
	caches[orderVec{
		p: matrix.Point[int]{
			X: 0,
			Y: 0,
		},
		l:  10,
		r:  10,
		u:  10,
		d:  10,
		di: right,
	}] = 0
	caches[orderVec{
		p: matrix.Point[int]{
			X: 0,
			Y: 0,
		},
		l:  10,
		r:  10,
		u:  10,
		d:  10,
		di: down,
	}] = 0
	v := &vec{
		Point: matrix.Point[int]{
			X:     0,
			Y:     0,
			Value: m.Get(0, 0),
		},
		d: map[direction]int{
			up:    10,
			down:  10,
			left:  10,
			right: 10,
		},
		di: right,
	}
	heap.Push(vq, v)
	v2 := &vec{
		Point: matrix.Point[int]{
			X:     0,
			Y:     0,
			Value: m.Get(0, 0),
		},
		d: map[direction]int{
			up:    10,
			down:  10,
			left:  10,
			right: 10,
		},
		di: down,
	}
	heap.Push(vq, v2)
	res := hl.fastPathWithRange()
	fmt.Printf("p2: %d\n", res)
}

func p1() {
	caches = make(map[orderVec]int)
	vq := &vecQueue{}
	heap.Init(vq)
	txt := input.NewTXTFile("input.txt")
	m := matrix.NewMatrix[int](row, col)
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		for j, v := range line {
			m.Add(i, j, int(v-'0'))
		}
		return nil
	})
	hl := &heatLossMap{
		m:  m,
		vq: vq,
	}
	caches[orderVec{
		p: matrix.Point[int]{
			X: 0,
			Y: 0,
		},
		l:  defaultLimit + 1,
		r:  defaultLimit,
		u:  defaultLimit + 1,
		d:  defaultLimit + 1,
		di: right,
	}] = m.Get(0, 0)
	caches[orderVec{
		p: matrix.Point[int]{
			X: 0,
			Y: 0,
		},
		l:  defaultLimit + 1,
		r:  defaultLimit + 1,
		u:  defaultLimit + 1,
		d:  defaultLimit,
		di: down,
	}] = m.Get(0, 0)
	v := &vec{
		Point: matrix.Point[int]{
			X:     0,
			Y:     0,
			Value: m.Get(0, 0),
		},
		d: map[direction]int{
			up:    defaultLimit + 1,
			down:  defaultLimit + 1,
			left:  defaultLimit + 1,
			right: defaultLimit,
		},
		di: right,
	}
	heap.Push(vq, v)
	v.d = map[direction]int{
		up:    defaultLimit + 1,
		down:  defaultLimit,
		left:  defaultLimit + 1,
		right: defaultLimit + 1,
	}
	v.di = down
	heap.Push(vq, v)
	res := hl.fastPath()
	fmt.Printf("p1: %d\n", res-m.Get(0, 0))
}

func main() {
	p1()
	p2()
}

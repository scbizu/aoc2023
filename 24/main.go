package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/magejiCoder/magejiAoc/grid"
	"github.com/magejiCoder/magejiAoc/input"
)

const (
	minArea = 200000000000000
	maxArea = 400000000000000
	// minArea = 7
	// maxArea = 27
)

type testArea struct {
	minX int
	maxX int
	minY int
	maxY int
}

type point3d struct {
	x int
	y int
	z int
}

type velocity3d struct {
	// per nanosecond
	x int
	y int
	z int
}

type hailStone struct {
	index    int
	position point3d
	velocity velocity3d
}

type point2d[T comparable] struct {
	x T
	y T
}

// normalVector 表示两点之间的法向量
func normalVector(p1, p2 point3d) point3d {
	// 平面叉乘
	return point3d{
		x: p1.y*p2.z - p1.z*p2.y,
		y: p1.z*p2.x - p1.x*p2.z,
		z: p1.x*p2.y - p1.y*p2.x,
	}
}

func (hs hailStone) normalVector() grid.XYZVec {
	// 由两点的法向量计算公式`normalVector(p1,p2)`化简而来
	//  - n = normalVector(p1, p2)
	// |
	// - p2 = p1 + v
	return grid.XYZVec{
		X: hs.position.y*hs.velocity.z - hs.position.z*hs.velocity.y,
		Y: hs.position.z*hs.velocity.x - hs.position.x*hs.velocity.z,
		Z: hs.position.x*hs.velocity.y - hs.position.y*hs.velocity.x,
	}
}

// slopeXY 表示XY平面上的斜率
func (hs hailStone) slopeXY() float64 {
	return float64(hs.velocity.y) / float64(hs.velocity.x)
}

// shiftXY 表示XY平面上的截距
func (hs hailStone) shiftXY() float64 {
	return float64(hs.position.y) - hs.slopeXY()*float64(hs.position.x)
}

func crossXYAt(h1, h2 hailStone) (point2d[float64], bool) {
	// fmt.Printf("[%v]:%v, [%v]:%v\n", h1, h1.slopeXY(), h2, h2.slopeXY())
	if h1.slopeXY() == h2.slopeXY() {
		return point2d[float64]{
			x: 0,
			y: 0,
		}, false
	}
	x := (h2.shiftXY() - h1.shiftXY()) / (h1.slopeXY() - h2.slopeXY())
	y := h1.slopeXY()*x + h1.shiftXY()
	if x < float64(minArea) || x > float64(maxArea) || y < float64(minArea) || y > float64(maxArea) {
		return point2d[float64]{
			x: 0,
			y: 0,
		}, false
	}
	time1 := (x - float64(h1.position.x)) / float64(h1.velocity.x)
	time2 := (x - float64(h2.position.x)) / float64(h2.velocity.x)
	// fmt.Printf("time1: %f, time2: %f\n", time1, time2)
	if time1 < 0 || time2 < 0 {
		return point2d[float64]{
			x: 0,
			y: 0,
		}, false
	}
	return point2d[float64]{
		x: x,
		y: y,
	}, true
}

type plane[T comparable] struct {
	// Ax + By + Cz + D = 0
	A T
	B T
	C T
	D T
}

type pointXYZ[T comparable] struct {
	x T
	y T
	z T
}

// 直线与平面的交点
func crossXYZAt(p plane[int], h hailStone) (pointXYZ[float64], int) {
	// 由直线轨迹方程与平面方程联立求解
	//  - 轨迹方程: x = p0.x + v.x*t, y = p0.y + v.y*t, z = p0.z + v.z*t
	//  - 平面方程: Ax + By + Cz = 0
	//  - 代入直线方程得: A*(p0.x + v.x*t) + B*(p0.y + v.y*t) + C*(p0.z + v.z*t) = 0
	//  - 整理得: t = - (A*p0.x + B*p0.y + C*p0.z) / (A*v.x + B*v.y + C*v.z)
	// t := -(p.A*p0.x + p.B*p0.y + p.C*p0.z) / (p.A*v.x + p.B*v.y + p.C*v.z)
	p0 := h.position
	v := h.velocity
	var t1 big.Int
	t1.Mul(big.NewInt(int64(p.A)), big.NewInt(int64(p0.x)))
	var t2 big.Int
	t2.Mul(big.NewInt(int64(p.B)), big.NewInt(int64(p0.y)))
	var t3 big.Int
	t3.Mul(big.NewInt(int64(p.C)), big.NewInt(int64(p0.z)))
	var t4 big.Int
	t4.Add(&t1, &t2)
	var t5 big.Int
	t5.Add(&t4, &t3)
	var tm1 big.Int
	tm1.Mul(big.NewInt(int64(p.A)), big.NewInt(int64(v.x)))
	var tm2 big.Int
	tm2.Mul(big.NewInt(int64(p.B)), big.NewInt(int64(v.y)))
	var tm3 big.Int
	tm3.Mul(big.NewInt(int64(p.C)), big.NewInt(int64(v.z)))
	var tm4 big.Int
	tm4.Add(&tm1, &tm2)
	var tm5 big.Int
	tm5.Add(&tm4, &tm3)
	var tm big.Int
	tm.Div(&t5, &tm5)
	t := -tm.Int64()
	// fmt.Printf("f: %v;p0: %v, v: %v, t: %d\n", p, p0, v, t)
	crossAt := pointXYZ[float64]{
		x: float64(p0.x) + float64(v.x)*float64(t),
		y: float64(p0.y) + float64(v.y)*float64(t),
		z: float64(p0.z) + float64(v.z)*float64(t),
	}
	fmt.Printf("cross: %v\n", crossAt)
	return crossAt, int(t)
}

func extractHailstone(raw string) hailStone {
	raw = strings.ReplaceAll(raw, " ", "")
	parts := strings.Split(raw, "@")
	position := strings.Split(parts[0], ",")
	pos := point3d{
		x: input.Atoi(position[0]),
		y: input.Atoi(position[1]),
		z: input.Atoi(position[2]),
	}
	velocity := strings.Split(parts[1], ",")
	v := velocity3d{
		x: input.Atoi(velocity[0]),
		y: input.Atoi(velocity[1]),
		z: input.Atoi(velocity[2]),
	}
	return hailStone{
		position: pos,
		velocity: v,
	}
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	var stones []hailStone
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		hs := extractHailstone(line)
		hs.index = i
		stones = append(stones, hs)
		return nil
	})
	// fmt.Printf("stones: %v\n", stones)
	var crossCount int
	for _, s1 := range stones {
		for _, s2 := range stones {
			if s1.index == s2.index {
				continue
			}
			if _, ok := crossXYAt(s1, s2); ok {
				// fmt.Printf("s1[%d], s2[%d], p[%v]\n", s1.index, s2.index, p)
				crossCount++
			}
		}
	}
	fmt.Printf("p1: %d\n", crossCount/2)
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	var stones []hailStone
	txt.ReadByLineEx(context.Background(), func(i int, line string) error {
		hs := extractHailstone(line)
		hs.index = i
		stones = append(stones, hs)
		return nil
	})
	origin := stones[0]
	stone0 := stones[0]
	// 设 stone0 是静止的
	// 以 stone0 为标准，所有其它的 stone 变为相对于 stone0 的位置和速度
	for i := range stones {
		if i == 0 {
			continue
		}
		stones[i] = hailStone{
			index: i,
			position: point3d{
				x: stones[i].position.x - stone0.position.x,
				y: stones[i].position.y - stone0.position.y,
				z: stones[i].position.z - stone0.position.z,
			},
			velocity: velocity3d{
				x: stones[i].velocity.x - stone0.velocity.x,
				y: stones[i].velocity.y - stone0.velocity.y,
				z: stones[i].velocity.z - stone0.velocity.z,
			},
		}
	}
	// 从相对位置来看，stone0是固定不动的，所以平面必定经过stone0，即(0,0,0)
	stone0 = hailStone{
		index: 0,
		position: point3d{
			x: 0,
			y: 0,
			z: 0,
		},
		velocity: velocity3d{
			x: 0,
			y: 0,
			z: 0,
		},
	}
	stone1 := stones[1]
	n := stone1.normalVector()
	// 以 stone0 作为平面上一点，n 为法向量，构建平面方程
	// 设: Ax + By + Cz + D = 0
	// 因为一定经过(0,0,0)，所以 D = 0
	// fmt.Printf("f(x,y,z) = %d*x + %d*y + %d*z = 0\n", n.X, n.Y, n.Z)
	stone2, stone3 := stones[2], stones[3]
	// 计算与 stone2 和 stone3 的交点与时间差
	cross2, t2 := crossXYZAt(plane[int]{A: n.X, B: n.Y, C: n.Z}, stone2)
	cross3, t3 := crossXYZAt(plane[int]{A: n.X, B: n.Y, C: n.Z}, stone3)
	// 设石头的速度为 vx,vy,vz
	// vx = delta(x) / delta(t)
	vx := (cross2.x - cross3.x) / float64(t2-t3)
	vy := (cross2.y - cross3.y) / float64(t2-t3)
	vz := (cross2.z - cross3.z) / float64(t2-t3)
	fmt.Printf("system:v: %.1f,%.1f,%.1f\n", vx, vy, vz)
	fmt.Printf("origin:v: %.1f,%.1f,%.1f\n", vx+float64(origin.velocity.x), vy+float64(origin.velocity.y), vz+float64(origin.velocity.z))
	// 所以可以得出初始位置的坐标
	x0 := cross3.x - float64(t3)*vx
	y0 := cross3.y - float64(t3)*vy
	z0 := cross3.z - float64(t3)*vz
	// 再转换回原来的坐标系
	fmt.Printf("system:coord: %.1f, %.1f, %.1f\n", x0, y0, z0)
	originX0 := float64(origin.position.x) + x0
	originY0 := float64(origin.position.y) + y0
	originZ0 := float64(origin.position.z) + z0
	fmt.Printf("origin:coord: %.1f, %.1f, %.1f\n", x0+float64(origin.position.x), y0+float64(origin.position.y), z0+float64(origin.position.z))
	fmt.Printf("p2: %.1f\n", originX0+originY0+originZ0)
}

func main() {
	p1()
	p2()
}

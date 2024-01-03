package main

import (
	"container/list"
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
	"github.com/magejiCoder/magejiAoc/queue"
	"github.com/magejiCoder/magejiAoc/set"
)

func main() {
	// p1()
	p2()
}

type moduleKind uint8

const (
	notDetermined moduleKind = iota
	broadcaster
	flipFlop
	conjunction
	output
)

func (k moduleKind) String() string {
	switch k {
	case broadcaster:
		return "broadcaster"
	case flipFlop:
		return "ff"
	case conjunction:
		return "co"
	case output:
		return "out"
	}
	return ""
}

type iModule interface {
	run(from string, pulse ...pulse) pulse
	stop(p pulse) bool
	String() string
}

type flipFlopModule struct {
	state bool
}

func (f *flipFlopModule) String() string {
	return fmt.Sprintf("%t", f.state)
}

func (f *flipFlopModule) run(from string, ps ...pulse) pulse {
	res := no
	for _, p := range ps {
		if p == low {
			if f.state {
				f.state = !f.state
				res = low
			} else {
				f.state = !f.state
				res = high
			}
		} else {
			continue
		}
	}
	return res
}

func (f *flipFlopModule) stop(p pulse) bool {
	return p == high
}

type conjunctionModule struct {
	mems *set.Set[string]
}

func (c *conjunctionModule) run(from string, pulses ...pulse) pulse {
	for _, pulse := range pulses {
		if pulse == high {
			c.mems.Remove(from)
		}
		if c.mems.Size() == 0 {
			return low
		}
	}
	return high
}

func (c *conjunctionModule) stop(p pulse) bool {
	return false
}

func (c *conjunctionModule) String() string {
	return fmt.Sprintf("%+v", c.mems)
}

type broadcasterModule struct{}

func (b *broadcasterModule) run(from string, pulses ...pulse) pulse {
	return pulses[0]
}

func (b *broadcasterModule) stop(p pulse) bool {
	return false
}

func (b *broadcasterModule) String() string {
	return "broadcaster"
}

type outputModule struct{}

func (o *outputModule) run(from string, pulses ...pulse) pulse {
	for _, p := range pulses {
		if p == low {
			return low
		}
	}
	return high
}

func (o *outputModule) stop(p pulse) bool {
	return false
}

func (o *outputModule) String() string {
	return "output"
}

type pulse uint8

const (
	low pulse = iota
	high
	no
)

func (p pulse) String() string {
	switch p {
	case low:
		return "low"
	case high:
		return "high"
	default:
		return "unknown"
	}
}

type moduleNode struct {
	name   string
	kind   moduleKind
	ins    []pulse
	outs   []*moduleNode
	state  iModule
	parent string
}

type moduleGraph struct {
	nodes     []*moduleNode
	q         queue.Queue[*moduleNode]
	states    map[string]iModule
	inits     map[string][]pulse
	pulses    map[pulse]int
	exitNodes *set.Set[string]
}

func (g *moduleGraph) hasNode(from string) bool {
	for _, n := range g.nodes {
		if n.name == from {
			return true
		}
	}
	return false
}

func (g *moduleGraph) getNode(from string) int {
	for i, n := range g.nodes {
		if n.name == from {
			return i
		}
	}
	return -1
	// fmt.Printf("add node: %s\n", from)
	// g.nodes = append(g.nodes, &moduleNode{
	// 	name:  from,
	// 	kind:  output,
	// 	state: &outputModule{},
	// })
	// g.states[from] = &outputModule{}
	// return len(g.nodes) - 1
}

func (g *moduleGraph) countPulses() (int, int) {
	return g.pulses[high], g.pulses[low]
}

func (g *moduleGraph) run() {
	for g.q.Len() > 0 {
		node := g.q.Pop()
		from := node.name
		if _, ok := g.inits[from]; ok {
			for _, p := range g.inits[from] {
				g.pulses[p]++
			}
			node.ins = append(node.ins, g.inits[from]...)
			g.inits[from] = []pulse{}
		}
		// fmt.Printf("[%d]pop: %s,kind: %s,ins: %+v\n", g.q.Len(), from, node.kind, node.ins)
		var ok bool
		node.state, ok = g.states[node.name]
		if !ok {
			continue
		}
		// var stop int
		// for _, in := range node.ins {
		// 	if node.state.stop(in) {
		// 		stop++
		// 	}
		// }
		// if stop == len(node.ins) {
		// 	continue
		// }
		var out pulse
		switch node.kind {
		case flipFlop:
			ff := node.state.(*flipFlopModule)
			// fmt.Printf("before: %t,ins: %+v\n", ff.state, node.ins)
			out = ff.run(node.parent, node.ins...)
			if out == no {
				continue
			}
			g.states[node.name] = ff
			// node.ins = []pulse{}
			// fmt.Printf("after: %t\n", ff.state)
		case conjunction:
			co := node.state.(*conjunctionModule)
			// fmt.Printf("[%s]before: %+v,parent: %s,ins: %+v\n", node.name, co.mems, node.parent, node.ins)
			out = co.run(node.parent, node.ins[0])
			// if node.ins[0] == low && node.name == "zb" {
			// 	fmt.Println("zb")
			// }
			g.states[node.name] = co
			// node.ins = node.ins[1:]
			// fmt.Printf("[%s]after: %+v\n", node.name, co.mems)
		default:
			panic("unknown module kind")
		}

		for _, n := range node.outs {
			// if out == low {
			// 	if (node.name == "lg" && n.name == "rr") || (node.name == "st" && n.name == "zb") ||
			// 		(node.name == "gr" && n.name == "js") || (node.name == "bn" && n.name == "bs") {
			// 		// fmt.Printf("found %s\n", node.name)
			// 		fmt.Printf("%s -> %s -> %s\n", node.name, out, n.name)
			// 		g.exitNodes.Add(node.name)
			// 	}
			// }
			fmt.Printf("%s -> %s -> %s\n", node.name, out, n.name)
			g.pulses[out]++
			n.parent = node.name
			if g.getNode(n.name) > 0 {
				n.outs = g.nodes[g.getNode(n.name)].outs
			}
			n.ins = []pulse{out}
			g.q.Push(n)
		}
	}
}

type edge struct {
	from string
	to   string
}

func copyMap(ins map[string][]pulse) map[string][]pulse {
	res := make(map[string][]pulse)
	for k, v := range ins {
		res[k] = v
	}
	return res
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	var edges []edge
	var start []string
	kindMap := make(map[string]moduleKind)
	insMap := make(map[string][]pulse)
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		parts := strings.Split(line, " -> ")
		md, down := parts[0], parts[1]
		if md == "broadcaster" {
			ds := strings.Split(down, ", ")
			for _, d := range ds {
				insMap[d] = append(insMap[d], low)
				start = append(start, d)
			}
			return nil
		}
		name := md
		switch {
		case strings.HasPrefix(md, "%"):
			kindMap[md[1:]] = flipFlop
			name = md[1:]
		case strings.HasPrefix(md, "&"):
			kindMap[md[1:]] = conjunction
			name = md[1:]
		case md == "broadcaster":
			kindMap[md] = broadcaster
		}
		for _, d := range strings.Split(down, ", ") {
			edges = append(edges, edge{
				from: name,
				to:   d,
			})
			if d == "output" {
				kindMap[d] = output
			}
		}
		return nil
	})
	g := &moduleGraph{
		q: queue.Queue[*moduleNode]{
			List: list.New(),
		},
		states:    map[string]iModule{},
		inits:     copyMap(insMap),
		pulses:    map[pulse]int{},
		exitNodes: set.New[string](),
	}
	for _, e := range edges {
		// for mermaid graph
		// fmt.Printf("%s-->%s\n", e.from, e.to)
		nd := &moduleNode{}
		if e.to == "output" {
			g.nodes = append(g.nodes, &moduleNode{
				name:  "output",
				kind:  output,
				state: &outputModule{},
			})
			g.states["output"] = &outputModule{}
		}
		if kindMap[e.to] == conjunction {
			if _, ok := g.states[e.to]; !ok {
				g.states[e.to] = &conjunctionModule{
					mems: set.New[string](e.from),
				}
			} else {
				g.states[e.to].(*conjunctionModule).mems.Add(e.from)
			}
		}
		if kindMap[e.from] == flipFlop {
			g.states[e.from] = &flipFlopModule{}
		}
		if kindMap[e.from] == broadcaster {
			g.states[e.from] = &broadcasterModule{}
			nd.parent = "broadcaster"
		}
		nd.name = e.to
		nd.kind = kindMap[e.to]
		if !g.hasNode(e.from) {
			var im iModule
			switch kindMap[e.from] {
			case flipFlop:
				im = &flipFlopModule{}
			case conjunction:
				im = &conjunctionModule{}
			default:
				panic("unknown module kind")
			}
			var ins []pulse
			n := &moduleNode{
				name:  e.from,
				kind:  kindMap[e.from],
				outs:  []*moduleNode{nd},
				ins:   ins,
				state: im,
			}
			g.nodes = append(g.nodes, n)
		} else {
			idx := g.getNode(e.from)
			if idx > 0 {
				g.nodes[idx].outs = append(g.nodes[idx].outs, nd)
			}
		}
	}
	var h, l int
	for i := 0; i < 1000; i++ {
		g.pulses[high] = 0
		g.pulses[low] = 1
		for _, s := range start {
			g.q.Push(g.nodes[g.getNode(s)])
		}
		g.inits = copyMap(insMap)
		for idx := range g.nodes {
			g.nodes[idx].ins = []pulse{}
		}
		g.run()
		ph, pl := g.countPulses()
		// fmt.Printf("high: %d, low: %d\n", ph, pl)
		h += ph
		l += pl
		// fmt.Println()
	}
	fmt.Printf("p1: %d\n", h*l)
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	var edges []edge
	var start []string
	kindMap := make(map[string]moduleKind)
	insMap := make(map[string][]pulse)
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		parts := strings.Split(line, " -> ")
		md, down := parts[0], parts[1]
		if md == "broadcaster" {
			ds := strings.Split(down, ", ")
			for _, d := range ds {
				insMap[d] = append(insMap[d], low)
				start = append(start, d)
			}
			return nil
		}
		name := md
		switch {
		case strings.HasPrefix(md, "%"):
			kindMap[md[1:]] = flipFlop
			name = md[1:]
		case strings.HasPrefix(md, "&"):
			kindMap[md[1:]] = conjunction
			name = md[1:]
		case md == "broadcaster":
			kindMap[md] = broadcaster
		}
		for _, d := range strings.Split(down, ", ") {
			edges = append(edges, edge{
				from: name,
				to:   d,
			})
			if d == "output" {
				kindMap[d] = output
			}
		}
		return nil
	})
	g := &moduleGraph{
		q: queue.Queue[*moduleNode]{
			List: list.New(),
		},
		states:    map[string]iModule{},
		inits:     copyMap(insMap),
		pulses:    map[pulse]int{},
		exitNodes: set.New[string](),
	}
	for _, e := range edges {
		// for mermaid graph
		// fmt.Printf("%s-->%s\n", e.from, e.to)
		nd := &moduleNode{}

		if kindMap[e.to] == conjunction {
			if _, ok := g.states[e.to]; !ok {
				g.states[e.to] = &conjunctionModule{
					mems: set.New[string](e.from),
				}
			} else {
				g.states[e.to].(*conjunctionModule).mems.Add(e.from)
			}
		}
		if kindMap[e.from] == flipFlop {
			g.states[e.from] = &flipFlopModule{}
		}
		if kindMap[e.from] == broadcaster {
			g.states[e.from] = &broadcasterModule{}
		}
		nd.name = e.to
		nd.kind = kindMap[e.to]
		if !g.hasNode(e.from) {
			var im iModule
			switch kindMap[e.from] {
			case flipFlop:
				im = &flipFlopModule{}
			case conjunction:
				im = &conjunctionModule{}
			default:
				panic("unknown module kind")
			}
			var ins []pulse
			n := &moduleNode{
				name:  e.from,
				kind:  kindMap[e.from],
				outs:  []*moduleNode{nd},
				ins:   ins,
				state: im,
			}
			g.nodes = append(g.nodes, n)
		} else {
			idx := g.getNode(e.from)
			if idx > 0 {
				g.nodes[idx].outs = append(g.nodes[idx].outs, nd)
			}
		}
	}

	var count int

	for i := 0; i < 4000; i++ {
		count++
		fmt.Printf("push: %d\n", count)
		for _, s := range start {
			g.q.Push(g.nodes[g.getNode(s)])
		}
		g.inits = copyMap(insMap)
		for idx := range g.nodes {
			g.nodes[idx].ins = []pulse{}
		}
		g.run()
		fmt.Printf("next\n")
	}

	// var count int
	// // totals := make([]int64, 0, 4)
	// for {
	// 	count++
	// 	for _, s := range start {
	// 		g.q.Push(s)
	// 	}
	// 	g.inits = copyMap(insMap)
	// 	for idx := range g.nodes {
	// 		g.nodes[idx].ins = []pulse{}
	// 	}
	// 	fmt.Printf("run: %d\n", count)
	// 	g.run()

	// 	// if len(totals) == 4 {
	// 	// 	break
	// 	// }
	// 	// if g.exitNodes.Has("lg") {
	// 	// 	lg := int64(count)
	// 	// 	fmt.Printf("found lg: %d\n", lg)
	// 	// 	totals = append(totals, lg)
	// 	// 	g.exitNodes.Remove("lg")
	// 	// }
	// 	// if g.exitNodes.Has("st") {
	// 	// 	st := int64(count)
	// 	// 	fmt.Printf("found st: %d\n", st)
	// 	// 	totals = append(totals, st)
	// 	// 	g.exitNodes.Remove("st")
	// 	// }
	// 	// if g.exitNodes.Has("gr") {
	// 	// 	gr := int64(count)
	// 	// 	fmt.Printf("found gr: %d\n", gr)
	// 	// 	totals = append(totals, gr)
	// 	// 	g.exitNodes.Remove("gr")
	// 	// }
	// 	// if g.exitNodes.Has("bn") {
	// 	// 	bn := int64(count)
	// 	// 	fmt.Printf("found bn: %d\n", bn)
	// 	// 	totals = append(totals, bn)
	// 	// 	g.exitNodes.Remove("bn")
	// 	// }
	// }
}

// lcm calculate the least common multiple
func lcm(a, b int64) int64 {
	return a * b / gcd(a, b)
}

// gcd calculate the greatest common divisor
func gcd(a, b int64) int64 {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
)

func main() {
	p1()
	p2()
}

const (
	defaultMin = 1
	defaultMax = 4000
)

type rng struct {
	min int
	max int
}

type rateRng struct {
	x *rng
	m *rng
	a *rng
	s *rng
}

func (rg rateRng) IsEmpty() bool {
	return rg.x == nil && rg.m == nil && rg.a == nil && rg.s == nil
}

func (rg rateRng) String() string {
	var buf bytes.Buffer
	if rg.x != nil {
		fmt.Fprintf(&buf, "x=%d-%d,", rg.x.min, rg.x.max)
	}
	if rg.m != nil {
		fmt.Fprintf(&buf, "m=%d-%d,", rg.m.min, rg.m.max)
	}
	if rg.a != nil {
		fmt.Fprintf(&buf, "a=%d-%d,", rg.a.min, rg.a.max)
	}
	if rg.s != nil {
		fmt.Fprintf(&buf, "s=%d-%d,", rg.s.min, rg.s.max)
	}
	return strings.TrimSuffix(buf.String(), ",")
}

func count(rgs []rng) int {
	if intersect(rgs).max-intersect(rgs).min+1 < 0 {
		return 0
	}
	return intersect(rgs).max - intersect(rgs).min + 1
}

func intersect(rgs []rng) rng {
	var ir rng
	for _, r := range rgs {
		if ir.min == 0 {
			ir.min = r.min
		} else {
			ir.min = max(ir.min, r.min)
		}
		if ir.max == 0 {
			ir.max = r.max
		} else {
			ir.max = min(ir.max, r.max)
		}
	}
	return ir
}

func merge(rgs []rateRng) rateRng {
	var rr rateRng
	for _, r := range rgs {
		if r.x != nil {
			if rr.x == nil {
				rr.x = r.x
			} else {
				x := intersect([]rng{*rr.x, *r.x})
				rr.x = &x
			}
		}
		if r.m != nil {
			if rr.m == nil {
				rr.m = r.m
			} else {
				m := intersect([]rng{*rr.m, *r.m})
				rr.m = &m
			}
		}
		if r.a != nil {
			if rr.a == nil {
				rr.a = r.a
			} else {
				a := intersect([]rng{*rr.a, *r.a})
				rr.a = &a
			}
		}
		if r.s != nil {
			if rr.s == nil {
				rr.s = r.s
			} else {
				s := intersect([]rng{*rr.s, *r.s})
				rr.s = &s
			}
		}
	}
	return rr
}

func (r rateRng) Reverse() rateRng {
	var rr rateRng
	if r.x != nil {
		if r.x.min == defaultMin {
			rr.x = &rng{
				min: r.x.max + 1,
				max: defaultMax,
			}
		} else {
			rr.x = &rng{
				min: defaultMin,
				max: r.x.min - 1,
			}
		}
	}
	if r.m != nil {
		if r.m.min == defaultMin {
			rr.m = &rng{
				min: r.m.max + 1,
				max: defaultMax,
			}
		} else {
			rr.m = &rng{
				min: defaultMin,
				max: r.m.min - 1,
			}
		}
	}
	if r.a != nil {
		if r.a.min == defaultMin {
			rr.a = &rng{
				min: r.a.max + 1,
				max: defaultMax,
			}
		} else {
			rr.a = &rng{
				min: defaultMin,
				max: r.a.min - 1,
			}
		}
	}
	if r.s != nil {
		if r.s.min == defaultMin {
			rr.s = &rng{
				min: r.s.max + 1,
				max: defaultMax,
			}
		} else {
			rr.s = &rng{
				min: defaultMin,
				max: r.s.min - 1,
			}
		}
	}
	return rr
}

type (
	condition func(x, m, a, s int) (string, bool)
)

type possibility struct {
	name string
	rg   rateRng
}

type workflow struct {
	name       string
	conditions []condition
	poss       []possibility
}

func (w workflow) String() string {
	pos := bytes.NewBuffer(nil)
	for _, p := range w.poss {
		fmt.Fprintf(pos, "%s:", p.name)
		if p.rg.a != nil {
			fmt.Fprintf(pos, "a=%d-%d,", p.rg.a.min, p.rg.a.max)
		}
		if p.rg.m != nil {
			fmt.Fprintf(pos, "m=%d-%d,", p.rg.m.min, p.rg.m.max)
		}
		if p.rg.s != nil {
			fmt.Fprintf(pos, "s=%d-%d,", p.rg.s.min, p.rg.s.max)
		}
		if p.rg.x != nil {
			fmt.Fprintf(pos, "x=%d-%d,", p.rg.x.min, p.rg.x.max)
		}
	}
	return fmt.Sprintf("%s{%s}", w.name, strings.TrimSuffix(pos.String(), ","))
}

type rating struct {
	x int
	m int
	a int
	s int
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	var flows []workflow
	txt.ReadByBlock(context.Background(), "\n\n", func(block []string) error {
		if len(block) != 2 {
			panic("invalid input")
		}
		workflows, _ := block[0], block[1]
		wks := strings.Split(workflows, "\n")
		for _, wk := range wks {
			w := parseWorkflow(wk)
			flows = append(flows, w)
		}
		return nil
	})
	flowMap := make(map[string]workflow)
	for _, f := range flows {
		flowMap[f.name] = f
	}
	routeV2(flowMap["in"], flowMap, []rateRng{})
	var total int
	for _, rgs := range totalRgs {
		var srng, arng, mrng, xrng []rng
		for _, rg := range rgs {
			// fmt.Printf("rg: %+v\n", rg)
			if rg.s != nil {
				srng = append(srng, *rg.s)
			} else {
				srng = append(srng, rng{
					min: defaultMin,
					max: defaultMax,
				})
			}
			if rg.a != nil {
				arng = append(arng, *rg.a)
			} else {
				arng = append(arng, rng{
					min: defaultMin,
					max: defaultMax,
				})
			}
			if rg.m != nil {
				mrng = append(mrng, *rg.m)
			} else {
				mrng = append(mrng, rng{
					min: defaultMin,
					max: defaultMax,
				})
			}
			if rg.x != nil {
				xrng = append(xrng, *rg.x)
			} else {
				xrng = append(xrng, rng{
					min: defaultMin,
					max: defaultMax,
				})
			}
		}
		total += count(srng) * count(arng) * count(mrng) * count(xrng)
	}
	fmt.Fprintf(os.Stdout, "p2: %d\n", total)
}

func copyRgs(rgs []rateRng) []rateRng {
	var nrgs []rateRng
	nrgs = append(nrgs, rgs...)
	return nrgs
}

var totalRgs [][]rateRng

func routeV2(f workflow,
	flowMap map[string]workflow,
	rgs []rateRng,
) {
	fmt.Fprintf(os.Stdout, "flow: %s,rgs: %+v\n", f, rgs)
	for _, p := range f.poss {
		if p.name == "R" {
			continue
		}
		if p.name == "A" {
			rng := append(copyRgs(rgs), p.rg)
			totalRgs = append(totalRgs, rng)
			// fmt.Printf("stored: %+v, total: %+v\n", rng, totalRgs)
			continue
		}
		rng := append(rgs, p.rg)
		routeV2(flowMap[p.name], flowMap, rng)
	}
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	var flows []workflow
	var rates []rating
	txt.ReadByBlock(context.Background(), "\n\n", func(block []string) error {
		if len(block) != 2 {
			panic("invalid input")
		}
		workflows, ratings := block[0], block[1]
		wks := strings.Split(workflows, "\n")
		for _, wk := range wks {
			w := parseWorkflow(wk)
			flows = append(flows, w)
		}
		rts := strings.Split(ratings, "\n")
		for _, rt := range rts {
			r := parseRate(rt)
			rates = append(rates, r)
		}
		return nil
	})
	flowMap := make(map[string]workflow)
	for _, f := range flows {
		flowMap[f.name] = f
	}

	var accepts []rating

	for _, r := range rates {
		f := flowMap["in"]
		ff := route(flowMap, f, r)
		if ff == "A" {
			accepts = append(accepts, r)
			continue
		}
		if ff == "R" {
			continue
		}
	}

	var total int
	for _, a := range accepts {
		total += a.a + a.m + a.s + a.x
	}
	fmt.Fprintf(os.Stdout, "p1: %d\n", total)
}

func route(fm map[string]workflow,
	in workflow,
	r rating,
) string {
	var out string
	for _, c := range in.conditions {
		res, ok := c(r.x, r.m, r.a, r.s)
		if ok {
			out = res
			break
		}
	}
	if out == "A" {
		return out
	}
	if out == "R" {
		return out
	}
	return route(fm, fm[out], r)
}

func parseWorkflow(str string) workflow {
	w := workflow{}
	read := bytes.NewBuffer(nil)
	for {
		if len(str) == 0 {
			break
		}
		b := str[0]
		if b == '{' {
			w.name = read.String()
			read.Reset()
			break
		}
		str = str[1:]
		read.WriteByte(b)
	}
	ruleStr := strings.TrimPrefix(strings.TrimSuffix(str, "}"), "{")
	rules := strings.Split(ruleStr, ",")
	var reverseOps []rateRng
	for _, r := range rules {
		parts := strings.Split(r, ":")
		switch {
		case len(parts) == 2 && strings.Contains(parts[0], "<"):
			condParts := strings.Split(parts[0], "<")
			c := func(x, m, a, s int) (string, bool) {
				switch condParts[0] {
				case "x":
					if x < input.Atoi(condParts[1]) {
						return parts[1], true
					}
				case "m":
					if m < input.Atoi(condParts[1]) {
						return parts[1], true
					}
				case "a":
					if a < input.Atoi(condParts[1]) {
						return parts[1], true
					}
				case "s":
					if s < input.Atoi(condParts[1]) {
						return parts[1], true
					}
				}
				return "", false
			}
			w.conditions = append(w.conditions, c)
			var rr rateRng
			switch condParts[0] {
			case "x":
				rr.x = &rng{
					min: defaultMin,
					max: input.Atoi(condParts[1]) - 1,
				}
			case "m":
				rr.m = &rng{
					min: defaultMin,
					max: input.Atoi(condParts[1]) - 1,
				}
			case "a":
				rr.a = &rng{
					min: defaultMin,
					max: input.Atoi(condParts[1]) - 1,
				}
			case "s":
				rr.s = &rng{
					min: defaultMin,
					max: input.Atoi(condParts[1]) - 1,
				}
			}
			p := possibility{
				name: parts[1],
				rg:   rr,
			}
			if len(reverseOps) > 0 {
				p.rg = merge(append(reverseOps, p.rg))
			}
			w.poss = append(w.poss, p)
			reverseOps = append(reverseOps, rr.Reverse())
		case len(parts) == 2 && strings.Contains(parts[0], ">"):
			condParts := strings.Split(parts[0], ">")
			c := func(x, m, a, s int) (string, bool) {
				switch condParts[0] {
				case "x":
					if x > input.Atoi(condParts[1]) {
						return parts[1], true
					}
				case "m":
					if m > input.Atoi(condParts[1]) {
						return parts[1], true
					}
				case "a":
					if a > input.Atoi(condParts[1]) {
						return parts[1], true
					}
				case "s":
					if s > input.Atoi(condParts[1]) {
						return parts[1], true
					}
				}
				return "", false
			}
			w.conditions = append(w.conditions, c)
			var rr rateRng
			switch condParts[0] {
			case "x":
				rr.x = &rng{
					min: input.Atoi(condParts[1]) + 1,
					max: defaultMax,
				}
			case "m":
				rr.m = &rng{
					min: input.Atoi(condParts[1]) + 1,
					max: defaultMax,
				}
			case "a":
				rr.a = &rng{
					min: input.Atoi(condParts[1]) + 1,
					max: defaultMax,
				}
			case "s":
				rr.s = &rng{
					min: input.Atoi(condParts[1]) + 1,
					max: defaultMax,
				}
			}
			p := possibility{
				name: parts[1],
				rg:   rr,
			}
			if len(reverseOps) > 0 {
				p.rg = merge(append(reverseOps, p.rg))
			}
			w.poss = append(w.poss, p)
			reverseOps = append(reverseOps, rr.Reverse())
		case len(parts) == 1 && parts[0] == "A":
			c := func(x, m, a, s int) (string, bool) {
				return "A", true
			}
			w.conditions = append(w.conditions, c)
			p := possibility{
				name: "A",
				rg:   merge(reverseOps),
			}
			w.poss = append(w.poss, p)
		case len(parts) == 1 && parts[0] == "R":
			c := func(x, m, a, s int) (string, bool) {
				return "R", true
			}
			w.conditions = append(w.conditions, c)
			p := possibility{
				name: "R",
				rg:   merge(reverseOps),
			}
			w.poss = append(w.poss, p)
		case len(parts) == 1:
			c := func(x, m, a, s int) (string, bool) {
				return parts[0], true
			}
			w.conditions = append(w.conditions, c)
			p := possibility{
				name: parts[0],
				rg:   merge(reverseOps),
			}
			w.poss = append(w.poss, p)
		default:
			panic("invalid rule")
		}
	}

	return w
}

func parseRate(rate string) rating {
	var r rating
	rate = strings.TrimPrefix(strings.TrimSuffix(rate, "}"), "{")
	parts := strings.Split(rate, ",")
	for _, p := range parts {
		rateParts := strings.Split(p, "=")
		switch rateParts[0] {
		case "x":
			r.x = input.Atoi(rateParts[1])
		case "m":
			r.m = input.Atoi(rateParts[1])
		case "a":
			r.a = input.Atoi(rateParts[1])
		case "s":
			r.s = input.Atoi(rateParts[1])
		default:
			panic("invalid rate")
		}
	}
	return r
}

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
)

func main() {}

type moduleKind uint8

const (
	notDetermined moduleKind = iota
	broadcaster
	flipFlop
	conjunction
)

type iModule interface {
	run(pulse bool) bool
}

type flipFlopModule struct {
	state bool
}

func (f *flipFlopModule) run(pulse bool) bool {
	if pulse {
		return pulse
	}
	if f.state {
		f.state = !f.state
		return false
	}
	f.state = !f.state
	return true
}

type conjunctionModule struct {
	rmInputs int
}

func (c *conjunctionModule) run(pulse bool) bool {
	c.rmInputs--
	return c.rmInputs == 0
}

type module struct {
	name       string
	kind       moduleKind
	mods       []iModule
	downstream []string
}

func newModuleFromRaw(
	raw string,
	downs []string,
) *module {
	m := &module{}
	switch {
	case raw == "broadcaster":
		m.kind = broadcaster
		m.name = raw
	case strings.HasPrefix(raw, "%"):
		m.kind = flipFlop
		m.name = strings.TrimPrefix(raw, "%")
	case strings.HasPrefix(raw, "&"):
		m.kind = conjunction
		m.name = strings.TrimPrefix(raw, "&")
	}
	m.downstream = append(m.downstream, downs...)
	return m
}

type workflow struct {
	modules    []*module
	highPulses int
	lowPulses  int
}

func (w *workflow) run() {
	for _, m := range w.modules {
		fmt.Printf("running %s\n", m.name)
		for _, mod := range m.mods {
			if mod.run(true) {
				w.highPulses++
			} else {
				w.lowPulses++
			}
		}
	}
}

func getInputs(mm map[string]*module, name string) int {
	var inputs int
	for _, m := range mm {
		for _, d := range m.downstream {
			if d == name {
				inputs++
				break
			}
		}
	}
	return inputs
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	var modules []*module
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		parts := strings.Split(line, " -> ")
		md, down := parts[0], parts[1]
		modules = append(modules, newModuleFromRaw(md, strings.Split(down, ", ")))
		return nil
	})
	moduleMap := make(map[string]*module)
	for _, m := range modules {
		moduleMap[m.name] = m
	}

	for _, mm := range moduleMap {
		for _, m := range mm.downstream {
			if _, ok := moduleMap[m]; !ok {
				panic("not found")
			} else {
				switch moduleMap[m].kind {
				case flipFlop:
					moduleMap[m].mods = append(moduleMap[m].mods, &flipFlopModule{
						state: false,
					})
				case conjunction:
					moduleMap[m].mods = append(moduleMap[m].mods, &conjunctionModule{
						rmInputs: getInputs(moduleMap, m),
					})
				}
			}
		}
	}

	wk := &workflow{
		modules: modules,
	}
	wk.run()
}

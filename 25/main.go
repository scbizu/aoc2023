package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/magejiCoder/magejiAoc/input"
)

type graph struct {
	g map[string]*unit
}

func (g graph) String() string {
	var sb strings.Builder
	for _, u := range g.g {
		for _, o := range u.outs {
			sb.WriteString(fmt.Sprintf("%s -> %s;\n", u.name, o))
		}
	}
	return sb.String()
}

// graph node
type unit struct {
	name string
	outs []string
	ins  []string
}

func main() {
	txt := input.NewTXTFile("input.txt")
	unitMap := make(map[string][]string)
	txt.ReadByLineEx(context.Background(), func(_ int, line string) error {
		parts := strings.Split(line, ":")
		un := parts[0]
		outParts := strings.Split(parts[1], " ")
		_ = outParts[0]
		outs := outParts[1:]
		unitMap[un] = outs
		return nil
	})
	gh := &graph{g: make(map[string]*unit)}
	for un, us := range unitMap {
		if _, ok := gh.g[un]; !ok {
			gh.g[un] = &unit{name: un}
		}
		gh.g[un].outs = us
		for _, u := range us {
			if _, ok := gh.g[u]; !ok {
				gh.g[u] = &unit{name: u}
			}
			gh.g[u].ins = append(gh.g[u].ins, un)
		}
	}
	fmt.Printf("%s", gh)
	// JUST FOR FUN
	// 看图说话吧，graphviz 牛逼
	// dot -Tsvg -Ksfdp graph.dot > graph.svg
	// 参见 graph.dot 和 graph.svg
}

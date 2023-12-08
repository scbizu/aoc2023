package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/scbizu/aoc2022/helper/input"
)

func parseFormat(raw string) (string, string, string) {
	r := regexp.MustCompile(`(.*)=\((.*),(.*)\)`)
	sub := r.FindStringSubmatch(raw)
	return sub[1], sub[2], sub[3]
}

type lr struct {
	l string
	r string
}

func parseNetwork2Map(raw string) map[string]lr {
	networks := make(map[string]lr)
	nets := strings.Split(raw, "\n")
	for _, net := range nets {
		if net == "" {
			continue
		}
		parentStr, leftStr, rightStr := parseFormat(removeBlank(net))
		networks[parentStr] = lr{
			l: leftStr,
			r: rightStr,
		}
	}
	return networks
}

func parseNetwork(raw string) ([]string, map[string]lr) {
	networks := make(map[string]lr)
	var froms []string
	nets := strings.Split(raw, "\n")
	for _, net := range nets {
		if net == "" {
			continue
		}
		parentStr, leftStr, rightStr := parseFormat(removeBlank(net))
		networks[parentStr] = lr{
			l: leftStr,
			r: rightStr,
		}
		if strings.HasSuffix(parentStr, "A") {
			froms = append(froms, parentStr)
		}
	}
	return froms, networks
}

func p1() {
	txt := input.NewTXTFile("input.txt")
	var instructions []byte
	var raw string
	txt.ReadByBlock(context.Background(), "\n\n", func(block []string) error {
		if len(block) != 2 {
			panic("invalid block")
		}
		instructions = []byte(block[0])
		raw = block[1]
		return nil
	})
	bt := parseNetwork2Map(raw)
	from := "AAA"
	var step int
	for {
		for _, in := range instructions {
			step++
			// fmt.Printf("from: %s\n", from)
			switch in {
			case 'L':
				if bt[from].l == "ZZZ" {
					fmt.Fprintf(os.Stdout, "p1: %d\n", step)
					return
				}
				from = bt[from].l
			case 'R':
				if bt[from].r == "ZZZ" {
					fmt.Fprintf(os.Stdout, "p1: %d\n", step)
					return
				}
				from = bt[from].r
			default:
				panic("invalid instruction")
			}
		}
	}
}

func p2() {
	txt := input.NewTXTFile("input.txt")
	var instructions []byte
	var raw string
	txt.ReadByBlock(context.Background(), "\n\n", func(block []string) error {
		if len(block) != 2 {
			panic("invalid block")
		}
		instructions = []byte(block[0])
		raw = block[1]
		return nil
	})
	froms, bt := parseNetwork(raw)
	var zSteps []int64
	for _, from := range froms {
		var step int64
		f := from
	INNER:
		for {
			for _, in := range instructions {
				step++
				fmt.Printf("from: %s\n", f)
				switch in {
				case 'L':
					if strings.HasSuffix(bt[f].l, "Z") {
						zSteps = append(zSteps, step)
						break INNER
					}
					f = bt[f].l
				case 'R':
					if strings.HasSuffix(bt[f].r, "Z") {
						zSteps = append(zSteps, step)
						break INNER
					}
					f = bt[f].r
				default:
					panic("invalid instruction")
				}
			}
		}
	}
	fmt.Printf("zSteps: %v\n", zSteps)
	allZ := zSteps[0]
	for i := 1; i < len(zSteps); i++ {
		allZ = lcm(allZ, zSteps[i])
	}
	fmt.Fprintf(os.Stdout, "p2: %d\n", allZ)
}

func main() {
	p1()
	p2()
}

func removeBlank(raw string) string {
	return strings.ReplaceAll(raw, " ", "")
}

// gcd calculate the greatest common divisor
func gcd(a, b int64) int64 {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

// lcm calculate the least common multiple
func lcm(a, b int64) int64 {
	return a * b / gcd(a, b)
}

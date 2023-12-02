package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scbizu/aoc2022/helper/input"
)

type gameLoader struct {
	id    uint32
	confs []gameConf
}

type gameConf struct {
	red, blue, green uint32
}

func (gl *gameLoader) parse(line string) error {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return errors.New("invalid line")
	}
	id, err := strconv.ParseInt(
		strings.TrimPrefix(parts[0], "Game"), 10, 64)
	if err != nil {
		return err
	}
	gl.id = uint32(id)
	setParts := strings.Split(parts[1], ";")
	if len(setParts) == 0 {
		return errors.New("invalid line")
	}
	for _, setPart := range setParts {
		conf := gameConf{}
		cubeParts := strings.Split(setPart, ",")
		if len(cubeParts) == 0 {
			continue
		}
		for _, cubePart := range cubeParts {
			switch {
			case strings.HasSuffix(cubePart, "blue"):
				blueCount, err := strconv.ParseInt(strings.TrimSuffix(cubePart, "blue"), 10, 64)
				if err != nil {
					return err
				}
				conf.blue = uint32(blueCount)
			case strings.HasSuffix(cubePart, "red"):
				redCount, err := strconv.ParseInt(strings.TrimSuffix(cubePart, "red"), 10, 64)
				if err != nil {
					return err
				}
				conf.red = uint32(redCount)
			case strings.HasSuffix(cubePart, "green"):
				greenCount, err := strconv.ParseInt(strings.TrimSuffix(cubePart, "green"), 10, 64)
				if err != nil {
					return err
				}
				conf.green = uint32(greenCount)
			default:
				return errors.New("unknown cube color")
			}
		}
		gl.confs = append(gl.confs, conf)
	}
	return nil
}

func (gl *gameLoader) canLoadConf(
	gc gameConf,
) bool {
	for _, conf := range gl.confs {
		if conf.blue > gc.blue || conf.red > gc.red || conf.green > gc.green {
			return false
		}
	}
	return true
}

func (gl *gameLoader) fewestGameConf() gameConf {
	gc := gameConf{}
	for _, conf := range gl.confs {
		if conf.red > gc.red {
			gc.red = conf.red
		}
		if conf.blue > gc.blue {
			gc.blue = conf.blue
		}
		if conf.green > gc.green {
			gc.green = conf.green
		}
	}
	return gc
}

func main() {
	txt := input.NewTXTFile("input.txt")

	ctx := context.Background()

	var idSum uint32
	var mSum uint32

	txt.ReadByLineEx(ctx, func(_ int, line string) error {
		if line == "" {
			return nil
		}
		line = removeBlank(line)
		gl := &gameLoader{}
		if err := gl.parse(line); err != nil {
			panic(err)
		}
		// part 1
		if gl.canLoadConf(gameConf{red: 12, blue: 14, green: 13}) {
			idSum += gl.id
		}
		// part 2
		fc := gl.fewestGameConf()
		mSum += fc.red * fc.blue * fc.green
		return nil
	})
	fmt.Fprintf(os.Stdout, "p1: idSum: %d\n", idSum)
	fmt.Fprintf(os.Stdout, "p2: mSum: %d\n", mSum)
}

func removeBlank(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

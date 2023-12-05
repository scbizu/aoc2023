package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/scbizu/aoc2022/helper/input"
)

type almanac struct {
	seeds                                           []int64
	soils, ferts, waters, lights, temps, hums, locs []int64
}

func (a *almanac) String() string {
	return fmt.Sprintf("seeds: %v\nsoils: %v\nferts: %v\nwaters: %v\nlights: %v\ntemps: %v\nhums: %v\nlocs: %v\n", a.seeds, a.soils, a.ferts, a.waters, a.lights, a.temps, a.hums, a.locs)
}

type i64Vec struct {
	from, to int64
}

type rangeAlmanac struct {
	seeds, soils, ferts, waters, lights, temps, hums, locs []i64Vec
	seedCount                                              int64
	almanacFn
}

func (ra *rangeAlmanac) Do() {
	for _, seed := range ra.seeds {
		var soils []int64
		for i := seed.from; i <= seed.to; i++ {
			soils = append(soils, ra.seed2Soil(i))
		}
		ra.soils = append(ra.soils, i64SliceToi64Vecs(soils)...)
		ra.soils = mergeI64Vecs(ra.soils)
	}
}

func mergeI64Vecs(vecs []i64Vec) []i64Vec {
	sort.Slice(vecs, func(i, j int) bool {
		return vecs[i].from < vecs[j].from
	})
	var merged []i64Vec
	merged = append(merged, vecs[0])
	for i := 1; i < len(vecs); i++ {
		if vecs[i].from <= merged[len(merged)-1].to {
			merged[len(merged)-1].to = vecs[i].to
		} else {
			merged = append(merged, vecs[i])
		}
	}
	return merged
}

func i64SliceToi64Vecs(s []int64) []i64Vec {
	var vecs []i64Vec
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
	from := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1]+1 {
			vecs = append(vecs, i64Vec{from, s[i-1]})
			from = s[i]
		}
	}
	vecs = append(vecs, i64Vec{from, s[len(s)-1]})
	return vecs
}

type almanacFn struct {
	seed2Soil, soil2Fert, fert2Water, water2Light, light2Temp, temp2Hum, hum2Loc func(int64) int64
}

func parseAlmanacFn(lines []string) (*almanacFn, error) {
	af := &almanacFn{}
	// soils
	sos := strings.Split(lines[1], "\n")
	rds := make([]rangeDirection, 0, len(sos)-1)
	for _, soilString := range sos[1:] {
		parts := strings.Split(soilString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid soil: %s", soilString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	af.seed2Soil = func(seed int64) int64 {
		rds := rds
		return doMapping(rds, seed)
	}
	// fertilizers
	ferts := strings.Split(lines[2], "\n")
	rds = make([]rangeDirection, 0, len(ferts)-1)
	for _, fertString := range ferts[1:] {
		parts := strings.Split(fertString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid fertilizer: %s", fertString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	af.soil2Fert = func(soil int64) int64 {
		rds := rds
		return doMapping(rds, soil)
	}
	// water
	waters := strings.Split(lines[3], "\n")
	rds = make([]rangeDirection, 0, len(waters)-1)
	for _, waterString := range waters[1:] {

		parts := strings.Split(waterString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid water: %s", waterString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}

	af.fert2Water = func(fert int64) int64 {
		rds := rds
		return doMapping(rds, fert)
	}

	// light
	lights := strings.Split(lines[4], "\n")
	rds = make([]rangeDirection, 0, len(lights)-1)
	for _, lightString := range lights[1:] {
		parts := strings.Split(lightString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid light: %s", lightString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}

	af.water2Light = func(water int64) int64 {
		rds := rds
		return doMapping(rds, water)
	}

	// temperature
	temperatures := strings.Split(lines[5], "\n")
	rds = make([]rangeDirection, 0, len(temperatures)-1)
	for _, temperatureString := range temperatures[1:] {
		parts := strings.Split(temperatureString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid temperature: %s", temperatureString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}

	af.light2Temp = func(light int64) int64 {
		rds := rds
		return doMapping(rds, light)
	}

	// humidity
	humidities := strings.Split(lines[6], "\n")
	rds = make([]rangeDirection, 0, len(humidities)-1)
	for _, humidityString := range humidities[1:] {
		parts := strings.Split(humidityString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid humidity: %s", humidityString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}

	af.temp2Hum = func(temp int64) int64 {
		rds := rds
		return doMapping(rds, temp)
	}

	// location
	locations := strings.Split(lines[7], "\n")
	rds = make([]rangeDirection, 0, len(locations)-1)
	for _, locationString := range locations[1:] {
		if locationString == "" {
			continue
		}
		parts := strings.Split(locationString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid location: %s", locationString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}

	af.hum2Loc = func(hum int64) int64 {
		rds := rds
		return doMapping(rds, hum)
	}
	return af, nil
}

func parseAlmanac(lines []string) (*almanac, error) {
	a := &almanac{}
	if len(lines) != 8 {
		return nil, fmt.Errorf("invalid almanac")
	}
	// seeds
	seedsString := strings.TrimPrefix(lines[0], "seeds: ")
	ss, err := splitString2Int64(seedsString, " ")
	if err != nil {
		return nil, fmt.Errorf("failed to parse seeds: %w", err)
	}
	a.seeds = ss
	// soils
	sos := strings.Split(lines[1], "\n")
	rds := make([]rangeDirection, 0, len(sos)-1)
	for _, soilString := range sos[1:] {
		parts := strings.Split(soilString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid soil: %s", soilString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}

	for _, seed := range a.seeds {
		a.soils = append(a.soils, doMapping(rds, seed))
	}
	// fertilizers
	ferts := strings.Split(lines[2], "\n")
	rds = make([]rangeDirection, 0, len(ferts)-1)
	for _, fertString := range ferts[1:] {
		parts := strings.Split(fertString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid fertilizer: %s", fertString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, soil := range a.soils {
		a.ferts = append(a.ferts, doMapping(rds, soil))
	}
	// water
	waters := strings.Split(lines[3], "\n")
	rds = make([]rangeDirection, 0, len(waters)-1)
	for _, waterString := range waters[1:] {
		parts := strings.Split(waterString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid water: %s", waterString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, fert := range a.ferts {
		a.waters = append(a.waters, doMapping(rds, fert))
	}
	// light
	lights := strings.Split(lines[4], "\n")
	rds = make([]rangeDirection, 0, len(lights)-1)
	for _, lightString := range lights[1:] {
		parts := strings.Split(lightString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid light: %s", lightString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, water := range a.waters {
		a.lights = append(a.lights, doMapping(rds, water))
	}
	// temperature
	temperatures := strings.Split(lines[5], "\n")
	rds = make([]rangeDirection, 0, len(temperatures)-1)
	for _, temperatureString := range temperatures[1:] {
		parts := strings.Split(temperatureString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid temperature: %s", temperatureString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, light := range a.lights {
		a.temps = append(a.temps, doMapping(rds, light))
	}
	// humidity
	humidities := strings.Split(lines[6], "\n")
	rds = make([]rangeDirection, 0, len(humidities)-1)
	for _, humidityString := range humidities[1:] {
		parts := strings.Split(humidityString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid humidity: %s", humidityString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, temperature := range a.temps {
		a.hums = append(a.hums, doMapping(rds, temperature))
	}
	// location
	locations := strings.Split(lines[7], "\n")
	rds = make([]rangeDirection, 0, len(locations)-1)
	for _, locationString := range locations[1:] {
		if locationString == "" {
			continue
		}
		parts := strings.Split(locationString, " ")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid location: %s", locationString)
		}
		rd := rangeDirection{
			source: atoi64(parts[1]),
			target: atoi64(parts[0]),
			rngLen: atoi64(parts[2]),
		}
		rds = append(rds, rd)
	}
	for _, humidity := range a.hums {
		a.locs = append(a.locs, doMapping(rds, humidity))
	}

	return a, nil
}

type rangeDirection struct {
	source int64
	target int64
	rngLen int64
}

func (rd rangeDirection) mapTo(source int64) int64 {
	if source > rd.source+rd.rngLen || source < rd.source {
		return source
	}
	return rd.target + (source - rd.source)
}

func doMapping(rds []rangeDirection, source int64) int64 {
	for _, rd := range rds {
		target := rd.mapTo(source)
		if target != source {
			return target
		}
	}
	return source
}

func main() {
	txt := input.NewTXTFile("input.txt")
	txt.ReadByBlock(context.Background(), "\n\n", func(lines []string) error {
		a, err := parseAlmanac(lines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse almanac: %v\n", err)
			return nil
		}
		minLoc := math.MaxInt64
		for _, loc := range a.locs {
			if loc < int64(minLoc) {
				minLoc = int(loc)
			}
		}
		fmt.Fprintf(os.Stdout, "p1: %d\n", minLoc)
		return nil
	})

	txt.ReadByBlock(context.Background(), "\n\n", func(lines []string) error {
		a, err := parseAlmanac(lines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse almanac: %v\n", err)
			return nil
		}
		minLoc := math.MaxInt64
		for _, loc := range a.locs {
			if loc < int64(minLoc) {
				minLoc = int(loc)
			}
		}
		fmt.Fprintf(os.Stdout, "p2: %d\n", minLoc)
		return nil
	})
}

func splitString2Int64(raw string, sep string) ([]int64, error) {
	raws := strings.Split(raw, sep)
	ints := make([]int64, len(raws))
	for i, raw := range raws {
		if _, err := fmt.Sscanf(raw, "%d", &ints[i]); err != nil {
			return nil, fmt.Errorf("failed to parse string to int64: %w", err)
		}
	}
	return ints, nil
}

func atoi64(raw string) int64 {
	var i int64
	if _, err := fmt.Sscanf(raw, "%d", &i); err != nil {
		panic(err)
	}
	return i
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [inputfile]\n", os.Args[0])
		os.Exit(2)
	}
	profileList, err := ParseProfiles(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %s\n Exits.\n", err.Error())
		os.Exit(2)
	}
	calcCoverage(profileList)
}

var IgnoredFiles = []string{"module.go", "registertxroutes.go"}

func calcCoverage(profileList []*Profile) {
	totalLines := 0
	coveredLines := 0
	for _, profile := range profileList {
		ignore := false
		for _, f := range IgnoredFiles {
			if strings.HasSuffix(profile.FileName, f) {
				ignore = true
				break
			}
		}
		if ignore {
			continue
		}
		for _, blk := range profile.Blocks {
			lineCount := blk.EndLine - blk.StartLine + 1
			totalLines += lineCount
			if blk.Count != 0 {
				coveredLines += lineCount
			}
		}
	}
	fmt.Printf("Coverage :%0.2f\n", float64(coveredLines)*100.0/float64(totalLines))
}

// Profile represents the profiling data for a specific file.
type Profile struct {
	FileName string
	Mode     string
	Blocks   []ProfileBlock
}

// ProfileBlock represents a single block of profiling data.
type ProfileBlock struct {
	StartLine, StartCol int
	EndLine, EndCol     int
	NumStmt, Count      int
}

type byFileName []*Profile

func (p byFileName) Len() int           { return len(p) }
func (p byFileName) Less(i, j int) bool { return p[i].FileName < p[j].FileName }
func (p byFileName) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// ParseProfiles parses profile data in the specified file and returns a
// Profile for each source file described therein.
func ParseProfiles(fileName string) ([]*Profile, error) {
	pf, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer pf.Close()

	files := make(map[string]*Profile)
	buf := bufio.NewReader(pf)
	// First line is "mode: foo", where foo is "set", "count", or "atomic".
	// Rest of file is in the format
	//	encoding/base64/base64.go:34.44,37.40 3 1
	// where the fields are: name.go:line.column,line.column numberOfStatements count
	s := bufio.NewScanner(buf)
	mode := ""
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			continue
		}
		const start = "mode: "
		if mode == "" {
			if !strings.HasPrefix(line, start) || line == start {
				return nil, fmt.Errorf("bad mode line: %v", line)
			}
			mode = line[len(start):]
			continue
		} else {
			if strings.HasPrefix(line, start) {
				continue
			}
		}
		m := lineRe.FindStringSubmatch(line)
		if m == nil {
			return nil, fmt.Errorf("line %q doesn't match expected format: %v", m, lineRe)
		}
		fn := m[1]
		p := files[fn]
		if p == nil {
			p = &Profile{
				FileName: fn,
				Mode:     mode,
			}
			files[fn] = p
		}
		p.Blocks = append(p.Blocks, ProfileBlock{
			StartLine: toInt(m[2]),
			StartCol:  toInt(m[3]),
			EndLine:   toInt(m[4]),
			EndCol:    toInt(m[5]),
			NumStmt:   toInt(m[6]),
			Count:     toInt(m[7]),
		})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	for _, p := range files {
		sort.Sort(blocksByStart(p.Blocks))
	}
	// Generate a sorted slice.
	profiles := make([]*Profile, 0, len(files))
	for _, profile := range files {
		profiles = append(profiles, profile)
	}
	sort.Sort(byFileName(profiles))
	return profiles, nil
}

type blocksByStart []ProfileBlock

func (b blocksByStart) Len() int      { return len(b) }
func (b blocksByStart) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b blocksByStart) Less(i, j int) bool {
	bi, bj := b[i], b[j]
	return bi.StartLine < bj.StartLine || bi.StartLine == bj.StartLine && bi.StartCol < bj.StartCol
}

var lineRe = regexp.MustCompile(`^(.+):([0-9]+).([0-9]+),([0-9]+).([0-9]+) ([0-9]+) ([0-9]+)$`)

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

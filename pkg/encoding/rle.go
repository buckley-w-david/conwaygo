package encoding

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// LoadFirstRoundFromRLE generates a new field from a rle-text-file
func LoadFieldFromFile(width, height int, filename string) *Field {
	var length int
	var field *Field
	x := 0
	y := 0
	finfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println(filename + " doesn't exist")
		return GenerateFirstRound(width, height)
	}
	if finfo.IsDir() {
		fmt.Println(filename + " is a directory")
		return GenerateFirstRound(width, height)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	xre := regexp.MustCompile("x ?= ?(\\d+)")
	yre := regexp.MustCompile("y ?= ?(\\d+)")
	lre := regexp.MustCompile(`(((\d*)([bo]+))([$!]*))`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			if line[0] == '#' {
				continue
			}
			xm := xre.FindStringSubmatch(line)
			ym := yre.FindStringSubmatch(line)

			if len(xm) == 2 && len(ym) == 2 {
				pint, _ := strconv.ParseInt(xm[1], 10, 64)
				width = int(pint)
				pint, _ = strconv.ParseInt(ym[1], 10, 64)
				height = int(pint)
				field = newField(width, height)
			}

			l := lre.FindAllStringSubmatch(line, -1)

			if len(l) > 0 {
				for _, sm := range l {
					if sm[3] == "" {
						length = 1
					} else {
						pint, _ := strconv.ParseInt(sm[3], 10, 64)
						length = int(pint)
					}
					for i := 1; i < length; i++ {
						if sm[4][0] == 'o' {
							field.setVitality(x, y, 1)
						}
						x++
					}
					for i := range sm[4] {
						if sm[4][i] == 'o' {
							field.setVitality(x, y, 1)
						}
						x++
					}
					if sm[5] != "" {
						x = 0
						for range sm[5] {
							y++
						}
					}

				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return field
}

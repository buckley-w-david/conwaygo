package encoding

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/buckley-w-david/conwaygo/pkg/conway"
)

var (
	xre *regexp.Regexp
	yre *regexp.Regexp
	lre *regexp.Regexp
)

func init() {
	xre = regexp.MustCompile("x ?= ?(\\d+)")
	yre = regexp.MustCompile("y ?= ?(\\d+)")
	lre = regexp.MustCompile(`(((\d*)([bo]+))([$!]*))`)
}

// LoadFieldFromFile generates a new field from a rle-text-file
func LoadFieldFromFile(filename string) (*conway.Field, error) {
	var (
		header string
		width  int
		height int
	)
	field := new(conway.Field)
	finfo, err := os.Stat(filename)

	if err != nil {
		return nil, err
	}
	if finfo.IsDir() {
		return nil, errors.New(filename + " is a directory")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		header = scanner.Text()
	} else {
		return nil, errors.New("No header found")
	}

	xm := xre.FindStringSubmatch(header)
	ym := yre.FindStringSubmatch(header)
	if len(xm) == 2 && len(ym) == 2 {
		pint, _ := strconv.ParseInt(xm[1], 10, 64)
		width = int(pint)
		pint, _ = strconv.ParseInt(ym[1], 10, 64)
		height = int(pint)
	} else {
		return nil, errors.New("Unable to parse header line")
	}
	field = conway.NewField(map[conway.Location]*conway.Cell{})
	fmt.Println(width, height)

	for scanner.Scan() {
		// line := scanner.Text()
	}
	//
	// 	xre := regexp.MustCompile("x ?= ?(\\d+)")
	// 	yre := regexp.MustCompile("y ?= ?(\\d+)")
	// 	lre := regexp.MustCompile(`(((\d*)([bo]+))([$!]*))`)
	//
	// 	for scanner.Scan() {
	// 		line := scanner.Text()
	// 		if len(line) > 0 {
	// 			if line[0] == '#' {
	// 				continue
	// 			}
	// 			xm := xre.FindStringSubmatch(line)
	// 			ym := yre.FindStringSubmatch(line)
	//
	// 			if len(xm) == 2 && len(ym) == 2 {
	// 				pint, _ := strconv.ParseInt(xm[1], 10, 64)
	// 				width = int(pint)
	// 				pint, _ = strconv.ParseInt(ym[1], 10, 64)
	// 				height = int(pint)
	// 				field = newField(width, height)
	// 			}
	//
	// 			l := lre.FindAllStringSubmatch(line, -1)
	//
	// 			if len(l) > 0 {
	// 				for _, sm := range l {
	// 					if sm[3] == "" {
	// 						length = 1
	// 					} else {
	// 						pint, _ := strconv.ParseInt(sm[3], 10, 64)
	// 						length = int(pint)
	// 					}
	// 					for i := 1; i < length; i++ {
	// 						if sm[4][0] == 'o' {
	// 							field.setVitality(x, y, 1)
	// 						}
	// 						x++
	// 					}
	// 					for i := range sm[4] {
	// 						if sm[4][i] == 'o' {
	// 							field.setVitality(x, y, 1)
	// 						}
	// 						x++
	// 					}
	// 					if sm[5] != "" {
	// 						x = 0
	// 						for range sm[5] {
	// 							y++
	// 						}
	// 					}
	//
	// 				}
	// 			}
	// 		}
	// 	}
	//
	// 	if err := scanner.Err(); err != nil {
	// 		fmt.Println(err)
	// 	}
	//
	// 	return field
	return field, nil
}

func SaveFieldToFile(field *conway.Field, filename string) error {
	return nil
}

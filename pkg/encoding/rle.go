package encoding

import (
	"bufio"
	"errors"
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

type rleScanner struct {
	*bufio.Scanner
}

func (r *rleScanner) Scan() bool {
	for r.Scanner.Scan() {
		line := r.Scanner.Text()
		if len(line) > 0 && line[0] != byte('#') {
			return true
		}
	}
	return false
}

// LoadFieldFromFile generates a new field from a rle-text-file
// Inspired by https://github.com/SimonWaldherr/cgolGo
func LoadFieldFromFile(filename string) (*conway.Field, error) {
	var (
		x      int
		y      int
		header string
	)
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
	scanner := rleScanner{bufio.NewScanner(file)}
	if scanner.Scan() {
		header = scanner.Text()
	} else {
		return nil, errors.New("Unable to find header")
	}
	xm := xre.FindStringSubmatch(header)
	ym := yre.FindStringSubmatch(header)

	if len(xm) == 2 && len(ym) == 2 {
	} else {
		return nil, errors.New("Unable to parse header")
	}
	var length int
	field := conway.NewField([]conway.Location{})
	for scanner.Scan() {
		line := scanner.Text()
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
						field.SetCell(conway.Location{X: x, Y: y}, true)
					}
					x++
				}
				for i := range sm[4] {
					if sm[4][i] == 'o' {
						field.SetCell(conway.Location{X: x, Y: y}, true)
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

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}
	return field, nil
}

// SaveFieldToFile exports the given field into a .rle file
func SaveFieldToFile(field *conway.Field, filename string) error {
	return nil
}

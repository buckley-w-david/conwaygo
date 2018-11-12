package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/buckley-w-david/conwaygo/pkg/conway"
	"github.com/buckley-w-david/conwaygo/pkg/encoding"

	tl "github.com/JoelOtter/termloop"
)

type conwayField struct {
	*conway.Field
}

var (
	live     map[int8]*tl.Cell
	cellChar *tl.Cell
	game     *tl.Game
	tick     float64
	start    bool
	level    *lifeLevel
	dirty    bool
	field    *conway.Field
	delay    float64
	debug    bool
	filename string
)

type lifeLevel struct {
	*tl.BaseLevel
	Bg       tl.Cell
	Entities []tl.Drawable
}

func newLifeLevel(bg tl.Cell) *lifeLevel {
	lv := tl.NewBaseLevel(bg)
	level := lifeLevel{Entities: make([]tl.Drawable, 0), Bg: bg, BaseLevel: lv}
	return &level
}

func (level *lifeLevel) SetBg(cell tl.Cell) {
	level.Bg = cell
}

// DrawBackground draws the background Cell bg to each Cell of the Screen s.
func (level *lifeLevel) DrawBackground(s *tl.Screen) {
	width, height := s.Size()

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			s.RenderCell(i, j, &level.Bg)
		}
	}
}

func init() {
	flag.BoolVar(&debug, "debug", false, "Display in debug mode")
	flag.Float64Var(&delay, "delay", 0.2, "Seconds between updates")
	flag.StringVar(&filename, "f", "", "Life RLE file")
	flag.Parse()

	live = map[int8]*tl.Cell{}
	var i int8
	for i = 0; i < 9; i++ {
		live[i] = &tl.Cell{Ch: []rune(strconv.Itoa(int(i)))[0], Fg: tl.ColorRed, Bg: tl.ColorBlack}
	}
	cellChar = &tl.Cell{Ch: '#', Fg: tl.ColorRed, Bg: tl.ColorBlack}
	tick = 0.0
	start = false
	dirty = true
}

func (f *conwayField) Tick(ev tl.Event) {
	// Enable arrow key movement
	switch ev.Type {
	case tl.EventKey:
		switch ev.Key {
		case tl.KeySpace:
			start = !start
		case tl.KeyCtrlS:
			encoding.SaveFieldToFile(field, "output.rle")
		case tl.KeyCtrlD:
			debug = !debug
			fg := tl.ColorDefault
			if debug {
				fg = tl.ColorBlue
			} else {
				fg = tl.ColorBlack
			}
			cell := tl.Cell{
				Bg: tl.ColorBlack,
				Fg: fg,
				Ch: '0',
			}
			level.SetBg(cell)
		case tl.KeyPgup:
			delay -= 0.01
		case tl.KeyPgdn:
			delay += 0.01
		case tl.KeyArrowLeft:
			x, y := level.Offset()
			level.SetOffset(x+1, y)
		case tl.KeyArrowRight:
			x, y := level.Offset()
			level.SetOffset(x-1, y)
		case tl.KeyArrowUp:
			x, y := level.Offset()
			level.SetOffset(x, y+1)
		case tl.KeyArrowDown:
			x, y := level.Offset()
			level.SetOffset(x, y-1)
		}
	case tl.EventMouse:
		x, y := level.Offset()
		if ev.Key == tl.MouseLeft {
			dirty = true
			l := conway.Location{X: ev.MouseX - x, Y: ev.MouseY - y}
			cell, exists := f.Cells[l]
			if exists {
				f.SetCell(l, !cell.State)
			} else {
				f.SetCell(l, true)
			}
		}
	}

	tick += game.Screen().TimeDelta()
	if tick > delay && start {
		tick = 0
		if dirty {
			f.Update()
			dirty = false
		}
		f.Commit()
		f.Count() //After a commit, the rc value does not reflect the state
		f.Update()
	}
}

func (f *conwayField) Draw(screen *tl.Screen) {
	var (
		cell *conway.Cell
	)
	for l := range f.Cells {
		cell, _ = f.Cells[l]
		if debug {
			screen.RenderCell(l.X, l.Y, live[cell.Rc])
		} else if cell.State {
			screen.RenderCell(l.X, l.Y, cellChar)
		}
	}
}

func main() {
	game = tl.NewGame()
	ch := '0'
	fg := tl.ColorBlack
	bg := tl.ColorBlack

	if debug {
		game.SetDebugOn(true)
		fg = tl.ColorBlue
	}

	var err error
	if filename == "" {
		field = conway.NewField([]conway.Location{})
	} else {
		field, err = encoding.LoadFieldFromFile(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		dirty = true
	}
	// spew.Dump(field)
	cField := &conwayField{field}
	level = newLifeLevel(tl.Cell{
		Bg: bg,
		Fg: fg,
		Ch: ch,
	})

	level.AddEntity(cField)

	game.Screen().SetLevel(level)
	game.Start()
}

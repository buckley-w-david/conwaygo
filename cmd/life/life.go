package main

import (
	"flag"
	"strconv"

	"github.com/buckley-w-david/conwaygo/pkg/conway"

	tl "github.com/JoelOtter/termloop"
)

type ConwayField struct {
	*conway.Field
}

var (
	live     map[int8]*tl.Cell
	cellChar *tl.Cell
	game     *tl.Game
	tick     float64
	start    bool
	level    *LifeLevel
	dirty    bool
	Delay    float64
	Debug    bool
)

type LifeLevel struct {
	*tl.BaseLevel
	Bg       tl.Cell
	Entities []tl.Drawable
}

func NewLifeLevel(bg tl.Cell) *LifeLevel {
	lv := tl.NewBaseLevel(bg)
	level := LifeLevel{Entities: make([]tl.Drawable, 0), Bg: bg, BaseLevel: lv}
	return &level
}

func (life *LifeLevel) SetBg(cell tl.Cell) {
	life.Bg = cell
}

// DrawBackground draws the background Cell bg to each Cell of the Screen s.
func (l *LifeLevel) DrawBackground(s *tl.Screen) {
	width, height := s.Size()

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			s.RenderCell(i, j, &l.Bg)
		}
	}
}

func init() {
	flag.BoolVar(&Debug, "debug", false, "Display in debug mode")
	flag.Float64Var(&Delay, "delay", 0.2, "Seconds between updates")
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

func (f *ConwayField) Tick(ev tl.Event) {
	// Enable arrow key movement
	switch ev.Type {
	case tl.EventKey:
		switch ev.Key {
		case tl.KeySpace:
			start = !start
		case tl.KeyCtrlD:
			Debug = !Debug
			fg := tl.ColorDefault
			if Debug {
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
			Delay -= 0.01
		case tl.KeyPgdn:
			Delay += 0.01
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
			l := conway.Location{ev.MouseX - x, ev.MouseY - y}
			cell, exists := f.Cells[l]
			if exists {
				f.SetCell(l, !cell.State)
			} else {
				f.SetCell(l, true)
			}
		}
	}

	tick += game.Screen().TimeDelta()
	if tick > Delay && start {
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

func (m *ConwayField) Draw(screen *tl.Screen) {
	var (
		cell *conway.Cell
	)
	for l := range m.Cells {
		cell, _ = m.Cells[l]
		if Debug {
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

	if Debug {
		game.SetDebugOn(true)
		fg = tl.ColorBlue
	}
	level = NewLifeLevel(tl.Cell{
		Bg: bg,
		Fg: fg,
		Ch: ch,
	})

	field := ConwayField{
		&conway.Field{
			map[conway.Location]*conway.Cell{},
		},
	}

	level.AddEntity(&field)

	game.Screen().SetLevel(level)
	game.Start()
}

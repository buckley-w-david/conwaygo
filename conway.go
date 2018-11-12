package main

import (
	"flag"
	"strconv"

	tl "github.com/JoelOtter/termloop"
	//"github.com/davecgh/go-spew/spew"
)

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

type Location struct {
	X int
	Y int
}

type Cell struct {
	state bool
	rc    int8
}

type Field struct {
	Cells map[Location]*Cell
}

func (l Location) Neighbours() [8]Location {
	loc := [8]Location{}
	adjust := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 1 {
				adjust = -1
			}
			loc[3*(i+1)+(j+1)+adjust] = Location{l.X + i, l.Y + j}
		}
	}
	return loc
}

func (f *Field) SetCell(l Location, state bool) {
	var (
		cell      *Cell
		neighbour *Cell
		exists    bool
	)

	neighbours := l.Neighbours()
	// game.Log("%v: %v", l, neighbours)
	if state {
		// If we're setting a cell to alive, track all adjacent cells
		for _, loc := range neighbours {
			neighbour, exists = f.Cells[loc]
			if !exists {
				f.SetCell(loc, false)
			}
		}
	}

	cell, exists = f.Cells[l]
	if exists {
		old := cell.state
		cell.state = state

		if !old && state {
			// Dead -> Living
			for _, loc := range neighbours {
				neighbour, exists = f.Cells[loc]
				if exists {
					neighbour.rc++
				} else {
					game.Log("Dead cell adjacent to living @ %v", loc)
				}
			}
		} else if old && !state {
			// Living -> Dead
			for _, loc := range neighbours {
				neighbour, exists = f.Cells[loc]
				if exists {
					neighbour.rc--
				} else {
					game.Log("Dead cell adjacent to previously living @ %v", loc)
				}
			}
		}
	} else {
		cell = &Cell{state: state, rc: 0}
		for _, loc := range neighbours {
			neighbour, exists = f.Cells[loc]
			if exists && neighbour.state {
				cell.rc++
			}
			if exists && state {
				neighbour.rc++
			}
		}
		f.Cells[l] = cell
	}
}

func (f *Field) commit() {
	// Update alive/dead status
	for _, cell := range f.Cells {
		switch cell.rc {
		case 2:
		case 3:
			cell.state = true
		default:
			cell.state = false
		}
	}

}

func (f *Field) clean() {
	var exists bool
	// Update relivant dead cells to track
	for l, cell := range f.Cells {
		for _, loc := range l.Neighbours() {
			_, exists = f.Cells[loc]
			// If we're not tracking a location and it has a living adjacent cell
			if !exists && cell.state {
				// start Tracking it
				f.SetCell(loc, false)
			}
		}

		// If we're tracking a dead cell with no living neighbours
		if !cell.state && cell.rc == 0 {
			// Stop
			delete(f.Cells, l)
		}
	}

}

func (f *Field) update() {
	var (
		neighbours int8
		exists     bool
		neighbour  *Cell
	)

	for l, cell := range f.Cells {
		neighbours = 0
		for _, loc := range l.Neighbours() {
			neighbour, exists = f.Cells[loc]
			if exists && neighbour.state {
				neighbours++
			}
		}
		cell.rc = neighbours
	}
}

func (f *Field) Tick(ev tl.Event) {
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
			l := Location{ev.MouseX - x, ev.MouseY - y}
			cell, exists := f.Cells[l]
			if exists {
				f.SetCell(l, !cell.state)
			} else {
				f.SetCell(l, true)
			}
		}
	}

	tick += game.Screen().TimeDelta()
	if tick > Delay && start {
		tick = 0
		if dirty {
			f.update()
			dirty = false
		}
		f.commit()
		f.update() //After a commit, the rc value does not reflect the state
		f.clean()
	}
}

func (m *Field) Draw(screen *tl.Screen) {
	var (
		cell *Cell
	)
	for l := range m.Cells {
		cell, _ = m.Cells[l]
		if Debug {
			screen.RenderCell(l.X, l.Y, live[cell.rc])
		} else if cell.state {
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

	field := Field{
		map[Location]*Cell{},
	}

	level.AddEntity(&field)

	game.Screen().SetLevel(level)
	game.Start()
}

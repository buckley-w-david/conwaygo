package conway

type Location struct {
	X int
	Y int
}

type Cell struct {
	State bool
	Rc    int8
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

func (f *Field) SetCell(l Location, State bool) {
	var (
		cell      *Cell
		neighbour *Cell
		exists    bool
	)

	neighbours := l.Neighbours()
	// game.Log("%v: %v", l, neighbours)
	if State {
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
		old := cell.State
		cell.State = State

		if !old && State {
			// Dead -> Living
			for _, loc := range neighbours {
				neighbour, exists = f.Cells[loc]
				if exists {
					neighbour.Rc++
				}
			}
		} else if old && !State {
			// Living -> Dead
			for _, loc := range neighbours {
				neighbour, exists = f.Cells[loc]
				if exists {
					neighbour.Rc--
				}
			}
		}
	} else {
		cell = &Cell{State: State, Rc: 0}
		for _, loc := range neighbours {
			neighbour, exists = f.Cells[loc]
			if exists && neighbour.State {
				cell.Rc++
			}
			if exists && State {
				neighbour.Rc++
			}
		}
		f.Cells[l] = cell
	}
}

func (f *Field) Commit() {
	// Update alive/dead status
	for _, cell := range f.Cells {
		switch cell.Rc {
		case 2:
		case 3:
			cell.State = true
		default:
			cell.State = false
		}
	}

}

func (f *Field) Update() {
	var exists bool
	// Update relivant dead cells to track
	for l, cell := range f.Cells {
		for _, loc := range l.Neighbours() {
			_, exists = f.Cells[loc]
			// If we're not tracking a location and it has a living adjacent cell
			if !exists && cell.State {
				// start Tracking it
				f.SetCell(loc, false)
			}
		}

		// If we're tracking a dead cell with no living neighbours
		if !cell.State && cell.Rc == 0 {
			// Stop
			delete(f.Cells, l)
		}
	}

}

func (f *Field) Count() {
	var (
		neighbours int8
		exists     bool
		neighbour  *Cell
	)

	for l, cell := range f.Cells {
		neighbours = 0
		for _, loc := range l.Neighbours() {
			neighbour, exists = f.Cells[loc]
			if exists && neighbour.State {
				neighbours++
			}
		}
		cell.Rc = neighbours
	}
}

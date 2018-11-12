package conway

// Location is an x, y coordinate
type Location struct {
	X int
	Y int
}

// Cell is a structure to containe the state of a location (alive or dead) and a count of it's living neighbours
type Cell struct {
	State bool
	Rc    int8
}

// Field is a structure containing all data needed to represent the Life grid
type Field struct {
	Cells map[Location]*Cell
}

// NewField creates a new Field struct given a list of locations to create live cells
func NewField(m []Location) *Field {
	f := &Field{make(map[Location]*Cell)}
	for _, l := range m {
		f.SetCell(l, true)
	}
	return f
}

// Neighbours are the surrounding locations of l
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

// SetCell updates the cell at location l to state, creating the cell if it does not already exist
// It then updates the Rc count of surrounding cells
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
		old := cell.State
		cell.State = state

		if !old && state {
			// Dead -> Living
			for _, loc := range neighbours {
				neighbour, exists = f.Cells[loc]
				if exists {
					neighbour.Rc++
				}
			}
		} else if old && !state {
			// Living -> Dead
			for _, loc := range neighbours {
				neighbour, exists = f.Cells[loc]
				if exists {
					neighbour.Rc--
				}
			}
		}
	} else {
		cell = &Cell{State: state, Rc: 0}
		for _, loc := range neighbours {
			neighbour, exists = f.Cells[loc]
			if exists && neighbour.State {
				cell.Rc++
			}
			if exists && state {
				neighbour.Rc++
			}
		}
		f.Cells[l] = cell
	}
}

// Commit changes the state of all cells to reflect their Rc counts
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

// Update handles adding untracked dead cells that have the potential to come alive in the next iteration
// And removes dead cells that do not have the potential to come alive
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

// Count updates the Rc count of each cell to reflect the number of living neighbours it has
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

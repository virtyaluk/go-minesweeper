package main

// Cell contains cell's state.
type Cell struct {
	State      byte
	DangerZone int
}

// Cell constants
const (
	Bomb     = byte(9)
	Unfolded = byte(1 << 4)
)

// IsBomb checks if a cell contain a bomb
func (c Cell) IsBomb() bool {
	return c.State&Bomb == Bomb
}

// IsUnfolded checks if a cell is unfolded
func (c Cell) IsUnfolded() bool {
	return c.State&Unfolded == Unfolded
}

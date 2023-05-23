package main

import "math"

/*
 * Part 1: the game state is represented as 2d array of Cell objects.
 * Each Cell objects represents either bomb, unfolded or untouched (default) game cell.
 * Each cell has a counter of neighbouring bomb cells.
 */

type BoardState int

const (
	BoardLost      BoardState = 0
	BoardWon       BoardState = 1
	BoardUndefined BoardState = 2
)

// Board contains board's state.
type Board struct {
	cells      [][]Cell
	state      BoardState
	rows, cols int
	diffc      float64
}

var visRows = []int{-1, -1, -1, 0, 1, 1, 1, 0}
var visCols = []int{-1, 0, 1, 1, 1, 0, -1, -1}

// NewBoard returns a new instance of Board object.
func NewBoard(rows, cols int, diffc float64) *Board {
	cells := make([][]Cell, rows)
	cellsNum := rows * cols
	bombs := int(math.Ceil(float64(cellsNum) * diffc))

	/*
	 * Part 2: A modified Fisher-Yates algorithm is used to generate a set of K elements
	 * distributed randomly and uniformly on the range of n elements.
	 */
	bombsPositions := RandomSample(rows*cols, bombs)

	for i := 0; i < cols; i++ {
		cells[i] = make([]Cell, cols)
	}

	for _, bomb := range bombsPositions {
		cells[bomb/rows][bomb%rows].State = Bomb
	}

	/*
	 * Part 3: Each cell hold a counter representing the number of 8-directional adjacent bomb cells.
	 */

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			for k := 0; k < 8; k++ {
				if isBombAtCell(cells, rows, cols, i+visRows[k], j+visCols[k]) {
					cells[i][j].DangerZone++
				}
			}
		}
	}

	return &Board{
		cells: cells,
		state: BoardUndefined,
		rows:  rows,
		cols:  cols,
		diffc: diffc,
	}
}

// BombsCount returns the number of bombs on the board.
func (b *Board) BombsCount() int {
	var bombsLeft int

	for _, r := range b.cells {
		for _, c := range r {
			if c.IsBomb() {
				bombsLeft++
			}
		}
	}

	return bombsLeft
}

// FreeCellsLeft returns the number of cells free to open.
func (b *Board) FreeCellsLeft() int {
	var freeCellsLeft int

	for _, r := range b.cells {
		for _, c := range r {
			if !c.IsBomb() && !c.IsUnfolded() {
				freeCellsLeft++
			}
		}
	}

	return freeCellsLeft
}

// GetCell returns a cell given coordinates
func (b *Board) GetCell(row, col int) *Cell {
	return &b.cells[row][col]
}

// UnfoldCell unfolds a cell and returns a tuple where first elements indicates whether a bomb was unfolded
// and the number of adjacent cells with a bomb and the list of discovered empty cells, if any.
func (b *Board) UnfoldCell(row, col int) (bool, int, [][]int) {
	if b.GetCell(row, col).IsUnfolded() {
		return false, b.GetCell(row, col).DangerZone, nil
	}

	if b.GetCell(row, col).IsBomb() {
		b.GetCell(row, col).State = Bomb
		b.state = BoardLost
		return true, 0, nil
	} else if b.GetCell(row, col).DangerZone == 0 {
		if b.FreeCellsLeft() == 0 {
			b.state = BoardWon
		}

		return false, 0, b.revealEmpty(row, col)
	} else {
		if b.FreeCellsLeft() == 0 {
			b.state = BoardWon
		}

		b.GetCell(row, col).State = Unfolded
		return false, b.GetCell(row, col).DangerZone, nil
	}
}

// IsLost returns true if a game in BoardLost state, false otherwise.
func (b *Board) IsLost() bool {
	return b.state == BoardLost
}

// IsWon returns true if a game in BoardWon state, false otherwise.
func (b *Board) IsWon() bool {
	return b.state == BoardWon
}

func (b *Board) revealEmpty(row, col int) [][]int {
	if !b.isValidCoordinates(row, col) || b.GetCell(row, col).IsBomb() || b.GetCell(row, col).IsUnfolded() || b.GetCell(row, col).DangerZone > 0 {
		return [][]int{}
	}

	res := make([][]int, 0)

	b.GetCell(row, col).State = Unfolded
	res = append(res, []int{row, col})

	for k := 0; k < 8; k++ {
		nextRow, nextCol := row+visRows[k], col+visCols[k]
		res = append(res, b.revealEmpty(nextRow, nextCol)...)
	}

	return res
}

func (b *Board) isValidCoordinates(row, col int) bool {
	if row < 0 || row >= b.rows || col < 0 || col >= b.cols {
		return false
	}

	return true
}

func isBombAtCell(cells [][]Cell, rows, cols, row, col int) bool {
	if row < 0 || row >= rows || col < 0 || col >= cols {
		return false
	}

	return cells[row][col].IsBomb()
}

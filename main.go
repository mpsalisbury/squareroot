package main

// Solve the Square Root sliding block puzzle.
// Starting position:
//  ____
// |abbc|
// |abbc|
// |deef|
// |dghf|
// |i  j|
//  ~~~~
//
// Each letter represents a piece that occupies the given spaces on the board.
// Pieces can slide left/right/up/down within the bounds of the frame. They cannot
// rotate. A piece can slide if its target space is open. The goal is to move piece b
// to the bottom middle spot on the board where it can slide out of the puzzle.
// Find a sequence of moves gets piece b to the bottom middle?

import (
	"fmt"
	"sort"
	"strings"
)

// Solution strategy:
//
// Maintain a queue of board configurations ordered by the number of moves taken to reach them.
// Also maintain a set of board configurations we've already seen.
// For the first board on the queue:
//   Collect all legal moves
//   For each remaining move:
//     Apply the move to the current board -> nextBoard (move piece, record new move)
//     If we've seen nextBoard before, skip it
//     If mark nextBoard as seen
//     If nextBoard is a winning configuration, print it, and we're done.
//     Add nextBoard to the queue of boards to consider
func main() {
	bs := []*Board{makeStartingBoard()}
	seenBoards := make(map[string]bool)
	numSkipped := 0
	for {
		if len(bs) == 0 {
			fmt.Print("Couldn't find solution\n")
			return
		}
		b := bs[0]
		bs = bs[1:]
		for _, m := range b.possibleMoves() {
			nb := b.move(m)
			nbConfig := nb.Config()
			if seenBoards[nbConfig] {
				numSkipped++
				continue
			}
			seenBoards[nbConfig] = true
			if b.isWin() {
				fmt.Printf("Found solution (%d moves, %d configurations, %d skipped):\n",
					len(b.mvs), len(seenBoards), numSkipped)
				printMoves(b.mvs)
				return
			}
			bs = append(bs, nb)
		}
	}
}

func printMoves(mvs []Move) {
	b := makeStartingBoard()
	fmt.Print(b.String())
	for i, m := range mvs {
		fmt.Printf("%d: %s\n", i+1, m.String())
		b = b.move(m)
		fmt.Print(b.String())
	}
}

// Returns the starting board configuration.
func makeStartingBoard() *Board {
	//    0123
	//    ____
	// 0 |abbc|
	// 1 |abbc|
	// 2 |deef|
	// 3 |dghf|
	// 4 |i  j|
	//    ~~~~
	ps := []Piece{
		Piece{"a", 1, 2, 0, 0},
		Piece{"b", 2, 2, 1, 0},
		Piece{"c", 1, 2, 3, 0},
		Piece{"d", 1, 2, 0, 2},
		Piece{"e", 2, 1, 1, 2},
		Piece{"f", 1, 2, 3, 2},
		Piece{"g", 1, 1, 1, 3},
		Piece{"h", 1, 1, 2, 3},
		Piece{"i", 1, 1, 0, 4},
		Piece{"j", 1, 1, 3, 4},
	}
	pm := make(map[string]Piece)
	for _, p := range ps {
		pm[p.id] = p
	}

	return &Board{4, 5, pm, []Move{}}
}

// Records the configuration of a board and how it got there (set of moves).
type Board struct {
	// The size of the board.
	w, h int

	// The pieces on the board.
	ps map[string]Piece

	// The moves used to get the pieces where they are.
	mvs []Move
}

// Is the given space unoccupied by a piece on this board.
func (b *Board) isOpen(s Space) bool {
	if s.x < 0 || s.y < 0 || s.x >= b.w || s.y >= b.h {
		return false
	}
	for _, p := range b.ps {
		if p.covers(s) {
			return false
		}
	}
	return true
}

// Returns the set of legal moves of pieces given this board configuration.
func (b *Board) possibleMoves() []Move {
	mvs := []Move{}
	for _, p := range b.ps {
		pmvs := p.possibleMoves(b)
		mvs = append(mvs, pmvs...)
	}
	return mvs
}

// Returns a new board the same as this one but with the given move applied.
func (b *Board) move(m Move) *Board {
	// The new pieces are the old pieces with one piece moved.
	nps := make(map[string]Piece)
	for pid, p := range b.ps {
		if pid == m.pid {
			nps[pid] = p.move(m.dir)
		} else {
			nps[pid] = p
		}
	}
	// The new moves are the old moves plus the new move.
	nmvs := []Move{}
	for _, m := range b.mvs {
		nmvs = append(nmvs, m)
	}
	nmvs = append(nmvs, m)

	return &Board{b.w, b.h, nps, nmvs}
}

// Is the current board position a winning configuration.
func (b *Board) isWin() bool {
	pb := b.ps["b"]
	return pb.x == 1 && pb.y == 3
}

// Config returns the configuration of the pieces on the given board.
// We use this to record which configurations we've already considered
// so that we don't consider them again.
func (b *Board) Config() string {
	pcs := []string{}
	for _, p := range b.ps {
		pcs = append(pcs, p.Config())
	}
	sort.Strings(pcs)
	return strings.Join(pcs, ";")
}

// Returns a spatial representation of the board. e.g.:
//  ____
// |abbc|
// |abbc|
// |deef|
// |dghf|
// |i  j|
//  ~~~~
func (b *Board) String() string {
	grid := makeGrid(b.w, b.h)
	for _, p := range b.ps {
		p.drawInto(grid)
	}

	var sb strings.Builder
	sb.WriteString(" ")
	sb.WriteString(strings.Repeat("_", b.w))
	sb.WriteString("\n")

	for i := 0; i < b.h; i++ {
		sb.WriteString("|")
		sb.WriteString(grid.row(i))
		sb.WriteString("|\n")
	}

	sb.WriteString(" ")
	sb.WriteString(strings.Repeat("~", b.w))
	sb.WriteString("\n")

	return sb.String()
}

// Piece records the id and configuration of a piece.
type Piece struct {
	id   string
	w, h int // size in squares
	x, y int // position of upper-left square
}

// A piece's configuration records its size and location.
// This is used to record which configurations of all pieces we've seen before
// so we don't consider them again. It ignores the id because we don't care
// which because any piece of the same shape is equivalent for the solution.
func (p Piece) Config() string {
	return fmt.Sprintf("%dx%d-%d,%d", p.w, p.h, p.x, p.y)
}

func (p Piece) drawInto(grid *Grid) {
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			grid.set(p.x+x, p.y+y, p.id[0])
		}
	}
}

// Is this piece free to move in the given direction on this board.
func (p Piece) canMove(b *Board, d Direction) bool {
	for _, ts := range p.targetSpaces(d) {
		if !b.isOpen(ts) {
			return false
		}
	}
	return true
}

// What are all of the possible legal moves this piece can move on this board.
func (p Piece) possibleMoves(b *Board) []Move {
	mvs := []Move{}
	for _, d := range Directions {
		if p.canMove(b, d) {
			mvs = append(mvs, Move{p.id, d})
		}
	}
	return mvs
}

// Returns this piece moved in the given direction.
func (p Piece) move(d Direction) Piece {
	switch d {
	case Up:
		return Piece{p.id, p.w, p.h, p.x, p.y - 1}
	case Down:
		return Piece{p.id, p.w, p.h, p.x, p.y + 1}
	case Left:
		return Piece{p.id, p.w, p.h, p.x - 1, p.y}
	case Right:
		return Piece{p.id, p.w, p.h, p.x + 1, p.y}
	}
	panic("Invalid directon")
}

// Space represents a 1x1 space on the board.
type Space struct {
	x, y int
}

// Which spaces will be moved into if this piece moves in the given direction.
func (p Piece) targetSpaces(d Direction) []Space {
	switch d {
	case Up:
		return hSpaces(p.y-1, p.x, p.x+p.w-1)
	case Down:
		return hSpaces(p.y+p.h, p.x, p.x+p.w-1)
	case Left:
		return vSpaces(p.x-1, p.y, p.y+p.h-1)
	case Right:
		return vSpaces(p.x+p.w, p.y, p.y+p.h-1)
	}
	panic("Invalid directon")
}

// hSpaces returns a horizontal set of spaces.
func hSpaces(y, x1, x2 int) []Space {
	ss := []Space{}
	for x := x1; x <= x2; x++ {
		ss = append(ss, Space{x, y})
	}
	return ss
}

// vSpaces returns a vertical set of spaces.
func vSpaces(x, y1, y2 int) []Space {
	ss := []Space{}
	for y := y1; y <= y2; y++ {
		ss = append(ss, Space{x, y})
	}
	return ss
}

// Does this piece cover the given space?
func (p Piece) covers(s Space) bool {
	return s.x >= p.x && s.y >= p.y && s.x < p.x+p.w && s.y < p.y+p.h
}

// Records a move of a piece in a direction for a single unit distance.
type Move struct {
	pid string
	dir Direction
}

func (m Move) String() string {
	return fmt.Sprintf("%s -> %s", m.pid, m.dir)
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

var Directions = []Direction{Up, Down, Left, Right}

func (d Direction) String() string {
	return []string{"Up", "Down", "Left", "Right"}[d]
}

// Grid holds a visual representation of a Board.
type Grid struct {
	w, h int
	c    [][]byte
}

func makeGrid(w, h int) *Grid {
	c := [][]byte{}
	for y := 0; y < h; y++ {
		row := make([]byte, w)
		for x := 0; x < w; x++ {
			row[x] = ' '
		}
		c = append(c, row)
	}
	return &Grid{w, h, c}
}

func (g *Grid) set(x, y int, c byte) {
	if x < 0 || y < 0 || x >= g.w || y >= g.h {
		return
	}
	g.c[y][x] = c
}

func (g *Grid) row(y int) string {
	return string(g.c[y])
}

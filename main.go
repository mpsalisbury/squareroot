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
				continue
			}
			//fmt.Printf("Checking board %s - %d\n", nbConfig, len(nb.mvs))
			seenBoards[nbConfig] = true
			if b.isWin() {
				fmt.Printf("Found solution (%d moves):\n", len(b.mvs))
				for _, m := range b.mvs {
					fmt.Printf("  %s\n", m.String())
				}
				return
			}
			bs = append(bs, nb)
		}
	}
}

func makePiece(id string, x, y, w, h int) Piece {
	return Piece{id, w, h, x, y}
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
		makePiece("a", 0, 0, 1, 2),
		makePiece("b", 1, 0, 2, 2),
		makePiece("c", 3, 0, 1, 2),
		makePiece("d", 0, 2, 1, 2),
		makePiece("e", 1, 2, 2, 1),
		makePiece("f", 3, 2, 1, 2),
		makePiece("g", 1, 3, 1, 1),
		makePiece("h", 2, 3, 1, 1),
		makePiece("i", 0, 4, 1, 1),
		makePiece("j", 3, 4, 1, 1),
	}
	pm := make(map[string]Piece)
	for _, p := range ps {
		pm[p.id] = p
	}

	return &Board{4, 5, pm, []Move{}}
}

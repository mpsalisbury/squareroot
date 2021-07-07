# Square Root
***Computes the solution for the Square Root sliding block puzzle***

(See [physical puzzle](http://squarerootgames.com/puzzles.html))

Starting position:
     ____
    |abbc|
    |abbc|
    |deef|
    |dghf|
    |i  j|
     ~~~~

Each letter represents a piece that occupies the given spaces on the board.
Pieces can slide left/right/up/down within the bounds of the frame. They cannot
rotate. A piece can slide if its target space is open. The goal is to move piece b
to the bottom middle spot on the board where it can slide out of the puzzle.

This program computes and prints the shortest solution using a breadth-first search.

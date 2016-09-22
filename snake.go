// Design:
// inLoop:
// Board: Game board
//
package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Snake
type Snake struct {
	x   int
	y   int
	len int
	dir int // 0, 1, 2, 3 => up, down, left, right
}

// Board
type Board struct {
	cells      [][]int
	size       int
	updateChan chan int
	exitChan   chan bool
	snake      Snake
}

func newBoard(size int) *Board {
	b := Board{}
	b.snake = Snake{x: 0, y: 0, len: 5, dir: 3}
	b.size = size
	b.cells = make([][]int, 0, size)
	b.updateChan = make(chan int)
	b.exitChan = make(chan bool)
	for i := 0; i < size; i++ {
		tmp := make([]int, size)
		b.cells = append(b.cells, tmp)
	}
	go b.displayLoop()
	return &b
}

func (b *Board) displayLoop() {
	for {
		select {
		case d := <-b.updateChan:
			b.updateBoard(d)
			clear()
			b.display()
		case <-time.After(time.Millisecond * 200):
			b.updateBoard(b.snake.dir)
			clear()
			b.display()
		}
	}
}

// display
func (b Board) display() {
	rowStr := make([]byte, len(b.cells))
	for row := range b.cells {
		for col := range b.cells[row] {
			if b.cells[row][col] > 0 {
				rowStr[col] = '@'
			} else {
				rowStr[col] = ' '
			}
		}
		fmt.Printf("|%s|\n", string(rowStr))
	}
}

func (b *Board) updateBoard(direct int) {
	b.snake.dir = direct
	switch direct {
	case 0: // w
		b.snake.x = (b.snake.x - 1 + b.size) % b.size
	case 1: // s
		b.snake.x = (b.snake.x + 1) % b.size
	case 2: // a
		b.snake.y = (b.snake.y - 1 + b.size) % b.size
	case 3: // d
		b.snake.y = (b.snake.y + 1) % b.size
	}
	if b.cells[b.snake.x][b.snake.y] == 0 {
		b.cells[b.snake.x][b.snake.y] = b.snake.len
	} else {
		b.exitChan <- true
	}
	for row := range b.cells {
		for col := range b.cells[row] {
			if b.cells[row][col] > 0 {
				b.cells[row][col]--
			}
		}
	}
}

func (b Board) update(direct int) {
	b.updateChan <- direct
}

// getDirection
func getDirection(directChan chan int, exitChan chan bool) {
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		switch c := b[0]; c {
		case 'w':
			directChan <- 0
		case 'a':
			directChan <- 2
		case 's':
			directChan <- 1
		case 'd':
			directChan <- 3
		case 'q':
			exitChan <- true
		}
	}
}

func clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

// inLoop
func inLoop(isComplete chan bool) {
	board := newBoard(30)
	directChan := make(chan int)
	go getDirection(directChan, board.exitChan)
	for {
		select {
		case direct := <-directChan:
			fmt.Println(direct)
			board.update(direct)
		case <-board.exitChan:
			goto exit
		}
	}
exit:
	fmt.Println("Leaving Game")
	isComplete <- true
}

// main
func main() {
	// set tty
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// init board

	isComplete := make(chan bool)
	go inLoop(isComplete)
	<-isComplete
}

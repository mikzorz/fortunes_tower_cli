package main

import (
	"bufio"
	"fmt"
	_ "io"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
  StateBetting = iota
  StatePlaying
)

// Game contains the deck and the tower
type Game struct {
	deck    []int
	counts  map[int]int
	tower   [][]int
	curRow  int
	balance int
	in      *bufio.Scanner
  state   int
  wager   int
}

// New() creates a new game
func (g *Game) New() {
	rand.Seed(time.Now().UnixNano())
  g.NewDeckAndTower()
  g.curRow=0
  g.state=StateBetting
  g.balance=300
  g.wager=15
}

// Set the deck, counts and tower to defaults
func (g *Game) NewDeckAndTower() {
	d := []int{}

	c := make(map[int]int)
	c[0] = 4
	for i := 1; i <= 7; i++ {
		c[i] = 8
		for j := 0; j < 8; j++ {
			d = append(d, i)
		}
	}
	g.counts = c

	for i := 0; i < 4; i++ {
		d = append(d, 0)
	}
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
	g.deck = d

	g.tower = make([][]int, 8)
}

// Deal() deals the next row of cards
func (g *Game) Deal() {
  if g.curRow == 0 {
    g.balance -= g.wager
  }
  if g.curRow < 8 {
    g.state = StatePlaying
    for i := 0; i < g.curRow+1; i++ {
      drawnCard := g.deck[0]
      g.deck = g.deck[1:] // ok because deck never empties completely
      g.counts[drawnCard]--
      g.tower[g.curRow] = append(g.tower[g.curRow], drawnCard)
    }
    g.curRow++
  } else {
    g.CashOut()
  }
}

// CashOut() adds the sum of the last row to the player's balance.
func (g *Game) CashOut() {
  if g.curRow > 0 {
    g.state = StateBetting
    sum := 0
    for _, v := range g.tower[g.curRow-1] {
      sum += v
    }
    g.balance += sum
    g.NewDeckAndTower()
    g.curRow=0
  } else {
  }
}

// Input() reads and processes user input.
// z deals a new row/confirms
// x cashes out at the current row if the player has not bust
// func (g *Game) Input(in io.Reader) {
func (g *Game) Input(in string) {
	in = strings.TrimSuffix(in, "\n")
	switch in {
	case "z":
		if g.curRow == 0 {
			g.Deal()
		}
		g.Deal()
	case "x":
		g.CashOut()
	}
}

// Balance() returns player's current cash amount.
func (g *Game) Balance() int {
	return g.balance
}

// State() returns the current game state.
func (g *Game) State() int {
  return g.state
}

// GetWager() returns the current wager.
func (g *Game) GetWager() int {
  return g.wager
}

// Print the current game state, with instructions
func (g *Game) Print() {
	for row := 0; row < g.curRow; row++ {
		fmt.Println(strings.Repeat(" ", 8-row), g.tower[row])
	}

  fmt.Println()

  switch g.State() {
  case StateBetting:
	  fmt.Println(`Type "z" to bet 15`)
  case StatePlaying:
	  fmt.Println(`"z" to deal the next row, "x" to cash out`)
  }

  fmt.Printf("Money: %d\n", g.Balance())
}

func main() {
	g := Game{}
	g.New()
	reader := bufio.NewReader(os.Stdin)
	for {
		g.Print()
		in, _ := reader.ReadString('\n')
		g.Input(in)
		time.Sleep(time.Second / 5)
		// if g.GameOver() {
		// 	break
		// }
	}
}

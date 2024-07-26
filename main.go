package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	StateBetting = iota
	StatePlaying
	StateGameOver

  maxRows = 8
)

// Game contains the deck and the tower
type Game struct {
	deck       []int
	counts     map[int]int
	tower      [][]int
	curRow     int
	balance    int
	in         *bufio.Scanner
	out        io.Writer
	state      int
	wager      int
	multiplier int
	gameover   bool
}

// NewGame() creates a new game with a fresh deck, tower and money
func NewGame() Game {
	rand.Seed(time.Now().UnixNano())

	g := Game{}
	g.NewRound()
	g.balance = 300
	g.wager = 15
	g.out = os.Stdout
	return g
}

func (g *Game) NewRound() {
	g.state = StateBetting
	g.multiplier = 1
	g.NewDeckAndTower()
  g.curRow = 0
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

	g.tower = make([][]int, maxRows)
	g.curRow = 0
}

// deal() deals the next row of cards
func (g *Game) deal() {
	if g.curRow == 0 {
		g.balance -= g.wager
		g.multiplier *= g.wager / 15
	}
	if !g.IsGameOver() {
		g.state = StatePlaying
		for i := 0; i <= g.curRow; i++ {
			drawnCard := g.deck[0]
			g.deck = g.deck[1:] // ok because deck never empties completely
			g.counts[drawnCard]--
			g.tower[g.curRow] = append(g.tower[g.curRow], drawnCard)
		}

		if g.curRow > 1 {
      if bust := g.handleBust(); bust {
        return
      }
		}
		g.checkMulti()

    if g.curRow < maxRows-1 {
		  g.curRow++
    } else {
      g.gameOver()
    }
	} else {
		g.cashOut()
	}
}

// dealX() repeates deal() x times
func (g *Game) dealX(x int) {
	for i := 0; i < x; i++ {
		g.deal()
	}
}

// handleBust() checks for a bust. If there is, try to replace the first occurence with the gate card.
// If gate doesn't exist, gameover, return true.
// Check for a bust again. If there is, gameover, return true.
// Else, return false
func (g *Game) handleBust() bool {
	for i := 0; i < 2; i++ {
		if bust, ci := g.IsBust(); bust {
			if len(g.tower[0]) > 0 {
				g.tower[g.curRow][ci] = g.tower[0][0]
				g.tower[0] = []int{}
			} else {
				g.gameOver()
        return true
			}
		}
	}
  return false
}

// IsBust() compares each card on the last dealt row with each card directly above it.
// If they match, return true and the index of the bust card.
// Else, return false, 0
func (g *Game) IsBust() (bool, int) {
  curRow := g.curRow
	for cardIndex, cardVal1 := range g.tower[curRow] {
		if cardVal1 == 0 {
			return false, 0
		}

		if cardIndex != len(g.tower[curRow])-1 {
			// compare currow[cardIndex] with lastrow[cardIndex]
			cardVal2 := g.tower[curRow-1][cardIndex]
			if cardVal1 == cardVal2 {
				return true, cardIndex
			}
		}

		if cardIndex != 0 {
			// compare currow[cardIndex] with lastrow[i - 1]
			cardVal2 := g.tower[curRow-1][cardIndex-1]
			if cardVal1 == cardVal2 {
				return true, cardIndex
			}
		}
	}

	return false, 0
}

func (g *Game) checkMulti() {
	cardsToCheck := g.tower[g.curRow]
	for i := 0; i < len(cardsToCheck)-1; i++ {
		if cardsToCheck[i] != cardsToCheck[i+1] {
			return
		}
	}
	g.multiplier *= len(cardsToCheck)
}

func (g *Game) getRowValue(row int) int {
	rv := 0
	for _, v := range g.tower[row] {
		rv += v
	}
	return rv
}

func (g *Game) getJackpotValue() int {
	sum := 0
	for r := maxRows-1; r > 0; r-- {
		sum += g.getRowValue(r)
	}
	return sum
}

// cashOut() adds the sum of the last row to the player's balance.
func (g *Game) cashOut() {
	if g.curRow > 0 {
		sum := 0
		if g.curRow == 7 && len(g.tower[7]) == 8 && len(g.tower[0]) == 1 {
			sum = g.getJackpotValue()
		} else {
			sum = g.getRowValue(g.curRow - 1)
		}
		g.balance += sum * g.multiplier
    g.NewRound()
	}
}

// gameOver() sets game state to StateGameOver.
func (g *Game) gameOver() {
	g.state = StateGameOver
}

func (g *Game) IsGameOver() bool {
	return g.State() == StateGameOver
}

// Input() reads and processes user input.
// z deals a new row/confirms
// x cashes out at the current row if the player has not bust
func (g *Game) Input(in string) {
	in = strings.TrimSuffix(in, "\n")
	switch in {
	case "z":
		if g.IsGameOver() {
			g.NewRound()
		} else {
			if g.curRow == 0 {
				g.deal()
			}
			g.deal()
		}
	case "x":
		if g.IsGameOver() {
			g.NewRound()
		} else {
			g.cashOut()
		}
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

func (g *Game) SetWager(w int) {
	g.wager = w
}

func (g *Game) PrintRow(row int) {
	spacing := strings.Repeat(" ", 8-row)
	if row == 0 {
		if len(g.tower[0]) == 0 {
			fmt.Fprint(g.out, spacing, "[ ]")
		} else {
			fmt.Fprint(g.out, spacing, "[?]")
		}
	} else {
		rv := 0
		if g.curRow == 7 && len(g.tower[7]) == 8 && len(g.tower[0]) == 1 {
			rv = g.getJackpotValue()
		} else {
			rv = g.getRowValue(row)
		}
		fmt.Fprint(g.out, spacing, g.tower[row], spacing, fmt.Sprintf("(%d)", rv))
	}
	fmt.Fprint(g.out, "\n")
}

func (g *Game) PrintTower() {
	for row := 0; row < g.curRow; row++ {
		g.PrintRow(row)
	}
  if g.IsGameOver() {
    g.PrintRow(g.curRow)
  }

	fmt.Println()
}

// Print the current game state, with instructions
func (g *Game) PrintText() {
	// fmt.Printf("\033[2K\r") -- Use this to replace rows of text (untested)
  // fmt.Print("\033[u\033[K") // restore the cursor position and clear the line
	switch g.State() {
	case StateBetting:
		fmt.Println(`Type "z" to bet 15`)
	case StatePlaying:
		fmt.Println(`"z" to deal the next row, "x" to cash out`)
	case StateGameOver:
		fmt.Println(`BUST! "z" or "x" to start a new round`)
	}

	fmt.Printf("Money: %d\n", g.Balance())
}

func main() {
	g := NewGame()
	reader := bufio.NewReader(os.Stdin)
	for {
    // fmt.Print("\033[s") // save the cursor position
		g.PrintText()
		in, _ := reader.ReadString('\n')
		g.Input(in)
		g.PrintTower()
		time.Sleep(time.Second / 5)
		// if g.GameOver() {
		// 	break
		// }
	}
}

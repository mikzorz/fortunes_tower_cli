package main

import (
  "bufio"
  "fmt"
  _ "io"
  "os"
  "strings"
  "time"
  "math/rand"
)

// Game contains the deck and the tower
type Game struct {
  deck []int
  counts map[int]int 
  tower [][]int
  curRow int
  balance int
  in *bufio.Scanner
}

// New() resets the deck and tower
func (g *Game) New() {
  rand.Seed(time.Now().UnixNano())
  c := make(map[int]int)
  d := []int{}
  c[0] = 4
  for i := 0; i < 4; i++ {
    d = append(d, 0)
  }
  for i := 1; i <= 7; i++ {
    c[i] = 8
    for j := 0; j < 8; j++ {
      d = append(d, i)
    }
  }
  g.counts = c
  rand.Shuffle(len(d), func(i, j int) {
    d[i], d[j] = d[j], d[i]
  })
  g.deck = d

  g.tower = make([][]int, 8)

}

// Deal() deals the next row of cards
func (g *Game) Deal() {
  for i := 0; i < g.curRow+1; i++ {
    drawnCard := g.deck[0]
    g.deck = g.deck[1:] // ok because deck never empties completely
    g.counts[drawnCard]--
    g.tower[g.curRow] = append(g.tower[g.curRow], drawnCard)
  }
  g.curRow++
}

// Input() reads and processes user input.
// z deals a new row/confirms
// x cashes out at the current row if the player has not bust
// func (g *Game) Input(in io.Reader) {
func (g *Game) Input(in string) {
  // g.in = bufio.NewScanner(in)
  // g.in.Scan()
  // switch g.in.Text() {
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

// CashOut() adds the sum of the last row to the player's balance.
func (g *Game) CashOut() {
  sum := 0
  for _, v := range g.tower[g.curRow-1] {
    sum += v
  }
  g.balance += sum
}

func (g *Game) Print() {
  for row := 0; row < g.curRow; row++ {
    fmt.Println(strings.Repeat(" ", 8-row), g.tower[row])
  }
}

func main() {
  g := Game{}
  g.New()
  reader := bufio.NewReader(os.Stdin)
  for i := 0; i < 7; i++ {
    // in := strings.NewReader("z\n")
    in, _ := reader.ReadString('\n')
    g.Input(in)
    time.Sleep(time.Second/5)
    g.Print()
  }
}


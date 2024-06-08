package main

import (
  "fmt"
  "time"
  "math/rand"
)

// Game contains the deck and the tower
type Game struct {
  deck map[int]int // deck as a map is good for tracking card amounts. A shuffled slice is nice for Deal() (no while true). Use both?
  tower [][]int
  curRow int
}

// New() resets the deck and tower
func (g *Game) New() {
  rand.Seed(time.Now().UnixNano())
  d := make(map[int]int)
  d[0] = 4
  for i := 1; i <= 7; i++ {
    d[i] = 8
  }
  g.deck = d

  g.tower = make([][]int, 8)
}

// Deal() deals the next row of cards
func (g *Game) Deal() {
  for i := 0; i < g.curRow+1; i++ {
    cardToDraw := rand.Intn(8)
    for g.deck[cardToDraw] <= 0 {
      cardToDraw = rand.Intn(8)
    }
    g.tower[g.curRow] = append(g.tower[g.curRow], cardToDraw)
    g.deck[cardToDraw]--
  }
  g.curRow++
}

func main() {
  g := Game{}
  g.New()
  for i := 0; i < 8; i++ {
    g.Deal()
  }
  for _, row := range g.tower {
    fmt.Println(row)
  }
}


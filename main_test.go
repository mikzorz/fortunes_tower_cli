package main

import (
  "fmt"
  _ "strings"
  "testing"
)

func TestNewGame (t *testing.T) {
  // When a new game is created, check contents of deck and tower
  g := Game{}
  g.New()

  // deck is a []int
  // counts is a map[int]int
  // tower is [][]int

  t.Run("g.counts should not be nil", func(t *testing.T) {
    if g.counts == nil {
      t.Fatalf("Game.New() did not create a map of card counts")
    }
  })

  t.Run("deck should contain 4 hero cards", func(t *testing.T) {
    deckCount := 0
    for _, v := range g.deck {
      if v == 0 {
        deckCount++
      }
    }
    c, ok := g.counts[0]
    if !ok {
      t.Fatalf("card with value %d should exist in g.counts but doesn't", 0)
    }

    if c != 4 || deckCount != 4 {
      t.Errorf("wrong amount of hero cards, want 4, got %d", c)
    }
  })

  t.Run("deck should contain 8 copies of 1-7", func(t *testing.T) {
    for i := 1; i <= 7; i++ {
      deckCount := 0
      for _, v := range g.deck {
        if v == i {
          deckCount++
        }
      }
      c, ok := g.counts[i]
      if !ok {
        t.Fatalf("card with value %d should exist in g.counts but doesn't", i)
      }

      if c != 8 {
        t.Errorf("wrong amount of %d-value cards in g.counts, want 8, got %d", i, c)
      }
      if deckCount != 8 {
        t.Errorf("wrong amount of %d-value cards in g.deck, want 8, got %d", i, c)
      }
    }
  })

  t.Run("tower should not be nil", func(t *testing.T) {
    if g.tower == nil {
      t.Fatalf("Game.New() did not create a tower")
    }
  })

  t.Run("tower should be empty", func(t *testing.T) {
    for _, row := range g.tower {
      if len(row) != 0 {
        t.Fatalf("row %d of tower should be empty, contains %d cards", row, len(row))
      }
    }
  })
}

func TestDealing (t *testing.T) {
  g := Game{}
  g.New()

  dealt := make(map[int]int)

  // Keep track of how many of each card have been dealt.
  countRow := func(row int) {
    for _, v := range g.tower[row] {
      dealt[v]++
    } 
  }

  // Check that the deck contains the right amount of each card.
  assertCounts := func(t *testing.T, g Game) {
    t.Helper()
    for v, count := range dealt {
      startCount := 8
      if v == 0 {
        // Hero card
        startCount = 4
      }
      want := startCount - count
      got := g.counts[v]
      if got != want {
        t.Fatalf("deck has wrong amount of cards of value %d, want %d, got %d", v, want, got)
      }
    }
  }

  // Deal the whole 36 card tower (ignore burned cards for now)
  t.Run("Deal whole tower and check counts", func(t *testing.T) {
    for i := 0; i < 8; i++ {
      g.Deal()
      countRow(i)
    }
    assertCounts(t, g)

    wantedNumOfCardsDealt := 0
    numOfCardsDealt := 0
    for i := 0; i < 8; i++ {
      numOfCardsDealt += len(g.tower[i])
      wantedNumOfCardsDealt += i+1
    }

    if wantedNumOfCardsDealt != numOfCardsDealt {
      t.Fatalf("Wrong amount of cards dealt, want %d, got %d", wantedNumOfCardsDealt, numOfCardsDealt)
    }
    
    if len(g.deck) != 60 - wantedNumOfCardsDealt {
      t.Fatalf("Wrong amount of cards left in deck, want %d, got %d", 60 - wantedNumOfCardsDealt, len(g.deck))
    }
  })
}

func TestCashOut (t *testing.T) {
  g := Game{}
  g.New()
  g.deck[1], g.deck[2] = 1, 1

  g.Deal()
  g.Deal()

  rowVal := 2

  g.CashOut()

  if g.Balance() != rowVal {
    t.Errorf("balance after cashing out should be %d, got %d", rowVal, g.Balance())
  }
}

func TestInput (t *testing.T) {

  t.Run("at game start, z deals first two rows", func (t *testing.T) {
    g := Game{}
    g.New()
  
    // in := strings.NewReader("z\n")
    in := "z"
    g.Input(in)

    if len(g.tower[0]) != 1 {
      t.Error("gate card should have been dealt")
    }
    if len(g.tower[1]) != 2 {
      t.Error("2nd row should have been dealt")
    }

    for deal := 2; deal <= 7; deal++ {
      t.Run(fmt.Sprintf("deal %d should only deal one row", deal), func (t *testing.T) {
        g.Input(in)

        if len(g.tower[deal]) != deal+1 {
          t.Errorf("row %d should have been dealt", deal)
        }
        
        if deal < 7 {
          if len(g.tower[deal+1]) != 0 {
            t.Errorf("row %d should not have been dealt", deal)
          }
        }
      })
    }
  })

  t.Run("after first deal, x cashes out", func (t *testing.T) {
    g := Game{}
    g.New()
    g.balance = 0

    g.Deal()
    g.Deal()

    rowVal := 0
    for _, v := range g.tower[1] {
      rowVal += v
    }

    // in := strings.NewReader("x\n")
    in := "x"
    g.Input(in)

    if g.Balance() != rowVal {
      t.Errorf("balance after cashing out should be %d, got %d", rowVal, g.Balance())
    }
  })

  t.Run("x does nothing if no cards have been dealt", func (t *testing.T) {
    t.Fail()
  })
}

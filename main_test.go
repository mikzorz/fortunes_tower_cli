package main

import "testing"

func TestNewGame (t *testing.T) {
  // When a new game is created, check contents of deck and tower
  g := Game{}
  g.New()

  // deck is a []int
  // tower is [][]int

  t.Run("deck should not be nil", func(t *testing.T) {
    if g.deck == nil {
      t.Fatalf("Game.New() did not create a deck")
    }
  })

  t.Run("deck should contain 4 hero cards", func(t *testing.T) {
    c, ok := g.deck[0]
    if !ok {
      t.Fatalf("card with value %d should exist but doesn't", 0)
    }

    if c != 4 {
      t.Errorf("wrong amount of hero cards, want 4, got %d", c)
    }
  })

  t.Run("deck should contain 8 copies of 1-7", func(t *testing.T) {
    for i := 1; i <= 7; i++ {
      c, ok := g.deck[i]
      if !ok {
        t.Fatalf("card with value %d should exist but doesn't", i)
      }

      if c != 8 {
        t.Errorf("wrong amount of %d-value cards, want 8, got %d", i, c)
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
      got := g.deck[v]
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
    for i := 1; i <= 8; i++ {
      wantedNumOfCardsDealt += i
    }
    numOfCardsDealt := 0
    for i := 0; i < 8; i++ {
      numOfCardsDealt += len(g.tower[i])
    }

    if wantedNumOfCardsDealt != numOfCardsDealt {
      t.Fatalf("Wrong amount of cards dealt, want %d, got %d", wantedNumOfCardsDealt, numOfCardsDealt)
    }
  })
}

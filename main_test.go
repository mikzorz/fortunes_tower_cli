package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestNewGame(t *testing.T) {
	// When a new game is created, check contents of deck and tower
	g := NewGame()

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

	t.Run("state should be StateBetting", func(t *testing.T) {
		if g.State() != StateBetting {
			t.Fatalf("state should be %d, got %d", g.State(), StateBetting)
		}
	})
}

func TestDealing(t *testing.T) {

	dealt := make(map[int]int)

	// Keep track of how many of each card have been dealt.
	countRow := func(g Game, row int) {
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
		g := NewGame()

		g.deck = safeDeck()

		for i := 0; i < 8; i++ {
			g.deal()
			countRow(g, i)
		}
		assertCounts(t, g)

		wantedNumOfCardsDealt := 0
		numOfCardsDealt := 0
		for i := 0; i < 8; i++ {
			numOfCardsDealt += len(g.tower[i])
			wantedNumOfCardsDealt += i + 1
		}

		if wantedNumOfCardsDealt != numOfCardsDealt {
			t.Fatalf("Wrong amount of cards dealt, want %d, got %d", wantedNumOfCardsDealt, numOfCardsDealt)
		}

		if len(g.deck) != 60-wantedNumOfCardsDealt {
			t.Fatalf("Wrong amount of cards left in deck, want %d, got %d", 60-wantedNumOfCardsDealt, len(g.deck))
		}
	})

	t.Run("dealing should change state to StatePlaying", func(t *testing.T) {
		g := NewGame()

		g.deal()

		if g.State() != StatePlaying {
			t.Fatalf("state is %d, should be %d", g.State(), StatePlaying)
		}
	})

	t.Run("first deal should subtract wager", func(t *testing.T) {
		g := NewGame()

		balBefore := g.Balance()
		g.deal()
		balAfter := g.Balance()

		if balAfter != balBefore-g.GetWager() {
			t.Errorf("wager not subtracted, got %d, want %d", balAfter, balBefore-g.GetWager())
		}
	})

}

func TestCashOut(t *testing.T) {
	t.Run("balance increases by last row value", func(t *testing.T) {
		g := NewGame()
		g.deck[1], g.deck[2] = 1, 2

		g.dealX(2)

		balBeforecashOut := g.Balance()
		rowVal := 3

		g.cashOut()

		want := balBeforecashOut + (rowVal * g.multiplier)
		if g.Balance() != want {
			t.Errorf("balance after cashing out should be %d, got %d", want, g.Balance())
		}
	})

	t.Run("round should end and return to betting state", func(t *testing.T) {
		g := NewGame()
		// Play a few rows
		g.dealX(3)

		// Cash out and return to betting state
		g.cashOut()

		if g.State() != StateBetting {
			t.Fatalf("cashing out should return state to betting, got %v", g.State())
		}
	})

	t.Run("round should stop and cash out after row 8 is played", func(t *testing.T) {
		g := NewGame()

		for i := 0; i < 8; i++ {
			g.deal()
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("should not try to deal 9th row")
			}
		}()

		g.deal()

		if g.State() != StateBetting {
			t.Errorf("want state %d, got %d", StateBetting, g.State())
		}
	})

	t.Run("cashOut() should do nothing if current row is 0", func(t *testing.T) {
		g := NewGame()

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("cashing out when tower is empty causes panic")
			}
		}()

		g.cashOut()
	})

	t.Run("cashing out should reset deck, counts and tower", func(t *testing.T) {
		g := NewGame()

		for i := 0; i < 4; i++ {
			g.deal()
		}

		g.cashOut()

		assertGameReset(t, g)
	})

	t.Run("set current row to 0", func(t *testing.T) {
		g := NewGame()

		g.deal()
		g.cashOut()

		if g.curRow != 0 {
			t.Fatalf("g.curRow wasn't reset to 0")
		}
	})

	t.Run("if last row is played without using gate, JACKPOT", func(t *testing.T) {
		g := NewGame()
		g.deck = deckNoMultis()

		t.Log(g.deck)
		g.dealX(8)
		t.Log(g.tower)

		want := 0
		for i := 1; i <= 7; i++ {
			for _, v := range g.tower[i] {
				want += v
			}
		}

		balBefore := g.Balance()
		g.cashOut()
		balAfter := g.Balance()

		diff := balAfter - balBefore

		if diff != want {
			t.Fatalf("last row did not give jackpot reward, got %d, want %d", diff, want)
		}
	})

}

func TestGetRowValue(t *testing.T) {
	g := NewGame()
	g.deck = []int{1, 1, 2}
	g.dealX(2)

	if g.getRowValue(1) != 3 {
		t.Fatalf("getRowValue() returned %d, want %d", g.getRowValue(1), 3)
	}
}

// These tests are retesting other functionality instead of just testing input. Increase complexity for purer tests?
func TestInput(t *testing.T) {

	t.Run("at game start, z deals first two rows", func(t *testing.T) {
		g := NewGame()

		g.deck = safeDeck()

		in := "z"
		g.Input(in)

		if len(g.tower[0]) != 1 {
			t.Error("gate card should have been dealt")
		}
		if len(g.tower[1]) != 2 {
			t.Error("2nd row should have been dealt")
		}

		for deal := 2; deal <= 7; deal++ {
			t.Run(fmt.Sprintf("deal %d should only deal one row", deal), func(t *testing.T) {
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

	t.Run("after first deal, x cashes out", func(t *testing.T) {
		g := NewGame()
		g.balance = 0

		g.dealX(2)

		rowVal := 0
		for _, v := range g.tower[1] {
			rowVal += v
		}

		balBeforecashOut := g.Balance()

		in := "x"
		g.Input(in)

		want := balBeforecashOut + rowVal*g.multiplier
		if g.Balance() != want {
			t.Errorf("balance after cashing out should be %d, got %d", want, g.Balance())
		}
	})

	for _, in := range []string{"z", "x"} {
		t.Run(fmt.Sprintf("after gameover, %s resets to betting state and empty tower", in), func(t *testing.T) {
			g := NewGame()

			g.gameOver()

			g.Input(in)

			if g.State() != StateBetting {
				t.Fatalf("wrong game state, got %d, want %d", g.State(), StateBetting)
			}

			// Not working?
			assertGameReset(t, g)
		})
	}
}

func TestBust(t *testing.T) {
	assertGateUsed := func(t *testing.T, g Game) {
		t.Helper()
		if len(g.tower[0]) > 0 {
			t.Fatalf("first row should be empty")
		}
	}

	t.Run("leftmost card busts and replaced with gate", func(t *testing.T) {
		g := NewGame()

		g.deck = []int{
			7,
			1, 1,
			1, 2, 2,
		}

		g.dealX(3)

		if g.tower[2][0] != 7 {
			t.Fatalf("gate card should replace burned card")
		}

		assertGateUsed(t, g)

		if g.IsGameOver() {
			t.Fatalf("game should continue")
		}
	})

	t.Run("rightmost card busts and replaced with gate", func(t *testing.T) {
		g := NewGame()

		g.deck = []int{
			7,
			1, 1,
			2, 2, 1,
		}

		g.dealX(3)

		if g.tower[2][2] != 7 {
			t.Fatalf("gate card should replace burned card")
		}

		assertGateUsed(t, g)

		if g.IsGameOver() {
			t.Fatalf("game should continue")
		}
	})

	t.Run("middle card busts and replaced with gate, game continues", func(t *testing.T) {
		g := NewGame()

		g.deck = []int{
			7,
			1, 1,
			2, 2, 2,
			3, 2, 3, 3,
		}

		g.dealX(4)

		if g.tower[3][1] != 7 {
			t.Fatalf("gate card should replace burned card")
		}

		assertGateUsed(t, g)

		if g.IsGameOver() {
			t.Fatalf("game should continue")
		}
	})

	t.Run("middle card busts and replaced with gate, game over", func(t *testing.T) {
		g := NewGame()

		g.deck = []int{
			7,
			1, 7,
			2, 1, 2,
		}

		g.dealX(3)

		if g.tower[2][1] != 7 {
			t.Fatalf("gate card should replace burned card")
		}

		assertGateUsed(t, g)

		if !g.IsGameOver() {
			t.Fatalf("game should end")
		}
	})

	t.Run("don't bust if row contains hero", func(t *testing.T) {
		t.Run("hero dealt directly from the deck", func(t *testing.T) {
			g := NewGame()
			g.deck = []int{1, 1, 2, 0, 2, 3}

			g.dealX(3)

			// check that player hasnt busted
			if bust, _ := g.IsBust(); bust {
				t.Fatalf("should not have bust")
			}
			// check that gate wasn't used
			if len(g.tower[0]) != 1 {
				t.Fatalf("should not have used gate card")
			}
		})

		t.Run("hero saves a bust row", func(t *testing.T) {
			g := NewGame()
			g.deck = []int{0, 1, 2, 1, 2, 3}

			g.dealX(3)

			// check that player hasnt bust
			if bust, _ := g.IsBust(); bust {
				t.Fatalf("should not have bust")
			}
			// check that gate was used
			if len(g.tower[0]) != 0 {
				t.Fatalf("should have used gate card")
			}
		})
	})
}

func TestMultiplier(t *testing.T) {
	t.Run("on round start, multiplier equals wager / 15", func(t *testing.T) {
		g := NewGame()
		g.SetWager(45)
		g.deal()

		if g.multiplier != 3 {
			t.Fatalf("multiplier should be wager / 15, want %d, got %d", 3, g.multiplier)
		}
	})

	t.Run("multipliers should compound", func(t *testing.T) {
		// put double 1 in second row of deck, multiplier should become x2.
		g := NewGame()
		g.tower[1] = []int{1, 1}
		g.curRow = 1
		g.checkMulti()

		if g.multiplier != 2 {
			t.Fatalf("multiplier should be 2, got %d", g.multiplier)
		}

		// row 3, triple 2, multi should be x6.
		g.tower[2] = []int{2, 2, 2}
		g.curRow = 2
		g.checkMulti()

		if g.multiplier != 6 {
			t.Fatalf("multiplier should be 6, got %d", g.multiplier)
		}
	})

	t.Run("multiplier should increase even after gate is used", func(t *testing.T) {
		g := NewGame()
		g.deck = []int{2, 1, 7, 1, 2, 2}

		g.dealX(3)

		if g.multiplier != 3 {
			t.Fatalf("multiplier should be %d, got %d", 3, g.multiplier)
		}
	})

	t.Run("deal() and cashOut() use multiplier", func(t *testing.T) {
		g := NewGame()
		g.deck = []int{0, 1, 1, 2, 2, 2, 3, 3, 3, 3}

		g.dealX(4)

		// want := 24
		// if g.multiplier != want {
		//   t.Fatalf("wrong multiplier, got %d, want %d", g.multiplier, want)
		// }
		balBefore := g.Balance()
		g.cashOut()
		balAfter := g.Balance()

		diff := balAfter - balBefore
		want := 12 * 24
		if diff != want {
			t.Fatalf("cashOut() changes balance by wrong amount, want %d, got %d", want, diff)
		}
	})
}

func TestPrinting(t *testing.T) {
	t.Run("Game.out should default to stdout", func(t *testing.T) {
		g := NewGame()
		if g.out != os.Stdout {
			t.Fatalf("g.out should be os.Stdout, got %v", g.out)
		}
	})

	t.Run("gate card should be shown as [?] until revealed", func(t *testing.T) {
		g := NewGame()
		g.deck = []int{1, 2, 3, 2, 4, 5}

		out := &bytes.Buffer{}
		g.out = out

		g.deal()

		g.PrintRow(0)
		got := strings.TrimSpace(out.String())
		want := "[?]"
		if got != want {
			t.Fatalf("gate card should be obscured, want %s, got %s", want, got)
		}

		out.Reset()
		g.dealX(2)
		g.PrintRow(0)
		got = strings.TrimSpace(out.String())
		want = "[ ]"
		if got != want {
			t.Fatalf("gate card should be blank, want %s, got %s", want, got)
		}
	})

	t.Run("each row should end with its value", func(t *testing.T) {
		g := NewGame()
		out := &bytes.Buffer{}
		g.out = out

		g.dealX(2)
		g.PrintRow(1)

		rowVal := g.getRowValue(1)

		txt := out.String()
		txt = txt[len(txt)-(2+lenOfNum(rowVal))-1 : len(txt)-1] // assumes value is enclosed in 1 delimiter each side and ends with \n

		want := fmt.Sprintf("(%d)", rowVal)
		if txt != want {
			t.Fatalf("row value should show %s, got %s", want, txt)
		}
	})

	t.Run("jackpot should show jackpot value", func(t *testing.T) {
		g := NewGame()
		g.deck = deckNoMultis()
		out := &bytes.Buffer{}
		g.out = out

		g.dealX(8)
		g.PrintRow(7)

		jackpotVal := g.getJackpotValue()

		txt := out.String()
		txt = txt[len(txt)-(2+lenOfNum(jackpotVal))-1 : len(txt)-1] // assumes value is enclosed in 1 delimiter each side and ends with \n

		want := fmt.Sprintf("(%d)", jackpotVal)
		if txt != want {
			t.Fatalf("row value should show %s, got %s", want, txt)
		}
	})

  // test the printed instructions and money
  // what if i replace lines in place? will that affect tests?
}

func lenOfNum(i int) int {
	return len(strconv.Itoa(i))
}

func safeDeck() []int {
	// create deck with no busts
	deck := []int{}
	for i := 0; i < 8; i++ {
		for j := 0; j <= i; j++ {
			deck = append(deck, i)
		}
	}
	for i := len(deck); i < 60; i++ {
		deck = append(deck, 7)
	}
	return deck
}

func deckNoMultis() []int {

	return []int{
		0,
		1, 2,
		3, 4, 5,
		1, 1, 7, 7,
		2, 2, 2, 6, 6,
		3, 3, 3, 4, 4, 4,
		2, 2, 2, 2, 7, 7, 7,
		5, 5, 5, 5, 5, 5, 5, 1,
	}
}

func assertGameReset(t *testing.T, g Game) {
	if len(g.deck) != 60 {
		t.Error("deck was not reset")
	}

	if g.counts[0] != 4 {
		t.Error("counts were not reset")
	} else {
		for v := 1; v <= 7; v++ {
			if g.counts[v] != 8 {
				t.Error("counts were not reset")
				break
			}
		}
	}

	for r := 0; r < 8; r++ {
		if len(g.tower[r]) > 0 {
			t.Error("tower was not reset")
			break
		}
	}
}

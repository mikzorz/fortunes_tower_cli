# CLI demake of Fortune's Tower from Fable 2.

## How to play

- The player bets a multiple of 15 gold.
- The deck contains 8 copies each of cards with values 1-7 along with 4 copies of the Hero card. *(This deck is called the Diamond Deck)*
- A Hero protects all cards on its row.
- A Gate card is played first, at the top, face down.
- Hitting deals a new row of cards below the previous.
- Each row contains 1 more card than the row above it, forming a triangle.
- If a card shares a value with one of the cards directly above it, it becomes burned.
- If any cards are still burned, and if the Gate card is still face down, replace the burned card with the Gate card.
- If, after that, any cards are still burned, or if the 8th row is played, the round ends.
- If you reach the bottom of the tower (8th row) without using the Gate card, the final score for that round is equal to the value of ALL cards in the tower.
- If all cards in a row have the same value (including Hero cards), the final score is multiplied by the amount of cards in that row.
- The prize is multiplied by the bet / 15.



todo

- BUG: some tests do not use predetermined decks, occasionally returns an error
- allow player to change wager (for accuracy)
- add other decks from F2 (accuracy, but not important to me)
- change printing to replace, not append (nice to have)
- custom tower sizes? (n2h)
- colours? (n2h)
- make code more nicerer (custom types, methods for read/writing to/from tower etc)

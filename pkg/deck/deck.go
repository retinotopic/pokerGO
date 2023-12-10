package deck

import (
	"math/rand"
)

type card struct {
	Value int
	Suit  int
}

// Deck represents a deck of cards
type deck []card

// NewDeck returns a new shuffled deck of 52 cards
func NewDeck() deck {
	deck := make(deck, 52)
	suits := []int{0, 1, 2, 3}
	values := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	for i := range deck {
		deck[i] = card{
			Value: values[i%13],
			Suit:  suits[i/13],
		}
	}
	return deck
}

// Shuffle randomly shuffles a deck of cards
func Shuffle(deck deck, r *rand.Rand) {
	for i := range deck {
		j := r.Intn(len(deck) - 1)
		deck[i], deck[j] = deck[j], deck[i]
	}

}

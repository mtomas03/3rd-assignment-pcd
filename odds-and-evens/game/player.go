package game

import (
	"math/rand/v2"
)

// PlayRequest acts as a signal from the referee to notify a player that it is their turn.
type PlayRequest struct{}

// PlayerChannels encapsulates a player's identity and their communication channels.
// Passing this struct ensures that the winning player's active goroutine and
// specific channels are seamlessly promoted to the next round
// without relying on shared state.
type PlayerChannels struct {
	ID        int
	RequestCh chan PlayRequest // referee -> player (signals it's time to play)
	MoveCh    chan int         // player -> referee (delivers the chosen number)
}

// Player represents the active goroutine for a participant.
// It idles until receiving a PlayRequest, generates a random integer between 1 and 5,
// and sends it back via the move channel. The goroutine automatically terminates
// when the request channel is closed.
func Player(requestCh <-chan PlayRequest, moveCh chan<- int) {
	for range requestCh {
		move := rand.IntN(5) + 1
		moveCh <- move
	}
}

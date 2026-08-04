package game

import (
	"fmt"
)

// Referee manages a single match between two players.
// Designed to run concurrently, it signals its completion to the caller using a WaitGroup.
//
// Communication protocol:
//  1. Sends a PlayRequest to both players via unbuffered channels.
//  2. Waits for a move from each player.
//  3. Calculates the sum of the moves to determine the parity (even/odd).
//  4. Sends the winning player's channels to the result channel.
func referee(p1, p2 PlayerChannels, resultCh chan<- PlayerChannels) {
	// Trigger both players' turns.
	p1.RequestCh <- PlayRequest{}
	p2.RequestCh <- PlayRequest{}

	// Retrieve moves from each player.
	m1 := <-p1.MoveCh
	m2 := <-p2.MoveCh

	// Player 1 wins on an odd sum, Player 2 wins on an even sum.
	sum := m1 + m2
	parityLabel := "even"
	winner := p2
	if sum%2 != 0 {
		parityLabel = "odd"
		winner = p1
	}

	fmt.Printf("Player %d chose %d && Player %d chose %d = sum %d (%s) -> Player %d wins\n",
		p1.ID, m1, p2.ID, m2, sum, parityLabel, winner.ID)

	// Report the winner back to the round coordinator.
	resultCh <- winner
}

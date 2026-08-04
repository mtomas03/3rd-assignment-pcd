/*
Concurrency architecture:
  - Each player is a long-lived goroutine, uniquely identified by an integer ID.
  - Zero shared memory; all synchronization relies on synchronous message passing via unbuffered channels.
  - Referees are short-lived goroutines spawned per match. They send out play requests,
    gather moves, determine the winner and report the result back to the coordinator.
  - Round coordinator manages the tournament round by round, spawning referees concurrently,
    and waiting for their completion using a sync.WaitGroup.
*/
package main

import (
	"fmt"
	"odds-and-evens/game"
)

func main() {
	const numRounds = 3
	numPlayers := 1 << numRounds // N = 2^m

	fmt.Printf("Odds-and-Evens tournament - %d players, %d rounds\n", numPlayers, numRounds)

	// Initialize N players for the tournament.
	// Players are instantiated once and persist throughout the entire tournament.
	// Unbuffered channels ensure synchronous message passing (rendezvous).
	players := make([]game.PlayerChannels, numPlayers)
	for i := range players {
		pc := game.PlayerChannels{
			ID:        i,
			RequestCh: make(chan game.PlayRequest),
			MoveCh:    make(chan int),
		}
		players[i] = pc
		go game.Player(pc.RequestCh, pc.MoveCh)
	}

	// Execute the tournament over m rounds.
	current := players
	for r := 1; r <= numRounds; r++ {
		current = game.RunRound(current, r)
	}

	fmt.Printf("\nTournament winner: Player %d\n", current[0].ID)

	// Graceful shutdown: closing the request channels signals all
	// player goroutines to terminate cleanly.
	for _, p := range players {
		close(p.RequestCh)
	}
}

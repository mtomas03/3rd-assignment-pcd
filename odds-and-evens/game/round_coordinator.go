package game

import (
	"fmt"
	"sync"
)

// RunRound manages all matches for a specific tournament round.
// It matches consecutive players in the slice, launches a referee goroutine for each match,
// and collects the winners to advance them to the next round.
func RunRound(players []PlayerChannels, round int) []PlayerChannels {
	games := len(players) / 2
	fmt.Printf("\n- Round %d\n", round)

	resultCh := make(chan PlayerChannels)
	var wg sync.WaitGroup

	// Spawn a referee for each pair of players.
	for i := 0; i < len(players); i += 2 {
		wg.Go(func() {
			referee(players[i], players[i+1], resultCh)
		})
	}

	// Spawn a goroutine to close the result channel once all referees finish.
	// This ensures the subsequent range loop can safely terminate.
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect and return the winning players.
	// The range loop retrieves each winner, triggering the corresponding referee one at a time.
	winners := make([]PlayerChannels, 0, games)
	for w := range resultCh {
		winners = append(winners, w)
	}
	return winners
}

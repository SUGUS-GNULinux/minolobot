package interaction

import (
	"math/rand"
	"time"
)

// StartPoleLogic starts an internal counter which increments every hour,
// the higher the value of the counter the more possible the emmission of a
// true value in the returned channel
func StartPoleLogic() <-chan bool {
	counter := 0
	tick := time.NewTicker(time.Hour)
	poleSignal := make(chan bool)
	go func() {
		for _ = range tick.C {
			counter++
			if rand.Intn(100) < counter {
				poleSignal <- true
				counter = 0
			}
		}
	}()
	return poleSignal
}

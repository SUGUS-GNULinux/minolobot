package interaction

import (
	"math/rand"
	"time"
)

// StartPoleLogic starts the logic to send a signal at 00:00, when it's
// started calculates the time until the next day and after that blocks in periods
// of 24h.
func StartPoleLogic() <-chan bool {
	poleSignal := make(chan bool)
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	untilTwelve := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		0, 0, 0, 0, tomorrow.Location()).Sub(now)
	firstSignal := time.After(untilTwelve)
	go func() {
		// block until 00:00
		<-firstSignal
		tick := time.NewTicker(time.Hour * 24)
		go delayActivation(poleSignal)
		for _ = range tick.C {
			go delayActivation(poleSignal)
		}
	}()
	return poleSignal
}

func delayActivation(c chan bool) {
	<-time.After(time.Millisecond * time.Duration(rand.Intn(4000)))
	c <- true
}

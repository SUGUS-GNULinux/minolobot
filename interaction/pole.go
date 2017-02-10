package interaction

import (
	"math/rand"
	"minolobot/config"
	"time"
)

// StartPoleLogic starts the logic to send a signal at 00:00, when it's
// started calculates the time until the next day and after that blocks in periods
// of 24h.
func StartPoleLogic() <-chan bool {
	poleSignal := make(chan bool)
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	untilTwelve := now.Sub(time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
		0, 0, 0, 0, tomorrow.Location()))
	firstSignal := time.After(untilTwelve)
	go func() {
		// block until 00:00
		<-firstSignal
		tick := time.NewTicker(time.Hour * 24)
		if config.Enabled {
			go delayActivation(poleSignal)
		}
		for _ = range tick.C {
			if config.Enabled {
				go delayActivation(poleSignal)
			}
		}
	}()
	return poleSignal
}

func delayActivation(c chan bool) {
	<-time.After(time.Millisecond * time.Duration(rand.Intn(4000)))
	c <- true
}

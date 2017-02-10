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

// TODO which id to do pole?
// actualTime := time.Now()
// if actualTime.Hour() == 0 && actualTime.Minute() == 0 {
// 	msg := tgbotapi.NewMessage(ADD ID HERE, "pole")
// 	go func(){
// 		wait := rand.Intn(5000)
// 		<-time.After(wait * time.Millisecond)
// 		bot.Send(msg)
// 	}
// }

// Copyright 2017-2018 SUGUS GNU/Linux <sugus@us.es>
//
// This file is part of Minolobot.
//
//     Minolobot is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     Minolobot is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with Minolobot.  If not, see <http://www.gnu.org/licenses/>.

package interaction

import (
	"log"
	"math/rand"
	"time"
)

// StartPoleLogic starts the logic to send a signal at 00:00, when it's
// started calculates the time until the next day and after that blocks in periods
// of 24h.
func StartPoleLogic() <-chan bool {
	poleSignal := make(chan bool)

	// Load timezone
	utc, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		log.Fatal(err)
	}
	// Load current time with the deviation of the zone
	now := time.Now().In(utc)
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

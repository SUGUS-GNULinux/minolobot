// Copyright 2017 Alejandro Sirgo Rica
// Copyright 2018 Manuel López Ruiz <manuellr.git@gmail.com>
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

// Package config holds everything related to configuration
package config

import (
	"encoding/json"
	"github.com/SUGUS-GNULinux/minolobot/utilities"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
	"time"
)

// ChatConfig contains the base configuration for a single user
type ChatConfig struct {
	// Enabled defines if the bot should answer
	Enabled bool
	// DisableUntil defines until bot should not answer
	DisableUntil time.Time
	// PercentAnswer is the probability of anwser to any update
	PercentAnswer int
	// IsGroup is true if the related chat is a group
	IsGroup bool
	// Pole is true if the bot should say "Pole" every day
	Pole bool
}

// bucketChatConfig define the bucketName in BoltDB
const bucketChatConfig = "chatConfig"

// Contains the timer for the automatic conversation enabled
var AutoEnableTimer map[int64]*time.Timer

func init() {
	AutoEnableTimer = make(map[int64]*time.Timer)
	reScheduleAutoEnable()
}

func CreateChatConfig(chatId int64, isGroup bool) *ChatConfig {
	conf := NewChatConfig(isGroup)
	err := utilities.Update(bucketChatConfig, strconv.FormatInt(chatId, 10), conf)
	if err != nil {
		log.Println(err)
	}
	return conf
}

func UpdateChatConfig(chatId int64, input *ChatConfig) error {
	err := utilities.Update(bucketChatConfig, strconv.FormatInt(chatId, 10), input)
	if err != nil {
		log.Println(err)
	}
	return err
}

func FindChatConfig(chatId int64) (res *ChatConfig, err error) {
	err = utilities.View(bucketChatConfig, strconv.FormatInt(chatId, 10), &res)
	return
}

func FindAllChatConfigWithId() (map[int64]ChatConfig, error) {
	res := make(map[int64]ChatConfig)

	err := utilities.BoltDB.View(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucketChatConfig))
		if b == nil {
			return bolt.ErrBucketNotFound
		}
		// Retrieve the records
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			key, err := strconv.ParseInt(string(k), 10, 0)
			if err != nil {
				log.Println(err)
			}

			// Decode the record
			var value ChatConfig
			err = json.Unmarshal(v, &value)
			if err != nil {
				log.Println(err)
			}

			res[key] = value
		}
		return nil
	})

	return res, err
}

// NewChatConfig creates a default chat configuration
func NewChatConfig(isGroup bool) *ChatConfig {
	return &ChatConfig{Enabled: true, PercentAnswer: 5, IsGroup: isGroup, Pole: isGroup}
}

// EnabledChatConfig enable/disable the bot answer
func EnabledChatConfig(enable bool, scheduleEnableFor *time.Time, chatId int64) error {
	log.Println(chatId, enable)
	chatConfig, err := FindChatConfig(chatId)
	if err != nil {
		return err
	}

	// Stop pending AutoEnableTimer if exist
	if a, ok := AutoEnableTimer[chatId]; ok {
		a.Stop()
	}

	// Modify chatConfig properties and schedule enable
	if scheduleEnableFor != nil {
		chatConfig.DisableUntil = *scheduleEnableFor
		chatConfig.Enabled = enable
		go ScheduleAutoEnable(chatId, scheduleEnableFor)
	} else {
		chatConfig.Enabled = enable
		chatConfig.DisableUntil = time.Time{}
	}

	// Persist in the database
	err = UpdateChatConfig(chatId, chatConfig)
	if err != nil {
		return err
	}
	return nil
}

// reScheduleAutoEnable reShedule the auto enable chatConfig after a bot restart
func reScheduleAutoEnable() {
	chatConfigs, _ := FindAllChatConfigWithId()
	now := time.Now()
	for k, v := range chatConfigs {
		tmpTime := v.DisableUntil // Prevent lazy pointers
		if tmpTime.After(now) {
			go ScheduleAutoEnable(k, &v.DisableUntil)
		} else if !tmpTime.Equal(time.Time{}) {
			v.DisableUntil = time.Time{}
			v.Enabled = true
			UpdateChatConfig(k, &v)
		}
	}
}

// ScheduleAutoEnable schedule the auto enable chatConfig
func ScheduleAutoEnable(chatId int64, scheduleFor *time.Time) {
	d := scheduleFor.Sub(time.Now())
	timer := time.NewTimer(d)
	AutoEnableTimer[chatId] = timer
	log.Println("ScheduleAutoEnable: ", chatId, "\\^/", scheduleFor)

	// Wait until timer ends
	<-timer.C

	// Delete timer from pending autoEnable
	delete(AutoEnableTimer, chatId)

	EnabledChatConfig(true, nil, chatId)
}

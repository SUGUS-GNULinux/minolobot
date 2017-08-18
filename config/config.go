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
	"bufio"
	"encoding/csv"
	"log"
	"os"
)

// ChatConfig contains the base configuration for a single user
type ChatConfig struct {
	// Enabled defines if the bot should answer
	Enabled bool
	// PercentAnswer is the probability of anwser to any update
	PercentAnswer int
	// IsGroup is true if the related chat is a group
	IsGroup bool
	// Pole is true if the bot should say "Pole" every day
	Pole bool
}

// NewChatConfig creastes a default chat configuration
func NewChatConfig(isGroup bool) *ChatConfig {
	return &ChatConfig{Enabled: true, PercentAnswer: 5, IsGroup: isGroup, Pole: isGroup}
}

var (
	// Token is the token value of the bot
	Token string
	// BotName contains the @name of the bot after initializing in an account
	BotName string
)

// init token
func init() {
	tokenFile, err := os.Open("datafiles/token")
	if err != nil {
		log.Fatal(err)
	}
	defer tokenFile.Close()
	scanner := bufio.NewScanner(tokenFile)
	if scanner.Scan() {
		Token = string(scanner.Text())
	} else {
		log.Fatal("invalid token in token file content")
	}
}

var (
	// Phrases is the list of predefined phrases to say randomly
	Phrases = []string{}
	// CionList is the list of words ending with "cion"
	CionList map[string]bool
	// ConfigList contains all the user configurations
	ConfigList map[int64]*ChatConfig
)

func init() {
	ConfigList = make(map[int64]*ChatConfig)
}

// init "cion" ended words list
func init() {
	cionFile, err := os.Open("datafiles/cion.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer cionFile.Close()
	rc := csv.NewReader(cionFile)
	cionSlice, err := rc.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// init cion ended words map
	CionList = make(map[string]bool)
	for _, line := range cionSlice {
		CionList[line[0]] = true
	}
}

// init phrases list
func init() {
	phrasesFile, err := os.Open("datafiles/phrases.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer phrasesFile.Close()
	rp := csv.NewReader(phrasesFile)
	phrasesSlice, err := rp.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	// init phrases list
	for _, line := range phrasesSlice {
		Phrases = append(Phrases, line[0])
	}
}

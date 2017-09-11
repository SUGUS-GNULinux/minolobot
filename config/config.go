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

// Package config holds everything related to configuration
package config

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
)

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
)

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

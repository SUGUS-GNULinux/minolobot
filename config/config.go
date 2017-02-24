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
}

// NewChatConfig creastes a default chat configuration
func NewChatConfig(isGroup bool) *ChatConfig {
	return &ChatConfig{Enabled: true, PercentAnswer: 5, IsGroup: isGroup}
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

// init "cion" ended words list and phrases list
func init() {
	cionFile, err := os.Open("datafiles/cion.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer cionFile.Close()
	phrasesFile, err := os.Open("datafiles/phrases.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer phrasesFile.Close()
	// init cion ended words map
	rc := csv.NewReader(cionFile)
	cionSlice, err := rc.Read()
	if err != nil {
		log.Fatal(err)
	}
	CionList = make(map[string]bool)
	for _, word := range cionSlice {
		CionList[word] = true
	}
	// init phrases list
	rp := csv.NewReader(phrasesFile)
	Phrases, err = rp.Read()
	if err != nil {
		log.Fatal(err)
	}
}

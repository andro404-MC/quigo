package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type configSTR struct {
	Apikey  string
	Prompts map[string]prompt
}

type prompt struct {
	Text string
}

var (
	config          configSTR
	unstagedChanges bool
)

const configPath = "/.config/quigo/quigo.conf"

func load(c *configSTR) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error getting user home directory:", err)
		os.Exit(3)
	}

	fullPath := homeDir + configPath
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.Apikey = "GEMINI-PRO"
		c.Prompts = make(map[string]prompt)
		c.Prompts["Ask"] = prompt{Text: "give the shortest respond MAX 50 words"}
		c.Prompts["Correct"] = prompt{
			Text: "Correct the grammar of the following sentence without any extra text just pure correction",
		}
		return
	}

	_, err = toml.DecodeFile(fullPath, &c)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	return
}

func save(c *configSTR) {
	fileString := fmt.Sprintf("apikey = \"%s\"\n", c.Apikey)
	for i, x := range c.Prompts {
		fileString += fmt.Sprintf("\n[Prompts.\"%s\"]\n   Text = \"%s\"\n", i, x.Text)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error getting user home directory:", err)
		os.Exit(3)
	}

	fullPath := homeDir + configPath

	dirPath := filepath.Dir(fullPath)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	err = os.WriteFile(fullPath, []byte(fileString), 0o644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	unstagedChanges = false
	return
}

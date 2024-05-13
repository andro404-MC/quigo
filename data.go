package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type configSTR struct {
	Apikey string
}

var (
	config  configSTR
	prompts = map[string]string{
		"Ask":     "give the shortest respond MAX 50 words",
		"Correct": "Correct the grammar of the following sentence without any extra text just pure correction",
	}
)

func load(c *configSTR) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error getting user home directory:", err)
		os.Exit(3)
	}

	fullPath := homeDir + configPath
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.Apikey = "GEMINI-PRO"
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
	fileString := fmt.Sprintf("apikey = \"%s\"", c.Apikey)

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
	return
}

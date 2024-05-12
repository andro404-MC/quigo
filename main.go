package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BurntSushi/toml"
	"github.com/sqweek/dialog"
	"golang.design/x/clipboard"
)

type configSTR struct {
	Apikey string
}

var config configSTR

const (
	askP     = "give the shortest respond MAX 50 words"
	correctP = "Correct the grammar of the following sentence without any extra text just pure correction : "

	url        = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key="
	configPath = "/.config/quigo/quigo.conf"
)

var APIKEY string

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Quigo")
	myWindow.Resize(fyne.NewSize(800, 400))
	load(&config)
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	// home
	confirmBtn := widget.NewButton("Go", nil)
	confirmBtn.Disable()

	combo := widget.NewSelect([]string{"Ask", "Correct"}, nil)

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")
	input.MultiLine = true
	input.Wrapping = fyne.TextWrapBreak

	aitext := widget.NewRichText()
	aitext.ParseMarkdown(`# Quigo
  1. Choose your option
  2. Write your prompt
  3. Go !!`)
	aitext.Wrapping = fyne.TextWrapBreak
	aitext.Scroll = container.ScrollVerticalOnly

	loading := widget.NewProgressBarInfinite()
	loading.Hidden = true

	copyText := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.Write(clipboard.FmtText, []byte(aitext.String()))
	})
	copyText.Disable()

	btnValide := func(_ string) {
		if combo.Selected != "" && utf8.RuneCountInString(input.Text) > 2 &&
			utf8.RuneCountInString(config.Apikey) > 30 {
			confirmBtn.Enable()
		} else {
			confirmBtn.Disable()
		}
	}
	btnAction := func() {
		loading.Hidden, aitext.Hidden = false, true
		input.Disable()
		confirmBtn.Disable()

		value := input.Text
		var prompt string
		switch combo.Selected {
		case "Ask":
			prompt = askP
			break
		case "Correct":
			prompt = correctP
			break
		}

		respond, err, merr := handle(value, prompt)

		if err != nil {
			log.Println(err, " : ", merr)
			dialog.Message("%s", err).Title("Error").Error()
		} else {
			aitext.ParseMarkdown(respond)
			copyText.Enable()
		}

		loading.Hidden, aitext.Hidden = true, false
		input.Enable()
		confirmBtn.Enable()
	}

	combo.OnChanged, input.OnChanged = btnValide, btnValide
	confirmBtn.OnTapped = func() { go btnAction() }

	main := container.NewBorder(
		combo,
		container.NewGridWithColumns(2, confirmBtn, copyText),
		nil,
		nil,
		container.New(
			layout.NewGridLayoutWithColumns(1),
			container.NewHSplit(input, container.NewGridWithColumns(1, loading, aitext)),
		),
	)

	// Settings
	apiLabel := widget.NewLabel("APIKEY")
	apiInput := widget.NewPasswordEntry()
	apiAplly := widget.NewButton("Save", func() { save(&config); load(&config) })
	apiInput.Text = config.Apikey

	setting := container.NewBorder(
		container.NewBorder(nil, nil, apiLabel, nil, apiInput),
		apiAplly,
		nil,
		nil,
		nil,
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("", theme.HomeIcon(), main),
		container.NewTabItemWithIcon("", theme.SettingsIcon(), setting),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

func AIBULLSHIT(text string) string {
	time.Sleep(2 * time.Second)
	return text
}

func handle(value string, prompt string) (respond string, err error, moreError error) {
	// Define request payload
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": prompt + " : " + value,
					},
				},
			},
		},
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", errors.New("Error marshalling JSON"), err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url+config.Apikey, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", errors.New("Error creating request"), err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("Error sending request"), err
	}
	defer resp.Body.Close()

	// Read response body
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", errors.New("Error decoding response"), err
	}

	// Access and print generated content
	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", errors.New("Error : No candidates found in response"), err
	}

	content, ok := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	if !ok {
		return "", errors.New("Error : No content found in response"), err
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", errors.New("Error : No parts found in response"), err
	}

	generatedText, ok := parts[0].(map[string]interface{})["text"].(string)
	if !ok {
		return "", errors.New("Error : No text found in response"), err
	}

	return generatedText, nil, nil
}

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

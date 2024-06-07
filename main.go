package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/sqweek/dialog"
	"golang.design/x/clipboard"
)

var APIKEY string

const url = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key="

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Quigo")
	myWindow.Resize(fyne.NewSize(800, 400))
	load(&config)
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("AI", theme.ComputerIcon(), mainTab(myWindow, myApp)),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), settingTab()),
	)

	myWindow.SetCloseIntercept(
		func() {
			if unstagedChanges {
				ok := dialog.Message("%s", "Changes not saved. Do you still want to Quit?").
					Title("Quit ?").
					YesNo()

				if ok {
					myApp.Quit()
				}

			} else {
				myApp.Quit()
			}
		},
	)

	tabs.SetTabLocation(container.TabLocationTop)
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
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

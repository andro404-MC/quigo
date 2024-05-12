package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func settingTab() *fyne.Container {
	apiLabel := widget.NewLabel("APIKEY")
	apiInput := widget.NewPasswordEntry()
	apiInput.Text = config.Apikey
	apiInput.OnChanged = func(s string) { config.Apikey = s }
	apiAplly := widget.NewButton("Save", func() { save(&config); load(&config) })

	setting := container.NewBorder(
		container.NewBorder(nil, nil, apiLabel, nil, apiInput),
		apiAplly,
		nil,
		nil,
		nil,
	)

	return setting
}

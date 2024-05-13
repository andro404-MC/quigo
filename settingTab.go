package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/maps"
)

func settingTab() *fyne.Container {
	data := [][]string{
		maps.Keys(prompts),
		maps.Values(prompts),
	}

	apiLabel := widget.NewLabel("APIKEY")
	apiInput := widget.NewPasswordEntry()
	apiInput.Text = config.Apikey
	apiInput.OnChanged = func(s string) { config.Apikey = s }

	promptNameLabel := widget.NewLabel("Pormpt name")
	promptNameInput := widget.NewEntry()
	promptLabel := widget.NewLabel("Prompt")
	promptInput := widget.NewEntry()
	promptAdd := widget.NewButton("Add", nil)

	promptTable := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(data[0]), len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Col][i.Row])
		},
	)
	promptTable.StickyColumnCount = 2
	promptTable.ShowHeaderColumn = false

	settingsAplly := widget.NewButton("Save", func() { save(&config); load(&config) })

	setting := container.NewBorder(
		container.NewVBox(
			container.NewBorder(
				nil,
				nil,
				apiLabel,
				nil,
				apiInput,
			),
			widget.NewSeparator(),
			container.NewBorder(nil, nil, promptNameLabel, nil, promptNameInput),
			container.NewBorder(nil, nil, promptLabel, nil, promptInput),
			promptAdd,
		),
		settingsAplly,
		nil,
		nil,
		promptTable,
	)

	return setting
}

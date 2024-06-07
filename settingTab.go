package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/maps"
)

func settingTab() *fyne.Container {
	apiLabel := widget.NewLabel("APIKEY")
	apiInput := widget.NewPasswordEntry()
	apiInput.Text = config.Apikey
	apiInput.OnChanged = func(s string) {
		config.Apikey = s
		unstagedChanges = true
	}

	promptNameInput := widget.NewEntry()
	promptInput := widget.NewMultiLineEntry()
	promptAdd := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)
	promptAdd.Disable()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: promptNameInput},
			{Text: "Prompt", Widget: promptInput},
		},
	}

	promptDelText := widget.NewLabel("Prompt")

	promptShowText := widget.NewRichText()
	promptShowText.Wrapping = fyne.TextWrapWord
	promptShowText.Scroll = container.ScrollVerticalOnly

	promptDelCombo := widget.NewSelect(maps.Keys(config.Prompts), nil)
	promptDelCombo.OnChanged = func(s string) {
		promptShowText.ParseMarkdown(config.Prompts[s].Text)
	}
	promptDel := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		delete(config.Prompts, promptDelCombo.Selected)
		promptDelCombo.Options = maps.Keys(config.Prompts)
		promptDelCombo.Selected = ""
		combo.Options = maps.Keys(config.Prompts)
		combo.Selected = ""
		promptShowText.ParseMarkdown("")

		unstagedChanges = true
	})

	promptAdd.OnTapped = func() {
		config.Prompts[promptNameInput.Text] = prompt{Text: promptInput.Text}

		promptNameInput.Text, promptInput.Text = "", ""

		promptDelCombo.Options = maps.Keys(config.Prompts)
		combo.Options = maps.Keys(config.Prompts)

		form.Refresh()
		promptAdd.Disable()

		unstagedChanges = true
	}

	dataChanged := func(_ string) {
		if promptInput.Text == "" || promptNameInput.Text == "" {
			promptAdd.Disable()
		} else {
			promptAdd.Enable()
		}
	}

	promptInput.OnChanged = dataChanged
	promptNameInput.OnChanged = dataChanged
	settingsAplly := widget.NewButton("Save", func() { save(&config); load(&config) })

	setting := container.NewBorder(
		container.NewVBox(
			widget.NewSeparator(),
			container.NewBorder(
				nil,
				nil,
				apiLabel,
				nil,
				apiInput,
			),
			widget.NewSeparator(),
			container.NewBorder(
				nil,
				nil,
				nil,
				promptAdd,
				form,
			),
			widget.NewSeparator(),
			container.NewBorder(
				nil,
				nil,
				promptDelText,
				promptDel,
				promptDelCombo,
			),
		),
		settingsAplly,
		nil,
		nil,
		promptShowText,
	)

	return setting
}

# Quigo-gui

A simple Google Gemini prompter made with [fyne](https://github.com/fyne-io/fyne/) and the Power of Go.

<p align="center">
  <img
    src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white"
  />
    <img
    src="https://img.shields.io/badge/Gemini-8E75B2?style=for-the-badge&logo=googlebard&logoColor=fff"
  /><br />
  <img
    src="https://github.com/andro404-MC/quigo-gui/actions/workflows/test.yml/badge.svg"
  />
</p>

### Screenshots

![main](asset/Screenshot_2024-05-12-12-12-04_1366x768.png)
![correct](asset/Screenshot_2024-05-12-12-17-14_1366x768.png)

## Build :

> You need a to have `GOPATH` added to `PATH`

```
$ git clone https://github.com/andro404-MC/quigo-gui
$ cd quigo-gui

// Run
$ go run .

// Build
$ go Build .

// Install
$ go install .
```

## Usage :

If the path is set up correctly, you can just run:

```
$ quigo-gui
```

> [!NOTE]
> You need to get an API key from [aistudio](https://aistudio.google.com/app/apikey)

## Todo :

- Provide a way to add/remove prompts.

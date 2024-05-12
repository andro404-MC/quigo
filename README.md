# Quigo-gui


A dead simple app that let you know when your internet connection is back. (with polybar support)

> [!WARNING]
> the notification mode is probably only usable with dunst

#### Why does this exist :

sooo there was a time when my internet was so good at disapearing when needed so I made this lil app to just tell me when to go back the moment its back.

and yea I just kept messing with it once a while and its public now. :)

> my internet still sucks so this thing is actually usefull.

## Requirement :

Nothing unless :

`libnotify` : for notification support.

`go` : if you are going to build from source.

## Build :

> You need a to have `GOPATH` added to `PATH`

```
$ git clone https://github/untemi/pingo
$ cd pingo

// Run
$ go run .

// Install
$ go install .
```

## Usage :

If go path is setuped correctly you can just run the terminal version :

```
$ pingo
```

and for the notification version :

```
$ pingo -n
```

The default ping target is google.com but it can be changed using the environment variable `PINGOIP` :

```
$ PINGOIP=gnu.org pingo
```

## Todo :

- Optional notification sound
- Compleat the terminal experience

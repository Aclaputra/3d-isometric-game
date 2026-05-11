package main

import (
	"fmt"
	"log"

	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/go-steamworks"
	"golang.org/x/text/language"
)

const (
	appID       = 480
	isPublished = true
)

func main() {
	if isPublished {
		if steamworks.RestartAppIfNecessary(appID) {
			os.Exit(1)
		}
		if err := steamworks.Init(); err != nil {
			panic(fmt.Sprintf("steamworks.Init failed: %v", err))
		}
	}

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Habitate v0.01")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func SystemLang() language.Tag {
	switch steamworks.SteamApps().GetCurrentGameLanguage() {
	case "english":
		return language.English
	case "japanese":
		return language.Japanese
	}
	return language.Und
}

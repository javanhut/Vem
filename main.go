package main

import (
	"log"
	"os"

	gioapp "gioui.org/app"
	"gioui.org/unit"

	"github.com/javanhut/ProjectVem/internal/appcore"
)

func main() {
	go func() {
		w := new(gioapp.Window)
		w.Option(
			gioapp.Title("Vem - Vim Emulator"),
			gioapp.Size(unit.Dp(960), unit.Dp(640)),
		)
		filePaths := os.Args[1:]
		if err := appcore.Run(w, filePaths); err != nil {
			log.Printf("app exited: %v", err)
		}
		os.Exit(0)
	}()
	gioapp.Main()
}

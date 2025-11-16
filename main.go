package main

import (
	"os"

	gioapp "gioui.org/app"
	"gioui.org/unit"

	"github.com/javanhut/vem/internal/appcore"
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
			// Silently handle app exit errors
		}
		os.Exit(0)
	}()
	gioapp.Main()
}

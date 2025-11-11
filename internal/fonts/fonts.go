package fonts

import (
	_ "embed"
	"fmt"

	"gioui.org/font"
	"gioui.org/font/opentype"
)

//go:embed JetBrainsMonoNerdFont-Regular.ttf
var jetbrainsMonoRegular []byte

//go:embed JetBrainsMonoNerdFont-Bold.ttf
var jetbrainsMonoBold []byte

// Collection returns a collection of JetBrains Mono Nerd Font faces.
func Collection() ([]font.FontFace, error) {
	var fonts []font.FontFace

	// Parse regular font
	regularFace, err := opentype.Parse(jetbrainsMonoRegular)
	if err != nil {
		return nil, fmt.Errorf("failed to parse regular font: %w", err)
	}
	fonts = append(fonts, font.FontFace{
		Font: font.Font{Typeface: "JetBrainsMono"},
		Face: regularFace,
	})

	// Parse bold font
	boldFace, err := opentype.Parse(jetbrainsMonoBold)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bold font: %w", err)
	}
	fonts = append(fonts, font.FontFace{
		Font: font.Font{
			Typeface: "JetBrainsMono",
			Weight:   font.Bold,
		},
		Face: boldFace,
	})

	return fonts, nil
}

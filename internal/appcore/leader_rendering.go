package appcore

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// drawLeaderBar renders the leader bar popup (centered modal)
func (s *appState) drawLeaderBar(gtx layout.Context) layout.Dimensions {
	// Overlay background (semi-transparent)
	overlayBg := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xcc}
	overlayRect := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, overlayBg)
	overlayRect.Pop()

	// Calculate centered popup dimensions
	popupWidth := gtx.Constraints.Max.X * 2 / 3
	if popupWidth > 600 {
		popupWidth = 600
	}
	if popupWidth < 300 {
		popupWidth = 300
	}

	// Height based on number of matches, with limits
	itemHeight := 28
	headerHeight := 80
	maxItems := 12
	numItems := len(s.leaderBarMatches)
	if numItems > maxItems {
		numItems = maxItems
	}
	popupHeight := headerHeight + (numItems * itemHeight) + 20
	if popupHeight > gtx.Constraints.Max.Y*2/3 {
		popupHeight = gtx.Constraints.Max.Y * 2 / 3
	}

	offsetX := (gtx.Constraints.Max.X - popupWidth) / 2
	offsetY := (gtx.Constraints.Max.Y - popupHeight) / 3

	// Colors
	boxBg := color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff}
	boxBorder := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}
	titleColor := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}
	sequenceColor := color.NRGBA{R: 0xff, G: 0xc8, B: 0x6d, A: 0xff}
	keyColor := color.NRGBA{R: 0xa1, G: 0xc6, B: 0xff, A: 0xff}
	nameColor := color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
	selectedBg := color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0xff}
	hintColor := color.NRGBA{R: 0x6d, G: 0x7d, B: 0x9d, A: 0xff}

	// Position the popup
	offset := op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops)
	defer offset.Pop()

	// Draw border
	borderRect := clip.Rect{Max: image.Pt(popupWidth, popupHeight)}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, boxBorder)
	borderRect.Pop()

	// Draw background (slightly inset for border effect)
	bgRect := clip.Rect{
		Min: image.Pt(2, 2),
		Max: image.Pt(popupWidth-2, popupHeight-2),
	}.Push(gtx.Ops)
	paint.Fill(gtx.Ops, boxBg)
	bgRect.Pop()

	// Constrain drawing to popup area
	gtx.Constraints.Max.X = popupWidth - 4
	gtx.Constraints.Max.Y = popupHeight - 4

	inset := layout.Inset{
		Top:    unit.Dp(12),
		Right:  unit.Dp(16),
		Bottom: unit.Dp(12),
		Left:   unit.Dp(16),
	}

	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// Title
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(s.theme, "Leader Commands")
				label.Font.Typeface = "JetBrainsMono"
				label.Color = titleColor
				return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, label.Layout)
			}),

			// Current sequence
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				seqText := "Sequence: "
				if s.leaderBarSequence == "" {
					seqText += "(waiting for input)"
				} else {
					seqText += s.leaderBarSequence
				}
				label := material.Body1(s.theme, seqText)
				label.Font.Typeface = "JetBrainsMono"
				label.Color = sequenceColor
				return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, label.Layout)
			}),

			// Matches list
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if len(s.leaderBarMatches) == 0 {
					label := material.Body2(s.theme, "No matching commands")
					label.Font.Typeface = "JetBrainsMono"
					label.Color = hintColor
					return label.Layout(gtx)
				}

				list := layout.List{Axis: layout.Vertical}
				return list.Layout(gtx, len(s.leaderBarMatches), func(gtx layout.Context, index int) layout.Dimensions {
					binding := s.leaderBarMatches[index]

					// Highlight selected item
					if index == s.leaderBarIndex {
						rect := clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(24)))}.Push(gtx.Ops)
						paint.Fill(gtx.Ops, selectedBg)
						rect.Pop()
					}

					return layout.Inset{
						Top:    unit.Dp(3),
						Bottom: unit.Dp(3),
						Left:   unit.Dp(4),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							// Keys in brackets
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								keyLabel := material.Body2(s.theme, fmt.Sprintf("[%s]", binding.Keys))
								keyLabel.Font.Typeface = "JetBrainsMono"
								keyLabel.Color = keyColor
								return layout.Inset{Right: unit.Dp(12)}.Layout(gtx, keyLabel.Layout)
							}),
							// Name
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								nameLabel := material.Body2(s.theme, binding.Name)
								nameLabel.Font.Typeface = "JetBrainsMono"
								if index == s.leaderBarIndex {
									nameLabel.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
								} else {
									nameLabel.Color = nameColor
								}
								return nameLabel.Layout(gtx)
							}),
						)
					})
				})
			}),

			// Hint at bottom
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				hint := "Press keys to filter, Enter to select, Esc to cancel"
				label := material.Caption(s.theme, hint)
				label.Font.Typeface = "JetBrainsMono"
				label.Color = hintColor
				return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, label.Layout)
			}),
		)
	})
}

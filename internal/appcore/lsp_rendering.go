package appcore

import (
	"image"
	"image/color"
	"path/filepath"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"

	"github.com/javanhut/vem/internal/editor"
	"github.com/javanhut/vem/internal/lsp"
)

// LSP UI colors
var (
	diagnosticErrorColor   = color.NRGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	diagnosticWarningColor = color.NRGBA{R: 0xff, G: 0xa5, B: 0x00, A: 0xff}
	diagnosticInfoColor    = color.NRGBA{R: 0x00, G: 0xbf, B: 0xff, A: 0xff}
	diagnosticHintColor    = color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}

	completionBgColor       = color.NRGBA{R: 0x2b, G: 0x2b, B: 0x2b, A: 0xf0}
	completionSelectedColor = color.NRGBA{R: 0x00, G: 0x7a, B: 0xcc, A: 0xff}
	completionBorderColor   = color.NRGBA{R: 0x40, G: 0x40, B: 0x40, A: 0xff}

	hoverBgColor     = color.NRGBA{R: 0x1e, G: 0x1e, B: 0x1e, A: 0xf5}
	hoverBorderColor = color.NRGBA{R: 0x45, G: 0x45, B: 0x45, A: 0xff}

	referencesBgColor = color.NRGBA{R: 0x1a, G: 0x1a, B: 0x1a, A: 0xf0}

	completionDocBgColor     = color.NRGBA{R: 0x1f, G: 0x22, B: 0x2b, A: 0xf5}
	completionDocBorderColor = color.NRGBA{R: 0x3a, G: 0x3f, B: 0x4b, A: 0xff}
)

// drawDiagnosticUnderlines draws diagnostic underlines for a line.
func (s *appState) drawDiagnosticUnderlines(gtx layout.Context, buf *editor.Buffer, lineIdx int, offsetX, offsetY, lineHeight int) {
	if buf == nil || buf.FilePath() == "" {
		return
	}

	diags := s.getLSPDiagnosticsForLine(buf.FilePath(), lineIdx)
	if len(diags) == 0 {
		return
	}

	line := buf.Line(lineIdx)
	runes := []rune(line)

	for _, diag := range diags {
		// Calculate start and end columns for this line
		startCol := 0
		endCol := len(runes)

		if diag.Range.Start.Line == lineIdx {
			startCol = diag.Range.Start.Character
		}
		if diag.Range.End.Line == lineIdx {
			endCol = diag.Range.End.Character
		}

		// Clamp to line bounds
		if startCol < 0 {
			startCol = 0
		}
		if endCol > len(runes) {
			endCol = len(runes)
		}
		if startCol >= endCol {
			continue
		}

		// Calculate pixel positions
		charWidth := s.measureCharWidth(gtx)
		startX := offsetX + startCol*charWidth
		endX := offsetX + endCol*charWidth
		y := offsetY + lineHeight - 2

		// Choose color based on severity
		underlineColor := diagnosticErrorColor
		switch diag.Severity {
		case lsp.DiagnosticSeverityWarning:
			underlineColor = diagnosticWarningColor
		case lsp.DiagnosticSeverityInformation:
			underlineColor = diagnosticInfoColor
		case lsp.DiagnosticSeverityHint:
			underlineColor = diagnosticHintColor
		}

		// Draw wavy underline
		s.drawWavyUnderline(gtx, startX, y, endX-startX, underlineColor)
	}
}

// drawWavyUnderline draws a wavy underline (squiggle).
func (s *appState) drawWavyUnderline(gtx layout.Context, x, y, width int, col color.NRGBA) {
	if width <= 0 {
		return
	}

	step := 4
	amplitude := 2

	// Draw using small rectangles to simulate a wave
	up := true
	for px := x; px < x+width; px += step {
		dy := 0
		if up {
			dy = -amplitude
		}

		// Draw a small segment
		rect := clip.Rect{
			Min: image.Pt(px, y+dy),
			Max: image.Pt(px+step, y+dy+2),
		}.Push(gtx.Ops)
		paint.ColorOp{Color: col}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		rect.Pop()

		up = !up
	}
}

// drawCompletionMenu draws the completion popup.
func (s *appState) drawCompletionMenu(gtx layout.Context, offsetX, offsetY int) {
	if !s.completionActive || len(s.completionItems) == 0 {
		return
	}

	maxWidth := 400
	itemHeight := 24
	maxVisible := 10
	padding := 4

	visibleItems := len(s.completionItems)
	if visibleItems > maxVisible {
		visibleItems = maxVisible
	}

	popupHeight := visibleItems*itemHeight + padding*2
	popupWidth := maxWidth

	// Position popup below cursor
	popupX := offsetX
	popupY := offsetY + 20

	// Ensure popup stays within window bounds
	if popupX+popupWidth > gtx.Constraints.Max.X {
		popupX = gtx.Constraints.Max.X - popupWidth
	}
	if popupY+popupHeight > gtx.Constraints.Max.Y {
		popupY = offsetY - popupHeight
	}

	// Draw background
	bgRect := clip.Rect{
		Min: image.Pt(popupX, popupY),
		Max: image.Pt(popupX+popupWidth, popupY+popupHeight),
	}.Push(gtx.Ops)
	paint.ColorOp{Color: completionBgColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	bgRect.Pop()

	// Draw border
	s.drawBorder(gtx, popupX, popupY, popupWidth, popupHeight, completionBorderColor)

	// Draw items
	startIdx := 0
	if s.completionIndex >= maxVisible {
		startIdx = s.completionIndex - maxVisible + 1
	}

	for i := 0; i < visibleItems && startIdx+i < len(s.completionItems); i++ {
		itemIdx := startIdx + i
		item := s.completionItems[itemIdx]
		itemY := popupY + padding + i*itemHeight

		// Highlight selected item
		if itemIdx == s.completionIndex {
			selRect := clip.Rect{
				Min: image.Pt(popupX+2, itemY),
				Max: image.Pt(popupX+popupWidth-2, itemY+itemHeight),
			}.Push(gtx.Ops)
			paint.ColorOp{Color: completionSelectedColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			selRect.Pop()
		}

		// Draw icon
		icon := getCompletionIcon(item.Kind)
		iconLabel := material.Body2(s.theme, icon)
		iconLabel.Font.Typeface = "JetBrainsMono Nerd Font"
		iconLabel.Color = getCompletionIconColor(item.Kind)
		iconLabel.TextSize = 12

		iconOffset := op.Offset(image.Pt(popupX+6, itemY+4)).Push(gtx.Ops)
		iconLabel.Layout(gtx)
		iconOffset.Pop()

		// Draw label
		label := material.Body2(s.theme, item.Label)
		label.Font.Typeface = "JetBrainsMono Nerd Font"
		label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		label.TextSize = 12

		labelOffset := op.Offset(image.Pt(popupX+26, itemY+4)).Push(gtx.Ops)
		label.Layout(gtx)
		labelOffset.Pop()

		// Draw detail (type info)
		if item.Detail != "" {
			detail := material.Caption(s.theme, truncateString(item.Detail, 30))
			detail.Font.Typeface = "JetBrainsMono Nerd Font"
			detail.Color = color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}
			detail.TextSize = 10

			detailOffset := op.Offset(image.Pt(popupX+200, itemY+6)).Push(gtx.Ops)
			detail.Layout(gtx)
			detailOffset.Pop()
		}
	}

	selected := s.completionItems[s.completionIndex]
	docText := completionDocText(selected)
	if docText == "" {
		return
	}

	docLines := wrapText(docText, 50)
	if len(docLines) > 10 {
		docLines = docLines[:10]
		docLines[len(docLines)-1] = truncateString(docLines[len(docLines)-1], 47)
	}

	lineHeight := 16
	docPadding := 6
	docHeight := len(docLines)*lineHeight + docPadding*2
	docWidth := 320

	docX := popupX + popupWidth + 8
	docY := popupY

	if docX+docWidth > gtx.Constraints.Max.X {
		docX = popupX
		docY = popupY + popupHeight + 8
	}
	if docY+docHeight > gtx.Constraints.Max.Y {
		docY = popupY - docHeight - 8
	}

	docRect := clip.Rect{
		Min: image.Pt(docX, docY),
		Max: image.Pt(docX+docWidth, docY+docHeight),
	}.Push(gtx.Ops)
	paint.ColorOp{Color: completionDocBgColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	docRect.Pop()

	s.drawBorder(gtx, docX, docY, docWidth, docHeight, completionDocBorderColor)

	for i, line := range docLines {
		label := material.Caption(s.theme, line)
		label.Font.Typeface = "JetBrainsMono"
		label.Color = color.NRGBA{R: 0xe6, G: 0xe6, B: 0xe6, A: 0xff}
		label.TextSize = 11

		offset := op.Offset(image.Pt(docX+docPadding, docY+docPadding+i*lineHeight)).Push(gtx.Ops)
		label.Layout(gtx)
		offset.Pop()
	}
}

// drawHoverTooltip draws the hover information tooltip.
func (s *appState) drawHoverTooltip(gtx layout.Context, cursorX, cursorY int) {
	if !s.hoverActive || s.hoverContent == "" {
		return
	}

	maxWidth := 500
	padding := 8

	// Wrap long lines
	lines := wrapText(s.hoverContent, 80)
	lineHeight := 16
	popupHeight := len(lines)*lineHeight + padding*2

	// Calculate width based on longest line
	popupWidth := 0
	charWidth := s.measureCharWidth(gtx)
	for _, line := range lines {
		lineWidth := len(line)*charWidth + padding*2
		if lineWidth > popupWidth {
			popupWidth = lineWidth
		}
	}
	if popupWidth > maxWidth {
		popupWidth = maxWidth
	}
	if popupWidth < 100 {
		popupWidth = 100
	}

	// Position tooltip above cursor
	popupX := cursorX
	popupY := cursorY - popupHeight - 5

	// Ensure tooltip stays within window bounds
	if popupX+popupWidth > gtx.Constraints.Max.X {
		popupX = gtx.Constraints.Max.X - popupWidth - 10
	}
	if popupY < 0 {
		popupY = cursorY + 20 // Show below cursor instead
	}

	// Draw background
	bgRect := clip.Rect{
		Min: image.Pt(popupX, popupY),
		Max: image.Pt(popupX+popupWidth, popupY+popupHeight),
	}.Push(gtx.Ops)
	paint.ColorOp{Color: hoverBgColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	bgRect.Pop()

	// Draw border
	s.drawBorder(gtx, popupX, popupY, popupWidth, popupHeight, hoverBorderColor)

	// Draw content
	for i, line := range lines {
		textY := popupY + padding + i*lineHeight

		label := material.Body2(s.theme, line)
		label.Font.Typeface = "JetBrainsMono Nerd Font"
		label.Color = color.NRGBA{R: 0xe0, G: 0xe0, B: 0xe0, A: 0xff}
		label.TextSize = 12

		offset := op.Offset(image.Pt(popupX+padding, textY)).Push(gtx.Ops)
		label.Layout(gtx)
		offset.Pop()
	}
}

// drawReferencesList draws the references list at the bottom of the screen.
func (s *appState) drawReferencesList(gtx layout.Context) {
	if !s.referencesActive || len(s.referencesItems) == 0 {
		return
	}

	itemHeight := 20
	maxVisible := 8
	padding := 4

	visibleItems := len(s.referencesItems)
	if visibleItems > maxVisible {
		visibleItems = maxVisible
	}

	listHeight := visibleItems*itemHeight + padding*2

	// Draw at bottom of screen
	listY := gtx.Constraints.Max.Y - listHeight
	listWidth := gtx.Constraints.Max.X

	// Draw background
	bgRect := clip.Rect{
		Min: image.Pt(0, listY),
		Max: image.Pt(listWidth, gtx.Constraints.Max.Y),
	}.Push(gtx.Ops)
	paint.ColorOp{Color: referencesBgColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	bgRect.Pop()

	// Draw separator
	sepRect := clip.Rect{
		Min: image.Pt(0, listY),
		Max: image.Pt(listWidth, listY+1),
	}.Push(gtx.Ops)
	paint.ColorOp{Color: paneSeparator}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	sepRect.Pop()

	// Calculate scroll position
	startIdx := 0
	if s.referencesIndex >= maxVisible {
		startIdx = s.referencesIndex - maxVisible + 1
	}

	// Draw items
	for i := 0; i < visibleItems && startIdx+i < len(s.referencesItems); i++ {
		itemIdx := startIdx + i
		loc := s.referencesItems[itemIdx]
		itemY := listY + padding + i*itemHeight

		// Highlight selected
		if itemIdx == s.referencesIndex {
			selRect := clip.Rect{
				Min: image.Pt(2, itemY),
				Max: image.Pt(listWidth-2, itemY+itemHeight),
			}.Push(gtx.Ops)
			paint.ColorOp{Color: completionSelectedColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			selRect.Pop()
		}

		// Format location
		filePath := lsp.URIToFilePath(loc.URI)
		shortPath := filepath.Base(filePath)
		line := loc.Range.Start.Line + 1
		col := loc.Range.Start.Character + 1
		text := strings.Repeat(" ", 3-len(string(rune('0'+itemIdx%10)))) +
			string(rune('0'+itemIdx%10)) + ". " + shortPath + ":" +
			strings.Repeat(" ", 4-len(string(rune('0'+line)))) +
			formatInt(line) + ":" + formatInt(col)

		label := material.Body2(s.theme, text)
		label.Font.Typeface = "JetBrainsMono Nerd Font"
		label.Font.Weight = font.Normal
		label.Color = color.NRGBA{R: 0xe0, G: 0xe0, B: 0xe0, A: 0xff}
		label.TextSize = 12

		offset := op.Offset(image.Pt(8, itemY+2)).Push(gtx.Ops)
		label.Layout(gtx)
		offset.Pop()
	}

	// Draw hint
	hint := "[j/k: navigate] [Enter: open] [Esc: close]"
	hintLabel := material.Caption(s.theme, hint)
	hintLabel.Font.Typeface = "JetBrainsMono Nerd Font"
	hintLabel.Color = color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xff}
	hintLabel.TextSize = 10

	hintOffset := op.Offset(image.Pt(listWidth-300, listY+padding+2)).Push(gtx.Ops)
	hintLabel.Layout(gtx)
	hintOffset.Pop()
}

// drawCodeActionsMenu draws the code actions menu.
func (s *appState) drawCodeActionsMenu(gtx layout.Context, cursorX, cursorY int) {
	if !s.codeActionsActive || len(s.codeActionItems) == 0 {
		return
	}

	itemHeight := 24
	maxWidth := 350
	padding := 4

	popupHeight := len(s.codeActionItems)*itemHeight + padding*2

	// Position menu near cursor
	popupX := cursorX
	popupY := cursorY + 20

	// Ensure menu stays within window bounds
	if popupX+maxWidth > gtx.Constraints.Max.X {
		popupX = gtx.Constraints.Max.X - maxWidth - 10
	}
	if popupY+popupHeight > gtx.Constraints.Max.Y {
		popupY = cursorY - popupHeight - 5
	}

	// Draw background
	bgRect := clip.Rect{
		Min: image.Pt(popupX, popupY),
		Max: image.Pt(popupX+maxWidth, popupY+popupHeight),
	}.Push(gtx.Ops)
	paint.ColorOp{Color: completionBgColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	bgRect.Pop()

	// Draw border
	s.drawBorder(gtx, popupX, popupY, maxWidth, popupHeight, completionBorderColor)

	// Draw items
	for i, action := range s.codeActionItems {
		itemY := popupY + padding + i*itemHeight

		// Highlight selected
		if i == s.codeActionIndex {
			selRect := clip.Rect{
				Min: image.Pt(popupX+2, itemY),
				Max: image.Pt(popupX+maxWidth-2, itemY+itemHeight),
			}.Push(gtx.Ops)
			paint.ColorOp{Color: completionSelectedColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			selRect.Pop()
		}

		// Draw title
		label := material.Body2(s.theme, truncateString(action.Title, 45))
		label.Font.Typeface = "JetBrainsMono Nerd Font"
		label.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
		label.TextSize = 12

		offset := op.Offset(image.Pt(popupX+8, itemY+4)).Push(gtx.Ops)
		label.Layout(gtx)
		offset.Pop()
	}
}

// drawBorder draws a 1px border around a rectangle.
func (s *appState) drawBorder(gtx layout.Context, x, y, width, height int, col color.NRGBA) {
	// Top
	rect := clip.Rect{Min: image.Pt(x, y), Max: image.Pt(x+width, y+1)}.Push(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	rect.Pop()

	// Bottom
	rect = clip.Rect{Min: image.Pt(x, y+height-1), Max: image.Pt(x+width, y+height)}.Push(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	rect.Pop()

	// Left
	rect = clip.Rect{Min: image.Pt(x, y), Max: image.Pt(x+1, y+height)}.Push(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	rect.Pop()

	// Right
	rect = clip.Rect{Min: image.Pt(x+width-1, y), Max: image.Pt(x+width, y+height)}.Push(gtx.Ops)
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	rect.Pop()
}

// measureCharWidth returns the approximate width of a character.
func (s *appState) measureCharWidth(gtx layout.Context) int {
	// Approximate character width for monospace font
	return 8
}

// getCompletionIcon returns an icon character for a completion item kind.
func getCompletionIcon(kind lsp.CompletionItemKind) string {
	switch kind {
	case lsp.CompletionItemKindText:
		return "T"
	case lsp.CompletionItemKindMethod:
		return "m"
	case lsp.CompletionItemKindFunction:
		return "f"
	case lsp.CompletionItemKindConstructor:
		return "C"
	case lsp.CompletionItemKindField:
		return "F"
	case lsp.CompletionItemKindVariable:
		return "v"
	case lsp.CompletionItemKindClass:
		return "c"
	case lsp.CompletionItemKindInterface:
		return "i"
	case lsp.CompletionItemKindModule:
		return "M"
	case lsp.CompletionItemKindProperty:
		return "p"
	case lsp.CompletionItemKindUnit:
		return "u"
	case lsp.CompletionItemKindValue:
		return "V"
	case lsp.CompletionItemKindEnum:
		return "E"
	case lsp.CompletionItemKindKeyword:
		return "k"
	case lsp.CompletionItemKindSnippet:
		return "s"
	case lsp.CompletionItemKindColor:
		return "#"
	case lsp.CompletionItemKindFile:
		return "F"
	case lsp.CompletionItemKindReference:
		return "r"
	case lsp.CompletionItemKindFolder:
		return "D"
	case lsp.CompletionItemKindEnumMember:
		return "e"
	case lsp.CompletionItemKindConstant:
		return "C"
	case lsp.CompletionItemKindStruct:
		return "S"
	case lsp.CompletionItemKindEvent:
		return "E"
	case lsp.CompletionItemKindOperator:
		return "o"
	case lsp.CompletionItemKindTypeParameter:
		return "T"
	default:
		return " "
	}
}

// getCompletionIconColor returns a color for a completion item kind.
func getCompletionIconColor(kind lsp.CompletionItemKind) color.NRGBA {
	switch kind {
	case lsp.CompletionItemKindMethod, lsp.CompletionItemKindFunction:
		return color.NRGBA{R: 0xdc, G: 0xdc, B: 0xaa, A: 0xff} // Yellow
	case lsp.CompletionItemKindVariable, lsp.CompletionItemKindField:
		return color.NRGBA{R: 0x9c, G: 0xdc, B: 0xfe, A: 0xff} // Light blue
	case lsp.CompletionItemKindClass, lsp.CompletionItemKindStruct, lsp.CompletionItemKindInterface:
		return color.NRGBA{R: 0x4e, G: 0xc9, B: 0xb0, A: 0xff} // Teal
	case lsp.CompletionItemKindModule, lsp.CompletionItemKindFolder:
		return color.NRGBA{R: 0xce, G: 0x91, B: 0x78, A: 0xff} // Orange
	case lsp.CompletionItemKindKeyword:
		return color.NRGBA{R: 0xc5, G: 0x86, B: 0xc0, A: 0xff} // Purple
	case lsp.CompletionItemKindSnippet:
		return color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff} // Gray
	case lsp.CompletionItemKindConstant:
		return color.NRGBA{R: 0x56, G: 0x9c, B: 0xd6, A: 0xff} // Blue
	default:
		return color.NRGBA{R: 0xd4, G: 0xd4, B: 0xd4, A: 0xff} // Light gray
	}
}

// truncateString truncates a string to a maximum length.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func completionDocText(item lsp.CompletionItem) string {
	var parts []string

	if item.Detail != "" {
		parts = append(parts, item.Detail)
	}

	doc := completionDocString(item.Documentation)
	if doc != "" {
		parts = append(parts, doc)
	}

	return strings.TrimSpace(strings.Join(parts, "\n"))
}

func completionDocString(doc interface{}) string {
	switch v := doc.(type) {
	case string:
		return v
	case lsp.MarkupContent:
		return v.Value
	case map[string]interface{}:
		if value, ok := v["value"].(string); ok {
			return value
		}
		if value, ok := v["documentation"].(string); ok {
			return value
		}
	}
	return ""
}

// wrapText wraps text to a maximum line length.
func wrapText(text string, maxLen int) []string {
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		if len(line) <= maxLen {
			lines = append(lines, line)
		} else {
			// Simple word wrap
			words := strings.Fields(line)
			current := ""
			for _, word := range words {
				if len(current)+len(word)+1 <= maxLen {
					if current != "" {
						current += " "
					}
					current += word
				} else {
					if current != "" {
						lines = append(lines, current)
					}
					current = word
				}
			}
			if current != "" {
				lines = append(lines, current)
			}
		}
	}
	return lines
}

// formatInt formats an integer as a string.
func formatInt(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return formatInt(n/10) + string(rune('0'+n%10))
}

package appcore

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/javanhut/vem/internal/lsp"
	"github.com/javanhut/vem/internal/syntax"
)

// Tab constants
const (
	tabGeneral = iota
	tabEditor
	tabLSP
)

// settingsModalState holds the persistent state for the settings modal.
type settingsModalState struct {
	active      bool
	selectedTab int
	focusedItem int // Currently focused item within the tab (for keyboard navigation)

	// Tab buttons
	tabGeneral widget.Clickable
	tabEditor  widget.Clickable
	tabLSP     widget.Clickable

	// General settings widgets
	fontSizeMinus widget.Clickable
	fontSizePlus  widget.Clickable
	fontSizeValue int

	// Editor settings widgets
	autoIndent        widget.Bool
	autoPairs         widget.Bool
	softWrap          widget.Bool
	formatOnSave      widget.Bool
	useSpaces         widget.Bool
	tabWidthMinus     widget.Clickable
	tabWidthPlus      widget.Clickable
	tabWidthValue     int
	scrollOffsetMinus widget.Clickable
	scrollOffsetPlus  widget.Clickable
	scrollOffsetValue int

	// LSP settings widgets
	lspEnabled    widget.Bool
	lspAutoDetect widget.Bool

	// LSP status widgets
	lspServerList  widget.List
	lspServers     []lsp.LSPServerStatus
	lspRestartBtns []widget.Clickable // One per server
	lspStopBtns    []widget.Clickable // One per server
	lspStopAllBtn  widget.Clickable   // Stop all servers

	// Theme selector widgets
	themeList     widget.List
	themeSelected int      // Currently selected theme index
	themes        []string // Available theme names
	themePreview  bool     // Whether preview mode is active
	themeOriginal string   // Original theme before preview started

	// Content scroll
	contentList widget.List

	// Action buttons
	saveButton   widget.Clickable
	cancelButton widget.Clickable
}

// initSettingsModal initializes the modal widget values from current settings.
func (s *appState) initSettingsModal() {
	// UI
	s.settingsModal.fontSizeValue = s.settings.UI.FontSize

	// Editor
	s.settingsModal.autoIndent.Value = s.settings.Editor.AutoIndent
	s.settingsModal.autoPairs.Value = s.settings.Editor.AutoPairs
	s.settingsModal.softWrap.Value = s.settings.Editor.SoftWrap
	s.settingsModal.formatOnSave.Value = s.settings.Editor.FormatOnSave
	s.settingsModal.useSpaces.Value = s.settings.Editor.UseSpaces
	s.settingsModal.tabWidthValue = s.settings.Editor.TabWidth
	s.settingsModal.scrollOffsetValue = s.settings.Editor.ScrollOffset

	// LSP
	s.settingsModal.lspEnabled.Value = s.settings.LSP.Enabled
	s.settingsModal.lspAutoDetect.Value = s.settings.LSP.AutoDetect
	s.settingsModal.lspServerList.List.Axis = layout.Vertical
	s.settingsModal.lspServers = nil // Will be populated when tab is viewed
	s.settingsModal.lspRestartBtns = nil
	s.settingsModal.lspStopBtns = nil

	// Theme
	s.settingsModal.themes = make([]string, len(syntax.PresetThemes))
	copy(s.settingsModal.themes, syntax.PresetThemes)
	s.settingsModal.themeList.List.Axis = layout.Vertical
	s.settingsModal.themePreview = false
	s.settingsModal.themeOriginal = s.settings.UI.Theme

	// Find current theme index
	currentTheme := s.settings.UI.Theme
	// Map "dark"/"light" to actual Chroma theme names for matching
	if currentTheme == "dark" {
		currentTheme = "monokai"
	} else if currentTheme == "light" {
		currentTheme = "solarized-light"
	}
	s.settingsModal.themeSelected = 0
	for i, theme := range s.settingsModal.themes {
		if theme == currentTheme {
			s.settingsModal.themeSelected = i
			break
		}
	}

	// Content scroll list
	s.settingsModal.contentList.List.Axis = layout.Vertical

	// Reset to first tab and first item
	s.settingsModal.selectedTab = tabGeneral
	s.settingsModal.focusedItem = 0
}

// drawSettingsModal renders the settings modal overlay.
func (s *appState) drawSettingsModal(gtx layout.Context) layout.Dimensions {
	// Colors
	overlayBg := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xcc}
	boxBg := color.NRGBA{R: 0x1a, G: 0x1f, B: 0x2e, A: 0xff}
	boxBorder := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}
	titleColor := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}

	// Overlay background (semi-transparent)
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
	popupHeight := 450
	if popupHeight > gtx.Constraints.Max.Y*2/3 {
		popupHeight = gtx.Constraints.Max.Y * 2 / 3
	}

	offsetX := (gtx.Constraints.Max.X - popupWidth) / 2
	offsetY := (gtx.Constraints.Max.Y - popupHeight) / 3

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
				label := material.H6(s.theme, "Settings")
				label.Color = titleColor
				return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, label.Layout)
			}),

			// Tab bar
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.drawSettingsTabs(gtx)
			}),

			// Content area (flexible, clipped and scrollable)
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					// Clip content to prevent overflow
					defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
					return s.drawTabContent(gtx)
				})
			}),

			// Button row at bottom
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.drawSettingsButtons(gtx)
			}),
		)
	})
}

// drawSettingsTabs renders the tab bar for settings navigation.
func (s *appState) drawSettingsTabs(gtx layout.Context) layout.Dimensions {
	tabNames := []string{"General", "Editor", "LSP"}
	tabButtons := []*widget.Clickable{
		&s.settingsModal.tabGeneral,
		&s.settingsModal.tabEditor,
		&s.settingsModal.tabLSP,
	}

	// Check for tab clicks
	for i, btn := range tabButtons {
		if btn.Clicked(gtx) {
			s.settingsModal.selectedTab = i
			s.settingsModal.focusedItem = 0
			s.settingsModal.contentList.List.Position.First = 0
			s.settingsModal.contentList.List.Position.Offset = 0
		}
	}

	// Colors
	selectedBg := color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0xff}
	unselectedBg := color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}
	selectedText := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	unselectedText := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}

	var children []layout.FlexChild
	for i, name := range tabNames {
		idx := i
		tabName := name
		btn := tabButtons[i]

		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			isSelected := s.settingsModal.selectedTab == idx

			return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					// Draw tab background
					bgColor := unselectedBg
					textColor := unselectedText
					if isSelected {
						bgColor = selectedBg
						textColor = selectedText
					}

					return layout.Background{}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							rr := gtx.Dp(unit.Dp(4))
							rect := clip.RRect{
								Rect: image.Rectangle{Max: gtx.Constraints.Min},
								NE:   rr, NW: rr, SE: rr, SW: rr,
							}.Push(gtx.Ops)
							paint.Fill(gtx.Ops, bgColor)
							rect.Pop()
							return layout.Dimensions{Size: gtx.Constraints.Min}
						},
						func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top: unit.Dp(6), Bottom: unit.Dp(6),
								Left: unit.Dp(12), Right: unit.Dp(12),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body2(s.theme, tabName)
								lbl.Color = textColor
								return lbl.Layout(gtx)
							})
						},
					)
				})
			})
		}))
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, children...)
}

// drawTabContent renders the content for the selected tab with scrolling.
func (s *appState) drawTabContent(gtx layout.Context) layout.Dimensions {
	itemCount := s.getTabContentItemCount()

	return material.List(s.theme, &s.settingsModal.contentList).Layout(gtx,
		itemCount,
		func(gtx layout.Context, index int) layout.Dimensions {
			switch s.settingsModal.selectedTab {
			case tabGeneral:
				return s.drawGeneralTabItem(gtx, index)
			case tabEditor:
				return s.drawEditorTabItem(gtx, index)
			case tabLSP:
				return s.drawLSPTabItem(gtx, index)
			default:
				return s.drawGeneralTabItem(gtx, index)
			}
		},
	)
}

// getTabContentItemCount returns the number of content items for the current tab.
func (s *appState) getTabContentItemCount() int {
	switch s.settingsModal.selectedTab {
	case tabGeneral:
		return 3 // Section header, Font Size, Theme
	case tabEditor:
		return 10 // Section + 4 toggles, Section + 2 items, Section + 1 item
	case tabLSP:
		return 10 // Section + toggle + hint + toggle + hint + section + servers + section + diagnostics + button
	default:
		return 1
	}
}

// drawGeneralTabItem renders a single item in the General tab.
func (s *appState) drawGeneralTabItem(gtx layout.Context, index int) layout.Dimensions {
	focused := s.settingsModal.focusedItem

	switch index {
	case 0:
		return s.drawSectionHeader(gtx, "Appearance")
	case 1:
		return s.drawNumberStepper(gtx, "Font Size",
			&s.settingsModal.fontSizeValue,
			&s.settingsModal.fontSizeMinus,
			&s.settingsModal.fontSizePlus,
			8, 32, focused == 0)
	case 2:
		return s.drawThemeSelector(gtx, focused == 1)
	default:
		return layout.Dimensions{}
	}
}

// drawEditorTabItem renders a single item in the Editor tab.
func (s *appState) drawEditorTabItem(gtx layout.Context, index int) layout.Dimensions {
	focused := s.settingsModal.focusedItem

	switch index {
	case 0:
		return s.drawSectionHeader(gtx, "Behavior")
	case 1:
		return s.drawBoolSetting(gtx, &s.settingsModal.autoIndent, "Auto Indent", focused == 0)
	case 2:
		return s.drawBoolSetting(gtx, &s.settingsModal.autoPairs, "Auto Pairs", focused == 1)
	case 3:
		return s.drawBoolSetting(gtx, &s.settingsModal.softWrap, "Soft Wrap", focused == 2)
	case 4:
		return s.drawBoolSetting(gtx, &s.settingsModal.formatOnSave, "Format on Save", focused == 3)
	case 5:
		return s.drawSectionHeader(gtx, "Indentation")
	case 6:
		return s.drawBoolSetting(gtx, &s.settingsModal.useSpaces, "Use Spaces", focused == 4)
	case 7:
		return s.drawNumberStepper(gtx, "Tab Width",
			&s.settingsModal.tabWidthValue,
			&s.settingsModal.tabWidthMinus,
			&s.settingsModal.tabWidthPlus,
			2, 8, focused == 5)
	case 8:
		return s.drawSectionHeader(gtx, "Scrolling")
	case 9:
		return s.drawNumberStepper(gtx, "Scroll Offset",
			&s.settingsModal.scrollOffsetValue,
			&s.settingsModal.scrollOffsetMinus,
			&s.settingsModal.scrollOffsetPlus,
			0, 20, focused == 6)
	default:
		return layout.Dimensions{}
	}
}

// drawLSPTabItem renders a single item in the LSP tab.
func (s *appState) drawLSPTabItem(gtx layout.Context, index int) layout.Dimensions {
	focused := s.settingsModal.focusedItem
	hintColor := color.NRGBA{R: 0x6d, G: 0x7d, B: 0x9d, A: 0xff}

	switch index {
	case 0:
		return s.drawSectionHeader(gtx, "Language Server Protocol")
	case 1:
		return s.drawBoolSetting(gtx, &s.settingsModal.lspEnabled, "Enable LSP", focused == 0)
	case 2:
		return layout.Inset{Left: unit.Dp(4), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Caption(s.theme, "Connect to language servers for completions, diagnostics, etc.")
			lbl.Color = hintColor
			return lbl.Layout(gtx)
		})
	case 3:
		return s.drawBoolSetting(gtx, &s.settingsModal.lspAutoDetect, "Auto Detect", focused == 1)
	case 4:
		return layout.Inset{Left: unit.Dp(4), Bottom: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Caption(s.theme, "Automatically start language server based on file type")
			lbl.Color = hintColor
			return lbl.Layout(gtx)
		})
	case 5:
		return s.drawSectionHeader(gtx, "Active Servers")
	case 6:
		return s.drawLSPServerStatus(gtx)
	case 7:
		return s.drawSectionHeader(gtx, "Diagnostics Summary")
	case 8:
		return s.drawDiagnosticsSummary(gtx)
	case 9:
		return s.drawStopAllButton(gtx)
	default:
		return layout.Dimensions{}
	}
}

// drawStopAllButton renders the Stop All Servers button.
func (s *appState) drawStopAllButton(gtx layout.Context) layout.Dimensions {
	// Handle click
	if s.settingsModal.lspStopAllBtn.Clicked(gtx) {
		if s.lspManager != nil {
			s.lspManager.StopAll()
			s.status = "All language servers stopped"
		}
	}

	return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return s.settingsModal.lspStopAllBtn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Background{}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					rr := gtx.Dp(unit.Dp(4))
					rect := clip.RRect{
						Rect: image.Rectangle{Max: gtx.Constraints.Min},
						NE:   rr, NW: rr, SE: rr, SW: rr,
					}.Push(gtx.Ops)
					paint.Fill(gtx.Ops, color.NRGBA{R: 0x8b, G: 0x00, B: 0x00, A: 0xff}) // Dark red
					rect.Pop()
					return layout.Dimensions{Size: gtx.Constraints.Min}
				},
				func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top: unit.Dp(6), Bottom: unit.Dp(6),
						Left: unit.Dp(12), Right: unit.Dp(12),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(s.theme, "Stop All Servers")
						lbl.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
						return lbl.Layout(gtx)
					})
				},
			)
		})
	})
}

// drawBoolSetting renders a consistent row layout for boolean switches with optional focus highlight.
func (s *appState) drawBoolSetting(gtx layout.Context, boolWidget *widget.Bool, label string, isFocused bool) layout.Dimensions {
	labelColor := color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
	focusBgColor := color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0x44}

	return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				if isFocused {
					rr := gtx.Dp(unit.Dp(4))
					rect := clip.RRect{
						Rect: image.Rectangle{Max: gtx.Constraints.Min},
						NE:   rr, NW: rr, SE: rr, SW: rr,
					}.Push(gtx.Ops)
					paint.Fill(gtx.Ops, focusBgColor)
					rect.Pop()
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4), Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body1(s.theme, label)
							lbl.Color = labelColor
							if isFocused {
								lbl.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
							}
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Switch(s.theme, boolWidget, "").Layout(gtx)
						}),
					)
				})
			},
		)
	})
}

// drawSectionHeader renders a section header with styled text.
func (s *appState) drawSectionHeader(gtx layout.Context, title string) layout.Dimensions {
	return layout.Inset{Top: unit.Dp(12), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Body2(s.theme, title)
		lbl.Color = color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}
		lbl.Font.Weight = font.Bold
		return lbl.Layout(gtx)
	})
}

// drawNumberStepper renders a label with +/- buttons and value display.
func (s *appState) drawNumberStepper(gtx layout.Context, label string, value *int, minus, plus *widget.Clickable, minVal, maxVal int, isFocused bool) layout.Dimensions {
	// Handle button clicks
	if minus.Clicked(gtx) && *value > minVal {
		*value--
	}
	if plus.Clicked(gtx) && *value < maxVal {
		*value++
	}

	// Colors
	labelColor := color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
	btnColor := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}
	valueColor := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	focusBgColor := color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0x44}

	if isFocused {
		labelColor = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	}

	return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Background{}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				if isFocused {
					rr := gtx.Dp(unit.Dp(4))
					rect := clip.RRect{
						Rect: image.Rectangle{Max: gtx.Constraints.Min},
						NE:   rr, NW: rr, SE: rr, SW: rr,
					}.Push(gtx.Ops)
					paint.Fill(gtx.Ops, focusBgColor)
					rect.Pop()
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			},
			func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(4), Right: unit.Dp(4), Top: unit.Dp(2), Bottom: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						// Label
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body1(s.theme, label)
							lbl.Color = labelColor
							return lbl.Layout(gtx)
						}),

						// Minus button
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return minus.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Body1(s.theme, "-")
									lbl.Color = btnColor
									return lbl.Layout(gtx)
								})
							})
						}),

						// Value
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body1(s.theme, fmt.Sprintf("%d", *value))
								lbl.Color = valueColor
								return lbl.Layout(gtx)
							})
						}),

						// Plus button
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return plus.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Body1(s.theme, "+")
									lbl.Color = btnColor
									return lbl.Layout(gtx)
								})
							})
						}),
					)
				})
			},
		)
	})
}

// drawThemeSelector renders the theme selection list with live preview.
func (s *appState) drawThemeSelector(gtx layout.Context, isFocused bool) layout.Dimensions {
	labelColor := color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
	selectedBg := color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0x66}
	focusBg := color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0x44}
	descColor := color.NRGBA{R: 0x6d, G: 0x7d, B: 0x9d, A: 0xff}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Label row
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body1(s.theme, "Theme")
				lbl.Color = labelColor
				if isFocused {
					lbl.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
				}
				return lbl.Layout(gtx)
			})
		}),

		// Theme list with constrained height
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Focus highlight for the entire list area
			if isFocused {
				rr := gtx.Dp(unit.Dp(4))
				rect := clip.RRect{
					Rect: image.Rectangle{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(unit.Dp(150)))},
					NE:   rr, NW: rr, SE: rr, SW: rr,
				}.Push(gtx.Ops)
				paint.Fill(gtx.Ops, focusBg)
				rect.Pop()
			}

			gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(150))
			gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(150))

			return material.List(s.theme, &s.settingsModal.themeList).Layout(gtx,
				len(s.settingsModal.themes),
				func(gtx layout.Context, index int) layout.Dimensions {
					themeName := s.settingsModal.themes[index]
					isSelected := index == s.settingsModal.themeSelected

					return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(2), Left: unit.Dp(4), Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Background{}.Layout(gtx,
							func(gtx layout.Context) layout.Dimensions {
								if isSelected {
									rr := gtx.Dp(unit.Dp(3))
									rect := clip.RRect{
										Rect: image.Rectangle{Max: gtx.Constraints.Min},
										NE:   rr, NW: rr, SE: rr, SW: rr,
									}.Push(gtx.Ops)
									paint.Fill(gtx.Ops, selectedBg)
									rect.Pop()
								}
								return layout.Dimensions{Size: gtx.Constraints.Min}
							},
							func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										// Theme name row
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											displayName := themeName
											if isSelected {
												displayName = themeName + " (selected)"
											}
											lbl := material.Body2(s.theme, displayName)
											if isSelected {
												lbl.Color = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
											} else {
												lbl.Color = labelColor
											}
											return lbl.Layout(gtx)
										}),
										// Description row
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											desc := syntax.GetThemeDescription(themeName)
											lbl := material.Caption(s.theme, desc)
											lbl.Color = descColor
											return lbl.Layout(gtx)
										}),
									)
								})
							},
						)
					})
				},
			)
		}),
	)
}

// drawLSPServerStatus renders the active servers status panel.
func (s *appState) drawLSPServerStatus(gtx layout.Context) layout.Dimensions {
	// Refresh server status
	if s.lspManager != nil {
		s.settingsModal.lspServers = s.lspManager.GetDetailedStatus()
	}

	labelColor := color.NRGBA{R: 0xdf, G: 0xe7, B: 0xff, A: 0xff}
	hintColor := color.NRGBA{R: 0x6d, G: 0x7d, B: 0x9d, A: 0xff}
	runningColor := color.NRGBA{R: 0x50, G: 0xfa, B: 0x7b, A: 0xff} // Green
	stoppedColor := color.NRGBA{R: 0xff, G: 0x55, B: 0x55, A: 0xff} // Red
	accentColor := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}  // Blue

	if len(s.settingsModal.lspServers) == 0 {
		return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Caption(s.theme, "No active language servers")
			lbl.Color = hintColor
			return lbl.Layout(gtx)
		})
	}

	// Resize button slices to match server count
	if len(s.settingsModal.lspRestartBtns) != len(s.settingsModal.lspServers) {
		s.settingsModal.lspRestartBtns = make([]widget.Clickable, len(s.settingsModal.lspServers))
		s.settingsModal.lspStopBtns = make([]widget.Clickable, len(s.settingsModal.lspServers))
	}

	// Constrain height
	gtx.Constraints.Max.Y = gtx.Dp(unit.Dp(150))

	return material.List(s.theme, &s.settingsModal.lspServerList).Layout(gtx,
		len(s.settingsModal.lspServers),
		func(gtx layout.Context, index int) layout.Dimensions {
			server := s.settingsModal.lspServers[index]

			// Handle button clicks
			if s.settingsModal.lspRestartBtns[index].Clicked(gtx) {
				s.lspManager.StopServer(server.WorkspaceRoot)
				s.status = fmt.Sprintf("Restarting %s...", server.Name)
			}
			if s.settingsModal.lspStopBtns[index].Clicked(gtx) {
				s.lspManager.StopServer(server.WorkspaceRoot)
				s.status = fmt.Sprintf("Stopped %s", server.Name)
			}

			// Status indicator
			statusText := "Stopped"
			statusColor := stoppedColor
			if server.Running {
				statusText = "Running"
				statusColor = runningColor
			}

			return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// Server name, status, and buttons
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body2(s.theme, server.Name)
								lbl.Color = labelColor
								return lbl.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									lbl := material.Caption(s.theme, statusText)
									lbl.Color = statusColor
									return lbl.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									diagText := fmt.Sprintf("%d diag", server.DiagCount)
									lbl := material.Caption(s.theme, diagText)
									lbl.Color = hintColor
									return lbl.Layout(gtx)
								})
							}),
							// Restart button
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Left: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return s.settingsModal.lspRestartBtns[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										lbl := material.Caption(s.theme, "[Restart]")
										lbl.Color = accentColor
										return lbl.Layout(gtx)
									})
								})
							}),
							// Stop button
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Left: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return s.settingsModal.lspStopBtns[index].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										lbl := material.Caption(s.theme, "[Stop]")
										lbl.Color = stoppedColor
										return lbl.Layout(gtx)
									})
								})
							}),
						)
					}),
					// Workspace path
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Caption(s.theme, server.WorkspaceRoot)
							lbl.Color = hintColor
							return lbl.Layout(gtx)
						})
					}),
				)
			})
		},
	)
}

// drawDiagnosticsSummary renders an overall diagnostics count summary.
func (s *appState) drawDiagnosticsSummary(gtx layout.Context) layout.Dimensions {
	errors, warnings, info := 0, 0, 0
	if s.lspManager != nil {
		errors, warnings, info = s.lspManager.GetDiagnosticsSummary()
	}

	errorColor := color.NRGBA{R: 0xff, G: 0x55, B: 0x55, A: 0xff}
	warnColor := color.NRGBA{R: 0xff, G: 0xaa, B: 0x00, A: 0xff}
	infoColor := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}
	hintColor := color.NRGBA{R: 0x6d, G: 0x7d, B: 0x9d, A: 0xff}

	total := errors + warnings + info
	if total == 0 {
		return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Caption(s.theme, "No diagnostics")
			lbl.Color = hintColor
			return lbl.Layout(gtx)
		})
	}

	return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(s.theme, fmt.Sprintf("%d errors", errors))
				lbl.Color = errorColor
				return lbl.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Caption(s.theme, fmt.Sprintf("%d warnings", warnings))
					lbl.Color = warnColor
					return lbl.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Caption(s.theme, fmt.Sprintf("%d info", info))
					lbl.Color = infoColor
					return lbl.Layout(gtx)
				})
			}),
		)
	})
}

// drawSettingsButtons renders the Save and Cancel buttons.
func (s *appState) drawSettingsButtons(gtx layout.Context) layout.Dimensions {
	// Check for button clicks
	if s.settingsModal.saveButton.Clicked(gtx) {
		s.saveAndCloseSettings()
	}
	if s.settingsModal.cancelButton.Clicked(gtx) {
		s.closeSettingsModal()
	}

	// Colors
	saveBg := color.NRGBA{R: 0x2b, G: 0x50, B: 0x8a, A: 0xff}
	saveText := color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	cancelText := color.NRGBA{R: 0x6d, G: 0xb3, B: 0xff, A: 0xff}

	return layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx,
			// Cancel button
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.settingsModal.cancelButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top: unit.Dp(8), Bottom: unit.Dp(8),
						Left: unit.Dp(16), Right: unit.Dp(16),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body2(s.theme, "Cancel")
						lbl.Color = cancelText
						return lbl.Layout(gtx)
					})
				})
			}),

			// Spacer
			layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),

			// Save button
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.settingsModal.saveButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Background{}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							rr := gtx.Dp(unit.Dp(4))
							rect := clip.RRect{
								Rect: image.Rectangle{Max: gtx.Constraints.Min},
								NE:   rr, NW: rr, SE: rr, SW: rr,
							}.Push(gtx.Ops)
							paint.Fill(gtx.Ops, saveBg)
							rect.Pop()
							return layout.Dimensions{Size: gtx.Constraints.Min}
						},
						func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{
								Top: unit.Dp(8), Bottom: unit.Dp(8),
								Left: unit.Dp(16), Right: unit.Dp(16),
							}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body2(s.theme, "Save")
								lbl.Color = saveText
								return lbl.Layout(gtx)
							})
						},
					)
				})
			}),
		)
	})
}

// handleSettingsModalKey handles key events when the settings modal is active.
func (s *appState) handleSettingsModalKey(ev key.Event) {
	// Special case: Tab events only come through on Release (consumed by focus system on Press)
	isTabRelease := ev.Name == key.NameTab && ev.State == key.Release

	// Only process Press events, except Tab which comes as Release
	if ev.State != key.Press && !isTabRelease {
		return
	}

	// Get max items for current tab
	maxItems := s.getTabItemCount()

	// Normalize key name for case-insensitive matching
	keyName := strings.ToLower(string(ev.Name))

	// Handle special keys first
	switch ev.Name {
	case key.NameEscape:
		s.closeSettingsModal()
		return
	case key.NameReturn, key.NameEnter:
		s.saveAndCloseSettings()
		return
	case key.NameTab:
		// Cycle tabs and reset focus and scroll
		s.settingsModal.selectedTab = (s.settingsModal.selectedTab + 1) % 3
		s.settingsModal.focusedItem = 0
		s.settingsModal.contentList.List.Position.First = 0
		s.settingsModal.contentList.List.Position.Offset = 0
		return
	case key.NameDownArrow:
		if s.settingsModal.focusedItem < maxItems-1 {
			s.settingsModal.focusedItem++
		}
		return
	case key.NameUpArrow:
		if s.settingsModal.focusedItem > 0 {
			s.settingsModal.focusedItem--
		}
		return
	case key.NameLeftArrow:
		s.adjustFocusedSetting(-1)
		return
	case key.NameRightArrow:
		s.adjustFocusedSetting(1)
		return
	case key.NameSpace:
		s.toggleFocusedSetting()
		return
	}

	// Handle character keys (case-insensitive)
	switch keyName {
	case "1":
		s.settingsModal.selectedTab = tabGeneral
		s.settingsModal.focusedItem = 0
		s.settingsModal.contentList.List.Position.First = 0
		s.settingsModal.contentList.List.Position.Offset = 0
	case "2":
		s.settingsModal.selectedTab = tabEditor
		s.settingsModal.focusedItem = 0
		s.settingsModal.contentList.List.Position.First = 0
		s.settingsModal.contentList.List.Position.Offset = 0
	case "3":
		s.settingsModal.selectedTab = tabLSP
		s.settingsModal.focusedItem = 0
		s.settingsModal.contentList.List.Position.First = 0
		s.settingsModal.contentList.List.Position.Offset = 0
	case "j":
		if s.settingsModal.focusedItem < maxItems-1 {
			s.settingsModal.focusedItem++
		}
	case "k":
		if s.settingsModal.focusedItem > 0 {
			s.settingsModal.focusedItem--
		}
	case "h":
		s.adjustFocusedSetting(-1)
	case "l":
		s.adjustFocusedSetting(1)
	case " ":
		s.toggleFocusedSetting()
	}
}

// getTabItemCount returns the number of focusable items in the current tab.
func (s *appState) getTabItemCount() int {
	switch s.settingsModal.selectedTab {
	case tabGeneral:
		return 2 // Font Size, Theme
	case tabEditor:
		return 7 // AutoIndent, AutoPairs, SoftWrap, FormatOnSave, UseSpaces, TabWidth, ScrollOffset
	case tabLSP:
		return 2 // Enabled, AutoDetect
	default:
		return 1
	}
}

// toggleFocusedSetting toggles a boolean setting at the focused position.
func (s *appState) toggleFocusedSetting() {
	switch s.settingsModal.selectedTab {
	case tabGeneral:
		// Font size is a stepper, space does nothing
	case tabEditor:
		switch s.settingsModal.focusedItem {
		case 0:
			s.settingsModal.autoIndent.Value = !s.settingsModal.autoIndent.Value
		case 1:
			s.settingsModal.autoPairs.Value = !s.settingsModal.autoPairs.Value
		case 2:
			s.settingsModal.softWrap.Value = !s.settingsModal.softWrap.Value
		case 3:
			s.settingsModal.formatOnSave.Value = !s.settingsModal.formatOnSave.Value
		case 4:
			s.settingsModal.useSpaces.Value = !s.settingsModal.useSpaces.Value
		// 5 and 6 are steppers
		}
	case tabLSP:
		switch s.settingsModal.focusedItem {
		case 0:
			s.settingsModal.lspEnabled.Value = !s.settingsModal.lspEnabled.Value
		case 1:
			s.settingsModal.lspAutoDetect.Value = !s.settingsModal.lspAutoDetect.Value
		}
	}
}

// adjustFocusedSetting adjusts a numeric stepper setting or cycles through theme list.
func (s *appState) adjustFocusedSetting(delta int) {
	switch s.settingsModal.selectedTab {
	case tabGeneral:
		switch s.settingsModal.focusedItem {
		case 0:
			// Font Size: 8-32
			newVal := s.settingsModal.fontSizeValue + delta
			if newVal >= 8 && newVal <= 32 {
				s.settingsModal.fontSizeValue = newVal
			}
		case 1:
			// Theme selector: use delta to navigate up/down
			s.adjustThemeSelection(delta)
		}
	case tabEditor:
		switch s.settingsModal.focusedItem {
		case 5:
			// Tab Width: 2-8
			newVal := s.settingsModal.tabWidthValue + delta
			if newVal >= 2 && newVal <= 8 {
				s.settingsModal.tabWidthValue = newVal
			}
		case 6:
			// Scroll Offset: 0-20
			newVal := s.settingsModal.scrollOffsetValue + delta
			if newVal >= 0 && newVal <= 20 {
				s.settingsModal.scrollOffsetValue = newVal
			}
		}
	}
}

// adjustThemeSelection changes the selected theme and triggers live preview.
func (s *appState) adjustThemeSelection(delta int) {
	if len(s.settingsModal.themes) == 0 {
		return
	}

	newIdx := s.settingsModal.themeSelected + delta
	if newIdx < 0 {
		newIdx = 0
	}
	if newIdx >= len(s.settingsModal.themes) {
		newIdx = len(s.settingsModal.themes) - 1
	}

	if newIdx != s.settingsModal.themeSelected {
		s.settingsModal.themeSelected = newIdx

		// Enable preview mode and store original theme
		if !s.settingsModal.themePreview {
			s.settingsModal.themePreview = true
			s.settingsModal.themeOriginal = s.currentTheme
		}

		// Apply theme for live preview
		s.loadTheme(s.settingsModal.themes[s.settingsModal.themeSelected])
	}
}

// closeSettingsModal closes the modal without saving changes.
func (s *appState) closeSettingsModal() {
	// Restore original theme if preview was active
	if s.settingsModal.themePreview {
		s.loadTheme(s.settingsModal.themeOriginal)
		s.settingsModal.themePreview = false
	}

	s.settingsModal.active = false
	s.status = ""
}

// saveAndCloseSettings saves settings and closes the modal.
func (s *appState) saveAndCloseSettings() {
	// UI settings
	s.settings.UI.FontSize = s.settingsModal.fontSizeValue

	// Theme setting - save the selected theme name
	if s.settingsModal.themeSelected >= 0 && s.settingsModal.themeSelected < len(s.settingsModal.themes) {
		s.settings.UI.Theme = s.settingsModal.themes[s.settingsModal.themeSelected]
	}
	// Clear preview mode (theme already applied via live preview)
	s.settingsModal.themePreview = false

	// Editor settings
	s.settings.Editor.AutoIndent = s.settingsModal.autoIndent.Value
	s.settings.Editor.AutoPairs = s.settingsModal.autoPairs.Value
	s.settings.Editor.SoftWrap = s.settingsModal.softWrap.Value
	s.settings.Editor.FormatOnSave = s.settingsModal.formatOnSave.Value
	s.settings.Editor.UseSpaces = s.settingsModal.useSpaces.Value
	s.settings.Editor.TabWidth = s.settingsModal.tabWidthValue
	s.settings.Editor.ScrollOffset = s.settingsModal.scrollOffsetValue

	// LSP settings
	s.settings.LSP.Enabled = s.settingsModal.lspEnabled.Value
	s.settings.LSP.AutoDetect = s.settingsModal.lspAutoDetect.Value

	// Save to file
	if err := s.settings.Save(); err != nil {
		s.status = fmt.Sprintf("Error saving settings: %v", err)
		return
	}

	// Apply to runtime
	s.applySettings()

	s.settingsModal.active = false
	s.status = "Settings saved"
}

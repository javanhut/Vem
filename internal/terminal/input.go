package terminal

import (
	"gioui.org/io/key"
)

// KeyToTerminalSequence converts Gio key event to terminal input
func KeyToTerminalSequence(ev key.Event) string {
	// Special keys (VT100 escape sequences)
	switch ev.Name {
	case key.NameReturn, key.NameEnter:
		return "\r"
	case key.NameTab:
		return "\t"
	case key.NameDeleteBackward:
		return "\x7f" // DEL (127)
	case key.NameDeleteForward:
		return "\x1b[3~" // Delete
	case key.NameUpArrow:
		return "\x1b[A"
	case key.NameDownArrow:
		return "\x1b[B"
	case key.NameRightArrow:
		return "\x1b[C"
	case key.NameLeftArrow:
		return "\x1b[D"
	case key.NameHome:
		return "\x1b[H"
	case key.NameEnd:
		return "\x1b[F"
	case key.NamePageUp:
		return "\x1b[5~"
	case key.NamePageDown:
		return "\x1b[6~"
	case key.NameEscape:
		return "\x1b"
	}

	// Function keys
	switch ev.Name {
	case key.NameF1:
		return "\x1bOP"
	case key.NameF2:
		return "\x1bOQ"
	case key.NameF3:
		return "\x1bOR"
	case key.NameF4:
		return "\x1bOS"
	case key.NameF5:
		return "\x1b[15~"
	case key.NameF6:
		return "\x1b[17~"
	case key.NameF7:
		return "\x1b[18~"
	case key.NameF8:
		return "\x1b[19~"
	case key.NameF9:
		return "\x1b[20~"
	case key.NameF10:
		return "\x1b[21~"
	case key.NameF11:
		return "\x1b[23~"
	case key.NameF12:
		return "\x1b[24~"
	}

	// Shift+Arrow combinations
	if ev.Modifiers.Contain(key.ModShift) {
		switch ev.Name {
		case key.NameUpArrow:
			return "\x1b[1;2A"
		case key.NameDownArrow:
			return "\x1b[1;2B"
		case key.NameRightArrow:
			return "\x1b[1;2C"
		case key.NameLeftArrow:
			return "\x1b[1;2D"
		}
	}

	// Ctrl+Arrow combinations
	if ev.Modifiers.Contain(key.ModCtrl) {
		switch ev.Name {
		case key.NameUpArrow:
			return "\x1b[1;5A"
		case key.NameDownArrow:
			return "\x1b[1;5B"
		case key.NameRightArrow:
			return "\x1b[1;5C"
		case key.NameLeftArrow:
			return "\x1b[1;5D"
		}

		// Ctrl+Key combinations
		if len(ev.Name) == 1 {
			r := rune(ev.Name[0])
			// Ctrl+A = 0x01, Ctrl+B = 0x02, ..., Ctrl+Z = 0x1A
			if r >= 'a' && r <= 'z' {
				return string(rune(r - 'a' + 1))
			}
			if r >= 'A' && r <= 'Z' {
				return string(rune(r - 'A' + 1))
			}
		}

		// Special Ctrl combinations
		switch ev.Name {
		case " ":
			return "\x00" // Ctrl+Space = NUL
		case "2", "@":
			return "\x00" // Ctrl+2 or Ctrl+@ = NUL
		case "3", "[":
			return "\x1b" // Ctrl+3 or Ctrl+[ = ESC
		case "4", "\\":
			return "\x1c" // Ctrl+4 or Ctrl+\ = FS
		case "5", "]":
			return "\x1d" // Ctrl+5 or Ctrl+] = GS
		case "6", "^":
			return "\x1e" // Ctrl+6 or Ctrl+^ = RS
		case "7", "_":
			return "\x1f" // Ctrl+7 or Ctrl+_ = US
		case "8":
			return "\x7f" // Ctrl+8 = DEL
		}
	}

	// Alt+Key combinations
	if ev.Modifiers.Contain(key.ModAlt) {
		if len(ev.Name) == 1 {
			// Alt+X sends ESC followed by X
			return "\x1b" + string(ev.Name)
		}
	}

	// Regular character input comes through EditEvent, not KeyEvent
	return ""
}

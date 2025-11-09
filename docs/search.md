# Search Feature

This document describes the search functionality in ProjectVem, including basic text search and the planned fuzzy file finder.

## Text Search (Implemented)

The text search feature allows you to find and navigate between occurrences of a pattern in the current buffer.

### Entering Search Mode

From **NORMAL mode**, press `/` to enter SEARCH mode. The status bar will show a `/` prompt where you can type your search pattern.

### Search Behavior

- **Case-insensitive**: Searches ignore case (e.g., `hello` matches `Hello`, `HELLO`, `HeLLo`)
- **Substring matching**: Pattern matches anywhere in the text (e.g., `the` matches `the`, `there`, `weather`)
- **Whole buffer**: Searches the entire buffer from start to finish
- **Multiple matches**: All occurrences are found and highlighted

### Using Search Mode

1. **Enter pattern**: Type your search text
   - Characters are added to the search pattern as you type
   - The status bar shows: `/your_pattern`
   
2. **Edit pattern**:
   - Press `Backspace` to delete the last character
   - Press `Esc` to cancel search and return to NORMAL mode
   
3. **Execute search**:
   - Press `Enter` to execute the search
   - The editor jumps to the first match from the current cursor position
   - The status bar shows: `/pattern [current/total]` (e.g., `/hello [1/5]`)

### Navigating Search Results

After executing a search, you can navigate between matches in NORMAL mode:

| Key | Action | Description |
|-----|--------|-------------|
| `n` | Next Match | Jump to next occurrence (wraps to start) |
| `Shift+N` | Previous Match | Jump to previous occurrence (wraps to end) |

### Visual Feedback

Search matches are highlighted with two different colors:

- **Yellow highlight** (rgba(255, 255, 0, 0.47)): All other matches
- **Orange highlight** (rgba(255, 165, 0, 0.67)): Current match where cursor is positioned

The status bar displays the search pattern and position:
```
/pattern [2/5]
```
This means you're at match 2 out of 5 total matches.

### Search Lifecycle

```
NORMAL mode
    |
    | Press '/'
    v
SEARCH mode (typing pattern)
    |
    | Press Enter (execute search)
    v
NORMAL mode (with active search)
    |
    | Use 'n' / 'Shift+N' to navigate
    |
    | Search remains active until:
    | - New search initiated with '/'
    | - Buffer is modified
```

### Example Workflow

**Example 1: Find all occurrences of "function"**

1. Press `/` in NORMAL mode
2. Type `function`
3. Press `Enter` to search
4. Status shows: `/function [1/12]`
5. Press `n` repeatedly to cycle through all 12 matches
6. Press `Shift+N` to go back to previous matches

**Example 2: Cancel search**

1. Press `/` in NORMAL mode
2. Type `hello`
3. Press `Esc` to cancel
4. Returns to NORMAL mode, no search executed

**Example 3: No matches found**

1. Press `/` in NORMAL mode
2. Type `xyz123notfound`
3. Press `Enter`
4. Status shows: `Pattern not found: xyz123notfound`
5. No highlighting appears

### Edge Cases

- **Empty search pattern**: If you press `Enter` with no pattern, the search is cancelled
- **Single match**: Cursor jumps to that match, `n` wraps to same match
- **No matches**: Status shows "Pattern not found: pattern"
- **Search wrapping**: `n` from last match wraps to first match; `Shift+N` from first wraps to last

### Implementation Details

**Location**: `internal/appcore/app.go`

**Key methods**:
- `enterSearchMode()`: Enters SEARCH mode and initializes state
- `exitSearchMode()`: Cancels search and returns to NORMAL mode
- `executeSearch()`: Finds all matches and jumps to first
- `findAllMatches(pattern)`: Returns all SearchMatch instances for pattern
- `jumpToNextMatch()`: Navigates to next match (wraps)
- `jumpToPrevMatch()`: Navigates to previous match (wraps)
- `drawSearchHighlights()`: Renders highlight rectangles for matches

**Data structures**:
```go
type SearchMatch struct {
    Line int  // Line number (0-indexed)
    Col  int  // Column position (rune-based)
    Len  int  // Length of match in runes
}
```

**State fields**:
- `searchPattern string`: Current search pattern
- `searchMatches []SearchMatch`: All matches found
- `currentMatchIdx int`: Index of current match in searchMatches
- `searchActive bool`: Whether search is active with highlights

### Keybindings

**Global**:
- `/` (NORMAL mode): Enter search mode

**SEARCH mode**:
- `Esc`: Cancel search, return to NORMAL
- `Enter`: Execute search
- `Backspace`: Delete last character from pattern
- Any printable character: Append to search pattern

**NORMAL mode (after search)**:
- `n`: Next match
- `Shift+N`: Previous match

## Fuzzy File Finder (Planned)

The fuzzy file finder will allow you to quickly open files by typing partial file names.

### Planned Features

- **Fuzzy matching**: Type partial paths, characters don't need to be consecutive
- **Recursive search**: Searches all files in the workspace recursively
- **Score-based ranking**: Best matches appear first
- **Visual feedback**: Shows matched characters highlighted
- **Quick navigation**: Up/Down arrows to select, Enter to open

### Planned Keybinding

- `Ctrl+P`: Open fuzzy finder (global keybinding)

### Planned UI

The fuzzy finder will overlay on top of the editor with:
- Input field at top showing the pattern
- Scrollable list of matching files
- Highlighted matched characters in file paths
- Match score and position indicator

### Implementation Status

**Current**: Structure defined, not yet implemented
**Planned**: Milestone 4 (Weeks 19-24) - User Fluency & Ergonomics

See [ROADMAP.md](../ROADMAP.md) for full implementation timeline.

## See Also

- [Keybindings Reference](keybindings.md) - Complete keybinding documentation
- [Tutorial](tutorial.md) - Step-by-step guide for new users
- [Navigation Guide](navigation.md) - Buffer and pane navigation

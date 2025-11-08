# Vem

A NeoVim-inspired text editor written in Go with Gio UI, designed to run cross-platform with zero external dependencies.

## Features

### Current Implementation (Phase 1 Spike)

- Vim-like modal editing (NORMAL, INSERT, VISUAL, DELETE, EXPLORER modes)
- File tree explorer with directory navigation
- Multi-buffer support
- Vim-style commands (`:e`, `:w`, `:q`, etc.)
- Pane navigation between file tree and editor
- Fullscreen mode toggle for distraction-free editing

### Key Bindings

#### Window Management
- `Shift+Enter`: Toggle fullscreen mode

#### Pane Navigation
- `Ctrl+H`: Jump to file tree explorer (opens if hidden)
- `Ctrl+L`: Jump back to text editor
- `Ctrl+T`: Toggle file tree visibility

#### Text Navigation
- `h/j/k/l`: Move cursor left/down/up/right
- `0`: Jump to start of line
- `$`: Jump to end of line
- `gg`: Jump to first line
- `G`: Jump to last line
- `<count>G`: Jump to specific line

#### Editing
- `i`: Enter INSERT mode
- `v`: Enter VISUAL line mode
- `d`: Enter DELETE mode
- `Esc`: Return to NORMAL mode

For complete documentation, see `docs/navigation.md` and `SPIKE_NOTES.md`.

## Building and Running

### Prerequisites
- Go 1.25.3+
- Linux: `libvulkan-dev` package (Debian/Ubuntu) or equivalent Vulkan SDK headers

### Run
```bash
GOCACHE="$(pwd)/.gocache" go run .
```

### Test
```bash
go test ./...
```

## Project Structure

- `internal/appcore/`: Main application loop and rendering
- `internal/editor/`: Buffer and text editing logic
- `internal/filesystem/`: File tree and directory navigation
- `docs/`: Documentation
- `SPIKE_NOTES.md`: Current spike implementation notes
- `ROADMAP.md`: Project roadmap and milestones

## License

GPLv2

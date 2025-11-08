# Gio Spike Notes

Early rendering spike for ProjectVem. Captures how to run the Gio prototype plus follow up work before merging into Phase 1 deliverables.

## Prerequisites
- **Go 1.25+** (repo initialized with Go 1.25.3).
- **Gio desktop deps**: Linux builds need `vulkan/vulkan.h`. Install `libvulkan-dev` (Debian/Ubuntu) or the equivalent Vulkan SDK headers before compiling (`sudo apt install libvulkan-dev`). Without it the build fails inside `gioui.org/internal/vk/vulkan_x11.go`.
- Optional: set a repo-local Go build cache to avoid `$HOME` permission issues when invoking `go run`:
  ```bash
  export GOCACHE="$(pwd)/.gocache"
  ```

## Run the spike
```bash
GOCACHE="$(pwd)/.gocache" go run .
```

What you should see:
- A 960×640 window titled **ProjectVem Gio Spike**.
- Sample buffer text, active line highlight, and a Vim-like status bar showing mode, cursor position, and last key.

If the build fails with `fatal error: vulkan/vulkan.h: No such file or directory`, install the Vulkan headers (see prerequisites) or run on a host where they already exist.

## Controls implemented

### Basic Navigation
- `h j k l` or arrow keys: move the cursor.
- `0` / `$`: jump to start/end of the current line.
- `gg`, `G`, and `<count>G`: `gg` jumps to the top of the buffer (counts supported, e.g. `5gg`), bare `G` jumps to the bottom, and `<count>G` jumps to the specified line (for example `42G`).

### Editing Modes
- `i`: enter INSERT mode. Normal text keys (letters, digits, punctuation) insert directly into the buffer; `<Space>` and `<Enter>` work as expected.
- `<Backspace>`/`<Delete>` while in INSERT mode remove characters (with line joins when you backspace at column 0).
- `Shift+V` (or `V`): enter VISUAL (line) mode. Use `j/k` or other motions (including `1` + `Shift+G`) to extend the selection across multiple lines. Press `d` to delete the highlighted block or `<Esc>` to cancel.
- `d`: enter DELETE mode. Type a line number (optional) and press `d` again to delete that line (e.g., `d 5 d` removes line 5). `<Esc>` cancels and returns to NORMAL.
- `<Esc>`: return to NORMAL mode.

### File Tree Explorer
- `Ctrl+T`: toggle the file tree explorer visibility.
- `Ctrl+H`: jump from editor to file tree (opens tree if hidden).
- `Ctrl+L`: jump from file tree back to editor.
- When in EXPLORER mode:
  - `j/k` or arrow keys: navigate up/down the file tree.
  - `h/l` or arrow keys: collapse/expand directories.
  - `Enter`: open file or toggle directory.
  - `r`: refresh the tree.
  - `u`: navigate to parent directory.
  - `q` or `Esc`: exit explorer mode.

### Command Mode
- `:` enters COMMAND mode where you can execute Vim-like commands.
- `:e <path>`: open a file.
- `:w [path]`: save the current buffer.
- `:q`: close buffer/quit.
- `:ex` or `:explore`: toggle file tree explorer.
- `:cd <path>`: change working directory.
- `:bn/:bp`: navigate between buffers.

### Visual Feedback
- Status bar updates mirror cursor moves and record the last key observed so you can confirm actions as you type.
- Focused pane is indicated by a blue border (file tree or editor).
- Mode indicator shows current mode (NORMAL, INSERT, VISUAL, DELETE, EXPLORER, COMMAND).

## Next steps
1. **Text editing**: extend `internal/editor.Buffer` with insert/delete operations and wire `key.EditEvent` data into mutations.
2. **Cursor rendering**: draw a caret at `buffer.Cursor()` instead of highlighting the entire line.
3. **Window splits + multiple buffers**: prototype window manager abstractions on top of Gio layout primitives.
4. **Input abstraction**: replace the spike’s inline key handling with a reusable command layer so future modal logic (motions, macros) can plug in cleanly.
5. **Validation**: once Vulkan headers are available locally, run `go build ./cmd/spike` (or `go run`) to confirm everything links before promoting the spike output into the Phase 1 architecture notes.

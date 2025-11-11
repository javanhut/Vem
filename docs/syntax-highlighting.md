# Syntax Highlighting

Vem includes built-in syntax highlighting powered by [Chroma](https://github.com/alecthomas/chroma), supporting over 200 programming languages and markup formats.

## Features

- **Automatic Language Detection**: Vem automatically detects the programming language based on file extension
- **200+ Languages Supported**: Including Go, Python, JavaScript, C, C++, Rust, Java, HTML, CSS, Markdown, and many more
- **Multiple Color Themes**: Choose from a variety of color schemes to match your preference
- **Smart Caching**: Syntax highlighting results are cached for performance
- **Per-Buffer Highlighting**: Each open buffer gets its own highlighter instance

## Supported Languages

Vem supports syntax highlighting for all languages recognized by Chroma, including but not limited to:

- **Programming Languages**: Go, Python, JavaScript, TypeScript, C, C++, Rust, Java, C#, Ruby, PHP, Swift, Kotlin, Scala
- **Web Technologies**: HTML, CSS, SCSS, SASS, JSX, TSX, Vue, Svelte
- **Shell Scripts**: Bash, Zsh, Fish, PowerShell
- **Configuration**: YAML, TOML, JSON, XML, INI
- **Markup**: Markdown, reStructuredText, AsciiDoc
- **Data Formats**: SQL, GraphQL, Protobuf
- **And many more...

## Color Themes

Vem comes with several pre-configured color themes:

### Available Themes

- **monokai** (default) - Dark theme with vibrant colors
- **dracula** - Dark theme with purple accents
- **github-dark** - GitHub's dark theme
- **nord** - Arctic-inspired dark theme
- **one-dark** - Atom One Dark theme
- **solarized-dark** - Precision colors for machines and people
- **solarized-light** - Light variant of Solarized
- **vim** - Classic Vim color scheme
- **catppuccin-mocha** - Warm, pastel dark theme
- **gruvbox** - Retro groove dark theme

### Changing Themes

To change the color theme, use the `:colorscheme` command in command mode:

```
:colorscheme dracula
:colorscheme nord
:colorscheme solarized-dark
```

### Listing Available Themes

To see all available themes:

```
:colorscheme list
```

## Usage

### Opening Files

When you open a file, syntax highlighting is automatically enabled based on the file extension:

```
:e myfile.go       # Opens with Go syntax highlighting
:e script.py       # Opens with Python syntax highlighting
:e index.html      # Opens with HTML syntax highlighting
```

### Toggling Syntax Highlighting

To disable syntax highlighting for the current buffer:

```
:syntax off
```

To re-enable syntax highlighting:

```
:syntax on
```

To toggle syntax highlighting:

```
:syntax toggle
```

## Performance

### Caching

Vem uses intelligent caching to ensure syntax highlighting doesn't impact performance:

- **Line-based Cache**: Each highlighted line is cached with a hash
- **Automatic Invalidation**: Cache is invalidated when lines are edited
- **Viewport Optimization**: Only visible lines are highlighted

### Large Files

For very large files (>10,000 lines), you may want to disable syntax highlighting:

```
:syntax off
```

## Technical Details

### Architecture

- **Highlighter**: `internal/syntax/highlighter.go`
  - Creates lexer based on file extension
  - Tokenizes lines of code
  - Manages per-line caching

- **Theme Management**: `internal/syntax/theme.go`
  - Converts Chroma color definitions to Gio colors
  - Provides theme selection and queries

- **Integration**: `internal/appcore/app.go`
  - Maintains highlighter instances per buffer
  - Renders tokens with appropriate colors
  - Handles cache invalidation on edits

### Token Types

Chroma recognizes various token types including:

- **Keywords**: `if`, `for`, `class`, `def`, etc.
- **Strings**: String literals
- **Comments**: Single-line and multi-line comments
- **Numbers**: Integer and float literals
- **Operators**: `+`, `-`, `*`, `/`, etc.
- **Identifiers**: Variable and function names
- **Types**: Built-in and user-defined types
- **Punctuation**: Brackets, parentheses, commas, etc.

Each token type is colored according to the selected theme.

## Troubleshooting

### Syntax Highlighting Not Working

If syntax highlighting isn't working for a file:

1. **Check File Extension**: Ensure the file has a recognized extension
2. **Verify Theme**: Try switching to a different theme with `:colorscheme monokai`
3. **Check if Disabled**: Run `:syntax on` to ensure it's enabled
4. **Reload File**: Try closing and reopening the file with `:e`

### Incorrect Colors

If colors don't look right:

1. **Try Different Theme**: Some themes may not define colors for all token types
2. **Check Terminal Colors**: Ensure your terminal supports 24-bit color
3. **Reset Theme**: Switch back to default with `:colorscheme monokai`

### Performance Issues

If experiencing slowdowns with syntax highlighting:

1. **Disable for Large Files**: Use `:syntax off` for files over 10,000 lines
2. **Reduce Viewport**: Fewer visible lines means less highlighting work
3. **Use Simpler Theme**: Some themes are more computationally expensive

## Examples

### Opening a Go File

```
:e main.go
```

The file opens with Go syntax highlighting:
- Keywords like `func`, `package`, `import` in one color
- Strings in another color
- Comments in a muted color
- Types and identifiers appropriately colored

### Switching Themes

```
:colorscheme dracula
```

All open buffers immediately update to use the Dracula theme.

### Disabling for a Specific File

```
:syntax off
```

The current buffer displays in plain text without syntax highlighting.

## Future Enhancements

Planned improvements for syntax highlighting:

- Per-buffer theme settings
- Custom theme creation
- Semantic highlighting (LSP integration)
- Treesitter-based parsing for improved accuracy
- Configurable token colors
- Theme preview before switching

## See Also

- [Navigation](navigation.md) - Moving around in files
- [Search](search.md) - Finding text in files
- [Keybindings](keybindings.md) - Keyboard shortcuts

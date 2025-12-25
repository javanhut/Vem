# Carrion Language Syntax Highlighting

Vem includes comprehensive syntax highlighting support for the Carrion programming language (.crl files). This document describes the implementation and features.

## Overview

The Carrion syntax highlighter is a custom-built Chroma lexer that provides complete syntax highlighting for all Carrion language features. The lexer is automatically registered when Vem starts and is applied to any file with a `.crl` extension.

## Supported Features

### Enhanced Highlighting

The Carrion lexer includes intelligent context-aware highlighting:

- **Function Definitions**: Function names after `spell` are highlighted distinctly (typically yellow/cyan)
- **Class Definitions**: Class names after `grim` are highlighted distinctly (typically green/blue)
- **Function Calls**: Any identifier followed by `(` is highlighted as a function call
- **Built-in Functions**: Standard library functions like `print`, `len`, `range` have their own highlighting

This makes your code more readable by distinguishing:
- `spell calculate(x):` - `calculate` is highlighted as a function name
- `grim Person:` - `Person` is highlighted as a class name
- `calculate(5)` - `calculate` is highlighted as a function call
- Regular variables remain in the default color

### Keywords

All Carrion keywords are properly highlighted, including:

- **Control Flow**: `if`, `otherwise`, `else`, `for`, `in`, `while`, `match`, `case`, `stop`, `skip`, `return`
- **Error Handling**: `attempt`, `ensnare`, `resolve`, `raise`, `check`
- **OOP**: `grim` (class), `spell` (function), `init`, `self`, `super`, `arcane`, `arcanespell`
- **Logical**: `and`, `or`, `not`, `not in`
- **Module System**: `import`, `as`
- **Special Blocks**: `main`, `diverge`, `converge`, `autoclose`
- **Other**: `var`, `global`, `ignore`

### Built-in Types

Recognized built-in types:
- `Array`, `Boolean`, `Float`, `Integer`, `String`
- `Map`, `Tuple`, `File`, `OS`
- `HTTPServer`, `WebServer`, `Error`

### Built-in Functions

All standard Carrion built-in functions are highlighted:
- I/O: `print`, `input`
- Collections: `len`, `range`, `enumerate`, `pairs`
- Type operations: `type`, `int`, `float`, `str`, `bool`, `list`, `tuple`
- Math: `max`, `min`, `abs`
- Character: `ord`, `chr`
- Type checking: `is_sametype`
- Hash/JSON: `parseHash`, `httpParseJSON`, `httpStringifyJSON`
- HTTP: `httpGet`, `httpPost`, `http_response`
- Introspection: `help`, `version`, `modules`

### Operators

All Carrion operators are recognized and highlighted:

- **Arithmetic**: `+`, `-`, `*`, `/`, `//`, `%`, `**`
- **Assignment**: `=`, `+=`, `-=`, `*=`, `/=`, `++`, `--`
- **Comparison**: `==`, `!=`, `<`, `>`, `<=`, `>=`
- **Bitwise**: `&`, `|`, `^`, `~`, `<<`, `>>`
- **Special**: `->`, `<-`, `!`

### String Types

Multiple string types are supported with proper highlighting:

1. **Regular Strings**: `"string"` or `'string'`
2. **Triple-Quoted Strings**: `"""multiline"""` or `'''multiline'''`
3. **F-Strings**: `f"Hello {name}"` with interpolation
4. **Interpolated Strings**: `i"Value: ${expression}"`

### Comments

All three comment styles are recognized:

1. **Single-line**: `# comment`
2. **Block comments**: `/* comment */`
3. **Triple-backtick**: ` ``` comment ``` `

### Numbers

Both integer and floating-point literals are properly highlighted:

- Integers: `42`, `-17`, `0`
- Floats: `3.14`, `-2.5`, `0.5`

### Special Constructs

The lexer properly handles Carrion's special block statements:

- **Main entry point**: `main:`
- **Goroutine creation**: `diverge:` or `diverge worker:`
- **Synchronization**: `converge` or `converge worker1, worker2`

## Implementation Details

### File Structure

```
internal/syntax/
├── lexers/
│   ├── carrion.go           # Carrion lexer implementation
│   └── carrion_test.go      # Lexer unit tests
├── highlighter.go            # Main highlighter (imports Carrion lexer)
├── highlighter_carrion_test.go  # Integration tests
└── theme.go                  # Theme management
```

### Lexer Registration

The Carrion lexer is registered with Chroma during package initialization:

```go
func init() {
    lexers.Register(chroma.MustNewLexer(...))
}
```

The highlighter automatically imports the lexer package:

```go
import _ "github.com/javanhut/vem/internal/syntax/lexers"
```

### Token Types

The lexer uses standard Chroma token types for consistency:

- `chroma.Keyword` - Reserved keywords
- `chroma.KeywordConstant` - `True`, `False`, `None`
- `chroma.KeywordType` - Built-in types
- `chroma.NameBuiltin` - Built-in functions
- `chroma.NameFunction` - Function names (definitions and calls)
- `chroma.NameClass` - Class names (after `grim`)
- `chroma.CommentSingle` - Single-line comments
- `chroma.CommentMultiline` - Block comments
- `chroma.StringDouble` - Double-quoted strings and f-strings
- `chroma.StringSingle` - Single-quoted strings
- `chroma.StringDoc` - Triple-quoted strings (docstrings)
- `chroma.StringInterpol` - Interpolated strings
- `chroma.NumberInteger` - Integer literals
- `chroma.NumberFloat` - Floating-point literals
- `chroma.Operator` - Operators
- `chroma.OperatorWord` - Word operators (`and`, `or`, `not`)
- `chroma.Punctuation` - Delimiters and punctuation
- `chroma.Name` - Regular identifiers (variables)
- `chroma.Text` - Whitespace

## Usage

### Opening Carrion Files

Simply open any `.crl` file and syntax highlighting is applied automatically:

```
:e program.crl
vem script.crl
```

### Example File

A comprehensive test file is available at `examples/test.crl` that demonstrates all supported syntax features.

### Highlighting Examples

Here are examples of how different elements are highlighted:

```carrion
# Function definition - 'greet' will be highlighted in function color
spell greet(name):
    message = f"Hello, {name}!"
    print(message)          # 'print' is a built-in (different color)
    return True

# Class definition - 'Person' will be highlighted in class color
grim Person:
    init(name, age):        # 'init' is a keyword
        self.name = name    # 'self' is a keyword, 'name' is a variable
    
    spell get_info():       # 'get_info' is highlighted as function
        return self.name

# Function calls - highlighted distinctly
result = greet("Alice")     # 'greet' is highlighted as function call
person = Person("Bob", 30)  # 'Person' is highlighted as class name
info = person.get_info()    # 'get_info' is highlighted as function call

# Abstract class - both 'arcane' and 'grim' are keywords, 'Animal' is class name
arcane grim Animal:
    arcanespell
    spell make_sound():
        ignore
```

In the above example:
- **Keywords** (`spell`, `grim`, `init`, `self`, `return`, `arcane`, etc.) - One color
- **Function names** (`greet`, `get_info`, `make_sound`) - Distinct function color (yellow/cyan)
- **Class names** (`Person`, `Animal`) - Distinct class color (green/blue)
- **Built-in functions** (`print`) - Built-in color (often magenta/purple)
- **Variables** (`message`, `name`, `age`, `info`) - Default identifier color (white/gray)
- **Constants** (`True`) - Constant color

### Changing Themes

All Chroma color themes work with Carrion syntax highlighting:

```
:colorscheme dracula
:colorscheme nord
:colorscheme monokai
```

### Disabling/Enabling

Use standard syntax highlighting commands:

```
:syntax off      # Disable highlighting
:syntax on       # Enable highlighting
:syntax toggle   # Toggle highlighting
```

## Testing

The implementation includes comprehensive tests:

### Lexer Tests

Located in `internal/syntax/lexers/carrion_test.go`:

- Registration verification
- File extension matching (.crl)
- Token generation for all syntax categories
- Keyword, built-in, comment, string, number, and operator tests

### Integration Tests

Located in `internal/syntax/highlighter_carrion_test.go`:

- Language detection for .crl files
- Highlighter creation and initialization
- Code highlighting verification
- File extension filtering

### Running Tests

```bash
# Test the lexer
go test ./internal/syntax/lexers/

# Test the highlighter integration
go test ./internal/syntax/

# Run all tests with verbose output
go test -v ./internal/syntax/...
```

All tests pass successfully:

```
PASS: TestCarrionLexerRegistration
PASS: TestCarrionLexerTokenization
PASS: TestCarrionHighlighting
PASS: TestCarrionKeywordHighlighting
PASS: TestCarrionShouldHighlight
```

## Technical Reference

### Regular Expression Patterns

The lexer uses carefully ordered regex patterns to ensure correct tokenization:

1. **Comments** (highest priority to prevent conflicts)
2. **Multi-word keywords** (`not in`)
3. **Single keywords**
4. **Built-in types and functions**
5. **Strings** (f-strings and interpolated before regular)
6. **Numbers** (floats before integers)
7. **Multi-character operators** (before single-character)
8. **Single-character operators**
9. **Delimiters**
10. **Identifiers**
11. **Whitespace**

### Pattern Matching Order

The order is critical for correct tokenization:

- Comments are matched first (highest priority)
- Function definitions (`spell name`) are matched before keywords
- Class definitions (`grim Name`) are matched before keywords
- Keywords are matched before identifiers
- F-strings are matched before regular strings
- Triple-quoted strings are matched before regular strings
- Float literals are matched before integers
- Multi-character operators (`==`, `**`, `//`) are matched before single characters
- Function calls (identifier followed by `(`) are matched before regular identifiers

## Future Enhancements

Potential improvements for Carrion syntax highlighting:

- Semantic highlighting for function/class definitions
- Context-aware highlighting for `self` and `super`
- Better string interpolation visualization
- Error highlighting for invalid syntax
- Integration with Carrion LSP (if available)

## See Also

- [Syntax Highlighting Guide](syntax-highlighting.md) - General syntax highlighting documentation
- [Carrion Language Reference](../Syntax-Highlighter-Reference.md) - Complete language specification
- [Keybindings Reference](keybindings.md) - Editor keyboard shortcuts

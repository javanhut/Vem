package lexers

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

func init() {
	lexers.Register(chroma.MustNewLexer(
		&chroma.Config{
			Name:      "Carrion",
			Aliases:   []string{"carrion", "crl"},
			Filenames: []string{"*.crl"},
			MimeTypes: []string{"text/x-carrion"},
		},
		func() chroma.Rules {
			return chroma.Rules{
				"root": {
					// Comments (highest priority)
					{`#[^\n]*`, chroma.CommentSingle, nil},
					{`/\*(.|\n)*?\*/`, chroma.CommentMultiline, nil},
					{"```(.|\n)*?```", chroma.CommentMultiline, nil},

					// Multi-word keywords
					{`\bnot\s+in\b`, chroma.OperatorWord, nil},

					// Function definitions - spell keyword followed by function name
					{`(\bspell\b)(\s+)([a-zA-Z_][a-zA-Z0-9_]*)`, chroma.ByGroups(chroma.Keyword, chroma.Text, chroma.NameFunction), nil},

					// Class definitions - grim keyword followed by class name
					{`(\bgrim\b)(\s+)([a-zA-Z_][a-zA-Z0-9_]*)`, chroma.ByGroups(chroma.Keyword, chroma.Text, chroma.NameClass), nil},

					// Abstract class definitions - arcane grim followed by class name
					{`(\barcane\b)(\s+)(\bgrim\b)(\s+)([a-zA-Z_][a-zA-Z0-9_]*)`, chroma.ByGroups(chroma.Keyword, chroma.Text, chroma.Keyword, chroma.Text, chroma.NameClass), nil},

					// Keywords (spell and grim removed since they're handled above)
					{`\b(and|arcane|arcanespell|as|attempt|autoclose|case|check|converge|diverge|else|ensnare|for|global|if|ignore|import|in|init|main|match|or|otherwise|raise|resolve|return|self|skip|stop|super|var|while)\b`, chroma.Keyword, nil},

					// Boolean and None literals
					{`\b(True|False|None)\b`, chroma.KeywordConstant, nil},

					// Logical operators (keyword-like)
					{`\bnot\b`, chroma.OperatorWord, nil},

					// Built-in types
					{`\b(Array|Boolean|Float|Integer|String|Map|Tuple|File|OS|HTTPServer|WebServer|Error)\b`, chroma.KeywordType, nil},

					// Built-in functions
					{`\b(print|input|len|type|range|max|min|abs|ord|chr|int|float|str|bool|list|tuple|enumerate|pairs|is_sametype|parseHash|httpGet|httpPost|httpParseJSON|httpStringifyJSON|http_response|help|version|modules)\b`, chroma.NameBuiltin, nil},

					// F-strings (must match before regular strings)
					{`f("""(.|\n)*?"""|'''(.|\n)*?'''|"([^"\\]|\\.)*"|'([^'\\]|\\.)*')`, chroma.StringDouble, nil},

					// Interpolated strings
					{`i("([^"\\$]|\\.|\$\{[^}]*\})*"|'([^'\\$]|\\.|\$\{[^}]*\})*')`, chroma.StringInterpol, nil},

					// Triple-quoted strings (docstrings)
					{`"""(.|\n)*?"""`, chroma.StringDoc, nil},
					{`'''(.|\n)*?'''`, chroma.StringDoc, nil},

					// Regular strings
					{`"([^"\\]|\\.)*"`, chroma.StringDouble, nil},
					{`'([^'\\]|\\.)*'`, chroma.StringSingle, nil},

					// Numbers (float before int)
					{`-?\d+\.\d+`, chroma.NumberFloat, nil},
					{`-?\d+`, chroma.NumberInteger, nil},

					// Multi-character operators
					{`(==|!=|<=|>=|\*\*|//|<<|>>|\+=|-=|\*=|/=|\+\+|--|->|<-)`, chroma.Operator, nil},

					// Single-character operators
					{`[+\-*/%=<>!&|^~]`, chroma.Operator, nil},

					// Delimiters
					{`[()[\]{},:;.]`, chroma.Punctuation, nil},

					// Function calls - identifier followed by (
					{`([a-zA-Z_][a-zA-Z0-9_]*)(\s*)(\()`, chroma.ByGroups(chroma.NameFunction, chroma.Text, chroma.Punctuation), nil},

					// Identifiers
					{`[a-zA-Z_][a-zA-Z0-9_]*`, chroma.Name, nil},

					// Whitespace
					{`\s+`, chroma.Text, nil},
				},
			}
		},
	))
}

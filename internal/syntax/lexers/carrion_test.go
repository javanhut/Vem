package lexers

import (
	"testing"

	"github.com/alecthomas/chroma/v2/lexers"
)

func TestCarrionLexerRegistration(t *testing.T) {
	// Test that the Carrion lexer is registered
	lexer := lexers.Get("carrion")
	if lexer == nil {
		t.Fatal("Carrion lexer not registered")
	}

	config := lexer.Config()
	if config == nil {
		t.Fatal("Carrion lexer has no config")
	}

	if config.Name != "Carrion" {
		t.Errorf("Expected lexer name 'Carrion', got '%s'", config.Name)
	}

	// Test file name matching
	lexer = lexers.Match("test.crl")
	if lexer == nil {
		t.Fatal("Carrion lexer does not match .crl files")
	}

	config = lexer.Config()
	if config.Name != "Carrion" {
		t.Errorf("Expected matched lexer to be 'Carrion', got '%s'", config.Name)
	}
}

func TestCarrionLexerTokenization(t *testing.T) {
	lexer := lexers.Get("carrion")
	if lexer == nil {
		t.Fatal("Carrion lexer not registered")
	}

	testCases := []struct {
		name  string
		input string
	}{
		{"Keywords", "spell grim attempt ensnare if for while"},
		{"Built-ins", "print len range type Array String"},
		{"Comments", "# This is a comment\n/* block */"},
		{"Strings", `"hello" 'world' f"test {x}" i"value ${y}"`},
		{"Numbers", "42 3.14 -17 -2.5"},
		{"Operators", "+ - * / == != <= >= ** //"},
		{"Special", "main: diverge: converge"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iterator, err := lexer.Tokenise(nil, tc.input)
			if err != nil {
				t.Fatalf("Failed to tokenize %s: %v", tc.name, err)
			}

			tokens := iterator.Tokens()
			if len(tokens) == 0 {
				t.Errorf("No tokens generated for %s", tc.name)
			}
		})
	}
}

func TestCarrionFunctionAndClassHighlighting(t *testing.T) {
	lexer := lexers.Get("carrion")
	if lexer == nil {
		t.Fatal("Carrion lexer not registered")
	}

	testCases := []struct {
		name          string
		input         string
		expectedToken string
		expectedType  string
	}{
		{
			name:          "Function definition",
			input:         "spell calculate(x, y):",
			expectedToken: "calculate",
			expectedType:  "NameFunction",
		},
		{
			name:          "Class definition",
			input:         "grim Person:",
			expectedToken: "Person",
			expectedType:  "NameClass",
		},
		{
			name:          "Abstract class definition",
			input:         "arcane grim Animal:",
			expectedToken: "Animal",
			expectedType:  "NameClass",
		},
		{
			name:          "Built-in function call",
			input:         "print(message)",
			expectedToken: "print",
			expectedType:  "NameBuiltin",
		},
		{
			name:          "Method call",
			input:         "greet(name)",
			expectedToken: "greet",
			expectedType:  "NameFunction",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			iterator, err := lexer.Tokenise(nil, tc.input)
			if err != nil {
				t.Fatalf("Failed to tokenize: %v", err)
			}

			tokens := iterator.Tokens()
			found := false
			for _, token := range tokens {
				if token.Value == tc.expectedToken {
					tokenTypeStr := token.Type.String()
					if tokenTypeStr == tc.expectedType {
						found = true
						break
					}
				}
			}

			if !found {
				t.Errorf("Expected token '%s' with type '%s' not found in: %s",
					tc.expectedToken, tc.expectedType, tc.input)
				t.Logf("Tokens found:")
				for _, token := range tokens {
					t.Logf("  %s: %s", token.Type, token.Value)
				}
			}
		})
	}
}

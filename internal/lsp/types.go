// Package lsp provides Language Server Protocol client functionality.
// It implements LSP 3.17 specification for communication with language servers.
package lsp

import "encoding/json"

// Position in a text document expressed as zero-based line and character offset.
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Range in a text document expressed as start and end positions.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Location represents a location inside a resource, such as a line inside a text file.
type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

// LocationLink represents a link between a source and a target location.
type LocationLink struct {
	OriginSelectionRange *Range      `json:"originSelectionRange,omitempty"`
	TargetURI            DocumentURI `json:"targetUri"`
	TargetRange          Range       `json:"targetRange"`
	TargetSelectionRange Range       `json:"targetSelectionRange"`
}

// DocumentURI is a document identifier using file:// URI scheme.
type DocumentURI string

// TextDocumentIdentifier identifies a text document.
type TextDocumentIdentifier struct {
	URI DocumentURI `json:"uri"`
}

// VersionedTextDocumentIdentifier identifies a specific version of a text document.
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

// TextDocumentItem is an item to transfer a text document from the client to the server.
type TextDocumentItem struct {
	URI        DocumentURI `json:"uri"`
	LanguageID string      `json:"languageId"`
	Version    int         `json:"version"`
	Text       string      `json:"text"`
}

// TextDocumentPositionParams is a parameter literal used in requests to pass
// a text document and a position inside that document.
type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// DiagnosticSeverity represents the severity of a diagnostic.
type DiagnosticSeverity int

const (
	// DiagnosticSeverityError reports an error.
	DiagnosticSeverityError DiagnosticSeverity = 1
	// DiagnosticSeverityWarning reports a warning.
	DiagnosticSeverityWarning DiagnosticSeverity = 2
	// DiagnosticSeverityInformation reports an information.
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	// DiagnosticSeverityHint reports a hint.
	DiagnosticSeverityHint DiagnosticSeverity = 4
)

// DiagnosticTag provides additional metadata about diagnostics.
type DiagnosticTag int

const (
	// DiagnosticTagUnnecessary indicates unused or unnecessary code.
	DiagnosticTagUnnecessary DiagnosticTag = 1
	// DiagnosticTagDeprecated indicates deprecated or obsolete code.
	DiagnosticTagDeprecated DiagnosticTag = 2
)

// Diagnostic represents a diagnostic, such as a compiler error or warning.
type Diagnostic struct {
	Range    Range              `json:"range"`
	Severity DiagnosticSeverity `json:"severity,omitempty"`
	Code     interface{}        `json:"code,omitempty"`
	Source   string             `json:"source,omitempty"`
	Message  string             `json:"message"`
	Tags     []DiagnosticTag    `json:"tags,omitempty"`
}

// PublishDiagnosticsParams is the parameters of the textDocument/publishDiagnostics notification.
type PublishDiagnosticsParams struct {
	URI         DocumentURI  `json:"uri"`
	Version     int          `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// CompletionItemKind defines the type of a completion item.
type CompletionItemKind int

const (
	CompletionItemKindText          CompletionItemKind = 1
	CompletionItemKindMethod        CompletionItemKind = 2
	CompletionItemKindFunction      CompletionItemKind = 3
	CompletionItemKindConstructor   CompletionItemKind = 4
	CompletionItemKindField         CompletionItemKind = 5
	CompletionItemKindVariable      CompletionItemKind = 6
	CompletionItemKindClass         CompletionItemKind = 7
	CompletionItemKindInterface     CompletionItemKind = 8
	CompletionItemKindModule        CompletionItemKind = 9
	CompletionItemKindProperty      CompletionItemKind = 10
	CompletionItemKindUnit          CompletionItemKind = 11
	CompletionItemKindValue         CompletionItemKind = 12
	CompletionItemKindEnum          CompletionItemKind = 13
	CompletionItemKindKeyword       CompletionItemKind = 14
	CompletionItemKindSnippet       CompletionItemKind = 15
	CompletionItemKindColor         CompletionItemKind = 16
	CompletionItemKindFile          CompletionItemKind = 17
	CompletionItemKindReference     CompletionItemKind = 18
	CompletionItemKindFolder        CompletionItemKind = 19
	CompletionItemKindEnumMember    CompletionItemKind = 20
	CompletionItemKindConstant      CompletionItemKind = 21
	CompletionItemKindStruct        CompletionItemKind = 22
	CompletionItemKindEvent         CompletionItemKind = 23
	CompletionItemKindOperator      CompletionItemKind = 24
	CompletionItemKindTypeParameter CompletionItemKind = 25
)

// CompletionItem represents a completion suggestion.
type CompletionItem struct {
	Label               string              `json:"label"`
	Kind                CompletionItemKind  `json:"kind,omitempty"`
	Tags                []CompletionItemTag `json:"tags,omitempty"`
	Detail              string              `json:"detail,omitempty"`
	Documentation       interface{}         `json:"documentation,omitempty"`
	Deprecated          bool                `json:"deprecated,omitempty"`
	Preselect           bool                `json:"preselect,omitempty"`
	SortText            string              `json:"sortText,omitempty"`
	FilterText          string              `json:"filterText,omitempty"`
	InsertText          string              `json:"insertText,omitempty"`
	InsertTextFormat    InsertTextFormat    `json:"insertTextFormat,omitempty"`
	InsertTextMode      InsertTextMode      `json:"insertTextMode,omitempty"`
	TextEdit            *TextEdit           `json:"textEdit,omitempty"`
	AdditionalTextEdits []TextEdit          `json:"additionalTextEdits,omitempty"`
	CommitCharacters    []string            `json:"commitCharacters,omitempty"`
	Command             *Command            `json:"command,omitempty"`
	Data                interface{}         `json:"data,omitempty"`
}

// CompletionItemTag provides extra annotations for completion items.
type CompletionItemTag int

const (
	CompletionItemTagDeprecated CompletionItemTag = 1
)

// InsertTextFormat defines whether the insert text is a plain string or a snippet.
type InsertTextFormat int

const (
	InsertTextFormatPlainText InsertTextFormat = 1
	InsertTextFormatSnippet   InsertTextFormat = 2
)

// InsertTextMode defines how whitespace and indentation is handled during insertion.
type InsertTextMode int

const (
	InsertTextModeAsIs              InsertTextMode = 1
	InsertTextModeAdjustIndentation InsertTextMode = 2
)

// CompletionList represents a collection of completion items.
type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

// CompletionParams are the parameters for a completion request.
type CompletionParams struct {
	TextDocumentPositionParams
	Context *CompletionContext `json:"context,omitempty"`
}

// CompletionContext contains additional information about the context of the completion.
type CompletionContext struct {
	TriggerKind      CompletionTriggerKind `json:"triggerKind"`
	TriggerCharacter string                `json:"triggerCharacter,omitempty"`
}

// CompletionTriggerKind defines how a completion was triggered.
type CompletionTriggerKind int

const (
	CompletionTriggerKindInvoked                         CompletionTriggerKind = 1
	CompletionTriggerKindTriggerCharacter                CompletionTriggerKind = 2
	CompletionTriggerKindTriggerForIncompleteCompletions CompletionTriggerKind = 3
)

// Hover represents the result of a hover request.
type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

// MarkupContent represents a string value with a content type.
type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}

// MarkupKind describes the content type.
type MarkupKind string

const (
	MarkupKindPlainText MarkupKind = "plaintext"
	MarkupKindMarkdown  MarkupKind = "markdown"
)

// TextEdit is a textual edit applicable to a text document.
type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

// WorkspaceEdit represents changes to many documents.
type WorkspaceEdit struct {
	Changes         map[DocumentURI][]TextEdit `json:"changes,omitempty"`
	DocumentChanges []TextDocumentEdit         `json:"documentChanges,omitempty"`
}

// TextDocumentEdit describes textual changes on a single text document.
type TextDocumentEdit struct {
	TextDocument VersionedTextDocumentIdentifier `json:"textDocument"`
	Edits        []TextEdit                      `json:"edits"`
}

// CodeAction represents a code action for fixing issues or refactoring.
type CodeAction struct {
	Title       string           `json:"title"`
	Kind        CodeActionKind   `json:"kind,omitempty"`
	Diagnostics []Diagnostic     `json:"diagnostics,omitempty"`
	IsPreferred bool             `json:"isPreferred,omitempty"`
	Edit        *WorkspaceEdit   `json:"edit,omitempty"`
	Command     *Command         `json:"command,omitempty"`
	Data        *json.RawMessage `json:"data,omitempty"`
}

// CodeActionKind defines the kind of a code action.
type CodeActionKind string

const (
	CodeActionKindEmpty                 CodeActionKind = ""
	CodeActionKindQuickFix              CodeActionKind = "quickfix"
	CodeActionKindRefactor              CodeActionKind = "refactor"
	CodeActionKindRefactorExtract       CodeActionKind = "refactor.extract"
	CodeActionKindRefactorInline        CodeActionKind = "refactor.inline"
	CodeActionKindRefactorRewrite       CodeActionKind = "refactor.rewrite"
	CodeActionKindSource                CodeActionKind = "source"
	CodeActionKindSourceOrganizeImports CodeActionKind = "source.organizeImports"
	CodeActionKindSourceFixAll          CodeActionKind = "source.fixAll"
)

// CodeActionParams are the parameters for a code action request.
type CodeActionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      CodeActionContext      `json:"context"`
}

// CodeActionContext contains additional diagnostic information about the context.
type CodeActionContext struct {
	Diagnostics []Diagnostic     `json:"diagnostics"`
	Only        []CodeActionKind `json:"only,omitempty"`
}

// Command represents a reference to a command.
type Command struct {
	Title     string        `json:"title"`
	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments,omitempty"`
}

// ReferenceParams are the parameters for a references request.
type ReferenceParams struct {
	TextDocumentPositionParams
	Context ReferenceContext `json:"context"`
}

// ReferenceContext contains additional information for references requests.
type ReferenceContext struct {
	IncludeDeclaration bool `json:"includeDeclaration"`
}

// RenameParams are the parameters for a rename request.
type RenameParams struct {
	TextDocumentPositionParams
	NewName string `json:"newName"`
}

// DocumentFormattingParams are the parameters for a formatting request.
type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

// FormattingOptions specifies formatting options.
type FormattingOptions struct {
	TabSize                int  `json:"tabSize"`
	InsertSpaces           bool `json:"insertSpaces"`
	TrimTrailingWhitespace bool `json:"trimTrailingWhitespace,omitempty"`
	InsertFinalNewline     bool `json:"insertFinalNewline,omitempty"`
	TrimFinalNewlines      bool `json:"trimFinalNewlines,omitempty"`
}

// DidOpenTextDocumentParams are the parameters for textDocument/didOpen notification.
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidChangeTextDocumentParams are the parameters for textDocument/didChange notification.
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// TextDocumentContentChangeEvent describes a content change event.
type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength int    `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

// DidSaveTextDocumentParams are the parameters for textDocument/didSave notification.
type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Text         string                 `json:"text,omitempty"`
}

// DidCloseTextDocumentParams are the parameters for textDocument/didClose notification.
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// InitializeParams are the parameters sent in the initialize request.
type InitializeParams struct {
	ProcessID             int                `json:"processId"`
	RootURI               DocumentURI        `json:"rootUri"`
	Capabilities          ClientCapabilities `json:"capabilities"`
	InitializationOptions interface{}        `json:"initializationOptions,omitempty"`
	Trace                 string             `json:"trace,omitempty"`
	WorkspaceFolders      []WorkspaceFolder  `json:"workspaceFolders,omitempty"`
}

// WorkspaceFolder represents a workspace folder.
type WorkspaceFolder struct {
	URI  DocumentURI `json:"uri"`
	Name string      `json:"name"`
}

// ClientCapabilities describe capabilities of the client.
type ClientCapabilities struct {
	Workspace    *WorkspaceClientCapabilities    `json:"workspace,omitempty"`
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	General      *GeneralClientCapabilities      `json:"general,omitempty"`
}

// WorkspaceClientCapabilities describe workspace capabilities.
type WorkspaceClientCapabilities struct {
	ApplyEdit              bool                       `json:"applyEdit,omitempty"`
	WorkspaceEdit          *WorkspaceEditCapabilities `json:"workspaceEdit,omitempty"`
	DidChangeConfiguration *DidChangeConfigCapability `json:"didChangeConfiguration,omitempty"`
	Symbol                 *WorkspaceSymbolCapability `json:"symbol,omitempty"`
	WorkspaceFolders       bool                       `json:"workspaceFolders,omitempty"`
}

// WorkspaceEditCapabilities describe workspace edit capabilities.
type WorkspaceEditCapabilities struct {
	DocumentChanges bool `json:"documentChanges,omitempty"`
}

// DidChangeConfigCapability describes didChangeConfiguration capabilities.
type DidChangeConfigCapability struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// WorkspaceSymbolCapability describes workspace/symbol capabilities.
type WorkspaceSymbolCapability struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// TextDocumentClientCapabilities describe text document capabilities.
type TextDocumentClientCapabilities struct {
	Synchronization    *TextDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
	Completion         *CompletionClientCapabilities       `json:"completion,omitempty"`
	Hover              *HoverClientCapabilities            `json:"hover,omitempty"`
	SignatureHelp      *SignatureHelpClientCapabilities    `json:"signatureHelp,omitempty"`
	Declaration        *DeclarationClientCapabilities      `json:"declaration,omitempty"`
	Definition         *DefinitionClientCapabilities       `json:"definition,omitempty"`
	TypeDefinition     *TypeDefinitionClientCapabilities   `json:"typeDefinition,omitempty"`
	Implementation     *ImplementationClientCapabilities   `json:"implementation,omitempty"`
	References         *ReferencesClientCapabilities       `json:"references,omitempty"`
	DocumentHighlight  *DocumentHighlightClientCaps        `json:"documentHighlight,omitempty"`
	DocumentSymbol     *DocumentSymbolClientCapabilities   `json:"documentSymbol,omitempty"`
	CodeAction         *CodeActionClientCapabilities       `json:"codeAction,omitempty"`
	CodeLens           *CodeLensClientCapabilities         `json:"codeLens,omitempty"`
	Formatting         *FormattingClientCapabilities       `json:"formatting,omitempty"`
	RangeFormatting    *RangeFormattingClientCaps          `json:"rangeFormatting,omitempty"`
	OnTypeFormatting   *OnTypeFormattingClientCaps         `json:"onTypeFormatting,omitempty"`
	Rename             *RenameClientCapabilities           `json:"rename,omitempty"`
	PublishDiagnostics *PublishDiagnosticsClientCaps       `json:"publishDiagnostics,omitempty"`
}

// TextDocumentSyncClientCapabilities describe sync capabilities.
type TextDocumentSyncClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	WillSave            bool `json:"willSave,omitempty"`
	WillSaveWaitUntil   bool `json:"willSaveWaitUntil,omitempty"`
	DidSave             bool `json:"didSave,omitempty"`
}

// CompletionClientCapabilities describe completion capabilities.
type CompletionClientCapabilities struct {
	DynamicRegistration bool                          `json:"dynamicRegistration,omitempty"`
	CompletionItem      *CompletionItemClientCaps     `json:"completionItem,omitempty"`
	CompletionItemKind  *CompletionItemKindClientCaps `json:"completionItemKind,omitempty"`
	ContextSupport      bool                          `json:"contextSupport,omitempty"`
}

// CompletionItemClientCaps describe completion item capabilities.
type CompletionItemClientCaps struct {
	SnippetSupport          bool     `json:"snippetSupport,omitempty"`
	CommitCharactersSupport bool     `json:"commitCharactersSupport,omitempty"`
	DocumentationFormat     []string `json:"documentationFormat,omitempty"`
	DeprecatedSupport       bool     `json:"deprecatedSupport,omitempty"`
	PreselectSupport        bool     `json:"preselectSupport,omitempty"`
}

// CompletionItemKindClientCaps describe completion item kind capabilities.
type CompletionItemKindClientCaps struct {
	ValueSet []CompletionItemKind `json:"valueSet,omitempty"`
}

// HoverClientCapabilities describe hover capabilities.
type HoverClientCapabilities struct {
	DynamicRegistration bool         `json:"dynamicRegistration,omitempty"`
	ContentFormat       []MarkupKind `json:"contentFormat,omitempty"`
}

// SignatureHelpClientCapabilities describe signature help capabilities.
type SignatureHelpClientCapabilities struct {
	DynamicRegistration  bool                     `json:"dynamicRegistration,omitempty"`
	SignatureInformation *SignatureInfoClientCaps `json:"signatureInformation,omitempty"`
	ContextSupport       bool                     `json:"contextSupport,omitempty"`
}

// SignatureInfoClientCaps describe signature info capabilities.
type SignatureInfoClientCaps struct {
	DocumentationFormat  []MarkupKind         `json:"documentationFormat,omitempty"`
	ParameterInformation *ParamInfoClientCaps `json:"parameterInformation,omitempty"`
}

// ParamInfoClientCaps describe parameter info capabilities.
type ParamInfoClientCaps struct {
	LabelOffsetSupport bool `json:"labelOffsetSupport,omitempty"`
}

// DeclarationClientCapabilities describe declaration capabilities.
type DeclarationClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	LinkSupport         bool `json:"linkSupport,omitempty"`
}

// DefinitionClientCapabilities describe definition capabilities.
type DefinitionClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	LinkSupport         bool `json:"linkSupport,omitempty"`
}

// TypeDefinitionClientCapabilities describe type definition capabilities.
type TypeDefinitionClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	LinkSupport         bool `json:"linkSupport,omitempty"`
}

// ImplementationClientCapabilities describe implementation capabilities.
type ImplementationClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	LinkSupport         bool `json:"linkSupport,omitempty"`
}

// ReferencesClientCapabilities describe references capabilities.
type ReferencesClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// DocumentHighlightClientCaps describe document highlight capabilities.
type DocumentHighlightClientCaps struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// DocumentSymbolClientCapabilities describe document symbol capabilities.
type DocumentSymbolClientCapabilities struct {
	DynamicRegistration               bool                  `json:"dynamicRegistration,omitempty"`
	SymbolKind                        *SymbolKindClientCaps `json:"symbolKind,omitempty"`
	HierarchicalDocumentSymbolSupport bool                  `json:"hierarchicalDocumentSymbolSupport,omitempty"`
}

// SymbolKindClientCaps describe symbol kind capabilities.
type SymbolKindClientCaps struct {
	ValueSet []SymbolKind `json:"valueSet,omitempty"`
}

// SymbolKind defines the kind of a symbol.
type SymbolKind int

const (
	SymbolKindFile          SymbolKind = 1
	SymbolKindModule        SymbolKind = 2
	SymbolKindNamespace     SymbolKind = 3
	SymbolKindPackage       SymbolKind = 4
	SymbolKindClass         SymbolKind = 5
	SymbolKindMethod        SymbolKind = 6
	SymbolKindProperty      SymbolKind = 7
	SymbolKindField         SymbolKind = 8
	SymbolKindConstructor   SymbolKind = 9
	SymbolKindEnum          SymbolKind = 10
	SymbolKindInterface     SymbolKind = 11
	SymbolKindFunction      SymbolKind = 12
	SymbolKindVariable      SymbolKind = 13
	SymbolKindConstant      SymbolKind = 14
	SymbolKindString        SymbolKind = 15
	SymbolKindNumber        SymbolKind = 16
	SymbolKindBoolean       SymbolKind = 17
	SymbolKindArray         SymbolKind = 18
	SymbolKindObject        SymbolKind = 19
	SymbolKindKey           SymbolKind = 20
	SymbolKindNull          SymbolKind = 21
	SymbolKindEnumMember    SymbolKind = 22
	SymbolKindStruct        SymbolKind = 23
	SymbolKindEvent         SymbolKind = 24
	SymbolKindOperator      SymbolKind = 25
	SymbolKindTypeParameter SymbolKind = 26
)

// CodeActionClientCapabilities describe code action capabilities.
type CodeActionClientCapabilities struct {
	DynamicRegistration      bool                      `json:"dynamicRegistration,omitempty"`
	CodeActionLiteralSupport *CodeActionLiteralSupport `json:"codeActionLiteralSupport,omitempty"`
	IsPreferredSupport       bool                      `json:"isPreferredSupport,omitempty"`
}

// CodeActionLiteralSupport describes code action literal support.
type CodeActionLiteralSupport struct {
	CodeActionKind *CodeActionKindLiteralSupport `json:"codeActionKind,omitempty"`
}

// CodeActionKindLiteralSupport describes code action kind literal support.
type CodeActionKindLiteralSupport struct {
	ValueSet []CodeActionKind `json:"valueSet,omitempty"`
}

// CodeLensClientCapabilities describe code lens capabilities.
type CodeLensClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// FormattingClientCapabilities describe formatting capabilities.
type FormattingClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// RangeFormattingClientCaps describe range formatting capabilities.
type RangeFormattingClientCaps struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// OnTypeFormattingClientCaps describe on-type formatting capabilities.
type OnTypeFormattingClientCaps struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// RenameClientCapabilities describe rename capabilities.
type RenameClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	PrepareSupport      bool `json:"prepareSupport,omitempty"`
}

// PublishDiagnosticsClientCaps describe publish diagnostics capabilities.
type PublishDiagnosticsClientCaps struct {
	RelatedInformation bool                  `json:"relatedInformation,omitempty"`
	TagSupport         *DiagnosticTagSupport `json:"tagSupport,omitempty"`
	VersionSupport     bool                  `json:"versionSupport,omitempty"`
}

// DiagnosticTagSupport describes supported diagnostic tags.
type DiagnosticTagSupport struct {
	ValueSet []DiagnosticTag `json:"valueSet,omitempty"`
}

// GeneralClientCapabilities describe general capabilities.
type GeneralClientCapabilities struct {
	RegularExpressions *RegularExpressionsCapability `json:"regularExpressions,omitempty"`
	Markdown           *MarkdownCapability           `json:"markdown,omitempty"`
}

// RegularExpressionsCapability describes regex capabilities.
type RegularExpressionsCapability struct {
	Engine  string `json:"engine"`
	Version string `json:"version,omitempty"`
}

// MarkdownCapability describes markdown capabilities.
type MarkdownCapability struct {
	Parser  string `json:"parser"`
	Version string `json:"version,omitempty"`
}

// InitializeResult is the result returned from the initialize request.
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

// ServerInfo contains information about the server.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// ServerCapabilities describe the capabilities of the language server.
type ServerCapabilities struct {
	TextDocumentSync             interface{}                  `json:"textDocumentSync,omitempty"`
	CompletionProvider           *CompletionOptions           `json:"completionProvider,omitempty"`
	HoverProvider                interface{}                  `json:"hoverProvider,omitempty"`
	SignatureHelpProvider        *SignatureHelpOptions        `json:"signatureHelpProvider,omitempty"`
	DeclarationProvider          interface{}                  `json:"declarationProvider,omitempty"`
	DefinitionProvider           interface{}                  `json:"definitionProvider,omitempty"`
	TypeDefinitionProvider       interface{}                  `json:"typeDefinitionProvider,omitempty"`
	ImplementationProvider       interface{}                  `json:"implementationProvider,omitempty"`
	ReferencesProvider           interface{}                  `json:"referencesProvider,omitempty"`
	DocumentHighlightProvider    interface{}                  `json:"documentHighlightProvider,omitempty"`
	DocumentSymbolProvider       interface{}                  `json:"documentSymbolProvider,omitempty"`
	CodeActionProvider           interface{}                  `json:"codeActionProvider,omitempty"`
	CodeLensProvider             *CodeLensOptions             `json:"codeLensProvider,omitempty"`
	DocumentFormattingProvider   interface{}                  `json:"documentFormattingProvider,omitempty"`
	DocumentRangeFormatProvider  interface{}                  `json:"documentRangeFormattingProvider,omitempty"`
	DocumentOnTypeFormatProvider *DocumentOnTypeFormatOpts    `json:"documentOnTypeFormattingProvider,omitempty"`
	RenameProvider               interface{}                  `json:"renameProvider,omitempty"`
	WorkspaceSymbolProvider      interface{}                  `json:"workspaceSymbolProvider,omitempty"`
	Workspace                    *ServerWorkspaceCapabilities `json:"workspace,omitempty"`
}

// CompletionOptions describe completion provider options.
type CompletionOptions struct {
	TriggerCharacters   []string `json:"triggerCharacters,omitempty"`
	AllCommitCharacters []string `json:"allCommitCharacters,omitempty"`
	ResolveProvider     bool     `json:"resolveProvider,omitempty"`
}

// SignatureHelpOptions describe signature help provider options.
type SignatureHelpOptions struct {
	TriggerCharacters   []string `json:"triggerCharacters,omitempty"`
	RetriggerCharacters []string `json:"retriggerCharacters,omitempty"`
}

// CodeLensOptions describe code lens provider options.
type CodeLensOptions struct {
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

// DocumentOnTypeFormatOpts describe on-type formatting options.
type DocumentOnTypeFormatOpts struct {
	FirstTriggerCharacter string   `json:"firstTriggerCharacter"`
	MoreTriggerCharacter  []string `json:"moreTriggerCharacter,omitempty"`
}

// ServerWorkspaceCapabilities describe server workspace capabilities.
type ServerWorkspaceCapabilities struct {
	WorkspaceFolders *WorkspaceFoldersServerCaps `json:"workspaceFolders,omitempty"`
}

// WorkspaceFoldersServerCaps describe workspace folders capabilities.
type WorkspaceFoldersServerCaps struct {
	Supported           bool        `json:"supported,omitempty"`
	ChangeNotifications interface{} `json:"changeNotifications,omitempty"`
}

// TextDocumentSyncKind defines how the host should sync document changes.
type TextDocumentSyncKind int

const (
	// TextDocumentSyncKindNone means documents should not be synced.
	TextDocumentSyncKindNone TextDocumentSyncKind = 0
	// TextDocumentSyncKindFull means documents are synced by sending the full content.
	TextDocumentSyncKindFull TextDocumentSyncKind = 1
	// TextDocumentSyncKindIncremental means documents are synced by sending incremental updates.
	TextDocumentSyncKindIncremental TextDocumentSyncKind = 2
)

// TextDocumentSyncOptions describe text document sync options.
type TextDocumentSyncOptions struct {
	OpenClose bool                 `json:"openClose,omitempty"`
	Change    TextDocumentSyncKind `json:"change,omitempty"`
	WillSave  bool                 `json:"willSave,omitempty"`
	Save      *SaveOptions         `json:"save,omitempty"`
}

// SaveOptions describe save options.
type SaveOptions struct {
	IncludeText bool `json:"includeText,omitempty"`
}

// SignatureHelp represents the signature of a callable.
type SignatureHelp struct {
	Signatures      []SignatureInformation `json:"signatures"`
	ActiveSignature int                    `json:"activeSignature,omitempty"`
	ActiveParameter int                    `json:"activeParameter,omitempty"`
}

// SignatureInformation represents the signature of something callable.
type SignatureInformation struct {
	Label         string                 `json:"label"`
	Documentation interface{}            `json:"documentation,omitempty"`
	Parameters    []ParameterInformation `json:"parameters,omitempty"`
}

// ParameterInformation represents a parameter of a callable.
type ParameterInformation struct {
	Label         interface{} `json:"label"`
	Documentation interface{} `json:"documentation,omitempty"`
}

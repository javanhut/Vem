package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Manager manages multiple language server instances.
type Manager struct {
	servers   map[string]*ServerInstance // key is workspace root
	serversMu sync.RWMutex

	configs map[string]ServerConfig // language -> config

	// Diagnostic callbacks
	onDiagnostics func(uri DocumentURI, diagnostics []Diagnostic)

	// Status callback
	onStatus func(message string)

	// Error callback
	onError func(err error)

	// Invalidate callback (for UI refresh)
	onInvalidate func()
}

// ServerInstance represents a running language server.
type ServerInstance struct {
	Config        ServerConfig
	Client        *Client
	WorkspaceRoot string
	Documents     map[DocumentURI]*DocumentState
	DocumentsMu   sync.RWMutex
	Capabilities  ServerCapabilities
	Diagnostics   map[DocumentURI][]Diagnostic
	DiagnosticsMu sync.RWMutex
}

// DocumentState tracks the state of an open document.
type DocumentState struct {
	URI        DocumentURI
	Version    int
	Content    string
	LanguageID string
}

// LSPServerStatus contains detailed status info for a language server.
type LSPServerStatus struct {
	Name          string
	Language      string
	WorkspaceRoot string
	Running       bool
	DiagCount     int
	DocCount      int
}

// NewManager creates a new LSP manager.
func NewManager() *Manager {
	return &Manager{
		servers: make(map[string]*ServerInstance),
		configs: DefaultServers(),
	}
}

// OnDiagnostics sets the callback for diagnostic notifications.
func (m *Manager) OnDiagnostics(callback func(uri DocumentURI, diagnostics []Diagnostic)) {
	m.onDiagnostics = callback
}

// OnStatus sets the callback for status messages.
func (m *Manager) OnStatus(callback func(message string)) {
	m.onStatus = callback
}

// OnError sets the callback for errors.
func (m *Manager) OnError(callback func(err error)) {
	m.onError = callback
}

// OnInvalidate sets the callback for UI invalidation.
func (m *Manager) OnInvalidate(callback func()) {
	m.onInvalidate = callback
}

// GetServerForFile returns or creates a server instance for a file.
func (m *Manager) GetServerForFile(filePath string) (*ServerInstance, error) {
	cfg := GetConfigForFile(filePath)
	if cfg == nil {
		return nil, fmt.Errorf("no language server configured for this file type")
	}

	if !IsServerAvailable(cfg) {
		return nil, fmt.Errorf("language server %s not installed", cfg.Name)
	}

	// Find project root
	root := FindProjectRoot(filePath, cfg.RootPatterns)
	if root == "" {
		return nil, fmt.Errorf("could not determine project root")
	}

	// Check if server already running for this root
	m.serversMu.RLock()
	if server, ok := m.servers[root]; ok {
		m.serversMu.RUnlock()
		return server, nil
	}
	m.serversMu.RUnlock()

	// Create new server instance
	return m.startServer(*cfg, root)
}

// startServer starts a new language server.
func (m *Manager) startServer(cfg ServerConfig, workspaceRoot string) (*ServerInstance, error) {
	if m.onStatus != nil {
		m.onStatus(fmt.Sprintf("Starting %s...", cfg.Name))
	}

	// Get full path to server command (handles ~/go/bin for gopls, etc.)
	cmdPath := FindServerCommand(&cfg)
	if cmdPath == "" {
		return nil, fmt.Errorf("language server %s not found", cfg.Name)
	}

	client, err := NewClient(cmdPath, cfg.Args, workspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	if err := client.Start(); err != nil {
		return nil, fmt.Errorf("start client: %w", err)
	}

	server := &ServerInstance{
		Config:        cfg,
		Client:        client,
		WorkspaceRoot: workspaceRoot,
		Documents:     make(map[DocumentURI]*DocumentState),
		Diagnostics:   make(map[DocumentURI][]Diagnostic),
	}

	// Set up error handling
	client.SetErrorCallback(func(err error) {
		if m.onError != nil {
			m.onError(fmt.Errorf("%s: %v", cfg.Name, err))
		}
	})

	// Register notification handlers before initialization
	client.OnNotification("textDocument/publishDiagnostics", func(method string, params json.RawMessage) {
		var diagParams PublishDiagnosticsParams
		if err := json.Unmarshal(params, &diagParams); err != nil {
			return
		}

		// Store diagnostics
		server.DiagnosticsMu.Lock()
		server.Diagnostics[diagParams.URI] = diagParams.Diagnostics
		server.DiagnosticsMu.Unlock()

		// Call callback
		if m.onDiagnostics != nil {
			m.onDiagnostics(diagParams.URI, diagParams.Diagnostics)
		}

		// Request UI refresh
		if m.onInvalidate != nil {
			m.onInvalidate()
		}
	})

	// Initialize the server
	initParams := InitializeParams{
		ProcessID: os.Getpid(),
		RootURI:   FilePathToURI(workspaceRoot),
		Capabilities: ClientCapabilities{
			Workspace: &WorkspaceClientCapabilities{
				ApplyEdit: true,
				WorkspaceEdit: &WorkspaceEditCapabilities{
					DocumentChanges: true,
				},
				WorkspaceFolders: true,
			},
			TextDocument: &TextDocumentClientCapabilities{
				Synchronization: &TextDocumentSyncClientCapabilities{
					DynamicRegistration: false,
					WillSave:            false,
					WillSaveWaitUntil:   false,
					DidSave:             true,
				},
				Completion: &CompletionClientCapabilities{
					DynamicRegistration: false,
					CompletionItem: &CompletionItemClientCaps{
						SnippetSupport:          true,
						CommitCharactersSupport: true,
						DocumentationFormat:     []string{string(MarkupKindMarkdown), string(MarkupKindPlainText)},
						DeprecatedSupport:       true,
						PreselectSupport:        true,
					},
					ContextSupport: true,
				},
				Hover: &HoverClientCapabilities{
					DynamicRegistration: false,
					ContentFormat:       []MarkupKind{MarkupKindMarkdown, MarkupKindPlainText},
				},
				Definition: &DefinitionClientCapabilities{
					DynamicRegistration: false,
					LinkSupport:         true,
				},
				References: &ReferencesClientCapabilities{
					DynamicRegistration: false,
				},
				CodeAction: &CodeActionClientCapabilities{
					DynamicRegistration: false,
					CodeActionLiteralSupport: &CodeActionLiteralSupport{
						CodeActionKind: &CodeActionKindLiteralSupport{
							ValueSet: []CodeActionKind{
								CodeActionKindQuickFix,
								CodeActionKindRefactor,
								CodeActionKindRefactorExtract,
								CodeActionKindRefactorInline,
								CodeActionKindRefactorRewrite,
								CodeActionKindSource,
								CodeActionKindSourceOrganizeImports,
								CodeActionKindSourceFixAll,
							},
						},
					},
					IsPreferredSupport: true,
				},
				Formatting: &FormattingClientCapabilities{
					DynamicRegistration: false,
				},
				Rename: &RenameClientCapabilities{
					DynamicRegistration: false,
					PrepareSupport:      true,
				},
				PublishDiagnostics: &PublishDiagnosticsClientCaps{
					RelatedInformation: true,
					TagSupport:         true,
					VersionSupport:     true,
				},
			},
		},
		WorkspaceFolders: []WorkspaceFolder{
			{
				URI:  FilePathToURI(workspaceRoot),
				Name: workspaceRoot,
			},
		},
	}

	var initResult InitializeResult
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.Call(ctx, "initialize", initParams, &initResult); err != nil {
		client.Stop()
		return nil, fmt.Errorf("initialize: %w", err)
	}

	// Store capabilities
	server.Capabilities = initResult.Capabilities
	client.SetCapabilities(initResult.Capabilities)

	// Send initialized notification
	if err := client.Notify("initialized", struct{}{}); err != nil {
		client.Stop()
		return nil, fmt.Errorf("initialized notification: %w", err)
	}

	// Store server
	m.serversMu.Lock()
	m.servers[workspaceRoot] = server
	m.serversMu.Unlock()

	// Start server monitor
	go m.monitorServer(server)

	if m.onStatus != nil {
		m.onStatus(fmt.Sprintf("%s ready", cfg.Name))
	}

	return server, nil
}

// monitorServer monitors a server for crashes.
func (m *Manager) monitorServer(server *ServerInstance) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !server.Client.IsRunning() {
			m.serversMu.Lock()
			delete(m.servers, server.WorkspaceRoot)
			m.serversMu.Unlock()

			if m.onStatus != nil {
				m.onStatus(fmt.Sprintf("%s stopped", server.Config.Name))
			}

			if err := server.Client.GetLastError(); err != nil && m.onError != nil {
				m.onError(err)
			}

			return
		}
	}
}

// StopAll stops all language servers.
func (m *Manager) StopAll() {
	m.serversMu.Lock()
	defer m.serversMu.Unlock()

	for root, server := range m.servers {
		server.Client.Stop()
		delete(m.servers, root)
	}
}

// StopServer stops a specific server by workspace root.
func (m *Manager) StopServer(workspaceRoot string) {
	m.serversMu.Lock()
	defer m.serversMu.Unlock()

	if server, ok := m.servers[workspaceRoot]; ok {
		server.Client.Stop()
		delete(m.servers, workspaceRoot)
	}
}

// RestartServerForFile restarts the server for a given file.
func (m *Manager) RestartServerForFile(filePath string) (*ServerInstance, error) {
	cfg := GetConfigForFile(filePath)
	if cfg == nil {
		return nil, fmt.Errorf("no language server configured for this file type")
	}

	root := FindProjectRoot(filePath, cfg.RootPatterns)

	// Stop existing server
	m.StopServer(root)

	// Start new server
	return m.GetServerForFile(filePath)
}

// GetDiagnostics returns diagnostics for a file.
func (m *Manager) GetDiagnostics(filePath string) []Diagnostic {
	uri := FilePathToURI(filePath)

	m.serversMu.RLock()
	defer m.serversMu.RUnlock()

	for _, server := range m.servers {
		server.DiagnosticsMu.RLock()
		if diags, ok := server.Diagnostics[uri]; ok {
			server.DiagnosticsMu.RUnlock()
			return diags
		}
		server.DiagnosticsMu.RUnlock()
	}

	return nil
}

// GetAllDiagnostics returns all diagnostics from all servers.
func (m *Manager) GetAllDiagnostics() map[string][]Diagnostic {
	result := make(map[string][]Diagnostic)

	m.serversMu.RLock()
	defer m.serversMu.RUnlock()

	for _, server := range m.servers {
		server.DiagnosticsMu.RLock()
		for uri, diags := range server.Diagnostics {
			filePath := URIToFilePath(uri)
			result[filePath] = diags
		}
		server.DiagnosticsMu.RUnlock()
	}

	return result
}

// GetServerStatus returns status information about running servers.
func (m *Manager) GetServerStatus() []string {
	m.serversMu.RLock()
	defer m.serversMu.RUnlock()

	status := make([]string, 0, len(m.servers))
	for root, server := range m.servers {
		running := "stopped"
		if server.Client.IsRunning() {
			running = "running"
		}
		status = append(status, fmt.Sprintf("%s (%s): %s", server.Config.Name, root, running))
	}

	return status
}

// GetDetailedStatus returns detailed status for all servers.
func (m *Manager) GetDetailedStatus() []LSPServerStatus {
	m.serversMu.RLock()
	defer m.serversMu.RUnlock()

	result := make([]LSPServerStatus, 0, len(m.servers))
	for root, server := range m.servers {
		status := LSPServerStatus{
			Name:          server.Config.Name,
			WorkspaceRoot: root,
			Running:       server.Client.IsRunning(),
		}

		// Get language from file extensions
		if len(server.Config.FileExtensions) > 0 {
			status.Language = server.Config.FileExtensions[0]
		}

		// Count diagnostics
		server.DiagnosticsMu.RLock()
		for _, diags := range server.Diagnostics {
			status.DiagCount += len(diags)
		}
		server.DiagnosticsMu.RUnlock()

		// Count open documents
		server.DocumentsMu.RLock()
		status.DocCount = len(server.Documents)
		server.DocumentsMu.RUnlock()

		result = append(result, status)
	}

	return result
}

// GetDiagnosticsSummary returns counts of errors, warnings, and info diagnostics.
func (m *Manager) GetDiagnosticsSummary() (errors, warnings, info int) {
	m.serversMu.RLock()
	defer m.serversMu.RUnlock()

	for _, server := range m.servers {
		server.DiagnosticsMu.RLock()
		for _, diags := range server.Diagnostics {
			for _, d := range diags {
				switch d.Severity {
				case 1: // Error
					errors++
				case 2: // Warning
					warnings++
				case 3, 4: // Information, Hint
					info++
				}
			}
		}
		server.DiagnosticsMu.RUnlock()
	}
	return
}

// HasServerForFile returns true if there's a server available for the file.
func (m *Manager) HasServerForFile(filePath string) bool {
	cfg := GetConfigForFile(filePath)
	return cfg != nil && IsServerAvailable(cfg)
}

// IsServerRunningForFile returns true if a server is running for the file.
func (m *Manager) IsServerRunningForFile(filePath string) bool {
	cfg := GetConfigForFile(filePath)
	if cfg == nil {
		return false
	}

	root := FindProjectRoot(filePath, cfg.RootPatterns)

	m.serversMu.RLock()
	defer m.serversMu.RUnlock()

	if server, ok := m.servers[root]; ok {
		return server.Client.IsRunning()
	}

	return false
}

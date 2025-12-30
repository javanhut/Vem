package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Client is an LSP client that communicates with a language server via JSON-RPC 2.0.
type Client struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	stderr io.ReadCloser

	// Request/response tracking
	nextID      int64
	pendingReqs map[int64]chan *Response
	pendingMu   sync.Mutex

	// Notification handlers
	handlers   map[string]NotificationHandler
	handlersMu sync.RWMutex

	// Lifecycle
	ctx       context.Context
	cancel    context.CancelFunc
	running   bool
	runningMu sync.RWMutex
	lastErr   error
	lastErrMu sync.RWMutex

	// Server capabilities (after initialization)
	capabilities   ServerCapabilities
	capabilitiesMu sync.RWMutex

	// Error callback
	onError func(error)

	writeMu sync.Mutex

	stderrMu       sync.RWMutex
	lastStderrLine string
}

// Request represents a JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      *int64      `json:"id,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *int64          `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

// ResponseError represents an error in a JSON-RPC response.
type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("LSP error %d: %s", e.Code, e.Message)
}

// NotificationHandler processes server notifications.
type NotificationHandler func(method string, params json.RawMessage)

// NewClient creates a new LSP client for the given command.
func NewClient(command string, args []string, workDir string) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = workDir
	cmd.Env = os.Environ()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("stderr pipe: %w", err)
	}

	return &Client{
		cmd:         cmd,
		stdin:       stdin,
		stdout:      bufio.NewReader(stdout),
		stderr:      stderr,
		pendingReqs: make(map[int64]chan *Response),
		handlers:    make(map[string]NotificationHandler),
		ctx:         ctx,
		cancel:      cancel,
	}, nil
}

// Start starts the language server process.
func (c *Client) Start() error {
	c.runningMu.Lock()
	defer c.runningMu.Unlock()

	if c.running {
		return fmt.Errorf("client already running")
	}

	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	c.running = true

	// Start reader goroutine
	go c.readLoop()

	// Start stderr reader (for debugging/logging)
	go c.readStderr()

	// Monitor process exit
	go c.monitorProcess()

	return nil
}

// Stop stops the language server gracefully.
func (c *Client) Stop() error {
	c.runningMu.Lock()
	if !c.running {
		c.runningMu.Unlock()
		return nil
	}
	c.running = false
	c.runningMu.Unlock()

	// Send shutdown request (best effort, with timeout)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer shutdownCancel()
	_ = c.Call(shutdownCtx, "shutdown", nil, nil)

	// Send exit notification
	_ = c.Notify("exit", nil)

	// Close stdin to signal server
	c.stdin.Close()

	// Cancel context
	c.cancel()

	// Wait for process with timeout
	done := make(chan error, 1)
	go func() {
		done <- c.cmd.Wait()
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		c.cmd.Process.Kill()
	}

	return nil
}

// IsRunning returns whether the client is running.
func (c *Client) IsRunning() bool {
	c.runningMu.RLock()
	defer c.runningMu.RUnlock()
	return c.running
}

// GetLastError returns the last error that occurred.
func (c *Client) GetLastError() error {
	c.lastErrMu.RLock()
	defer c.lastErrMu.RUnlock()
	return c.lastErr
}

func (c *Client) setLastStderrLine(line string) {
	c.stderrMu.Lock()
	c.lastStderrLine = line
	c.stderrMu.Unlock()
}

func (c *Client) getLastStderrLine() string {
	c.stderrMu.RLock()
	defer c.stderrMu.RUnlock()
	return c.lastStderrLine
}

// SetErrorCallback sets the callback for errors.
func (c *Client) SetErrorCallback(callback func(error)) {
	c.onError = callback
}

// Call sends a request and waits for a response.
func (c *Client) Call(ctx context.Context, method string, params interface{}, result interface{}) error {
	c.runningMu.RLock()
	if !c.running {
		c.runningMu.RUnlock()
		return fmt.Errorf("client not running")
	}
	c.runningMu.RUnlock()

	id := atomic.AddInt64(&c.nextID, 1)

	respChan := make(chan *Response, 1)
	c.pendingMu.Lock()
	c.pendingReqs[id] = respChan
	c.pendingMu.Unlock()

	defer func() {
		c.pendingMu.Lock()
		delete(c.pendingReqs, id)
		c.pendingMu.Unlock()
	}()

	// Send request
	req := Request{
		JSONRPC: "2.0",
		ID:      &id,
		Method:  method,
		Params:  params,
	}

	if err := c.writeMessage(req); err != nil {
		return fmt.Errorf("write request: %w", err)
	}

	// Wait for response
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.ctx.Done():
		return fmt.Errorf("client stopped")
	case resp, ok := <-respChan:
		if !ok || resp == nil {
			if err := c.GetLastError(); err != nil {
				return err
			}
			return fmt.Errorf("client stopped")
		}
		if resp.Error != nil {
			return resp.Error
		}
		if result != nil && resp.Result != nil {
			return json.Unmarshal(resp.Result, result)
		}
		return nil
	}
}

// Notify sends a notification (no response expected).
func (c *Client) Notify(method string, params interface{}) error {
	c.runningMu.RLock()
	if !c.running {
		c.runningMu.RUnlock()
		return fmt.Errorf("client not running")
	}
	c.runningMu.RUnlock()

	req := Request{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	return c.writeMessage(req)
}

// OnNotification registers a handler for server notifications.
func (c *Client) OnNotification(method string, handler NotificationHandler) {
	c.handlersMu.Lock()
	c.handlers[method] = handler
	c.handlersMu.Unlock()
}

// Capabilities returns the negotiated server capabilities.
func (c *Client) Capabilities() ServerCapabilities {
	c.capabilitiesMu.RLock()
	defer c.capabilitiesMu.RUnlock()
	return c.capabilities
}

// SetCapabilities sets the server capabilities after initialization.
func (c *Client) SetCapabilities(caps ServerCapabilities) {
	c.capabilitiesMu.Lock()
	c.capabilities = caps
	c.capabilitiesMu.Unlock()
}

// writeMessage writes an LSP message with Content-Length header.
func (c *Client) writeMessage(msg interface{}) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))

	// Write header and body atomically
	_, err = c.stdin.Write([]byte(header))
	if err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	_, err = c.stdin.Write(data)
	if err != nil {
		return fmt.Errorf("write body: %w", err)
	}

	return nil
}

// readLoop continuously reads messages from the server.
func (c *Client) readLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		msg, err := c.readMessage()
		if err != nil {
			c.lastErrMu.Lock()
			c.lastErr = err
			c.lastErrMu.Unlock()

			if c.onError != nil {
				c.onError(err)
			}

			// Check if we should stop
			c.runningMu.RLock()
			running := c.running
			c.runningMu.RUnlock()
			if !running {
				return
			}

			// EOF or pipe closed indicates server exited
			if err == io.EOF || strings.Contains(err.Error(), "closed") {
				return
			}

			continue
		}

		c.handleMessage(msg)
	}
}

// readMessage reads a single LSP message.
func (c *Client) readMessage() (json.RawMessage, error) {
	// Read headers
	var contentLength int

	for {
		line, err := c.stdout.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break // End of headers
		}

		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				length, err := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err == nil {
					contentLength = length
				}
			}
		}
		// Ignore other headers like Content-Type
	}

	if contentLength == 0 {
		return nil, fmt.Errorf("missing or invalid Content-Length header")
	}

	// Read body
	body := make([]byte, contentLength)
	_, err := io.ReadFull(c.stdout, body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return body, nil
}

// handleMessage routes an incoming message to the appropriate handler.
func (c *Client) handleMessage(data json.RawMessage) {
	// Try to determine message type
	var msg struct {
		ID     *int64          `json:"id"`
		Method string          `json:"method"`
		Result json.RawMessage `json:"result"`
		Error  *ResponseError  `json:"error"`
		Params json.RawMessage `json:"params"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	// Response (has ID and result/error, no method)
	if msg.ID != nil && msg.Method == "" {
		c.pendingMu.Lock()
		if ch, ok := c.pendingReqs[*msg.ID]; ok {
			resp := &Response{
				JSONRPC: "2.0",
				ID:      msg.ID,
				Result:  msg.Result,
				Error:   msg.Error,
			}
			select {
			case ch <- resp:
			default:
			}
		}
		c.pendingMu.Unlock()
		return
	}

	// Notification (has method, no ID)
	if msg.Method != "" && msg.ID == nil {
		c.handlersMu.RLock()
		handler, ok := c.handlers[msg.Method]
		c.handlersMu.RUnlock()

		if ok {
			// Handle in goroutine to not block read loop
			go handler(msg.Method, msg.Params)
		}
		return
	}

	// Request from server (has method and ID) - we should respond
	// For now, we don't handle server-initiated requests
	if msg.Method != "" && msg.ID != nil {
		// Send empty response to acknowledge
		resp := Response{
			JSONRPC: "2.0",
			ID:      msg.ID,
			Result:  json.RawMessage("null"),
		}
		_ = c.writeMessage(resp)
	}
}

// readStderr reads and discards stderr output.
// This prevents the server from blocking on stderr writes.
func (c *Client) readStderr() {
	scanner := bufio.NewScanner(c.stderr)
	for scanner.Scan() {
		c.setLastStderrLine(scanner.Text())
	}
}

// monitorProcess monitors the server process for unexpected exits.
func (c *Client) monitorProcess() {
	err := c.cmd.Wait()

	c.runningMu.Lock()
	wasRunning := c.running
	c.running = false
	c.runningMu.Unlock()

	if wasRunning {
		c.lastErrMu.Lock()
		lastStderr := c.getLastStderrLine()
		if err != nil {
			if lastStderr != "" {
				c.lastErr = fmt.Errorf("server exited: %w (stderr: %s)", err, lastStderr)
			} else {
				c.lastErr = fmt.Errorf("server exited: %w", err)
			}
		} else {
			if lastStderr != "" {
				c.lastErr = fmt.Errorf("server exited unexpectedly (stderr: %s)", lastStderr)
			} else {
				c.lastErr = fmt.Errorf("server exited unexpectedly")
			}
		}
		c.lastErrMu.Unlock()

		if c.onError != nil {
			c.onError(c.lastErr)
		}

		// Cancel any pending requests
		c.pendingMu.Lock()
		for id, ch := range c.pendingReqs {
			close(ch)
			delete(c.pendingReqs, id)
		}
		c.pendingMu.Unlock()
	}
}

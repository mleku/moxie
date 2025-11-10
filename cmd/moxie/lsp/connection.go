// Copyright 2024 The Moxie Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
)

// Connection represents a JSON-RPC 2.0 connection
type Connection struct {
	reader   *bufio.Reader
	writer   io.Writer
	handlers map[string]Handler
	mu       sync.Mutex
	logger   *log.Logger
}

// Handler is a function that handles LSP requests
type Handler func(ctx context.Context, req *Request) (interface{}, error)

// Request represents a JSON-RPC 2.0 request
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC 2.0 error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Notification represents a JSON-RPC 2.0 notification
type Notification struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// NewConnection creates a new JSON-RPC connection
func NewConnection(r io.Reader, w io.Writer, logger *log.Logger) *Connection {
	return &Connection{
		reader:   bufio.NewReader(r),
		writer:   w,
		handlers: make(map[string]Handler),
		logger:   logger,
	}
}

// Handle registers a handler for a method
func (c *Connection) Handle(method string, handler Handler) {
	c.mu.Lock()
	c.handlers[method] = handler
	c.mu.Unlock()
}

// Run runs the connection loop
func (c *Connection) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Read message
		msg, err := c.readMessage()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			c.logger.Printf("Read error: %v", err)
			continue
		}

		// Parse request
		var req Request
		if err := json.Unmarshal(msg, &req); err != nil {
			c.logger.Printf("Parse error: %v", err)
			continue
		}

		// Handle request
		go c.handleRequest(ctx, &req)
	}
}

// readMessage reads a JSON-RPC message using the LSP header format
func (c *Connection) readMessage() ([]byte, error) {
	// Read headers
	var contentLength int
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			// Empty line marks end of headers
			break
		}

		// Parse Content-Length header
		if strings.HasPrefix(line, "Content-Length: ") {
			lengthStr := strings.TrimPrefix(line, "Content-Length: ")
			contentLength, err = strconv.Atoi(lengthStr)
			if err != nil {
				return nil, fmt.Errorf("invalid Content-Length: %v", err)
			}
		}
	}

	if contentLength == 0 {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	// Read content
	content := make([]byte, contentLength)
	_, err := io.ReadFull(c.reader, content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// writeMessage writes a JSON-RPC message using the LSP header format
func (c *Connection) writeMessage(msg []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Write headers
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(msg))
	if _, err := c.writer.Write([]byte(header)); err != nil {
		return err
	}

	// Write content
	if _, err := c.writer.Write(msg); err != nil {
		return err
	}

	return nil
}

// handleRequest handles a JSON-RPC request
func (c *Connection) handleRequest(ctx context.Context, req *Request) {
	c.mu.Lock()
	handler, ok := c.handlers[req.Method]
	c.mu.Unlock()

	if !ok {
		c.logger.Printf("No handler for method: %s", req.Method)
		if req.ID != nil {
			// Send error response for requests (not notifications)
			c.sendError(req.ID, -32601, "Method not found")
		}
		return
	}

	// Call handler
	result, err := handler(ctx, req)

	// Send response (only for requests, not notifications)
	if req.ID != nil {
		if err != nil {
			c.sendError(req.ID, -32603, err.Error())
		} else {
			c.sendResult(req.ID, result)
		}
	}
}

// sendResult sends a successful response
func (c *Connection) sendResult(id interface{}, result interface{}) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		c.logger.Printf("Marshal error: %v", err)
		return
	}

	if err := c.writeMessage(data); err != nil {
		c.logger.Printf("Write error: %v", err)
	}
}

// sendError sends an error response
func (c *Connection) sendError(id interface{}, code int, message string) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		c.logger.Printf("Marshal error: %v", err)
		return
	}

	if err := c.writeMessage(data); err != nil {
		c.logger.Printf("Write error: %v", err)
	}
}

// Notify sends a notification to the client
func (c *Connection) Notify(method string, params interface{}) {
	notif := Notification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(notif)
	if err != nil {
		c.logger.Printf("Marshal error: %v", err)
		return
	}

	if err := c.writeMessage(data); err != nil {
		c.logger.Printf("Write error: %v", err)
	}
}

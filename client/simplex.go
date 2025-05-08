package client

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// SimplexRequest is a struct that represents a request to the Simplex websocket.
type SimplexRequest struct {
	Id      string `json:"corrId"` // Unique identifier for the request
	Command string `json:"cmd"`    // Command to be executed
}

// NewSimplexRequest creates a new SimplexRequest with the given command
func NewSimplexRequest(command string) *SimplexRequest {
	return &SimplexRequest{
		Id:      fmt.Sprintf("simplex-client-%08d", rand.Int32N(10000000)), // Generate a random ID for the request with padding
		Command: command,
	}
}

// NewSimplexMessageRequest creates a new SimplexRequest for sending a message
func NewSimplexMessageRequest(recipient, message string) *SimplexRequest {
	// Prepare the recipient
	prefix, recipient := PrepareRecipient(recipient)
	// Create a new SimplexRequest with the formatted message
	return NewSimplexRequest(fmt.Sprintf("%s%s %s", prefix, recipient, message))
}

// PrepareRecipient prepares the recipient string by removing any prefixes
// and returning the prefix and the cleaned recipient string
func PrepareRecipient(recipient string) (string, string) {
	switch true {
	case strings.HasPrefix(recipient, "@"):
		return "@", strings.TrimPrefix(recipient, "@")
	case strings.HasPrefix(recipient, "#"):
		return "#", strings.TrimPrefix(recipient, "#")
	}
	// Fall back to default case without any prefix
	return "@", recipient
}

// SimplexClient is a struct that represents a client for the Simplex websocket.
type SimplexClient struct {
	Address string          // WebSocket address for the simplex client to connect to
	con     *websocket.Conn // WebSocket connection to the server
	context context.Context // Context for managing the lifecycle of the client
}

// NewSimplexClient creates a new instance of SimplexClient with the given address.
func NewSimplexClient(ctx context.Context, address string) *SimplexClient {
	return &SimplexClient{
		Address: address,
		context: ctx,
	}
}

// Connect establishes a WebSocket connection
func (c *SimplexClient) Connect() error {
	// Check if the connection is already established
	if c.con == nil {
		var err error
		// Connect to the websocket server
		c.con, _, err = websocket.Dial(c.context, c.Address, nil)
		if err != nil {
			// Return wrapped error if connection fails
			return fmt.Errorf("failed to connect to websocket: %w", err)
		}
	}
	// Return no errors
	return nil
}

// Close closes the WebSocket connection
func (c *SimplexClient) Close() error {
	if c.con != nil {
		err := c.con.Close(websocket.StatusNormalClosure, "simplex client closed")
		// If connection close fails, return the error and don't set con to nil
		if err != nil {
			return fmt.Errorf("failed to close websocket: %w", err)
		}
		// Set the connection to nil after closing (to avoid reusing the closed connection)
		c.con = nil
	}
	// Return no errors
	return nil
}

// Send sends a SimplexRequest to the websocket server
func (c *SimplexClient) Send(req *SimplexRequest) error {
	// Check if connection is nil before sending any messages
	err := c.Connect()
	if err != nil {
		// Return already wrapped error if connection fails
		return err
	}
	// Write the request to the WebSocket connection
	err = wsjson.Write(c.context, c.con, req)
	if err != nil {
		// Return wrapped error if sending message fails
		return fmt.Errorf("failed to send message: %w", err)
	}
	// Return no errors
	return nil
}

// Detects if a group already exists and client is a member of it
// TODO: Implement this function
func (c *SimplexClient) IsGroupMember(group string) bool {
	return false
}

// Detects if a contact already exists
// TODO: Implement this function
func (c *SimplexClient) IsContact(contact string) bool {
	return false
}

// SendMessage sends a message to a specific recipient
func (c *SimplexClient) SendMessage(recipient string, message string) error {
	slog.Info("Sending message", "recipient", recipient, "message", message)
	return c.Send(NewSimplexMessageRequest(recipient, message))
}

// ChangeDisplayName changes the display name of the client
func (c *SimplexClient) ChangeDisplayName(name string) error {
	slog.Info("Changing display name", "name", name)
	return c.Send(NewSimplexRequest(fmt.Sprintf("/p %s", name)))
}

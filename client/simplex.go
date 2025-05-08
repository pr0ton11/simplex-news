package client

import (
	"context"
	"fmt"
	"strings"

	"math/rand/v2"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// SimplexRequest is a struct that represents a request to the Simplex websocket.
type SimplexRequest struct {
	Id      string `json:"corrId"` // Unique identifier for the request
	Command string `json:"cmd"`    // Command to be executed
}

func NewSimplexRequest(command string) *SimplexRequest {
	return &SimplexRequest{
		Id:      fmt.Sprintf("id-%d", rand.Int32N(1000000)), // Generate a random ID for the request
		Command: command,
	}
}

// SimplexClient is a struct that represents a client for the Simplex websocket.
type SimplexClient struct {
	Address string          // WebSocket address for the simplex client to connect to
	con     *websocket.Conn // WebSocket connection to the server
	context context.Context // Context for managing the lifecycle of the client
}

// NewSimplexClient creates a new instance of SimplexClient with the given address.
func NewSimplexClient(address string, ctx context.Context) *SimplexClient {
	return &SimplexClient{
		Address: address,
		context: ctx,
	}
}

// Connect establishes a WebSocket connection
func (c *SimplexClient) Connect() error {
	conn, _, err := websocket.Dial(c.context, c.Address, nil)
	if err != nil {
		return err
	}
	c.con = conn
	return nil
}

// Close closes the WebSocket connection
func (c *SimplexClient) Close() error {
	if c.con != nil {
		err := c.con.Close(websocket.StatusNormalClosure, "Disconnection of Simplex client")
		if err != nil {
			return err
		}
		c.con = nil
	}
	return nil
}

// Send sends a message to the WebSocket server
func (c *SimplexClient) Send(req *SimplexRequest) error {
	if c.con == nil {
		return fmt.Errorf("websocket connection is not established")
	}
	err := wsjson.Write(c.context, c.con, req)
	if err != nil {
		return err
	}
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
	prefix := "@" // Default prefix for contacts
	switch true {
	case strings.HasPrefix(recipient, "@"):
		// If the recipient starts with '@', it's a user
		prefix = "@"
		// Remove the '@' from the recipient
		recipient = strings.TrimPrefix(recipient, "@")
	case strings.HasPrefix(recipient, "#"):
		// If the recipient starts with '#', it's a group
		prefix = "#"
		// Remove the '#' from the recipient
		recipient = strings.TrimPrefix(recipient, "#")
	default:
		// Detect if the recipient is a group or user by checking for existing groups and chats
		if c.IsGroupMember(recipient) {
			// Switch prefix to '#' if the recipient is a confirmed group
			prefix = "#"
		}
	}

	return c.Send(NewSimplexRequest(fmt.Sprintf("%s%s %s", prefix, recipient, message)))
}

// ChangeDisplayName changes the display name of the client
func (c *SimplexClient) ChangeDisplayName(name string) error {
	return c.Send(NewSimplexRequest(fmt.Sprintf("/p %s", name)))
}

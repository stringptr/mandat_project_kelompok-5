package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type MailTrapMessage struct {
	ID             int    `json:"id"`
	Subject        string `json:"subject"`
	ToEmail        string `json:"to_email"`
	CreatedAt      string `json:"created_at"`
	HTMLSourcePath string `json:"html_source_path"`
}

type mailTrapMessageListItem struct {
	ID        int    `json:"id"`
	Subject   string `json:"subject"`
	ToEmail   string `json:"to_email"`
	CreatedAt string `json:"created_at"`
}

type MailTrapClient struct {
	apiToken   string
	accountID  string
	inboxID    string
	baseURL    string
	httpClient *http.Client
}

func NewMailTrapClient(apiToken, accountID, inboxID string) *MailTrapClient {
	return &MailTrapClient{
		apiToken:   apiToken,
		accountID:  accountID,
		inboxID:    inboxID,
		baseURL:    "https://mailtrap.io/api",
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *MailTrapClient) doRequest(t *testing.T, method, path string) []byte {
	t.Helper()
	req, err := http.NewRequest(method, c.baseURL+path, nil)
	if err != nil {
		t.Fatalf("MailTrap: failed to create request: %v", err)
	}
	req.Header.Set("Api-Token", c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		t.Fatalf("MailTrap: request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("MailTrap: failed to read response body: %v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("MailTrap: %s %s returned %d: %s", method, path, resp.StatusCode, string(body))
	}
	return body
}

func (c *MailTrapClient) GetMessages(t *testing.T) []MailTrapMessage {
	t.Helper()
	path := fmt.Sprintf("/accounts/%s/inboxes/%s/messages", c.accountID, c.inboxID)
	body := c.doRequest(t, http.MethodGet, path)

	var items []mailTrapMessageListItem
	if err := json.Unmarshal(body, &items); err != nil {
		t.Fatalf("MailTrap: failed to decode messages: %v", err)
	}

	messages := make([]MailTrapMessage, len(items))
	for i, item := range items {
		messages[i] = MailTrapMessage{
			ID:        item.ID,
			Subject:   item.Subject,
			ToEmail:   item.ToEmail,
			CreatedAt: item.CreatedAt,
		}
	}
	return messages
}

func (c *MailTrapClient) GetMessageDetail(t *testing.T, messageID int) MailTrapMessage {
	t.Helper()
	path := fmt.Sprintf("/accounts/%s/inboxes/%s/messages/%d", c.accountID, c.inboxID, messageID)
	body := c.doRequest(t, http.MethodGet, path)

	var msg MailTrapMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		t.Fatalf("MailTrap: failed to decode message detail: %v", err)
	}
	return msg
}

func (c *MailTrapClient) GetHTMLBody(t *testing.T, msg MailTrapMessage) string {
	t.Helper()
	path := strings.TrimPrefix(msg.HTMLSourcePath, "/api")
	body := c.doRequest(t, http.MethodGet, path)
	return string(body)
}

func (c *MailTrapClient) DeleteMessage(t *testing.T, messageID int) {
	t.Helper()
	path := fmt.Sprintf("/accounts/%s/inboxes/%s/messages/%d", c.accountID, c.inboxID, messageID)
	c.doRequest(t, http.MethodDelete, path)
}

func (c *MailTrapClient) WaitForNewMessage(t *testing.T, initialCount int, timeout time.Duration) MailTrapMessage {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		messages := c.GetMessages(t)
		if len(messages) > initialCount {
			newMsg := c.GetMessageDetail(t, messages[0].ID)
			return newMsg
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("MailTrap: timeout waiting for new message after %v", timeout)
	return MailTrapMessage{}
}

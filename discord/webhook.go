package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Red color for Discord embeds (danger/alert).
const ColorRed = 16711680

// AlertType is the same as parser.AlertType; we use string to avoid coupling.
const (
	AlertTypeBinaryExec = "BINARY_EXEC"
	AlertTypeFileRead   = "FILE_READ"
)

// EmbedField represents a single field in a Discord embed.
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// DiscordEmbed is the payload structure for a Discord webhook embed.
type DiscordEmbed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Color       int          `json:"color"`
	Fields      []EmbedField `json:"fields,omitempty"`
}

// WebhookPayload is the top-level body Discord expects.
type WebhookPayload struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

// SendAlert builds a Discord embed from the alert and POSTs it to the webhook URL.
// alertType is BINARY_EXEC or FILE_READ; filePath is only used for FILE_READ.
func SendAlert(webhookURL string, alertType string, binary string, args string, pid uint32, filePath string) error {
	var title string
	var description string
	var fields []EmbedField

	switch alertType {
	case AlertTypeFileRead:
		title = "🚨 UNAUTHORIZED FILE ACCESS DETECTED"
		description = "Tetragon Sentinel intercepted an attempt to read a protected file."
		fields = []EmbedField{
			{Name: "File path", Value: filePath, Inline: false},
			{Name: "Process", Value: binary, Inline: true},
			{Name: "PID", Value: fmt.Sprintf("%d", pid), Inline: true},
			{Name: "Arguments", Value: truncate(args, 1024), Inline: false},
		}
	default:
		title = "🚨 eBPF Kernel Alert"
		description = "Tetragon Sentinel detected a monitored binary execution."
		fields = []EmbedField{
			{Name: "Binary", Value: binary, Inline: true},
			{Name: "PID", Value: fmt.Sprintf("%d", pid), Inline: true},
			{Name: "Arguments", Value: truncate(args, 1024), Inline: false},
		}
	}

	payload := WebhookPayload{
		Embeds: []DiscordEmbed{
			{
				Title:       title,
				Description: description,
				Color:       ColorRed,
				Fields:      fields,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal webhook payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

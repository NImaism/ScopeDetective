package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NImaism/ScopeDetective/model"
	"net/http"
	"time"
)

type Messager struct {
	Webhook string
}

// NewMessager function creates and returns a new instance of the Messager struct.
func NewMessager(webhook string) *Messager {
	return &Messager{Webhook: webhook}
}

// sendMessage function sends a message to a Discord channel using a webhook.
func (m *Messager) sendMessage(message model.Message) {
	msg := model.DiscordMessage{
		Content:   "",
		Username:  "ScopeDetective",
		AvatarUrl: "https://media.discordapp.net/attachments/996196305711943801/1144225219880423464/logo.png?width=631&height=631",
		Embeds: []model.DiscordEmbed{
			{
				Title:       message.Owner,
				Description: fmt.Sprintf("```yaml\n - 💣 Max Serverity: %s \n - 🏷 Url: %s \n```", message.MaxSeverity, message.SubDomain),
				Url:         message.Url,
				Color:       0xADD8E6,
				Timestamp:   time.Now().Format(time.RFC3339),
			},
		},
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", m.Webhook, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		return
	}
}

func (m *Messager) sendLog(message string) {
	msg := model.DiscordMessage{
		Content:   "",
		Username:  "ScopeDetective",
		AvatarUrl: "https://media.discordapp.net/attachments/996196305711943801/1144225219880423464/logo.png?width=631&height=631",
		Embeds: []model.DiscordEmbed{
			{
				Title:       "",
				Description: message,
				Color:       0xADD8E6,
				Timestamp:   time.Now().Format(time.RFC3339),
			},
		},
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", m.Webhook, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if _, err = client.Do(req); err != nil {
		return
	}
}

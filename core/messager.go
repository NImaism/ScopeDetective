package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Messager struct {
	Webhook string
}

type DiscordMessage struct {
	Content   string         `json:"content,omitempty"`
	Username  string         `json:"username,omitempty"`
	AvatarUrl string         `json:"avatar_url,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Url         string              `json:"url,omitempty"`
	Color       int                 `json:"color,omitempty"`
	Timestamp   string              `json:"timestamp,omitempty"`
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// NewMessager function creates and returns a new instance of the Messager struct.
func NewMessager(webhook string) *Messager {
	return &Messager{Webhook: webhook}
}

// sendMessage function sends a message to a Discord channel using a webhook.
func (m *Messager) sendMessage(message Message) {
	msg := DiscordMessage{
		Content:   "",
		Username:  "ScopeDetective",
		AvatarUrl: "https://cdn.discordapp.com/attachments/996196305711943801/1133578939684638780/manja-vitolic-gKXKBY-C-Dk-unsplash-scaled.jpg",
		Embeds: []DiscordEmbed{
			{
				Title:       message.Owner,
				Description: "",
				Url:         message.Url,
				Color:       0xADD8E6,
				Timestamp:   time.Now().Format(time.RFC3339),
				Fields: []DiscordEmbedField{
					{
						Name:   "Address",
						Value:  message.SubDomain,
						Inline: false,
					},
				},
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

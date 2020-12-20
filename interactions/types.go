package interactions

import "time"

type InteractionType int

const (
	_ InteractionType = iota
	Ping
	ApplicationCommand
)

type InteractionResponseType int

const (
	_ InteractionResponseType = iota
	Pong
	Acknowledge
	ChannelMessage
	ChannelMessageWithSource
	AcknowledgeWithSource
)

type InteractionResponseFlags int64

const Ephemeral InteractionResponseFlags = 1 << 6

type Data struct {
	Type   InteractionType `json:"type"`
	Token  string          `json:"token"`
	Member struct {
		User struct {
			ID            int64  `json:"id"`
			Username      string `json:"username"`
			Avatar        string `json:"avatar"`
			Discriminator string `json:"discriminator"`
			PublicFlags   int64  `json:"public_flags"`
		} `json:"user"`
		Roles        []string    `json:"roles"`
		PremiumSince time.Time   `json:"premium_since"`
		Permissions  string      `json:"permissions"`
		Pending      bool        `json:"pending"`
		Nick         interface{} `json:"nick"`
		Mute         bool        `json:"mute"`
		JoinedAt     time.Time   `json:"joined_at"`
		IsPending    bool        `json:"is_pending"`
		Deaf         bool        `json:"deaf"`
	} `json:"member"`
	ID      string `json:"id"`
	GuildID string `json:"guild_id"`
	Data    struct {
		Options []ApplicationCommandInteractionDataOption `json:"options"`
		Name    string                                    `json:"name"`
		ID      string                                    `json:"id"`
	} `json:"data"`
	ChannelID string `json:"channel_id"`
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                    `json:"name"`
	Value   interface{}                               `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

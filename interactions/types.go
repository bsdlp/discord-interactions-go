package interactions

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

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
			ID            string `json:"id"`
			Username      string `json:"username"`
			Avatar        string `json:"avatar"`
			Discriminator string `json:"discriminator"`
			PublicFlags   int64  `json:"public_flags"`
		} `json:"user"`
		Roles        []string  `json:"roles"`
		PremiumSince time.Time `json:"premium_since"`
		Permissions  string    `json:"permissions"`
		Pending      bool      `json:"pending"`
		Nick         string    `json:"nick"`
		Mute         bool      `json:"mute"`
		JoinedAt     time.Time `json:"joined_at"`
		IsPending    bool      `json:"is_pending"`
		Deaf         bool      `json:"deaf"`
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

func (data *Data) ResponseURL() string {
	return fmt.Sprintf("https://discord.com/api/v8/interactions/%s/%s/callback", data.ID, data.Token)
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                    `json:"name"`
	Value   interface{}                               `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

type InteractionResponse struct {
	Type InteractionResponseType                    `json:"type"`
	Data *InteractionApplicationCommandCallbackData `json:"data,omitempty"`
}

type InteractionApplicationCommandCallbackData struct {
	TTS             *bool                             `json:"tts,omitempty"`
	Content         string                            `json:"content"`
	Embeds          []*discordgo.MessageEmbed         `json:"embeds,omitempty"`
	AllowedMentions *discordgo.MessageAllowedMentions `json:"allowed_mentions,omitempty"`
}

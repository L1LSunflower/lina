package entities

import (
	"encoding/json"
	"strings"
)

type MediaGroup struct {
	BusinessConnectionId string        `json:"business_connection_id,omitempty"`
	ChatId               string        `json:"chat_id"`
	MessageThreadId      int           `json:"message_thread_id,omitempty"`
	Media                []*MediaPhoto `json:"media"`
	DisableNotification  bool          `json:"disable_notification,omitempty"`
	ProtectContent       bool          `json:"protect_content,omitempty"`
}

func (mg *MediaGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(mg)
}

func (mg *MediaGroup) SetImages(imageLinks []string) {
	mg.Media = make([]*MediaPhoto, len(imageLinks))
	for i, imageLink := range imageLinks {
		mg.Media[i] = &MediaPhoto{
			Type:  "photo",
			Media: imageLink,
		}
	}
}

type MediaPhoto struct {
	Type            string `json:"type"`
	Media           string `json:"media"`
	Caption         string `json:"caption,omitempty"`
	ParseMode       string `json:"parse_mode,omitempty"`
	CaptionEntities string `json:"caption_entities,omitempty"`
	HasSpoiler      bool   `json:"has_spoiler,omitempty"`
}

type MessageEntity struct {
	Type          string `json:"type"`
	Offset        int    `json:"offset"`
	Length        int    `json:"length"`
	Url           string `json:"url,omitempty"`
	User          string `json:"user,omitempty"`
	Language      string `json:"language,omitempty"`
	CustomEmojiId string `json:"custom_emoji_id,omitempty"`
}

type User struct {
	Id                      string `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name,omitempty"`
	LastName                string `json:"last_name,omitempty"`
	Username                string `json:"username,omitempty"`
	LanguageCode            string `json:"language_code,omitempty"`
	IsPremium               bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
	CanConnectToBusiness    bool   `json:"can_connect_to_business,omitempty"`
}

type Message struct {
	MessageId            int           `json:"message_id" yaml:"message_id,omitempty"`
	From                 *User         `json:"from" yaml:"from,omitempty"`
	SenderChat           *Chat         `json:"sender_chat" yaml:"sender_chat,omitempty"`
	BusinessConnectionId string        `json:"business_connection_id" yaml:"business_connection_id,omitempty"`
	ChatId               string        `json:"chat_id" yaml:"chat_id"`
	MessageThreadId      int           `json:"message_thread_id,omitempty" yaml:"message_thread_id,omitempty"`
	Text                 string        `json:"text" yaml:"text"`
	ParseMode            string        `json:"parse_mode,omitempty" yaml:"parse_mode,omitempty"`
	Entities             MessageEntity `json:"entities,omitempty" yaml:"entities,omitempty"`
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

func NewMsg(item *Item, chatId string) *Message {
	sb := new(strings.Builder)
	sb.WriteString("Name: " + item.Name + "\n")
	sb.WriteString("Url: " + item.Url + "\n")
	sb.WriteString("Article: " + item.Article + "\n")
	sb.WriteString("Expected Price: " + item.Name + " " + item.Currency + "\n")
	sb.WriteString("Actual Price: " + item.Name + " " + item.Currency + "\n")
	sb.WriteString("Sizes: " + strings.Join(item.Sizes, ","+"\n"))
	return &Message{
		ChatId: chatId,
		Text:   sb.String(),
	}
}

type Chat struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	IsForum   bool   `json:"is_forum,omitempty"`
	Date      int    `json:"date,omitempty"`
}

type UpdateResult struct {
	Ok     bool `json:"ok"`
	result []Message
}

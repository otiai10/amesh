package slack

import "github.com/otiai10/spell"

// Payload ...
// https://api.slack.com/events/app_mention
type Payload struct {
	Token       string   `json:"token"`
	TeamID      string   `json:"team_id"`
	APIAppID    string   `json:"api_app_id"`
	Type        string   `json:"type"`
	ID          string   `json:"event_id"`
	Timestamp   int64    `json:"event_time"`
	AuthedUsers []string `json:"authed_users"`
	Event       Event    `json:"event"`
	// ONLY for verification.
	Challenge string `json:"challenge"`

	// Ext is an extension for amesh-bot framework
	Ext struct {
		Words spell.Words `json:"-"`
	} `json:"-"`
}

// Event ...
type Event struct {
	Type      string `json:"type"`
	User      string `json:"user"`
	Text      string `json:"text"`
	Channel   string `json:"channel"`
	Timestamp string `json:"event_ts"`
}

// Block ...
type Block struct {
	Type     string    `json:"type"`
	ImageURL string    `json:"image_url,omitempty"`
	AltText  string    `json:"alt_text,omitempty"`
	Elements []Element `json:"elements,omitempty"`
}

// Element ...
type Element struct {
	Type     string `json:"type,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	AltText  string `json:"alt_text,omitempty"`
	Text     string `json:"text,omitempty"`
}

// Message ...
type Message struct {
	Channel string `json:"channel"`
	Text    string `json:"text,omitempty"`
	// https://api.slack.com/messaging/composing/layouts#sending-messages
	Blocks []Block `json:"blocks,omitempty"`
	// Markdown bool   `json:"mrkdwn"`
}

// OAuthResponse ...
// https://api.slack.com/methods/oauth.v2.access#response
type OAuthResponse struct {
	OK         bool   `json:"ok"     firestore:"ok"`
	AppID      string `json:"app_id" firestore:"app_id"`
	AuthedUser struct {
		ID string `json:"id" firestore:"id"`
	} `json:"authed_user" firestore:"authed_user"`
	Scope       string `json:"scope"        firestore:"scope"`
	TokenType   string `json:"token_type"   firestore:"token_type"`
	AccessToken string `json:"access_token" firestore:"access_token"`
	BotUserID   string `json:"bot_user_id"  firestore:"bot_user_id"`
	Team        struct {
		ID   string `json:"id"   firestore:"id"`
		Name string `json:"name" firestore:"name"`
	} `json:"team" firestore:"team"`
	Enterprise interface{} `json:"enterprise" firestore:"-"`
}

// Team is a wrapper of OAuthResponse to manage collection on Firestore.
type Team OAuthResponse

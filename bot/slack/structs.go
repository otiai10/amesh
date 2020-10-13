package slack

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
	Event       struct {
		Type      string `json:"type"`
		User      string `json:"user"`
		Text      string `json:"text"`
		Channel   string `json:"channel"`
		Timestamp string `json:"event_ts"`
	} `json:"event"`
	// ONLY for verification.
	Challenge string `json:"challenge"`
}

// Block ...
type Block struct {
	Type     string `json:"type"`
	ImageURL string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

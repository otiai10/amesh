package platforms

// Slack ...
type Slack struct {
}

// Validate ...
func (slack *Slack) Validate() error {
	return nil
}

// Upload ...
func (slack *Slack) Upload() error {
	return nil
}

// Message ...
func (slack *Slack) Message(message string) error {
	return nil
}

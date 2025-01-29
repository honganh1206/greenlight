package mailer

import "embed"

// Embed templates directly into the compiled binary
// These templates will be read during compile time

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *Dialer
	sender string
}

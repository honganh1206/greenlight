package mailer

type authTest struct {
	auths     []string
	methods   []string // aka challenges-list of auth mechanisms the server supports
	tls       bool
	wantData  []string
	wantError bool
}

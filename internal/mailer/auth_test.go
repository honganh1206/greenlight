package mailer

import (
	"net/smtp"
	"testing"
)

type authTest struct {
	auths     []string
	methods   []string // aka challenges-list of auth mechanisms the server supports
	tls       bool
	wantData  []string
	wantError bool
}

const (
	testUser = "user"
	testPwd  = "pwd"
	testHost = "smtp.example.com"
)

func testLoginAuth(t *testing.T, test *authTest) {
	auth := &loginAuth{
		username: testUser,
		password: testPwd,
		host:     testHost,
	}
	server := &smtp.ServerInfo{
		Name: testHost,
		TLS:  test.tls,
		Auth: test.auths,
	}

	proto, toServer, err := auth.Start(server)
	if err != nil && !test.wantError {
		t.Fatalf("loginAuth.Start(): %v", err)
	}

	if err != nil && test.wantError {
		return
	}

	if proto != "LOGIN" {
		t.Errorf("invalid protocol, got: %q, want LOGIN", proto)
	}

	i := 0

	got := string(toServer)
	if got != test.wantData[i] {
		t.Errorf("invalid response, got %q, want %q", got, test.wantData[i])
	}

	for _, method := range test.methods {
		i++
		if i >= len(test.wantData) {
			t.Fatalf("unexpected method: %q", method)
		}

		toServer, err := auth.Next([]byte(method), true)

		if err != nil {
			t.Fatalf("loginAuth.Auth(): %v", err)
		}
		got = string(toServer)
		if got != test.wantData[i] {
			t.Errorf("invalid response, got %q, want %q", got, test.wantData[i])
		}

	}

}

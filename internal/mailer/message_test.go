package mailer

import "testing"

type message struct {
	from    string
	to      []string
	content string
}

func TestMessage(t *testing.T) {
	m := NewMessage()
	m.SetAddressHeader("From", "from@example.com", "Señor From")
	m.SetHeader("To", m.FormatAddress("to@example.com", "Señor To"), "tobis@example.com")
	m.SetAddressHeader("Cc", "cc@example.com", "A, B")
	m.SetAddressHeader("X-To", "ccbis@example.com", "à, b")
	m.SetDateHeader("X-Date", now())
	m.SetHeader("X-Date-2", m.FormatDate(now()))
	m.SetHeader("Subject", "¡Hola, señor!")
	m.SetHeaders(map[string][]string{
		"X-Headers": {"Test", "Café"},
	})
	m.SetBody("text/plain", "¡Hola, señor!")

	want := &message{
		from: "from@example.com",
		to: []string{
			"to@example.com",
			"tobis@example.com",
			"cc@example.com",
		},
		content: "From: =?UTF-8?q?Se=C3=B1or_From?= <from@example.com>\r\n" +
			"To: =?UTF-8?q?Se=C3=B1or_To?= <to@example.com>, tobis@example.com\r\n" +
			"Cc: \"A, B\" <cc@example.com>\r\n" +
			"X-To: =?UTF-8?b?w6AsIGI=?= <ccbis@example.com>\r\n" +
			"X-Date: Wed, 25 Jun 2014 17:46:00 +0000\r\n" +
			"X-Date-2: Wed, 25 Jun 2014 17:46:00 +0000\r\n" +
			"X-Headers: Test, =?UTF-8?q?Caf=C3=A9?=\r\n" +
			"Subject: =?UTF-8?q?=C2=A1Hola,_se=C3=B1or!?=\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"Content-Transfer-Encoding: quoted-printable\r\n" +
			"\r\n" +
			"=C2=A1Hola, se=C3=B1or!",
	}

	testMessage(t, m, 0, want)
}

// HELPER FUNCTIONS

func testMessage(t *testing.T, m *Message, bCount int, want *message) {
	err := Send(stubSendMail(t, bCount, want), m)
	if err != nil {
		t.Error(err)
	}
}

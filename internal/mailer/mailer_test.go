package mailer

// import (
// 	"errors"
// 	"os"
// 	"strings"
// 	"testing"
// 	"time"
// )

// // Mock template for testing
// const testTemplate = `
// {{define "subject"}}Test Subject{{end}}
// {{define "plainBody"}}Hello, {{.Name}}!{{end}}
// {{define "htmlBody"}}<h1>Hello, {{.Name}}!</h1>{{end}}
// `

// // Mock Dialer for testing
// type mockDialer struct {
//     dialAndSendCalled bool
//     lastMsg          *Message
//     shouldError      bool
// }

// func (md *mockDialer) DialAndSend(msg *Message) error {
//     md.dialAndSendCalled = true
//     md.lastMsg = msg
//     if md.shouldError {
//         return errors.New("mock dial error")
//     }
//     return nil
// }

// func TestNew(t *testing.T) {
//     tests := []struct {
//         name     string
//         host     string
//         port     int
//         username string
//         password string
//         sender   string
//     }{
//         {
//             name:     "Valid configuration",
//             host:     "smtp.example.com",
//             port:     587,
//             username: "test@example.com",
//             password: "password123",
//             sender:   "sender@example.com",
//         },
//     }

//     for _, tt := range tests {
//         t.Run(tt.name, func(t *testing.T) {
//             mailer := New(tt.host, tt.port, tt.username, tt.password, tt.sender)

//             if mailer == nil {
//                 t.Fatal("Expected non-nil mailer")
//             }

//             if mailer.sender != tt.sender {
//                 t.Errorf("Expected sender %s, got %s", tt.sender, mailer.sender)
//             }

//             if mailer.dialer.Timeout != 5*time.Second {
//                 t.Errorf("Expected timeout 5s, got %v", mailer.dialer.Timeout)
//             }
//         })
//     }
// }

// func TestMailer_Send(t *testing.T) {
//     tests := []struct {
//         name       string
//         recipient  string
//         template   string
//         data       interface{}
//         wantError  bool
//         mockDialer *mockDialer
//     }{
//         {
//             name:      "Valid email",
//             recipient: "recipient@example.com",
//             template:  "test.tmpl",
//             data: struct {
//                 Name string
//             }{
//                 Name: "John",
//             },
//             wantError:  false,
//             mockDialer: &mockDialer{shouldError: false},
//         },
//         {
//             name:      "Dialer error",
//             recipient: "recipient@example.com",
//             template:  "test.tmpl",
//             data: struct {
//                 Name string
//             }{
//                 Name: "John",
//             },
//             wantError:  true,
//             mockDialer: &mockDialer{shouldError: true},
//         },
//     }

//     for _, tt := range tests {
//         t.Run(tt.name, func(t *testing.T) {
//             // Create test template file
//             err := os.WriteFile("templates/test.tmpl", []byte(testTemplate), 0666)
//             if err != nil {
//                 t.Fatal(err)
//             }
//             defer os.Remove("templates/test.tmpl")

//             mailer := &Mailer{
//                 dialer: tt.mockDialer,
//                 sender: "sender@example.com",
//             }

//             err = mailer.Send(tt.recipient, tt.template, tt.data)

//             if (err != nil) != tt.wantError {
//                 t.Errorf("Send() error = %v, wantError %v", err, tt.wantError)
//                 return
//             }

//             if !tt.wantError {
//                 if !tt.mockDialer.dialAndSendCalled {
//                     t.Error("Expected DialAndSend to be called")
//                 }

//                 msg := tt.mockDialer.lastMsg
//                 if msg == nil {
//                     t.Fatal("Expected non-nil message")
//                 }

//                 // Check message properties
//                 if !strings.Contains(msg.GetHeader("To"), tt.recipient) {
//                     t.Errorf("Expected recipient %s, got %s", tt.recipient, msg.GetHeader("To"))
//                 }

//                 if !strings.Contains(msg.GetHeader("Subject"), "Test Subject") {
//                     t.Errorf("Unexpected subject: %s", msg.GetHeader("Subject"))
//                 }

//                 plainBody := msg.GetBody("text/plain")
//                 if !strings.Contains(plainBody, "Hello, John!") {
//                     t.Errorf("Unexpected plain body: %s", plainBody)
//                 }

//                 htmlBody := msg.GetBody("text/html")
//                 if !strings.Contains(htmlBody, "<h1>Hello, John!</h1>") {
//                     t.Errorf("Unexpected HTML body: %s", htmlBody)
//                 }
//             }
//         })
//     }
// }

// // Helper function to test template parsing errors
// func TestMailer_SendTemplateError(t *testing.T) {
//     mailer := &Mailer{
//         dialer: &mockDialer{},
//         sender: "sender@example.com",
//     }

//     err := mailer.Send("recipient@example.com", "nonexistent.tmpl", nil)
//     if err == nil {
//         t.Error("Expected error for nonexistent template")
//     }
// }

// // Test with invalid template content
// func TestMailer_SendInvalidTemplate(t *testing.T) {
//     invalidTemplate := `
// {{define "subject"}}{{.InvalidField}}{{end}}
// {{define "plainBody"}}Test{{end}}
// {{define "htmlBody"}}Test{{end}}
// `

//     // Create test template file
//     err := os.WriteFile("templates/invalid.tmpl", []byte(invalidTemplate), 0666)
//     if err != nil {
//         t.Fatal(err)
//     }
//     defer os.Remove("templates/invalid.tmpl")

//     mailer := &Mailer{
//         dialer: &mockDialer{},
//         sender: "sender@example.com",
//     }

//     err = mailer.Send("recipient@example.com", "invalid.tmpl", struct{}{})
//     if err == nil {
//         t.Error("Expected error for invalid template")
//     }
// }

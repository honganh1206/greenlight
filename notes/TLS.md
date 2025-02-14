# TLS - Transport Layer System

More modern and secure than [[SSL]]

1. **TLS is Required Because**:

```go
// Modern SMTP server typically starts on port 587 (STARTTLS)
server := &smtp.Server{
    Addr:    ":587",
    TLSConfig: &tls.Config{
        Certificates: []tls.Certificate{cert},
        MinVersion:   tls.VersionTLS12, // Minimum TLS version
    },
}
```

2. **Common SMTP Security Ports**:

- Port 587: STARTTLS (recommended)
- Port 465: Implicit TLS/SSL (legacy)
- Port 25: Plain text (unsafe, often blocked)

3. **STARTTLS Implementation Example**:

```go
func (s *SMTPServer) handleClient(conn net.Conn) {
    // Initial plain connection
    c := smtp.NewConn(conn, s.Server)

    // Handle STARTTLS command
    if err := c.TLSHandler(s.TLSConfig); err != nil {
        log.Printf("TLS error: %v", err)
        return
    }

    // Continue with secure communication
}
```

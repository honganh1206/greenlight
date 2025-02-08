# MIME (Multipurpose Internet Mail Extensions)

The `mime.WordEncoder` encodes the body of the email according to MIME standards, including multiplart messages e.g., plain text and HTML content, attachments and character encoding

## MIME boundaries

Used to separate different part of a multipart email message (like when an email has both text and attachment)

## Different levels of MIME

The three MIME multipart levels serve different purposes in email composition:

1. MIME Level 1:

- Basic MIME support
- Handles text/plain and text/html content types
- Supports US-ASCII and basic character encodings
- Good for simple text-based emails

2. MIME Level 2:

- Everything in Level 1, plus:
- Supports multipart messages (multipart/mixed, multipart/alternative)
- Handles attachments
- Supports more character encodings
- Can handle non-text content (images, audio, etc.)

3. MIME Level 3:

- Everything in Level 1 and 2, plus:
- Supports message security features
- Handles digital signatures
- Supports encryption
- Includes S/MIME capabilities
- More advanced content handling and security features

Example of the hierarchy of MIME types:

```go
// multipart/mixed
//   ├── multipart/related
//   │     ├── multipart/alternative
//   │     │     ├── text/plain
//   │     │     └── text/html
//   │     └── image (logo.png)
//   └── application/pdf (report.pdf)
message := NewMessage()
message.SetBody("text/plain", "Hello")           // alternative
message.AddAlternative("text/html", "<p>Hello with <img src='cid:logo.png'></p>")  // alternative
message.Embed("logo.png")                        // related
message.Attach("report.pdf")                     // mixed
```

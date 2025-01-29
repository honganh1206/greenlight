# Writing my own SMTP server from scratch

I aim to keep it minimal, so at this point of writing, I will ignore TLS configurations, SSL and only use simple authentication method

## Dialer

Establish and maintain connections to remote mail servers or clients

## MIME (Multipurpose Internet Mail Extensions)

The `mime.WordEncoder` encodes the body of the email according to MIME standards, including multiplart messages e.g., plain text and HTML content, attachments and character encoding

## Encoding

### Base64 encoding vs. Quoted-printable encoding

| Base64                                                                      | Quoted-printable                  |
| --------------------------------------------------------------------------- | --------------------------------- |
| Encode binary data (e.g., attachments, images) or non-ASCII text in headers | Encode ASCII-compatible text data |

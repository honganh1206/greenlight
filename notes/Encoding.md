# Encoding

## Base64 encoding vs. Quoted-printable encoding

| Base64                                                                      | Quoted-printable                  |
| --------------------------------------------------------------------------- | --------------------------------- |
| Encode binary data (e.g., attachments, images) or non-ASCII text in headers | Encode ASCII-compatible text data |

Sometimes we need to create a `*base64LineWriter` to limit text encoded in base64 to 76 characters specified in RFC 2045 (MIME specifications)

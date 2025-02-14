# Dialer

Establish and maintain connections to remote mail servers or clients

`LocalName` as the identity your client presents to the SMTP server. Used in the **initial handshake** when connecting to an SMTP server

Example: `"mycompany.com"` - A valid domain from which we are sending emails

## HELO/EHLO command

HELO (Original command per RFC 821)/EHLO (Extended version per RFC 2821) are SMTP commands used during the _initial handshake_ between an SMTP client and server.

1. **Basic Commands**:

```
HELO (Hello) - Basic greeting
EHLO (Extended Hello) - Modern version supporting extensions
```

2. **Typical SMTP Conversation**:

```
Client -> Server: EHLO client.example.com
Server -> Client: 250-smtp.server.com
Server -> Client: 250-SIZE 14680064
Server -> Client: 250-8BITMIME
Server -> Client: 250-STARTTLS
Server -> Client: 250 AUTH LOGIN PLAIN
```

[[TLS]]

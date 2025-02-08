# Headers

There are **top-level headers** (From/To/Subject/...) to _handle routing/identification_ and there are nested headers/MIME part headers (Content-Type/Content-Disposition/...) to _handle multiple content types_

Structure:

```text
Email Message
├── Top-level Headers
│   ├── Routing information (From, To)
│   ├── Subject
│   └── Overall MIME structure
│
└── Body (MIME Parts)
    ├── Part 1
    │   ├── Part Headers (Content-Type, etc.)
    │   └── Content
    └── Part 2
        ├── Part Headers
        └── Content
```

More detailed version:

```
[Top-level Headers]       # Message envelope/metadata
From: sender@example.com
To: recipient@example.com
Subject: Hello
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="boundary1"

--boundary1              # Start of MIME part
[MIME Part Headers]      # Part-specific metadata
Content-Type: text/plain
Content-Disposition: inline

This is the text content.
--boundary1              # Another MIME part
[MIME Part Headers]      # Different part headers
Content-Type: image/jpeg
Content-Disposition: attachment; filename="photo.jpg"

[Binary data...]
--boundary1--            # End of multipart
```

The parent multipart is responsible for managing the boundaries between itself and its child multipart

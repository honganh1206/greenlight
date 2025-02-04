# Sending Emails

I use MailTrap SMTP service now (for easier deployment?), but adding Maildev as a docker container just in case

Our user welcome email template (as a `.tmpl` file)will have 3 parts:

- `subject` contains the subject line of the email
- `plainBody` contains the plain-text variant of the email message body
- `htmlBody` contains the HTML variant of the email message body

[Revisiting tmpl file](Revisiting%20tmpl%20file.md)

Also, we will [write a small SMTP server from scratch](./Write my own SMTP server from scratch.md) (mostly copy it from [gomail](https://github.com/go-gomail/gomail) )



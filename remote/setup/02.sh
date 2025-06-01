 #!/bin/bash

 read -p "Enter host name for Mailtrap SMTP server: " MAILTRAP_SMTP_HOST
 read -p "Enter user name for Mailtrap SMTP server: " MAILTRAP_SMTP_USERNAME
 read -p "Enter password for Mailtrap SMTP server: " MAILTRAP_SMTP_HOST

 echo "MAILTRAP_SMTP_HOST=${MAILTRAP_SMTP_HOST}" >> /etc/environment
 echo "MAILTRAP_SMTP_USERNAME=${MAILTRAP_SMTP_USERNAME}" >> /etc/environment
 echo "MAILTRAP_SMTP_PASSWORD=${MAILTRAP_SMTP_PASSWORD}" >> /etc/environment

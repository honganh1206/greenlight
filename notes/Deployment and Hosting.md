# Deployment and Hosting

Hosting provider: Digital Ocean - $5 a month

Server: Ubuntu Linux

Reverse proxy: Caddy

Things to do:

- Update packages
- Set server timezone
- Create a `greenlight` user on the server for day-to-day maintenance
- Copy the `root` user HOME to `greenlight` user directory
- Configure firewall to allow only 22 (SSH), 80 (HTTP) and 443 (HTTPS)
- Install PostgreSQL, migrate, and Caddy

```bash
# Setup script for the server
$ rsync -rP --delete ./remote/setup greenlight@45.55.49.87:~
# 02.sh could be a script to add Mailtrap host, username and password to /etc/environment
$ ssh -t greenlight@45.55.49.87 "sudo bash /home/greenlight/setup/02.sh"
```

Remember to update `greenlight` database owner to our user:

```bash
# Login as postgres user
sudo -u postgres psql
```

Then update the owner

```sql
alter database greenlight owner to greenlight;
```

Update the port for the app (as `greenlight` user)

```bash
# For greenlight
sudo ufw allow 4000/tcp
# For mailtrap
sudo ufw allow 2525/tcp
```

We will make our API run as a background service using `systemd` and an unit file to tell how `systemd` can run the service

View the log with `sudo journalctl -u api`

[[Using Caddy as reverse proxy]]

[[Scaling up]]

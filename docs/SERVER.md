# Server Administration Guide

Complete guide for deploying, configuring, and managing the Quotes API server.

## Installation

### Binary Installation

#### Linux

```bash
# Download latest release
wget https://github.com/apimgr/quotes/releases/latest/download/quotes-linux-amd64

# Make executable
chmod +x quotes-linux-amd64

# Move to system path
sudo mv quotes-linux-amd64 /usr/local/bin/quotes

# Create user and directories
sudo useradd -r -s /bin/false quotes
sudo mkdir -p /var/lib/quotes/{config,data,logs}
sudo chown -R quotes:quotes /var/lib/quotes

# Run
sudo -u quotes quotes --port 8080
```

#### macOS

```bash
# Download latest release
wget https://github.com/apimgr/quotes/releases/latest/download/quotes-macos-arm64

# Make executable
chmod +x quotes-macos-arm64

# Move to system path
sudo mv quotes-macos-arm64 /usr/local/bin/quotes

# Run
quotes --port 8080
```

#### Windows

```powershell
# Download from releases page
# https://github.com/apimgr/quotes/releases/latest/download/quotes-windows-amd64.exe

# Run
.\quotes-windows-amd64.exe --port 8080
```

### Docker Installation

#### Docker Run

```bash
# Pull image
docker pull ghcr.io/apimgr/quotes:latest

# Run container
docker run -d \
  --name quotes \
  --restart unless-stopped \
  -p 8080:80 \
  -v /var/lib/quotes/config:/config \
  -v /var/lib/quotes/data:/data \
  -v /var/lib/quotes/logs:/logs \
  -e ADMIN_USER=administrator \
  -e ADMIN_PASSWORD=changeme123 \
  ghcr.io/apimgr/quotes:latest
```

#### Docker Compose (Production)

Create `docker-compose.yml`:

```yaml
services:
  quotes:
    image: ghcr.io/apimgr/quotes:latest
    container_name: quotes
    restart: unless-stopped

    environment:
      - CONFIG_DIR=/config
      - DATA_DIR=/data
      - LOGS_DIR=/logs
      - PORT=80
      - ADDRESS=0.0.0.0
      - DB_PATH=/data/db/quotes.db
      - ADMIN_USER=administrator
      - ADMIN_PASSWORD=changeme123

    volumes:
      - ./rootfs/config/quotes:/config
      - ./rootfs/data/quotes:/data
      - ./rootfs/logs/quotes:/logs

    ports:
      - "172.17.0.1:64180:80"

    networks:
      - quotes

    healthcheck:
      test: ["CMD", "/usr/local/bin/quotes", "--status"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 10s

networks:
  quotes:
    name: quotes
    external: false
    driver: bridge
```

Start services:

```bash
docker compose up -d
```

#### Docker Compose (Development)

Create `docker-compose.test.yml`:

```yaml
services:
  quotes:
    image: quotes:dev
    container_name: quotes-test
    restart: "no"

    environment:
      - PORT=80
      - ADMIN_USER=administrator
      - ADMIN_PASSWORD=testpass123
      - DEV=true

    volumes:
      - /tmp/quotes/rootfs/config/quotes:/config
      - /tmp/quotes/rootfs/data/quotes:/data
      - /tmp/quotes/rootfs/logs/quotes:/logs

    ports:
      - "64181:80"

    networks:
      - quotes

networks:
  quotes:
    name: quotes
    external: false
    driver: bridge
```

Start services:

```bash
docker compose -f docker-compose.test.yml up -d
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `ADDRESS` | `0.0.0.0` | Bind address |
| `CONFIG_DIR` | Platform-specific | Configuration directory |
| `DATA_DIR` | Platform-specific | Data directory |
| `LOGS_DIR` | Platform-specific | Logs directory |
| `DB_PATH` | `{DATA_DIR}/db/quotes.db` | SQLite database path |
| `ADMIN_USER` | `administrator` | Admin username (first run) |
| `ADMIN_PASSWORD` | - | Admin password (first run) |
| `ADMIN_TOKEN` | Auto-generated | Admin API token (first run) |

### Directory Locations

#### Linux
- Config: `/var/lib/quotes/config`
- Data: `/var/lib/quotes/data`
- Logs: `/var/lib/quotes/logs`
- Database: `/var/lib/quotes/data/db/quotes.db`

#### macOS
- Config: `~/Library/Application Support/quotes/config`
- Data: `~/Library/Application Support/quotes/data`
- Logs: `~/Library/Application Support/quotes/logs`

#### Windows
- Config: `%APPDATA%\quotes\config`
- Data: `%APPDATA%\quotes\data`
- Logs: `%APPDATA%\quotes\logs`

#### Docker
- Config: `/config`
- Data: `/data`
- Logs: `/logs`
- Database: `/data/db/quotes.db`

### Command-Line Options

```bash
quotes [OPTIONS]

Options:
  --port PORT           Server port (default: 8080)
  --address ADDRESS     Bind address (default: 0.0.0.0)
  --config DIR          Config directory
  --data DIR            Data directory
  --logs DIR            Logs directory
  --status              Check server status (exit code 0 if healthy)
  --version             Print version and exit
  --help                Show help message
```

Examples:

```bash
# Run on port 9000
quotes --port 9000

# Run on localhost only
quotes --address 127.0.0.1

# Custom directories
quotes --config /etc/quotes --data /opt/quotes/data --logs /var/log/quotes

# Check status
quotes --status
```

## Systemd Service

### Create Service File

Create `/etc/systemd/system/quotes.service`:

```ini
[Unit]
Description=Quotes API Server
Documentation=https://quotes.readthedocs.io
After=network.target

[Service]
Type=simple
User=quotes
Group=quotes
WorkingDirectory=/var/lib/quotes

ExecStart=/usr/local/bin/quotes --port 8080
ExecReload=/bin/kill -HUP $MAINPID

Restart=always
RestartSec=10

# Environment
Environment="CONFIG_DIR=/var/lib/quotes/config"
Environment="DATA_DIR=/var/lib/quotes/data"
Environment="LOGS_DIR=/var/lib/quotes/logs"

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/quotes

# Limits
LimitNOFILE=65535
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

### Manage Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service (start on boot)
sudo systemctl enable quotes

# Start service
sudo systemctl start quotes

# Stop service
sudo systemctl stop quotes

# Restart service
sudo systemctl restart quotes

# Check status
sudo systemctl status quotes

# View logs
sudo journalctl -u quotes -f
```

## Reverse Proxy Configuration

### Nginx

Create `/etc/nginx/sites-available/quotes`:

```nginx
server {
    listen 80;
    server_name quotes.example.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name quotes.example.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/quotes.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/quotes.example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Logging
    access_log /var/log/nginx/quotes-access.log;
    error_log /var/log/nginx/quotes-error.log;

    # Proxy settings
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support (if needed)
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=quotes:10m rate=10r/s;
    limit_req zone=quotes burst=20 nodelay;
}
```

Enable and restart:

```bash
sudo ln -s /etc/nginx/sites-available/quotes /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Caddy

Create `Caddyfile`:

```caddy
quotes.example.com {
    reverse_proxy localhost:8080

    # Automatic HTTPS
    tls {
        protocols tls1.2 tls1.3
    }

    # Logging
    log {
        output file /var/log/caddy/quotes.log
        format json
    }

    # Security headers
    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains"
        X-Content-Type-Options "nosniff"
        X-Frame-Options "DENY"
        Referrer-Policy "no-referrer-when-downgrade"
    }

    # Rate limiting
    rate_limit {
        zone quotes {
            key {remote_host}
            events 100
            window 1m
        }
    }
}
```

Run Caddy:

```bash
caddy run --config Caddyfile
```

## Monitoring

### Health Checks

#### HTTP Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2025-10-14T12:00:00Z"
}
```

#### Status Check (Binary)

```bash
quotes --status
echo $?  # 0 if healthy, 1 if not
```

#### Docker Health Check

```bash
docker exec quotes quotes --status
```

### Metrics

Get server statistics via admin API:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/admin/stats
```

Response:
```json
{
  "success": true,
  "data": {
    "server": {
      "version": "0.0.1",
      "uptime": "2h30m15s"
    },
    "collections": {
      "total": 27500
    },
    "requests": {
      "total": 12345,
      "rate": "124.5 req/s"
    },
    "memory": {
      "allocated": "52.4 MB"
    }
  }
}
```

### Logging

#### Log Locations

- Linux: `/var/lib/quotes/logs/quotes.log`
- Docker: `/logs/quotes.log`
- Systemd: `journalctl -u quotes`

#### Log Levels

- `INFO`: Normal operations
- `WARN`: Non-critical issues
- `ERROR`: Critical errors

#### View Logs

```bash
# Tail logs
tail -f /var/lib/quotes/logs/quotes.log

# Systemd logs
journalctl -u quotes -f

# Docker logs
docker logs -f quotes
```

## Backup & Restore

### Backup

#### Full Backup

```bash
# Create backup directory
mkdir -p /backup/quotes-$(date +%Y%m%d)

# Backup data
sudo cp -r /var/lib/quotes/config /backup/quotes-$(date +%Y%m%d)/
sudo cp -r /var/lib/quotes/data /backup/quotes-$(date +%Y%m%d)/

# Create archive
cd /backup
sudo tar -czf quotes-$(date +%Y%m%d).tar.gz quotes-$(date +%Y%m%d)/
```

#### Database Only

```bash
# Stop service
sudo systemctl stop quotes

# Backup database
sudo cp /var/lib/quotes/data/db/quotes.db /backup/quotes-db-$(date +%Y%m%d).db

# Start service
sudo systemctl start quotes
```

#### Docker Backup

```bash
# Backup volumes
docker run --rm \
  -v quotes_config:/config \
  -v quotes_data:/data \
  -v $(pwd):/backup \
  alpine tar -czf /backup/quotes-backup-$(date +%Y%m%d).tar.gz /config /data
```

### Restore

#### Full Restore

```bash
# Stop service
sudo systemctl stop quotes

# Extract backup
cd /backup
sudo tar -xzf quotes-20251014.tar.gz

# Restore data
sudo rm -rf /var/lib/quotes/config /var/lib/quotes/data
sudo cp -r quotes-20251014/config /var/lib/quotes/
sudo cp -r quotes-20251014/data /var/lib/quotes/

# Fix permissions
sudo chown -R quotes:quotes /var/lib/quotes

# Start service
sudo systemctl start quotes
```

#### Database Restore

```bash
# Stop service
sudo systemctl stop quotes

# Restore database
sudo cp /backup/quotes-db-20251014.db /var/lib/quotes/data/db/quotes.db
sudo chown quotes:quotes /var/lib/quotes/data/db/quotes.db

# Start service
sudo systemctl start quotes
```

## Security

### Admin Credentials

#### First Run Setup

Set admin credentials via environment variables:

```bash
ADMIN_USER=administrator \
ADMIN_PASSWORD=secure-password-here \
quotes --port 8080
```

Credentials are saved to: `{CONFIG_DIR}/credentials.txt`

#### Change Admin Password

1. Stop the service
2. Delete the database: `rm /var/lib/quotes/data/db/quotes.db`
3. Start with new credentials:
   ```bash
   ADMIN_USER=administrator \
   ADMIN_PASSWORD=new-password \
   quotes --port 8080
   ```

### API Token Management

Admin API token is auto-generated on first run and saved in:
- Database: `/var/lib/quotes/data/db/quotes.db`
- Credentials file: `/var/lib/quotes/config/credentials.txt`

To regenerate token:
1. Stop service
2. Delete database
3. Restart service with new credentials

### Firewall Configuration

#### Linux (UFW)

```bash
# Allow HTTP
sudo ufw allow 8080/tcp

# Allow HTTPS (if using reverse proxy)
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable
```

#### Linux (firewalld)

```bash
# Allow HTTP
sudo firewall-cmd --permanent --add-port=8080/tcp

# Reload
sudo firewall-cmd --reload
```

### SSL/TLS

Always use HTTPS in production. Use a reverse proxy (nginx, Caddy) with Let's Encrypt certificates.

#### Let's Encrypt with Certbot

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d quotes.example.com

# Auto-renewal
sudo systemctl enable certbot.timer
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port
sudo lsof -i :8080

# Kill process
sudo kill -9 PID

# Or use different port
quotes --port 8081
```

### Database Locked

```bash
# Check for stale lock files
ls -la /var/lib/quotes/data/db/

# Remove lock files
sudo rm /var/lib/quotes/data/db/quotes.db-wal
sudo rm /var/lib/quotes/data/db/quotes.db-shm

# Restart service
sudo systemctl restart quotes
```

### Permission Denied

```bash
# Fix ownership
sudo chown -R quotes:quotes /var/lib/quotes

# Fix permissions
sudo chmod -R 755 /var/lib/quotes
sudo chmod 644 /var/lib/quotes/data/db/quotes.db
```

### High Memory Usage

```bash
# Check memory
free -h

# View process memory
ps aux | grep quotes

# Restart service to clear cache
sudo systemctl restart quotes
```

### Service Won't Start

```bash
# Check logs
sudo journalctl -u quotes -n 50

# Check config
quotes --help

# Verify binary
which quotes
quotes --version

# Test run
sudo -u quotes quotes --port 8080
```

## Performance Tuning

### System Limits

Edit `/etc/security/limits.conf`:

```
quotes soft nofile 65535
quotes hard nofile 65535
quotes soft nproc 4096
quotes hard nproc 4096
```

### Kernel Parameters

Edit `/etc/sysctl.conf`:

```
# Network
net.core.somaxconn = 65535
net.ipv4.tcp_max_syn_backlog = 8192
net.ipv4.tcp_tw_reuse = 1

# File descriptors
fs.file-max = 2097152
```

Apply changes:

```bash
sudo sysctl -p
```

### Database Optimization

SQLite is already optimized for read-heavy workloads. No additional tuning needed.

## Upgrading

### Binary Upgrade

```bash
# Stop service
sudo systemctl stop quotes

# Backup current binary
sudo cp /usr/local/bin/quotes /usr/local/bin/quotes.bak

# Download new version
wget https://github.com/apimgr/quotes/releases/latest/download/quotes-linux-amd64

# Replace binary
sudo mv quotes-linux-amd64 /usr/local/bin/quotes
sudo chmod +x /usr/local/bin/quotes

# Start service
sudo systemctl start quotes

# Verify version
quotes --version
```

### Docker Upgrade

```bash
# Pull new image
docker pull ghcr.io/apimgr/quotes:latest

# Stop and remove old container
docker stop quotes
docker rm quotes

# Start new container
docker compose up -d

# Or using docker run
docker run -d \
  --name quotes \
  --restart unless-stopped \
  -p 8080:80 \
  -v /var/lib/quotes/config:/config \
  -v /var/lib/quotes/data:/data \
  -v /var/lib/quotes/logs:/logs \
  ghcr.io/apimgr/quotes:latest
```

## Uninstallation

### Binary Uninstall

```bash
# Stop service
sudo systemctl stop quotes
sudo systemctl disable quotes

# Remove service file
sudo rm /etc/systemd/system/quotes.service
sudo systemctl daemon-reload

# Remove binary
sudo rm /usr/local/bin/quotes

# Remove data (optional)
sudo rm -rf /var/lib/quotes

# Remove user
sudo userdel quotes
```

### Docker Uninstall

```bash
# Stop and remove containers
docker compose down

# Remove images
docker rmi ghcr.io/apimgr/quotes:latest

# Remove volumes (optional)
sudo rm -rf ./rootfs
```

## Support

- **Documentation**: [https://quotes.readthedocs.io](https://quotes.readthedocs.io)
- **Issues**: [https://github.com/apimgr/quotes/issues](https://github.com/apimgr/quotes/issues)
- **Repository**: [https://github.com/apimgr/quotes](https://github.com/apimgr/quotes)

---

**Version**: 0.0.1
**Last Updated**: 2025-10-14

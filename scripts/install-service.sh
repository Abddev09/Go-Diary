#!/bin/bash

# Blog service o'rnatish skripti

# 1. Binary faylni yaratish
echo "Building application..."
go build -o /usr/local/bin/blog-server

# 2. Service faylini yaratish
sudo tee /etc/systemd/system/blog-server.service > /dev/null <<EOF
[Unit]
Description=Blog Server
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/blog
ExecStart=/usr/local/bin/blog-server
Restart=always
RestartSec=5
Environment=DATABASE_URL=postgres://blog_user:abdulloh_009@localhost/blog_system?sslmode=disable
Environment=PORT=8080

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=blog-server

[Install]
WantedBy=multi-user.target
EOF

# 3. Static fayllar uchun papka yaratish
sudo mkdir -p /var/www/blog/static
sudo mkdir -p /var/www/blog/uploads
sudo chown -R www-data:www-data /var/www/blog

# 4. Static fayllarni nusxalash
sudo cp -r static/* /var/www/blog/static/

# 5. Service ni yoqish
sudo systemctl daemon-reload
sudo systemctl enable blog-server
sudo systemctl start blog-server

echo "Service o'rnatildi! Status ko'rish uchun: sudo systemctl status blog-server"

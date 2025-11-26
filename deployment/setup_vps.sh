#!/bin/bash
set -e

echo "Starting VPS Setup..."

# Update system
echo "Updating system packages..."
sudo apt-get update
sudo apt-get upgrade -y

# Install Docker
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    sudo apt-get install -y docker.io
    sudo systemctl start docker
    sudo systemctl enable docker
    sudo usermod -aG docker $USER
    echo "Docker installed."
else
    echo "Docker already installed."
fi

# Install Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "Installing Docker Compose..."
    sudo apt-get install -y docker-compose
    echo "Docker Compose installed."
else
    echo "Docker Compose already installed."
fi

# Install Nginx
if ! command -v nginx &> /dev/null; then
    echo "Installing Nginx..."
    sudo apt-get install -y nginx
    sudo systemctl start nginx
    sudo systemctl enable nginx
    echo "Nginx installed."
else
    echo "Nginx already installed."
fi

# Install Certbot
if ! command -v certbot &> /dev/null; then
    echo "Installing Certbot..."
    sudo apt-get install -y certbot python3-certbot-nginx
    echo "Certbot installed."
else
    echo "Certbot already installed."
fi

# Setup Firewall
echo "Configuring Firewall..."
sudo ufw allow OpenSSH
sudo ufw allow 'Nginx Full'
# Only enable if not already enabled to avoid locking out (though OpenSSH allowed above)
if ! sudo ufw status | grep -q "Status: active"; then
    echo "y" | sudo ufw enable
fi

echo "✅ VPS Setup Complete!"
echo "⚠️  Please logout and login again to apply Docker group changes."

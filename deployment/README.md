# Event Campus - VPS Deployment Guide

This guide explains how to deploy the Event Campus API to a VPS (Virtual Private Server) running Ubuntu or Debian.

## Prerequisites
- A VPS (DigitalOcean, AWS, Google Cloud, etc.)
- Ubuntu 20.04/22.04 or Debian 10/11
- Root access or sudo privileges
- A domain name pointing to your VPS IP address

## Step 1: Initial Setup

1. SSH into your VPS:
   ```bash
   ssh root@your_vps_ip
   ```

2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/event-campus-backend.git
   cd event-campus-backend
   ```

3. Run the setup script to install dependencies (Docker, Nginx, Certbot):
   ```bash
   chmod +x deployment/setup_vps.sh
   ./deployment/setup_vps.sh
   ```
   *Note: You may need to logout and login again for Docker permission changes to take effect.*

## Step 2: Configuration

1. Create `.env` file from example:
   ```bash
   cp .env.example .env
   nano .env
   ```
   *Fill in your database credentials, JWT secret, etc.*

2. Configure Nginx:
   ```bash
   sudo cp deployment/nginx.conf /etc/nginx/sites-available/event-campus
   sudo ln -s /etc/nginx/sites-available/event-campus /etc/nginx/sites-enabled/
   sudo rm /etc/nginx/sites-enabled/default
   ```

3. Edit Nginx config to set your domain:
   ```bash
   sudo nano /etc/nginx/sites-available/event-campus
   # Change 'server_name your_domain.com;' to your actual domain
   ```

4. Test and restart Nginx:
   ```bash
   sudo nginx -t
   sudo systemctl restart nginx
   ```

## Step 3: SSL Setup (HTTPS)

Run Certbot to get a free SSL certificate:
```bash
sudo certbot --nginx -d your_domain.com
```

## Step 4: Deployment

Run the deployment script to build and start the application:
```bash
chmod +x deployment/deploy.sh
./deployment/deploy.sh
```

## Maintenance

To update the application later, just run the deployment script again:
```bash
./deployment/deploy.sh
```

This will:
1. Pull the latest code from Git
2. Rebuild the Docker image
3. Restart the container with the new version

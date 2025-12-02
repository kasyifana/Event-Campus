# 🚀 Jenkins + Docker Deployment Guide for Event Campus

Panduan lengkap untuk deployment aplikasi Event Campus ke VPS menggunakan Docker dan Jenkins CI/CD.

**Setup:** Single VPS dengan Jenkins + Docker + Direct Port Expose

---

## 📋 Prasyarat

- **VPS Ubuntu** (Fresh Install atau sudah ada)
- **Akses Root/Sudo**
- **Jenkins** (akan di-install di VPS yang sama)
- **Git Repository** (GitHub/GitLab) untuk source code

---

## 🛠️ Step 1: Setup VPS & Install Dependencies

### 1.1 Install Docker

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Start Docker
sudo systemctl start docker
sudo systemctl enable docker

# Verify
docker --version
```

### 1.2 Install Docker Compose

```bash
# Install docker-compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify
docker-compose --version
```

### 1.3 Install Jenkins

```bash
# Install Java (Required for Jenkins)
sudo apt install -y openjdk-17-jre

# Add Jenkins repository
curl -fsSL https://pkg.jenkins.io/debian-stable/jenkins.io-2023.key | sudo tee \
  /usr/share/keyrings/jenkins-keyring.asc > /dev/null
echo deb [signed-by=/usr/share/keyrings/jenkins-keyring.asc] \
  https://pkg.jenkins.io/debian-stable binary/ | sudo tee \
  /etc/apt/sources.list.d/jenkins.list > /dev/null

# Install Jenkins
sudo apt update
sudo apt install -y jenkins

# Start Jenkins
sudo systemctl start jenkins
sudo systemctl enable jenkins

# Get initial admin password
sudo cat /var/lib/jenkins/secrets/initialAdminPassword
## 7895545a3ab94a6bad7690f1293a47d6
## 47effd3a2dd64300ada5e2ca4eebf55b bawang
## 1c1dbaedf4cb486eb7d256f04839b2ba ali backup
```

**Akses Jenkins:** `http://YOUR_VPS_IP:8080`

### 1.4 Configure Jenkins User

```bash
# Add jenkins user to docker group (important!)
sudo usermod -aG docker jenkins

# Restart Jenkins
sudo systemctl restart jenkins
```

---

## ⚙️ Step 2: Setup Application Directory & Environment

### 2.1 Create Deployment Directory

```bash
# Create app directory
sudo mkdir -p /opt/event-campus/storage/posters
sudo mkdir -p /opt/event-campus/storage/documents

# Give ownership to Jenkins
sudo chown -R jenkins:jenkins /opt/event-campus
```

### 2.2 Create .env File

```bash
# Create .env file
sudo nano /opt/event-campus/.env
```

**Isi dengan konfigurasi Production** (gunakan `.env.example` sebagai template):

```env
PORT=8080
ENV=production
ALLOWED_ORIGINS=http://YOUR_VPS_IP:3000

JWT_SECRET=generate-with-openssl-rand-base64-32
JWT_EXPIRATION=168h

# Supabase Database
POSTGRES_HOST=aws-0-ap-southeast-1.pooler.supabase.com
POSTGRES_PORT=6543
POSTGRES_USER=postgres.xxxxxxxxxxxxx
POSTGRES_PASSWORD=your-supabase-password
POSTGRES_DB=postgres
POSTGRES_SSLMODE=require

# Email (Gmail SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-gmail-app-password

# File Upload
UPLOAD_PATH=storage
MAX_UPLOAD_SIZE=10485760
```

**Generate JWT Secret:**
```bash
openssl rand -base64 32
##T0q+vCESyFAzdjtiIVGXucnEdUBLm64GN8MrZH/W3uQ=
```

---

## 🤖 Step 3: Configure Jenkins Pipeline

### 3.1 Install Required Jenkins Plugins

1. Buka Jenkins Dashboard (`http://YOUR_VPS_IP:8080`)
2. Ke **Manage Jenkins** → **Manage Plugins**
3. Install plugins:
   - Git Plugin
   - Docker Plugin
   - Pipeline Plugin
   - Credentials Plugin

### 3.2 Add Git Credentials (Jika Private Repo)

1. **Manage Jenkins** → **Credentials** → **Global**
2. **Add Credentials**:
   - Kind: Username with password
   - Username: Git username
   - Password: Personal Access Token
   - ID: `git-credentials`

### 3.3 Create Pipeline Job

1. **New Item** → Masukkan nama: `event-campus-deployment`
2. Pilih **Pipeline**
3. Di **Pipeline** section:
   - Definition: **Pipeline script from SCM**
   - SCM: **Git**
   - Repository URL: `https://github.com/yourusername/event-campus-backend.git`
   - Credentials: Pilih credentials yang dibuat tadi (jika private)
   - Branch Specifier: `*/main`
   - Script Path: `Jenkinsfile`
4. **Save**

---

## 🚀 Step 4: Deploy untuk Pertama Kali

### 4.1 Trigger Build

1. Di Jenkins Dashboard, pilih pipeline `event-campus-deployment`
2. Klik **Build Now**
3. Monitor progress di **Console Output**

### 4.2 Verify Deployment

```bash
# Check container status
docker ps | grep event-campus-api

# Check logs
docker logs -f event-campus-api

# Test health endpoint
curl http://localhost:8080/health
```

**Expected Response:**
```json
{"status":"ok"}
```

### 4.3 Test API dari luar VPS

```bash
# Dari komputer lokal
curl http://YOUR_VPS_IP:8080/health
```

**Jika gagal**, pastikan firewall allow port 8080:
```bash
# Ubuntu UFW
sudo ufw allow 8080/tcp
sudo ufw reload
```

---

## 🔄 Step 5: Continuous Deployment

Setelah setup awal, setiap kali ada perubahan code:

1. **Push code** ke Git repository
2. **Trigger Jenkins** manual dengan klik **Build Now**, ATAU
3. **Setup Webhook** untuk auto-trigger (optional):
   - Di GitHub: Settings → Webhooks → Add webhook
   - Payload URL: `http://YOUR_VPS_IP:8080/github-webhook/`
   - Content type: `application/json`
   - Events: `Just the push event`

---

## 🆘 Troubleshooting

### Permission Denied (Docker)

**Symptom:** Jenkins build gagal dengan error `permission denied` saat run Docker

**Solution:**
```bash
sudo usermod -aG docker jenkins
sudo systemctl restart jenkins
# Verify
sudo -u jenkins docker ps
```

### .env File Not Found

**Symptom:** Container fail to start, logs show config error

**Solution:**
```bash
# Check if .env exists
ls -la /opt/event-campus/.env

# Verify ownership
sudo chown jenkins:jenkins /opt/event-campus/.env
```

### Port Already in Use

**Symptom:** Error `port 8080 already allocated`

**Solution:**
```bash
# Find process using port 8080
sudo lsof -i :8080

# Kill old container
docker stop event-campus-api
docker rm event-campus-api
```

### Health Check Failed

**Symptom:** Deployment fails at health check stage

**Solution:**
```bash
# Check container logs
docker logs event-campus-api

# Common issues:
# 1. Database connection failed → Check .env POSTGRES_* variables
# 2. Migration error → Check migration files
# 3. Port conflict → Check if port 8080 is free
```

### Rollback to Previous Version

Jika deployment baru bermasalah:

```bash
cd /path/to/repo
./deployment/rollback.sh
```

Script ini akan otomatis mengembalikan ke versi sebelumnya.

---

## 📊 Monitoring & Maintenance

### Check Application Status

```bash
# Container status
docker ps -a | grep event-campus

# Resource usage
docker stats event-campus-api

# Recent logs
docker logs --tail 100 -f event-campus-api
```

### Backup Database

Karena pakai Supabase, backup sudah automatic. Untuk manual backup:

1. Login ke Supabase Dashboard
2. Database → Backups
3. Download backup atau restore dari point-in-time

### Clean Up Old Images

```bash
# Remove dangling images
docker image prune -f

# Remove old builds (keep last 3)
docker images event-campus-api --format "{{.Tag}}" | \
    grep -v 'latest' | grep -v 'previous' | tail -n +4 | \
    xargs -r -I {} docker rmi event-campus-api:{}
```

---

## 🔒 Security Best Practices

1. ✅ **Never commit .env** to Git
2. ✅ **Use strong JWT_SECRET** (min 32 characters)
3. ✅ **Use App Password** untuk Gmail SMTP (bukan password asli)
4. ✅ **Update sistem** secara rutin:
   ```bash
   sudo apt update && sudo apt upgrade -y
   ```
5. ✅ **Setup firewall** (allow hanya port yang diperlukan):
   ```bash
   sudo ufw allow 22/tcp    # SSH
   sudo ufw allow 8080/tcp  # API
   sudo ufw allow 8080/tcp  # Jenkins
   sudo ufw enable
   ```

---

## 📁 File Structure Summary

```
event-campus-backend/
├── Dockerfile              # Multi-stage Docker build
├── .dockerignore          # Exclude unnecessary files
├── Jenkinsfile            # CI/CD pipeline definition
├── .env.example           # Environment template
├── deployment/
│   ├── docker-compose.yml # Compose configuration
│   ├── deploy.sh          # Manual deployment script
│   └── rollback.sh        # Rollback script
└── /opt/event-campus/     # VPS deployment directory
    ├── .env               # Production environment
    └── storage/           # Uploaded files
```

---

**Selamat! Setup deployment Anda sudah lengkap! 🎉**

Untuk pertanyaan atau issues, check logs dengan:
```bash
docker logs -f event-campus-api
```


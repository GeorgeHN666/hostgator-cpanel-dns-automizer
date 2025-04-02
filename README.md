# Hostgator / cPanel DNS Automizer 🛰️

A lightweight DNS automation tool designed to keep your A records in sync with your dynamic public IP address — specifically for setups where you’re hosting a server locally (behind a reverse proxy) but managing your domain through **Hostgator or any cPanel-based provider**.

## 🧠 Why this project exists

Many hosting providers like **Hostgator** that use **cPanel** do not allow complex backend server setups (like Golang, Node.js, etc.) or give SSH/root access. So, developers often **host their actual server locally** and expose it to the internet using tools like:

- [Cloudflare Tunnel](https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/)
- [Ngrok](https://ngrok.com/)
- [NGINX](https://nginx.org/)
- [Apache](https://httpd.apache.org/)

However, here's the catch:\
Your **local public IP changes frequently** (dynamic IP), and **manually updating DNS A records every time is annoying and error-prone**.

This tool solves that.

## ⚙️ How it works

1. **Periodically checks your public IP** (using a public IP resolver like `https://api.ipify.org`)
2. **Saves a persistent file on your machine** that stores the last known public IP, and compares it with the current public IP to detect any changes.
3. If there's a difference:
   - Logs in to your **cPanel account** (via its native HTTP interface or API)
   - **Automatically updates the A record** with the new IP for `https://yourdomain...` and `https://www.yourdomain...`.
4. Repeat — works as a background service on your machine.

## 🧪 Example Use Case

You're running:

- A **Golang API** on your local Linux, macOS, or Windows machine.
- Your domain is managed on **Hostgator** or plain cPanel.
- You expose your server via **Cloudflare Tunnel**, **NGINX**, or **Apache**.
- Your public IP changes because you're on a home/office connection.

This tool will:

- Detect the IP change.
- Automatically update your `api.example.com` A record in cPanel DNS settings.
- Ensure traffic always resolves to your live server, even if your IP changes overnight.

## 🖥️ Requirements

- Any OS that supports Go (Linux, macOS, Windows)
- Hostgator or any hosting provider with **cPanel DNS access**
- Domain/subdomain pointed to your cPanel account
- cPanel **API token** with DNS edit permissions

## 🔐 cPanel API Token (Required)

To allow this tool to update DNS records via cPanel, you'll need to create an API token:

1. Log in to your **cPanel account**.
2. Search for **Manage API Tokens** (usually under "Security").
3. Click **Create Token**.
4. Set a name (e.g., `dns-automizer`) and give it **Zone Editor (A, CNAME, etc)** permissions.
5. Copy and save the generated token securely.

You’ll use this token in your `.env` config.

---

## 🚀 Setup Instructions (Per OS)

### 🐧 Linux Setup (Ubuntu/Debian)

1. **Install Golang:**
```bash
sudo apt update
sudo apt install golang-go
```

2. **Clone the Project:**
```bash
git clone https://github.com/yourusername/hostgator-cpanel-dns-automizer.git
cd hostgator-cpanel-dns-automizer

go install
```

3. **Set Up Environment Files**
Create an `.env` file on the root folder `.env file already created. create it if needed`

4. **.env Configuration**
*Location: .env*
```env
VERSION=1.2.0 #current version
INTERVAL=20

PUBLIC_IP_CHECKER="https://ipv4.icanhazip.com" #url to check public ipv4 change to your taste.
CONFIG_IP_REGISTRY_FILE="/Add/Your/Path/hostgator-cpanel-dns-automizer/config/ip.cfg"
CONFIG_IP_REGISTRY_PATH="/Add/Your/Path/hostgator-cpanel-dns-automizer/config"
```

5. **Adding domains and configuration**
Create a json file in the `./domains` folder. Inside the is going to be a example.txt with basic fields.
```json

{
    "domain_name" : "your_domain",
    "main_domain_url" : "your_domain.com",
    "target_domains" : ["example.your_domain.com.","www.example.your_domain.com."],
    "auth" : {
        "username" : "cpanel_username",
        "token" : "cpanel_token"
    },
    "cpanel_port" : 2083 // usually 2083 but check on DNS Zone url
}
```

6. **Create a Makefile**
```Makefile
SERVICE_NAME=dns-automizer

run:build
	@ ./bin/dns-automizer-artifact
	

build:
	@ echo --- BUILDING ARTIFACT ---
	@ go build -o ./bin/dns-automizer-artifact ./cmd/main.go
	@ echo --- ARTIFACT BUILT ---

start:
	@ sudo systemctl start $(SERVICE_NAME)

restart:
	@ sudo systemctl restart $(SERVICE_NAME)

deploy: build
	sudo systemctl daemon-reload
	sudo systemctl reload $(SERVICE_NAME)
	@ echo --- SERVICE DEPLOYED AND RUNNING ---

```

8. **Create a .service File**
```ini
# /etc/systemd/system/dns-automizer.service
[Unit]
Description=Hostgator DNS Automizer
After=network.target

[Service]
ExecStart=/path/to/hostgator-cpanel-dns-automizer/bin/dns-automizer-artifact
WorkingDirectory=/path/to/hostgator-cpanel-dns-automizer/
Restart=always
user=AddYouUser #optional

[Install]
WantedBy=multi-user.target
```
Enable and start it:
```bash
sudo systemctl daemon-reexec
sudo systemctl daemon-reload
sudo systemctl enable dns-automizer
sudo systemctl start dns-automizer
```

9. **Improve Makefile for Deployment**
```bash
# make sure to be in the project folder
make deploy 
```

---

### 🍎 macOS Setup

1. **Install Go (via Homebrew):**
```bash
brew install go
```

2. **Clone the repo:**
```bash
git clone https://github.com/yourusername/hostgator-cpanel-dns-automizer.git
cd hostgator-cpanel-dns-automizer

go install
```

3. **Set up your `.env` file and add your domains to be watch** (same as Linux)

4. **Run it manually:**
```bash
make run
```

5. *(Optional)* Create a `launchd` `.plist` file to run it as a background service.

---

### 🪟 Windows Setup

1. **Install Go:**
   - Download and install from [golang.org](https://golang.org/dl/)

2. **Install Git Bash or use PowerShell**

3. **Clone the repo:**
```bash
git clone https://github.com/yourusername/hostgator-cpanel-dns-automizer.git
cd hostgator-cpanel-dns-automizer

go install
```

4. **Set up your `.env` file and add your domains to be watch** like on Linux/macOS.

5. **Run the app (from Git Bash or terminal):**
```bash
make run
```

6. *(Optional)* Use Task Scheduler to run the compiled `.exe` as a background service on startup.

---

### 🛠 Enhancements

- Support for multiple domains.
- Environment variables are already set up
- Plug and Play style setup.
- Variable watch interval.(you can modify the watch interval from the .env file)

---

## 🤝 Contribute

Feel free to contribute! Submit PRs for new features (like GUI support, alternative DNS providers, etc.) or improve documentation.

Open an issue or fork the repo to start collaborating!

Tests & Updates will be added periodically
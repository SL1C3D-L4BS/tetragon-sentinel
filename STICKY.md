# 📌 Start here — Tetragon Sentinel

**One webhook. One pipe. Kernel alerts in Discord.**

### 1. Get a Discord webhook URL
Create a channel (e.g. `#ai-alerts`) → **Channel settings → Integrations → Webhooks → New Webhook** → copy URL.

### 2. Build and run
```bash
export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/YOUR_ID/YOUR_TOKEN"
go build -o tetragon-sentinel .
```

### 3. Live pipeline (Tetragon → Discord)
```bash
sudo tetragon observe --output json | ./tetragon-sentinel --alert-binary=/bin/bash --alert-file=/etc/shadow --alert-file=.env
```

### 4. Test without Tetragon (mock)
```bash
echo '{"process_exec":{"process":{"binary":"/bin/bash","arguments":"-c id","pid":12345}}}' | ./tetragon-sentinel --alert-binary=/bin/bash
```

When the red embed hits Discord, you’re live.  
Full docs → [README.md](README.md)

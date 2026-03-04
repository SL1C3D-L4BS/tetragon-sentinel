# 👁️ Tetragon Sentinel Bot

**Real-time eBPF Kernel Alerts for AI Infrastructure.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/vericore/tetragon-sentinel)](https://golang.org/)
[![eBPF Powered](https://img.shields.io/badge/Kernel-eBPF%20%2F%20Tetragon-black)](https://tetragon.io/)

Tetragon Sentinel is a lightweight, zero-dependency Go daemon that pipes eBPF kernel security events from [Tetragon](https://tetragon.io/) directly into your Discord or Slack channels.

If a sandboxed AI agent escapes its runtime and attempts to execute `/bin/bash` or read sensitive system files, Sentinel catches the kernel syscall in microseconds and drops a rich alert into your chat before the agent can complete the exfiltration.

## 🚀 Quickstart

**1. Set your Discord Webhook URL:** (or use `DISCORD_WEBHOOK_SENTINEL` if you use [Discord_Server_Config](docs/DISCORD-SETUP.md))
```bash
export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/..."
```

**2. Pipe Tetragon events directly into Sentinel:**
By adhering to the Unix philosophy, Sentinel reads standard input (`stdin`). You can pipe Tetragon's JSON output directly into the bot.

```bash
sudo tetragon observe --output json | ./tetragon-sentinel --alert-binary="/bin/bash" --alert-binary="/usr/bin/curl"
```

**3. (Phase 2) Trap sensitive file reads:** Alert when a process reads protected files (e.g. `.env`, `/etc/shadow`, AWS credentials):

```bash
sudo tetragon observe --output json | ./tetragon-sentinel \
  --alert-binary="/bin/bash" \
  --alert-file="/etc/shadow" --alert-file=".env" --alert-file="/root/.aws/credentials"
```

**Test without Tetragon (mock JSON):**
```bash
# Binary exec alert
echo '{"process_exec":{"process":{"binary":"/bin/bash","arguments":"-c cat /etc/shadow","pid":12345}}}' | ./tetragon-sentinel --alert-binary=/bin/bash

# File read alert (process_kprobe)
echo '{"process_kprobe":{"process":{"binary":"/usr/bin/cat","pid":101},"args":[{"file_arg":{"path":"/etc/shadow"}}]}}' | ./tetragon-sentinel --alert-file=/etc/shadow
```

## 🧠 Architecture

* **eBPF Tracing:** Relies on Tetragon for low-overhead, kernel-level enforcement.
* **Stream Parsing:** Uses Go's `bufio.Scanner` to parse high-throughput JSON-lines without blowing up system memory.
* **Rich Webhooks:** Formats kernel events (PID, Binary, Arguments, Namespace) into actionable Discord Embeds.

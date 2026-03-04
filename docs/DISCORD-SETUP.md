# Discord setup for Tetragon Sentinel

Tetragon Sentinel sends alerts to a Discord channel via an **Incoming Webhook**. You can either create the webhook manually or use the SL1C3D Discord server config so `npm run provision` creates it for you.

## Option 1: Manual webhook

1. In your Discord server, open **Server settings → Integrations → Webhooks** (or edit a channel → Integrations → Webhooks).
2. Create an **Incoming Webhook** in the channel where you want alerts (e.g. `#bot-logs` or a dedicated `#sentinel-alerts`).
3. Copy the webhook URL and set:
   ```bash
   export DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/..."
   ```

## Option 2: Use Discord_Server_Config (SL1C3D)

If you use the **Discord_Server_Config** repo (e.g. `~/Desktop/Discord_Server_Config`) to provision your server, a **Tetragon Sentinel** webhook is defined there. It posts to the **🤖・bot-logs** channel (09 ▸ 🔐 INTERNAL).

1. In `Discord_Server_Config`, run:
   ```bash
   npm run provision
   ```
2. The script creates the webhook and prints its URL (or writes it to `.env` / `.env.webhooks`).
3. Set **one** of these in the environment when running tetragon-sentinel:
   - `DISCORD_WEBHOOK_URL` — paste the printed webhook URL, or
   - `DISCORD_WEBHOOK_SENTINEL` — if your `.env` from Discord_Server_Config is loaded (e.g. `source ../Discord_Server_Config/.env`), Sentinel will use this if `DISCORD_WEBHOOK_URL` is unset.

### Webhook entry in Discord_Server_Config

In `config/webhooks.ts` the Sentinel webhook is:

| Name              | Channel     | Description                          |
|-------------------|------------|--------------------------------------|
| Tetragon Sentinel | 🤖・bot-logs | eBPF kernel alerts from Tetragon Sentinel |

Env key: `DISCORD_WEBHOOK_SENTINEL`.

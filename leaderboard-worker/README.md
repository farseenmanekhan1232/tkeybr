# Leaderboard Worker

A Cloudflare Worker that serves the tkeybr leaderboard API.

## Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Set up secrets:
   ```bash
   # Generate a random token for API authentication
   wrangler secret put API_SECRET_TOKEN
   
   # Your Firebase project ID
   wrangler secret put FIREBASE_PROJECT_ID
   
   # Base64-encoded service account JSON
   # Get from Firebase Console > Project Settings > Service Accounts > Generate New Private Key
   # Then run: cat service-account.json | base64 | wrangler secret put FIREBASE_SERVICE_ACCOUNT
   wrangler secret put FIREBASE_SERVICE_ACCOUNT
   ```

3. Deploy:
   ```bash
   npm run deploy
   ```

## Usage

The worker exposes a single endpoint:

```
GET /api/leaderboard
Authorization: Bearer <API_SECRET_TOKEN>
```

Returns JSON with daily, weekly, and all-time leaderboard rankings.

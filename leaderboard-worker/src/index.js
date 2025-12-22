/**
 * tkeybr Leaderboard API Worker
 * 
 * This Cloudflare Worker provides a secure API to fetch leaderboard data
 * from Firestore. It's called by GitHub Actions to generate static JSON.
 * 
 * Environment Variables (set via wrangler secret):
 * - API_SECRET_TOKEN: Bearer token for authentication
 * - FIREBASE_PROJECT_ID: Firebase project ID
 * - FIREBASE_SERVICE_ACCOUNT: Base64-encoded service account JSON
 */

export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);
    
    // CORS headers for preflight requests
    if (request.method === 'OPTIONS') {
      return new Response(null, {
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Methods': 'GET, OPTIONS',
          'Access-Control-Allow-Headers': 'Authorization, Content-Type',
        },
      });
    }

    // Only handle /api/leaderboard endpoint
    if (url.pathname !== '/api/leaderboard') {
      return new Response('Not Found', { status: 404 });
    }

    // Verify authentication
    const authHeader = request.headers.get('Authorization');
    if (!authHeader || authHeader !== `Bearer ${env.API_SECRET_TOKEN}`) {
      return new Response('Unauthorized', { status: 401 });
    }

    try {
      // Get access token for Firestore
      const accessToken = await getFirestoreAccessToken(env);
      
      // Fetch leaderboard data from Firestore
      const leaderboardData = await fetchLeaderboardData(env.FIREBASE_PROJECT_ID, accessToken);
      
      // Process and categorize the data
      const categorizedData = categorizeLeaderboard(leaderboardData);
      
      return new Response(JSON.stringify(categorizedData, null, 2), {
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
          'Cache-Control': 'public, max-age=3600',
        },
      });
    } catch (error) {
      console.error('Error fetching leaderboard:', error);
      return new Response(JSON.stringify({ error: error.message }), {
        status: 500,
        headers: { 'Content-Type': 'application/json' },
      });
    }
  },
};

/**
 * Get an access token for Firestore using a service account
 */
async function getFirestoreAccessToken(env) {
  const serviceAccount = JSON.parse(atob(env.FIREBASE_SERVICE_ACCOUNT));
  
  // Create JWT for service account
  const now = Math.floor(Date.now() / 1000);
  const jwt = await createJWT({
    iss: serviceAccount.client_email,
    sub: serviceAccount.client_email,
    aud: 'https://oauth2.googleapis.com/token',
    iat: now,
    exp: now + 3600,
    scope: 'https://www.googleapis.com/auth/datastore',
  }, serviceAccount.private_key);

  // Exchange JWT for access token
  const tokenResponse = await fetch('https://oauth2.googleapis.com/token', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({
      grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
      assertion: jwt,
    }),
  });

  const tokenData = await tokenResponse.json();
  if (!tokenResponse.ok) {
    throw new Error(`Token error: ${JSON.stringify(tokenData)}`);
  }
  
  return tokenData.access_token;
}

/**
 * Create a signed JWT using the Web Crypto API
 */
async function createJWT(payload, privateKeyPem) {
  const header = { alg: 'RS256', typ: 'JWT' };
  
  // Encode header and payload
  const headerB64 = btoa(JSON.stringify(header)).replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_');
  const payloadB64 = btoa(JSON.stringify(payload)).replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_');
  const unsignedToken = `${headerB64}.${payloadB64}`;
  
  // Import the private key
  const pemContents = privateKeyPem
    .replace('-----BEGIN PRIVATE KEY-----', '')
    .replace('-----END PRIVATE KEY-----', '')
    .replace(/\s/g, '');
  
  const binaryKey = Uint8Array.from(atob(pemContents), c => c.charCodeAt(0));
  
  const cryptoKey = await crypto.subtle.importKey(
    'pkcs8',
    binaryKey,
    { name: 'RSASSA-PKCS1-v1_5', hash: 'SHA-256' },
    false,
    ['sign']
  );
  
  // Sign the token
  const signature = await crypto.subtle.sign(
    'RSASSA-PKCS1-v1_5',
    cryptoKey,
    new TextEncoder().encode(unsignedToken)
  );
  
  const signatureB64 = btoa(String.fromCharCode(...new Uint8Array(signature)))
    .replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_');
  
  return `${unsignedToken}.${signatureB64}`;
}

/**
 * Fetch all leaderboard entries from Firestore
 */
async function fetchLeaderboardData(projectId, accessToken) {
  const firestoreUrl = `https://firestore.googleapis.com/v1/projects/${projectId}/databases/(default)/documents/leaderboard`;
  
  const response = await fetch(firestoreUrl, {
    headers: {
      'Authorization': `Bearer ${accessToken}`,
    },
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(`Firestore error: ${error}`);
  }

  const data = await response.json();
  
  if (!data.documents) {
    return [];
  }

  // Parse Firestore document format
  return data.documents.map(doc => {
    const fields = doc.fields || {};
    return {
      id: doc.name.split('/').pop(),
      displayName: fields.displayName?.stringValue || 'Anonymous',
      bestWPM: parseFloat(fields.bestWPM?.doubleValue || fields.bestWPM?.integerValue || 0),
      averageWPM: parseFloat(fields.averageWPM?.doubleValue || fields.averageWPM?.integerValue || 0),
      totalTimeSeconds: parseFloat(fields.totalTimeSeconds?.doubleValue || fields.totalTimeSeconds?.integerValue || 0),
      totalCharacters: parseInt(fields.totalCharacters?.integerValue || 0, 10),
      totalSessions: parseInt(fields.totalSessions?.integerValue || 0, 10),
      lastUpdated: fields.lastUpdated?.timestampValue || new Date().toISOString(),
    };
  });
}

/**
 * Categorize leaderboard into daily, weekly, and all-time rankings
 */
function categorizeLeaderboard(entries) {
  const now = new Date();
  const oneDayAgo = new Date(now - 24 * 60 * 60 * 1000);
  const oneWeekAgo = new Date(now - 7 * 24 * 60 * 60 * 1000);

  // Filter entries by time period
  const dailyEntries = entries.filter(e => new Date(e.lastUpdated) >= oneDayAgo);
  const weeklyEntries = entries.filter(e => new Date(e.lastUpdated) >= oneWeekAgo);

  // Sort by best WPM descending
  const sortByBestWPM = (a, b) => b.bestWPM - a.bestWPM;
  const sortByAvgWPM = (a, b) => b.averageWPM - a.averageWPM;
  const sortByPracticeTime = (a, b) => b.totalTimeSeconds - a.totalTimeSeconds;

  // Create ranked lists (top 50)
  const createRankedList = (list, sortFn) => {
    return list
      .sort(sortFn)
      .slice(0, 50)
      .map((entry, index) => ({
        rank: index + 1,
        displayName: entry.displayName,
        bestWPM: Math.round(entry.bestWPM * 10) / 10,
        avgWPM: Math.round(entry.averageWPM * 10) / 10,
        practiceHours: Math.round(entry.totalTimeSeconds / 3600 * 10) / 10,
        totalCharacters: entry.totalCharacters,
      }));
  };

  return {
    lastUpdated: now.toISOString(),
    categories: {
      allTime: {
        byBestWPM: createRankedList([...entries], sortByBestWPM),
        byAvgWPM: createRankedList([...entries], sortByAvgWPM),
        byPracticeTime: createRankedList([...entries], sortByPracticeTime),
      },
      weekly: {
        byBestWPM: createRankedList([...weeklyEntries], sortByBestWPM),
        byAvgWPM: createRankedList([...weeklyEntries], sortByAvgWPM),
      },
      daily: {
        byBestWPM: createRankedList([...dailyEntries], sortByBestWPM),
        byAvgWPM: createRankedList([...dailyEntries], sortByAvgWPM),
      },
    },
  };
}

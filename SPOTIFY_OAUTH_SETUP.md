# Spotify OAuth Setup Guide

This guide walks you through setting up Spotify OAuth for your Muse application following the official Spotify Authorization Code Flow.

## 1. Create Spotify Application

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Log in with your Spotify account
3. Click "Create an app"
4. Fill in:
   - **App name**: `Muse Music App` (or your preferred name)
   - **App description**: `A music review and playlist management application`
   - **Website**: `https://127.0.0.1:3000` (for development)
   - **Redirect URI**: `https://127.0.0.1:8080/auth/spotify/callback`
5. Check the boxes for Developer Terms of Service
6. Click "Create"

## 2. Configure Redirect URIs

⚠️ **Critical**: The redirect URI must match exactly between your Spotify app settings and your `.env` file.

### For Development:

- Spotify App Setting: `https://127.0.0.1:8080/auth/spotify/callback`
- Your `.env` file: `SPOTIFY_REDIRECT_URL=https://127.0.0.1:8080/auth/spotify/callback`

### For Production:

- Spotify App Setting: `https://yourdomain.com/auth/spotify/callback`
- Your `.env` file: `SPOTIFY_REDIRECT_URL=https://yourdomain.com/auth/spotify/callback`

## 3. Environment Variables Setup

Create a `.env` file in the root directory with these variables:

```env
# Environment Configuration
ENVIRONMENT=development
PORT=8080

# Spotify OAuth Configuration
# Get these from https://developer.spotify.com/dashboard
SPOTIFY_CLIENT_ID=your_spotify_client_id_here
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret_here

# This must match exactly with what you configure in your Spotify app settings
SPOTIFY_REDIRECT_URL=https://127.0.0.1:8080/auth/spotify/callback

# Database Configuration
DATABASE_URL=postgresql://postgres:password@localhost:5432/muse
DB_HOST=localhost
DB_PORT=5432
DB_NAME=muse
DB_USER=postgres
DB_PASSWORD=your_db_password
DB_SSL_MODE=prefer

# Redis Configuration
REDIS_URL=redis://localhost:6379
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Security
JWT_SECRET=your_super_secure_jwt_secret_key_here

# Frontend URL (for CORS)
FRONTEND_URL=https://127.0.0.1:3000
```

## 4. OAuth Flow Implementation

Your implementation follows the Spotify Authorization Code Flow:

### Step 1: Request User Authorization

```typescript
// Frontend: User clicks "Connect to Spotify"
const response = await apollo.query({
  query: GET_SPOTIFY_AUTH_URL,
});

// User is redirected to Spotify for authorization
window.location.href = response.data.spotifyAuthURL.url;
```

### Step 2: User Grants Permission

- User logs into Spotify (if not already logged in)
- User sees permission dialog with requested scopes:
  - `user-read-private`: Read user profile
  - `user-read-email`: Read user email
  - `playlist-read-private`: Read private playlists
  - `playlist-read-collaborative`: Read collaborative playlists
  - `user-library-read`: Read saved tracks

### Step 3: Spotify Redirects Back

- User accepts: `https://127.0.0.1:8080/auth/spotify/callback?code=AUTHORIZATION_CODE&state=STATE_PARAMETER`
- User denies: `https://127.0.0.1:8080/auth/spotify/callback?error=access_denied&state=STATE_PARAMETER`

### Step 4: Exchange Code for Token

```go
// Backend automatically exchanges authorization code for access token
token, err := spotifyClient.ExchangeCode(ctx, code)
```

### Step 5: Store User Credentials

```go
// Backend stores tokens in database for future API calls
user, err := authService.updateUserWithSpotifyCredentials(ctx, userID, spotifyUser, token)
```

## 5. Security Features

### State Parameter Protection

- ✅ Cryptographically secure random state generation
- ✅ User ID embedded in state for validation
- ✅ Timestamp-based expiry (15 minutes)
- ✅ Protection against CSRF attacks

### Token Management

- ✅ Secure token storage in PostgreSQL
- ✅ Automatic token refresh when expired
- ✅ Scope validation on token receipt

### Error Handling

- ✅ Comprehensive error mapping for all OAuth error cases
- ✅ User-friendly error messages
- ✅ Detailed logging for debugging

## 6. Scopes Explained

Your application requests these scopes:

| Scope                         | Description                  | Usage in Muse                        |
| ----------------------------- | ---------------------------- | ------------------------------------ |
| `user-read-private`           | Read user profile data       | Display user's Spotify profile       |
| `user-read-email`             | Read user's email            | Link Spotify account to Muse account |
| `playlist-read-private`       | Read private playlists       | Import user's private playlists      |
| `playlist-read-collaborative` | Read collaborative playlists | Import collaborative playlists       |
| `user-library-read`           | Read saved tracks            | Access user's liked songs            |

## 7. Testing the Integration

### Test Successful Flow:

1. Start your backend: `cd backend && go run .`
2. Start your frontend: `cd frontend && npm run dev`
3. Navigate to playlists page
4. Click "Connect to Spotify"
5. Grant permissions on Spotify
6. Should redirect back with success message

### Test Error Cases:

- **User denies permission**: Should show access denied message
- **Invalid state**: Should show security error message
- **Expired state**: Wait 15+ minutes, should show expiry message

## 8. Production Considerations

### Security:

- Use HTTPS for all redirect URIs
- Store client secret securely (environment variables, not in code)
- Implement rate limiting for OAuth endpoints
- Monitor for suspicious OAuth activity

### Monitoring:

- Log all OAuth flows for debugging
- Monitor token refresh rates
- Track authentication success/failure rates
- Set up alerts for OAuth errors

### PKCE (Optional but Recommended):

For enhanced security, especially for mobile apps, consider implementing PKCE:

```go
// Generate PKCE parameters
codeVerifier := generateCodeVerifier()
codeChallenge := generateCodeChallenge(codeVerifier)

// Use PKCE auth URL
authURL := client.GetAuthURLWithPKCE(state, codeChallenge)

// Exchange with PKCE
token, err := client.ExchangeCodeWithPKCE(ctx, code, codeVerifier)
```

## 9. Troubleshooting

### Common Issues:

1. **"INVALID_CLIENT: Invalid redirect URI"**

   - Solution: Ensure redirect URI matches exactly in Spotify app and `.env`

2. **"State parameter mismatch"**

   - Solution: Check state generation and validation logic

3. **"Token expired"**

   - Solution: Implement automatic token refresh (already included)

4. **"Insufficient scopes"**

   - Solution: Update requested scopes in client configuration

5. **"CORS errors"**
   - Solution: Ensure frontend URL is properly configured in CORS middleware

### Debug Mode:

Enable debug logging by setting:

```env
SPOTIFY_DEBUG=true
```

This will log detailed OAuth flow information for troubleshooting.

## 10. Next Steps

After successful OAuth setup:

1. Test playlist import functionality
2. Test token refresh mechanism
3. Implement additional Spotify API features
4. Set up monitoring and error tracking
5. Prepare for production deployment

For more information, refer to:

- [Spotify Web API Documentation](https://developer.spotify.com/documentation/web-api/)
- [OAuth 2.0 RFC](https://tools.ietf.org/html/rfc6749)
- [PKCE RFC](https://tools.ietf.org/html/rfc7636)

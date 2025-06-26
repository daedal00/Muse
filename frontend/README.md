# Muse Frontend - Backend Testing Interface

This is a simple Next.js frontend created to test your Muse backend GraphQL API. It provides a user-friendly interface to test all the major backend functionality.

## Features

- 🔍 **Search**: Test Spotify API integration for albums and artists
- 🔐 **Authentication**: User registration and login functionality
- 💿 **Albums**: Browse stored albums with pagination
- ⭐ **Reviews**: View and create album reviews
- 🎵 **Playlists**: Browse user playlists
- 📊 **GraphQL Integration**: Full Apollo Client setup with error handling

## Prerequisites

- Node.js (v18 or higher)
- Your Muse backend running on port 8080
- Ensure your backend is properly configured with:
  - PostgreSQL database
  - Redis (optional)
  - Spotify API credentials

## Installation

1. Install dependencies:

```bash
cd frontend
npm install
```

2. Set up environment variables (optional):

```bash
# Create .env.local file
echo "NEXT_PUBLIC_GRAPHQL_URL=http://localhost:8080/query" > .env.local
```

3. Start the development server:

```bash
npm run dev
```

The frontend will be available at `http://localhost:3000`

## Testing Your Backend

### 1. Start Your Backend First

Make sure your backend is running:

```bash
cd backend
go run .
```

Your backend should be accessible at `http://localhost:8080`

### 2. Test Authentication

1. Visit `/auth` to create a new user account
2. Login with your credentials
3. Check that the header shows your authenticated status

### 3. Test Search Functionality

1. Visit `/search`
2. Search for albums or artists (e.g., "Taylor Swift", "The Beatles")
3. Verify that results are returned from Spotify

### 4. Test Data Browsing

1. Visit `/albums`, `/reviews`, and `/playlists`
2. These may be empty initially - that's expected!
3. Check for proper error handling and empty states

### 5. Test GraphQL Playground

- Visit `http://localhost:8080` for the GraphQL playground
- Try running queries directly

## Available Pages

- `/` - Homepage with overview and backend connection status
- `/search` - Search albums and artists using Spotify API
- `/auth` - User registration and login
- `/albums` - Browse stored albums with pagination
- `/reviews` - View album reviews
- `/playlists` - Browse user playlists

## GraphQL Operations

The frontend includes pre-built GraphQL operations for:

### Queries

- `me` - Get current user information
- `searchAlbums` - Search albums via Spotify
- `searchArtists` - Search artists via Spotify
- `albums` - Get stored albums with pagination
- `reviews` - Get reviews with pagination
- `playlists` - Get playlists with pagination

### Mutations

- `createUser` - Register new user
- `login` - Authenticate user
- `createReview` - Create album review
- `createPlaylist` - Create new playlist

## Development Notes

### Apollo Client Configuration

- Configured to connect to `http://localhost:8080/query`
- Includes authentication headers
- Error handling for network issues

### Styling

- Uses Tailwind CSS for styling
- Responsive design
- Modern UI components

### TypeScript

- Full TypeScript support
- Type-safe GraphQL operations
- Proper error handling

## Troubleshooting

### Backend Connection Issues

1. Ensure backend is running on port 8080
2. Check CORS settings in your backend
3. Verify GraphQL endpoint is `/query`

### Authentication Issues

1. Check JWT secret configuration
2. Verify token storage in localStorage
3. Check browser console for errors

### Empty Data

- Initially, albums, reviews, and playlists will be empty
- Use the search functionality to populate data
- Create user accounts to test reviews and playlists

## Testing Checklist

Use this checklist to verify your backend functionality:

- [ ] Backend starts without errors
- [ ] GraphQL playground accessible at `http://localhost:8080`
- [ ] Frontend connects to backend successfully
- [ ] User registration works
- [ ] User login works and updates UI
- [ ] Search returns results from Spotify
- [ ] Album/artist search works for both types
- [ ] Error handling works for network issues
- [ ] Pagination works for data browsing
- [ ] Authentication state persists on page refresh

## Next Steps

This frontend provides a solid foundation for testing your backend. You can extend it by:

1. Adding review creation forms
2. Adding playlist creation/management
3. Implementing track playback
4. Adding more advanced search filters
5. Adding user profile management

The code is well-structured and commented, making it easy to understand and extend for your specific testing needs.

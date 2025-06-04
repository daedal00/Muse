// MongoDB initialization script for Muse project
// Development environment with sample data

// Load the main initialization script
load('/docker-entrypoint-initdb.d/init.js');

// Switch to the development database
db = db.getSiblingDB('muse_dev');

// Run the same initialization but for dev database
eval(cat('/docker-entrypoint-initdb.d/init.js').replace(/db = db\.getSiblingDB\('muse'\);/, "db = db.getSiblingDB('muse_dev');"));

// Insert sample data for development

// Sample user activity logs
db.user_activity_logs.insertMany([
    {
        user_id: 'user-1',
        activity_type: 'login',
        timestamp: new Date(Date.now() - 86400000), // 1 day ago
        ip_address: '192.168.1.100',
        user_agent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36'
    },
    {
        user_id: 'user-1',
        activity_type: 'song_play',
        timestamp: new Date(Date.now() - 3600000), // 1 hour ago
        metadata: {
            song_id: 'song-1',
            duration_played: 180,
            source: 'playlist'
        },
        ip_address: '192.168.1.100',
        user_agent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36'
    },
    {
        user_id: 'user-2',
        activity_type: 'playlist_create',
        timestamp: new Date(Date.now() - 7200000), // 2 hours ago
        metadata: {
            playlist_id: 'playlist-1',
            playlist_name: 'My Workout Mix'
        },
        ip_address: '192.168.1.101',
        user_agent: 'Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15'
    }
]);

// Sample song play history
db.song_play_history.insertMany([
    {
        user_id: 'user-1',
        song_id: 'song-1',
        played_at: new Date(Date.now() - 3600000),
        duration_listened: 180,
        total_duration: 219,
        completion_percentage: 82.2,
        source: 'playlist',
        context_id: 'playlist-1'
    },
    {
        user_id: 'user-1',
        song_id: 'song-2',
        played_at: new Date(Date.now() - 3000000),
        duration_listened: 259,
        total_duration: 259,
        completion_percentage: 100.0,
        source: 'album',
        context_id: 'album-1'
    },
    {
        user_id: 'user-2',
        song_id: 'song-1',
        played_at: new Date(Date.now() - 1800000),
        duration_listened: 120,
        total_duration: 219,
        completion_percentage: 54.8,
        source: 'search',
        context_id: null
    }
]);

// Sample playlist analytics
db.playlist_analytics.insertMany([
    {
        playlist_id: 'playlist-1',
        owner_id: 'user-1',
        date: new Date(Date.now() - 86400000),
        plays: 15,
        shares: 3,
        likes: 8,
        unique_listeners: ['user-1', 'user-2', 'user-3', 'user-4']
    },
    {
        playlist_id: 'playlist-2',
        owner_id: 'user-2',
        date: new Date(Date.now() - 86400000),
        plays: 7,
        shares: 1,
        likes: 12,
        unique_listeners: ['user-2', 'user-3', 'user-5']
    }
]);

// Sample user recommendations cache
db.user_recommendations_cache.insertMany([
    {
        user_id: 'user-1',
        recommendations: [
            {
                song_id: 'song-3',
                score: 0.87,
                reason: 'Based on your high rating of similar pop songs'
            },
            {
                song_id: 'song-4',
                score: 0.73,
                reason: 'Other users with similar taste also liked this'
            },
            {
                song_id: 'song-5',
                score: 0.65,
                reason: 'From artists you frequently listen to'
            }
        ],
        generated_at: new Date(Date.now() - 3600000),
        expires_at: new Date(Date.now() + 82800000), // 23 hours from now
        algorithm_version: 'v2.1.0'
    },
    {
        user_id: 'user-2',
        recommendations: [
            {
                song_id: 'song-1',
                score: 0.92,
                reason: 'Perfect match for your workout playlists'
            },
            {
                song_id: 'song-6',
                score: 0.78,
                reason: 'High energy songs like your favorites'
            }
        ],
        generated_at: new Date(Date.now() - 1800000),
        expires_at: new Date(Date.now() + 84600000), // 23.5 hours from now
        algorithm_version: 'v2.1.0'
    }
]);

// Sample search analytics
db.search_analytics.insertMany([
    {
        query: 'taylor swift',
        user_id: 'user-1',
        timestamp: new Date(Date.now() - 7200000),
        results_count: 25,
        clicked_results: [
            {
                result_id: 'song-1',
                result_type: 'song',
                position: 1
            },
            {
                result_id: 'artist-1',
                result_type: 'artist',
                position: 3
            }
        ]
    },
    {
        query: 'workout music',
        user_id: 'user-2',
        timestamp: new Date(Date.now() - 5400000),
        results_count: 42,
        clicked_results: [
            {
                result_id: 'playlist-workout-1',
                result_type: 'playlist',
                position: 2
            }
        ]
    },
    {
        query: 'billie eilish bad guy',
        timestamp: new Date(Date.now() - 3600000), // Anonymous search
        results_count: 8,
        clicked_results: [
            {
                result_id: 'song-4',
                result_type: 'song',
                position: 1
            }
        ]
    }
]);

print('Sample data inserted successfully into MongoDB for development environment'); 
// MongoDB initialization script for Muse project
// Production environment

// Switch to the muse database
db = db.getSiblingDB('muse');

// Create collections with validation schemas

// User activity logs collection
db.createCollection('user_activity_logs', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['user_id', 'activity_type', 'timestamp'],
            properties: {
                user_id: {
                    bsonType: 'string',
                    description: 'User ID from PostgreSQL'
                },
                activity_type: {
                    bsonType: 'string',
                    enum: ['login', 'logout', 'song_play', 'playlist_create', 'rating_add', 'playlist_share'],
                    description: 'Type of user activity'
                },
                timestamp: {
                    bsonType: 'date',
                    description: 'When the activity occurred'
                },
                metadata: {
                    bsonType: 'object',
                    description: 'Additional activity-specific data'
                },
                ip_address: {
                    bsonType: 'string',
                    description: 'User IP address'
                },
                user_agent: {
                    bsonType: 'string',
                    description: 'User agent string'
                }
            }
        }
    }
});

// Playlist analytics collection
db.createCollection('playlist_analytics', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['playlist_id', 'owner_id', 'date'],
            properties: {
                playlist_id: {
                    bsonType: 'string',
                    description: 'Playlist ID from PostgreSQL'
                },
                owner_id: {
                    bsonType: 'string',
                    description: 'Playlist owner user ID'
                },
                date: {
                    bsonType: 'date',
                    description: 'Date of the analytics record'
                },
                plays: {
                    bsonType: 'int',
                    minimum: 0,
                    description: 'Number of times playlist was played'
                },
                shares: {
                    bsonType: 'int',
                    minimum: 0,
                    description: 'Number of times playlist was shared'
                },
                likes: {
                    bsonType: 'int',
                    minimum: 0,
                    description: 'Number of likes/favorites'
                },
                unique_listeners: {
                    bsonType: 'array',
                    items: {
                        bsonType: 'string'
                    },
                    description: 'Array of unique user IDs who listened'
                }
            }
        }
    }
});

// Song play history collection
db.createCollection('song_play_history', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['user_id', 'song_id', 'played_at'],
            properties: {
                user_id: {
                    bsonType: 'string',
                    description: 'User ID from PostgreSQL'
                },
                song_id: {
                    bsonType: 'string',
                    description: 'Song ID from PostgreSQL'
                },
                played_at: {
                    bsonType: 'date',
                    description: 'When the song was played'
                },
                duration_listened: {
                    bsonType: 'int',
                    minimum: 0,
                    description: 'Duration listened in seconds'
                },
                total_duration: {
                    bsonType: 'int',
                    minimum: 0,
                    description: 'Total song duration in seconds'
                },
                completion_percentage: {
                    bsonType: 'double',
                    minimum: 0,
                    maximum: 100,
                    description: 'Percentage of song completed'
                },
                source: {
                    bsonType: 'string',
                    enum: ['playlist', 'search', 'recommendation', 'album'],
                    description: 'How the song was discovered'
                },
                context_id: {
                    bsonType: 'string',
                    description: 'ID of playlist, album, etc. that was the source'
                }
            }
        }
    }
});

// User recommendations cache collection
db.createCollection('user_recommendations_cache', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['user_id', 'recommendations', 'generated_at'],
            properties: {
                user_id: {
                    bsonType: 'string',
                    description: 'User ID from PostgreSQL'
                },
                recommendations: {
                    bsonType: 'array',
                    items: {
                        bsonType: 'object',
                        required: ['song_id', 'score'],
                        properties: {
                            song_id: {
                                bsonType: 'string'
                            },
                            score: {
                                bsonType: 'double',
                                minimum: 0,
                                maximum: 1
                            },
                            reason: {
                                bsonType: 'string'
                            }
                        }
                    },
                    description: 'Array of recommended songs with scores'
                },
                generated_at: {
                    bsonType: 'date',
                    description: 'When recommendations were generated'
                },
                expires_at: {
                    bsonType: 'date',
                    description: 'When recommendations expire'
                },
                algorithm_version: {
                    bsonType: 'string',
                    description: 'Version of recommendation algorithm used'
                }
            }
        }
    }
});

// Search analytics collection
db.createCollection('search_analytics', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required: ['query', 'timestamp'],
            properties: {
                query: {
                    bsonType: 'string',
                    description: 'Search query string'
                },
                user_id: {
                    bsonType: 'string',
                    description: 'User ID who performed the search (optional for anonymous)'
                },
                timestamp: {
                    bsonType: 'date',
                    description: 'When the search was performed'
                },
                results_count: {
                    bsonType: 'int',
                    minimum: 0,
                    description: 'Number of search results returned'
                },
                clicked_results: {
                    bsonType: 'array',
                    items: {
                        bsonType: 'object',
                        properties: {
                            result_id: { bsonType: 'string' },
                            result_type: { 
                                bsonType: 'string',
                                enum: ['song', 'album', 'artist', 'playlist']
                            },
                            position: { bsonType: 'int' }
                        }
                    },
                    description: 'Which search results were clicked'
                }
            }
        }
    }
});

// Create indexes for performance
db.user_activity_logs.createIndex({ 'user_id': 1, 'timestamp': -1 });
db.user_activity_logs.createIndex({ 'activity_type': 1, 'timestamp': -1 });

db.playlist_analytics.createIndex({ 'playlist_id': 1, 'date': -1 });
db.playlist_analytics.createIndex({ 'owner_id': 1, 'date': -1 });

db.song_play_history.createIndex({ 'user_id': 1, 'played_at': -1 });
db.song_play_history.createIndex({ 'song_id': 1, 'played_at': -1 });
db.song_play_history.createIndex({ 'source': 1, 'played_at': -1 });

db.user_recommendations_cache.createIndex({ 'user_id': 1 }, { unique: true });
db.user_recommendations_cache.createIndex({ 'expires_at': 1 }, { expireAfterSeconds: 0 });

db.search_analytics.createIndex({ 'timestamp': -1 });
db.search_analytics.createIndex({ 'query': 1, 'timestamp': -1 });
db.search_analytics.createIndex({ 'user_id': 1, 'timestamp': -1 });

print('MongoDB collections and indexes created successfully for Muse project');
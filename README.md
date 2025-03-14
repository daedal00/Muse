# Project Technical Document

## 1. Introduction

Our project aims to create a music-focused platform where users can:

- Rate songs and albums (similar to Letterboxd’s rating of films).  
- Curate and share playlists with other users.  
- Convert playlists across different streaming services (Spotify, Apple Music, etc.).  
- Maintain a personal “My Muse” page displaying ratings, top favorites, and recent activity.  
- Receive personalized recommendations based on ratings and listening habits.

---

## 2. Overall Workflow

### 2.1 User Journey

1. **Sign Up / Login**  
   - Users create accounts or log in using either standard credentials or via third-party OAuth (Spotify, Apple Music, etc.).

2. **My Muse Page**  
   - Upon login, users land on their “My Muse” page, which displays their ratings, favorite artists, and playlists.

3. **Rating & Reviewing**  
   - Users search for songs/albums, leave ratings (and possibly short reviews or comments).

4. **Playlist Creation & Sharing**  
   - Users create or import playlists.  
   - Playlists can be shared publicly, privately, or with specific friends.  
   - Other users can view and clone these playlists.

5. **Playlist Conversion**  
   - A dedicated feature allows users to convert a playlist to another streaming service.

6. **Recommendations**  
   - The system analyzes ratings and listening patterns to provide recommended songs/albums.

### 2.2 Data Flow Summary

1. **Front-End**  
   - Makes HTTP calls or opens WebSocket connections to the backend via an API Gateway/Reverse Proxy.  
   - Renders real-time data (e.g., updated ratings, new recommendations) using WebSockets or server-sent events.

2. **Backend Services**  
   - **User, Playlist, Rating, Recommendation, and Conversion** microservices handle specific domains.  
   - Each service reads/writes to the databases (PostgreSQL, MongoDB) and uses Redis for caching.

3. **External Integrations**  
   - The Conversion service uses OAuth 2.0 to communicate with external music streaming APIs.

---

## 3. Front-End Changes

### 3.1 Technology & Framework

- **React (TypeScript)** for the web front-end.  
- **React Native (TypeScript)** or **Expo** for cross-platform mobile apps.

### 3.2 Components / Pages

1. **Landing / Login Page**  
   - Allows user registration, login, or OAuth sign-in.  

2. **My Muse Page**  
   - Displays user’s top-rated songs, albums, or artists, plus a quick overview of recently created playlists.  

3. **Search & Discovery Page**  
   - Lets users search for music, filter by genre/artist, and discover trending or recommended songs/albums.  

4. **Playlist Management**  
   - Create new playlists, view existing ones, and share or edit them.  

5. **Playlist Conversion Modal**  
   - UI to convert an existing playlist to another service (e.g., Spotify → Apple Music).

### 3.3 State Management

- **Redux** or **React Query** for data fetching, caching, and managing global state.  
- **WebSocket Integration** (e.g., using Socket.io client) for real-time updates on rating changes or recommendation notifications.

---

## 4. Back-End Changes

### 4.1 Microservices Overview

1. **User Service**  
   - Manages user accounts, authentication, profiles, and the “My Muse” page data.  

2. **Playlist Service**  
   - Handles creation, reading, updating, and deletion of playlists.  

3. **Rating Service**  
   - Receives and stores ratings for songs/albums; updates the user’s “My Muse” page.  

4. **Recommendation Service**  
   - Periodically processes rating data to generate personalized recommendations.  
   - Could use machine learning libraries (TensorFlow, PyTorch, scikit-learn) for advanced similarity grouping.  

5. **Conversion/Integration Service**  
   - Interacts with external music streaming APIs (Spotify, Apple Music, etc.) using OAuth.  
   - Converts playlists into the appropriate format for each platform.

### 4.2 API Gateway / Routing

- **Nginx** or a dedicated gateway (e.g., **Kong**, **AWS API Gateway**) sits in front of the microservices.  
- Routes requests to the correct microservice based on the URL path.  
- Manages SSL termination, rate limiting, and load balancing.

### 4.3 Real-Time Updates

- **WebSockets** (e.g., **Socket.io**) for pushing live notifications (e.g., “Playlist converted successfully” or “New recommended songs available”).

---

## 5. APIs

### 5.1 Core Endpoints

1. **User Service**  
   - `POST /users` → Create a new user  
   - `POST /auth/login` → Authenticate a user (JWT-based or session-based)  
   - `GET /users/{id}` → Fetch user profile, including “My Muse” data  

2. **Playlist Service**  
   - `GET /playlists` → List user’s playlists  
   - `POST /playlists` → Create a new playlist  
   - `PUT /playlists/{id}` → Update a playlist  
   - `DELETE /playlists/{id}` → Remove a playlist  

3. **Rating Service**  
   - `POST /ratings` → Add a rating for a song/album  
   - `GET /ratings/{songId}` → Retrieve ratings for a particular song  

4. **Recommendation Service**  
   - `GET /recommendations` → Get user-specific recommended songs/albums  

5. **Conversion/Integration Service**  
   - `POST /conversion/playlist/{playlistId}` → Convert playlist to another streaming service

### 5.2 External Streaming APIs

- **Spotify, Apple Music, etc.**  
  - Use **OAuth 2.0** to securely retrieve or modify user playlists.  
  - Map track identifiers from one platform to another.

---

## 6. Tech Stack & Rationale

| Layer                      | Technology                         | Reason                                                                                   |
|----------------------------|------------------------------------|------------------------------------------------------------------------------------------|
| **Frontend**               | React (TypeScript)                 | Mature ecosystem, strong community, type safety with TS.                                 |
| **Mobile**                 | React Native (TypeScript)          | Cross-platform development with a shared codebase.                                       |
| **Backend Framework**      | Node.js + NestJS/Express (TS)      | Fast, flexible, and great for microservices. TS provides type safety.                    |
| **Database (Relational)**  | PostgreSQL                         | ACID compliance, strong for structured data (user profiles, ratings).                    |
| **Database (NoSQL)**       | MongoDB                            | Flexible schema for playlist details, activity logs, or large dynamic documents.         |
| **Caching**                | Redis                              | In-memory data store for sessions and frequently accessed data, reducing database load.  |
| **Containerization**       | Docker                             | Consistent deployment environment across dev, staging, and production.                   |
| **Orchestration**          | Kubernetes (K8s)                   | Automated deployment, scaling, and management of containerized services.                 |
| **CI/CD**                  | GitHub Actions                     | Automate testing, building, and deployment pipelines.                                    |
| **Monitoring & Logging**   | Prometheus, Grafana, ELK Stack     | Real-time metrics, performance monitoring, and log aggregation.                          |
| **Reverse Proxy / Gateway**| Nginx or Kong                      | SSL termination, routing, load balancing for microservices.                              |

---

## 7. Data Modeling

### 7.1 Relational Diagram (PostgreSQL)
<pre>
   ┌────────────┐         ┌─────────────┐       ┌─────────────┐
   │   USERS    │         │   RATINGS   │       │   ALBUMS    │
   ├────────────┤         ├─────────────┤       ├─────────────┤
   │ user_id(PK)│1───────*│rating_id(PK)│       │ album_id(PK)│
   │ name       │         │user_id (FK) ├*─────1┤ title       │
   │ email      │         │album_id (FK)│       │ artist      │
   │ password   │         │rating       │       │ ...         │
   └────────────┬         └─────────────┘       └─────────────┘
                1
                │
                │
                *
   ┌────────────▼────────────┐
   │       PLAYLISTS         │
   ├─────────────────────────┤
   │ playlist_id (PK)        │
   │ user_id (FK)            │
   │ title                   │
   │ description             │
   └─────────────────────────┘
</pre>


- **Users** table: Stores user information and references to their playlists and ratings.  
- **Ratings** table: Associates a user with an album (or song) and stores the rating value.  
- **Albums** table: Contains basic metadata about albums (which could be extended to include track-level details).  
- **Playlists** table: Ties a playlist to its owner.

### 7.2 NoSQL (MongoDB) Usage

1. **Playlist Documents**  
   - Store detailed track information for each playlist, particularly useful if track data is large or frequently changing.  
   - **Example Structure:**
     ```json
     {
       "_id": "<playlistObjectId>",
       "title": "My Chill Vibes",
       "owner": "<userId>",
       "tracks": [
         {
           "songId": "<someId>",
           "title": "<songTitle>",
           "artist": "<artistName>",
           // Additional song metadata
         },
         // More tracks...
       ]
     }
     ```

2. **Activity Logs**  
   - Store user activity in a flexible schema (e.g., logging when a user creates a playlist or rates a song).

---

## 8. Explanation of Decisions

1. **Microservices Architecture**  
   - **Reason:** Allows each domain (user, rating, playlist, etc.) to be developed and scaled independently, enhancing maintainability and fault isolation.

2. **Node.js with NestJS/Express (TypeScript)**  
   - **Reason:** Ensures consistency across the front-end and back-end through TypeScript, and provides a structured approach to building scalable server-side applications.

3. **PostgreSQL + MongoDB**  
   - **Reason:**  
     - **PostgreSQL** is used for transactional and relational data (users, ratings, album information).  
     - **MongoDB** offers flexibility for handling unstructured data (detailed playlist information, activity logs).

4. **Redis**  
   - **Reason:** Caches frequently requested data (e.g., top-rated songs, user sessions) to reduce latency and load on the primary databases.

5. **WebSockets (Socket.io)**  
   - **Reason:** Facilitates real-time notifications (e.g., updates on ratings, playlist conversion status) to improve user experience and engagement.

6. **Docker & Kubernetes**  
   - **Reason:**  
     - **Docker** ensures consistent environments across development, staging, and production.  
     - **Kubernetes** automates deployment, scaling, and management of containerized services.

7. **OAuth 2.0**  
   - **Reason:** Provides a secure, standard method for integrating with external music streaming APIs without exposing user credentials.

8. **CI/CD with GitHub Actions**  
   - **Reason:** Automates testing, building, and deployment processes, ensuring code quality and accelerating release cycles.

9. **Nginx / Kong**  
   - **Reason:** Serves as a reverse proxy or API gateway for SSL termination, request routing, and load balancing across microservices.

10. **Prometheus / Grafana / ELK**  
    - **Reason:** Essential for real-time monitoring, performance metrics, and centralized log aggregation to facilitate troubleshooting and scaling.

---

## Next Steps

1. **Detailed API Contracts:**  
   - Define request/response formats for each endpoint (e.g., using OpenAPI/Swagger documentation).

2. **Security & Authentication:**  
   - Finalize token-based authentication, user roles, and permissions.

3. **Deployment Strategy:**  
   - Outline staging and production deployment pipelines using Docker, Kubernetes, etc.

4. **Performance Considerations:**  
   - Implement caching, indexing, and consider read replicas for PostgreSQL to enhance performance.

5. **Testing Plan:**  
   - Establish end-to-end testing for critical features like playlist conversion, rating updates, and recommendation accuracy.

---

*This document provides a foundational outline of the architecture, data model, and technology choices for our Letterboxd-style music platform. It is intended to be a living document that evolves alongside the project.*

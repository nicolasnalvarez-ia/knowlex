# Twitter Bookmarks API (Go Backend)

## Setup Instructions

### Prerequisites
- Go 1.21 or higher
- PostgreSQL database (Supabase account)
- Twitter Developer Account with OAuth 2.0 app

### Environment Variables

Copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
```

Required variables:
- `SUPABASE_URL`: Your Supabase PostgreSQL connection string
- `SUPABASE_SERVICE_KEY`: Your Supabase service role key
- `JWT_SECRET`: A secure random string for JWT signing
- `TWITTER_CLIENT_ID`: Your Twitter OAuth 2.0 client ID
- `TWITTER_CLIENT_SECRET`: Your Twitter OAuth 2.0 client secret
- `BACKEND_URL`: Your backend URL (e.g., http://localhost:8080)
- `FRONTEND_URL`: Your frontend URL (e.g., http://localhost:5173)
- `PORT`: Server port (default: 8080)

### Database Setup

1. Create a Supabase project at https://supabase.com
2. Go to SQL Editor in your Supabase dashboard
3. Copy and paste the contents of `schema.sql`
4. Execute the SQL to create all tables and indexes

### Twitter OAuth Setup

1. Go to https://developer.twitter.com/en/portal/dashboard
2. Create a new OAuth 2.0 application
3. Add callback URL: `http://localhost:8080/api/auth/twitter/callback` (for development)
4. Enable "Read" permissions for tweets and bookmarks
5. Copy your Client ID and Client Secret

### Installation

```bash
# Install dependencies
go mod download

# Run the server
go run main.go
```

### API Endpoints

#### Authentication
- `GET /api/auth/twitter` - Initiate Twitter OAuth
- `GET /api/auth/twitter/callback` - OAuth callback
- `GET /api/auth/me` - Get current user (protected)
- `POST /api/auth/logout` - Logout (protected)

#### Bookmarks
- `GET /api/bookmarks` - Get all bookmarks (protected)
  - Query params: `page`, `page_size`, `category_id`
- `POST /api/bookmarks/import` - Import bookmarks from JSON (protected)
- `DELETE /api/bookmarks/:id` - Delete bookmark (protected)
- `GET /api/bookmarks/search?q=query` - Search bookmarks (protected)
- `POST /api/bookmarks/:id/category` - Assign category (protected)
- `DELETE /api/bookmarks/:id/category/:categoryId` - Remove category (protected)

#### Categories
- `GET /api/categories` - Get all categories (protected)
- `POST /api/categories` - Create category (protected)
- `PUT /api/categories/:id` - Update category (protected)
- `DELETE /api/categories/:id` - Delete category (protected)

#### Export
- `GET /api/export/bookmarks` - Export all bookmarks (protected)
- `GET /api/export/category/:id` - Export category bookmarks (protected)

#### User
- `DELETE /api/user/account` - Delete account and all data (protected)

### Development

```bash
# Run with auto-reload (install air first: go install github.com/cosmtrek/air@latest)
air
```

### Project Structure

```
backend/
├── main.go           # Entry point
├── auth/             # OAuth and JWT logic
│   ├── oauth.go
│   └── jwt.go
├── database/         # Database connection and queries
│   ├── db.go
│   └── queries.go
├── handlers/         # HTTP request handlers
│   ├── auth.go
│   ├── bookmarks.go
│   ├── categories.go
│   ├── export.go
│   └── user.go
├── middleware/       # HTTP middleware
│   ├── auth.go
│   └── logger.go
├── models/           # Data structures
│   └── models.go
└── schema.sql        # Database schema
```

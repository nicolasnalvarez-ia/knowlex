# Twitter Bookmarks SaaS - Quick Setup Guide

## Project Overview

This is a complete full-stack Twitter Bookmarks management application with:
- **Go Backend** (REST API with Twitter OAuth)
- **React Frontend** (SPA with Tailwind CSS)
- **Chrome Extension** (Automatic bookmark sync)

## Directory Structure

```
/
├── backend/          # Go REST API
├── frontend/         # React SPA
├── extension/        # Chrome Extension
└── README.md         # Main documentation
```

## Quick Start

### 1. Database Setup (Supabase)

1. Create a Supabase account at https://supabase.com
2. Create a new project
3. Go to SQL Editor
4. Run the SQL from `backend/schema.sql`
5. Get your connection string from Settings > Database

### 2. Twitter OAuth Setup

1. Go to https://developer.twitter.com/en/portal/dashboard
2. Create a new OAuth 2.0 application
3. Add callback URL: `http://localhost:8080/api/auth/twitter/callback`
4. Enable read permissions
5. Copy Client ID and Secret

### 3. Backend Setup

```bash
cd backend

# Copy and edit environment file
cp .env.example .env
# Edit .env with your credentials

# Download dependencies (requires internet)
go mod download

# Run the server
go run main.go
```

Backend runs on http://localhost:8080

### 4. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Copy and edit environment file
cp .env.example .env
# Edit VITE_API_URL if needed

# Run dev server
npm run dev
```

Frontend runs on http://localhost:5173

### 5. Chrome Extension Setup

1. Navigate to `chrome://extensions/`
2. Enable "Developer mode" (toggle in top-right)
3. Click "Load unpacked"
4. Select the `extension/` directory
5. Copy your auth token from Settings page in web app
6. Paste token into extension popup
7. Click "Save & Start Sync"

## Usage Flow

1. **Login**: Open http://localhost:5173 and click "Login with Twitter"
2. **Setup Extension**: 
   - Install Chrome extension (load unpacked)
   - Go to Settings page and copy your auth token
   - Open extension popup and paste token
   - Click "Save & Start Sync"
3. **Auto-Sync**: 
   - Extension automatically opens Twitter bookmarks page
   - Auto-scrolls to load all bookmarks
   - Syncs them to your Knowlex account
   - View bookmarks in dashboard!
4. **Organize**:
   - Create categories with colors and icons
   - Assign bookmarks to categories
   - Search across all bookmarks
5. **Export**: Download your organized bookmarks anytime

## Environment Variables

### Backend (.env)
```
SUPABASE_URL=postgresql://postgres:[PASSWORD]@db.[PROJECT].supabase.co:5432/postgres
SUPABASE_SERVICE_KEY=your-service-key
JWT_SECRET=your-secret-key
TWITTER_CLIENT_ID=your-client-id
TWITTER_CLIENT_SECRET=your-client-secret
BACKEND_URL=http://localhost:8080
FRONTEND_URL=http://localhost:5173
PORT=8080
```

### Frontend (.env)
```
VITE_API_URL=http://localhost:8080/api
```

## Features

### Backend (Go)
- Twitter OAuth 2.0 authentication
- JWT token-based sessions
- PostgreSQL with parameterized queries
- RESTful API with Gin framework
- CORS middleware
- Request logging

### Frontend (React)
- Landing page with Twitter login
- Dashboard with sidebar navigation
- Category management (CRUD)
- Bookmark grid with pagination
- Full-text search
- Export functionality
- Dark/Light mode
- Responsive design
- Toast notifications
- Token display for extension setup

### Chrome Extension
- Auto-sync from Twitter/X bookmarks page
- Extract tweet data directly from DOM
- Progress tracking with visual overlay
- Background worker for API calls
- Works on twitter.com and x.com
- Manifest V3 compatible

## API Endpoints

**Authentication:**
- GET `/api/auth/twitter` - Initiate OAuth
- GET `/api/auth/twitter/callback` - OAuth callback
- GET `/api/auth/me` - Get current user
- POST `/api/auth/logout` - Logout

**Bookmarks:**
- GET `/api/bookmarks` - List bookmarks
- POST `/api/bookmarks/import` - Import from JSON
- DELETE `/api/bookmarks/:id` - Delete bookmark
- GET `/api/bookmarks/search?q=query` - Search

**Categories:**
- GET `/api/categories` - List categories
- POST `/api/categories` - Create category
- PUT `/api/categories/:id` - Update category
- DELETE `/api/categories/:id` - Delete category

**Export:**
- GET `/api/export/bookmarks` - Export all
- GET `/api/export/category/:id` - Export category

**User:**
- DELETE `/api/user/account` - Delete account

## Troubleshooting

### Backend won't start
- Check SUPABASE_URL is correct
- Verify database tables are created
- Ensure port 8080 is available

### Frontend can't connect to backend
- Verify backend is running
- Check VITE_API_URL in .env
- Check CORS settings

### Extension not working
- Make sure you've entered a valid auth token
- Test connection in extension popup before syncing
- Ensure you're on twitter.com/i/bookmarks or x.com/i/bookmarks
- Check browser console for errors
- Try reloading the extension

### Twitter OAuth fails
- Verify callback URL matches exactly
- Check Client ID and Secret are correct
- Ensure Twitter app has read permissions

## Production Deployment

### Backend
- Deploy to Railway, Fly.io, or DigitalOcean
- Set environment variables
- Use production Supabase URL
- Update FRONTEND_URL to production domain

### Frontend
- Deploy to Vercel or Netlify
- Set VITE_API_URL to production backend
- Run `npm run build` for production build

### Extension (Future)
- Create proper icon files
- Update manifest with production API URL
- Zip extension folder
- Submit to Chrome Web Store

## Security Notes

- JWT tokens expire after 7 days
- CORS restricted to frontend origin
- All queries use user_id checks
- Parameterized SQL queries prevent injection
- No RLS bypass in application layer

## Tech Stack

- **Backend**: Go 1.21+, Gin, pgx, JWT, OAuth2
- **Frontend**: React 18, Vite, Tailwind CSS, Zustand, React Router
- **Database**: PostgreSQL (Supabase)
- **Extension**: Vanilla JavaScript, Manifest V3

## Next Steps

1. Set up your Supabase database
2. Configure Twitter OAuth
3. Start backend and frontend
4. Load Chrome extension
5. Login to the app and copy auth token
6. Configure extension with token
7. Extension will automatically sync your bookmarks!

For more details, see README.md files in each directory.

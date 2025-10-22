# Knowlex - X Bookmarks Manager

A full-stack web application for managing, organizing, and searching your X (formerly Twitter) bookmarks with AI-powered categorization and export capabilities.

## 🚀 Features

- **X OAuth Authentication** - Secure login with your X account
- **Chrome Extension** - Automatic bookmark sync from X
- **Smart Search** - Full-text search across all your bookmarks
- **Categories** - Organize bookmarks with custom colors and icons
- **Import/Export** - JSON-based data portability
- **Dark Mode** - Beautiful UI with light and dark themes
- **Responsive Design** - Works on desktop, tablet, and mobile

## 📁 Project Structure

```
twitter-bookmarks-saas/
├── backend/          # Go REST API
│   ├── auth/         # OAuth & JWT logic
│   ├── database/     # PostgreSQL queries
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # Auth & CORS middleware
│   ├── models/       # Data structures
│   └── main.go       # Entry point
├── frontend/         # React SPA
│   ├── src/
│   │   ├── components/  # UI components
│   │   ├── pages/       # Route pages
│   │   ├── lib/         # API & state
│   │   └── hooks/       # Custom hooks
│   └── package.json
└── extension/        # Chrome Extension
    ├── manifest.json  # Extension config
    ├── popup.html/.js # Extension UI
    ├── content.js     # Bookmark extraction
    └── background.js  # API communication
```

## 🛠️ Tech Stack

### Backend
- **Go 1.21+** - Fast, compiled backend
- **Gin** - Web framework
- **pgx** - PostgreSQL driver
- **JWT** - Token-based authentication
- **OAuth 2.0** - Twitter integration

### Frontend
- **React 18** - UI library
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **Zustand** - State management
- **React Router** - Navigation
- **Axios** - HTTP client

### Database
- **Supabase PostgreSQL** - Cloud database with RLS

### Extension
- **Manifest V3** - Modern Chrome extension
- **Auto-sync** - Automatic bookmark extraction
- **Background worker** - API communication

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- Supabase account
- Twitter Developer account
- Modern web browser (Chrome, Firefox, Safari, Edge)

### 1. Database Setup

1. Create a Supabase project at [supabase.com](https://supabase.com)
2. Copy `backend/schema.sql` to Supabase SQL Editor
3. Execute the SQL to create tables
4. Get your connection string from Settings > Database

### 2. Twitter OAuth Setup

1. Go to [Twitter Developer Portal](https://developer.twitter.com)
2. Create a new OAuth 2.0 application
3. Add callback URL: `http://localhost:8080/api/auth/twitter/callback`
4. Enable read permissions
5. Copy Client ID and Secret

### 3. Backend Setup

```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env with your credentials
# SUPABASE_URL, TWITTER_CLIENT_ID, etc.

# Install dependencies
go mod download

# Run server
go run main.go
```

Server runs on http://localhost:8080

### 4. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.example .env

# Edit .env
# VITE_API_URL=http://localhost:8080/api

# Run dev server
npm run dev
```

App runs on http://localhost:5173

### 5. Chrome Extension (Optional but Recommended)

See `extension/README.md` for detailed setup instructions.

Quick setup:
1. Open `chrome://extensions/`
2. Enable "Developer mode"
3. Click "Load unpacked" and select the `extension/` folder
4. Copy your auth token from Settings page
5. Paste token in extension popup
6. Extension will automatically sync bookmarks!

## 📖 Usage

1. **Login** - Click "Login with Twitter" on landing page
2. **Install Extension** - Load the Chrome extension
3. **Configure** - Copy token from Settings and paste in extension
4. **Auto-sync** - Extension automatically syncs your Twitter bookmarks
5. **Organize** - Create categories and assign bookmarks
6. **Search** - Find bookmarks instantly with search
7. **Export** - Download your data anytime

## 🔐 Security

- JWT tokens with 7-day expiration
- Password-less authentication via Twitter OAuth
- Row-level security equivalent in application layer
- CORS restricted to frontend origin
- Parameterized SQL queries prevent injection
- Extension uses secure token storage

## 🌐 Deployment

### Backend
- **Railway** / **Fly.io** / **DigitalOcean**
- Set environment variables
- Deploy Go binary

### Frontend
- **Vercel** / **Netlify**
- Set `VITE_API_URL` to production backend
- Deploy from `frontend/` directory

### Database
- Already hosted on Supabase
- Configure production connection string

## 📝 API Documentation

### Authentication
- `GET /api/auth/twitter` - Initiate OAuth
- `GET /api/auth/twitter/callback` - OAuth callback
- `GET /api/auth/me` - Get current user
- `POST /api/auth/logout` - Logout

### Bookmarks
- `GET /api/bookmarks` - List bookmarks (paginated)
- `POST /api/bookmarks/import` - Import from JSON
- `DELETE /api/bookmarks/:id` - Delete bookmark
- `GET /api/bookmarks/search?q=query` - Search bookmarks

### Categories
- `GET /api/categories` - List categories
- `POST /api/categories` - Create category
- `PUT /api/categories/:id` - Update category
- `DELETE /api/categories/:id` - Delete category

### Export
- `GET /api/export/bookmarks` - Export all
- `GET /api/export/category/:id` - Export category

## 🤝 Contributing

This is a personal project. Feel free to fork and customize!

## 📄 License

MIT License - see LICENSE file for details

## 🐛 Known Issues

- Twitter's DOM structure may change, requiring extension updates
- Large bookmark collections (10,000+) may take time to sync
- Extension requires manual token configuration (first time only)

## 🔮 Future Enhancements

- AI-powered auto-categorization
- Bookmark sharing and collections
- Browser notifications for new bookmarks
- Bulk operations (delete, move)
- Advanced filtering and sorting
- Analytics dashboard

## 💡 Support

For issues or questions:
1. Check the README files in each directory
2. Review the code comments
3. Open an issue on GitHub

---

Built with ❤️ for better Twitter bookmark management

# AI-Powered Bookmark Categorization

Your X bookmarks manager now includes intelligent AI categorization powered by Claude AI!

## ✨ What's New

### 1. **Automatic AI Categorization**
- Click one button to categorize all uncategorized bookmarks
- AI analyzes tweet content and suggests relevant categories
- Automatically creates new categories when needed
- Reuses your existing categories when appropriate

### 2. **Smart Category Management**
- Filter bookmarks by category in the sidebar
- Search across all bookmarks with keyword search
- Categories show bookmark counts
- Beautiful color-coded organization

### 3. **Full Text Search**
- Search by tweet content, author, or username
- Instant results as you type
- Works across all your bookmarks

## 🚀 Quick Start

### Backend Setup (Required for AI)

1. **Get Claude API Key**
   ```bash
   # Visit: https://console.anthropic.com/
   # Sign up and create an API key
   ```

2. **Add to Environment**
   Create/update `backend/.env`:
   ```env
   ANTHROPIC_API_KEY=sk-ant-api03-your-key-here
   DATABASE_URL=your_database_url
   JWT_SECRET=your_jwt_secret
   TWITTER_CLIENT_ID=your_twitter_client_id
   TWITTER_CLIENT_SECRET=your_twitter_client_secret
   FRONTEND_URL=http://localhost:5173
   ```

3. **Start Backend**
   ```bash
   cd backend
   go run main.go
   ```

### Frontend Setup

1. **Install Dependencies**
   ```bash
   cd frontend
   npm install
   ```

2. **Start Frontend**
   ```bash
   npm run dev
   ```

## 🎯 How to Use

### Auto-Categorize Your Bookmarks

1. Import your X bookmarks via the extension
2. Go to **Settings** page
3. Toggle **"Auto-categorize new imports"** to enable automatic categorization
4. (Optional) Click **"Auto-Categorize Bookmarks"** to process existing uncategorized items
5. Wait a few seconds for AI magic ✨
6. See your bookmarks organized into categories!

### Filter & Search

**Filter by Category:**
- Use the sidebar to select a category
- Only bookmarks in that category will show
- Click "All Bookmarks" to see everything

**Search by Keyword:**
- Type in the search bar at the top
- Search works on tweet text, author, username
- Results update as you type (300ms debounce)

## 💰 Cost

Using Claude 3 Haiku (cheapest, fast, good quality):
- **~$0.0001 per bookmark** (1/100th of a penny)
- **~$1 for 10,000 bookmarks**

Extremely affordable! Most users will spend < $0.10/month.

## 🏗️ Technical Details

### Backend APIs

**Auto-categorize endpoint:**
```
POST /api/ai/categorize
Authorization: Bearer <token>

Response:
{
  "message": "Auto-categorization completed",
  "data": {
    "categorized": 45,
    "new_categories": 5,
    "total_processed": 50
  }
}
```

**Search endpoint:**
```
GET /api/bookmarks/search?q=keyword&page=1&page_size=20
Authorization: Bearer <token>
```

**Filter by category:**
```
GET /api/bookmarks?category_id=uuid&page=1&page_size=20
Authorization: Bearer <token>
```

### Database Schema

The existing schema already supports everything:
- `bookmarks` - stores tweets
- `categories` - user-created categories
- `bookmark_categories` - many-to-many relationship

No migrations needed! ✅

### AI Categorization Logic

1. Fetches uncategorized bookmarks (no existing categories)
2. Sends tweet text to Claude with existing category context
3. Claude suggests 1-2 relevant categories
4. If category exists → assigns it
5. If new → creates with smart defaults (color/icon)
6. Batch processes up to 50 bookmarks per request

## 🎨 Category Defaults

When AI creates new categories, it assigns smart defaults:

| Category | Color | Icon |
|----------|-------|------|
| Tech | Blue #3B82F6 | 💻 |
| AI | Purple #8B5CF6 | 🤖 |
| Design | Pink #EC4899 | 🎨 |
| Coding | Green #10B981 | ⚡ |
| Business | Amber #F59E0B | 💼 |
| Marketing | Red #EF4444 | 📈 |
| News | Indigo #6366F1 | 📰 |
| Tutorial | Teal #14B8A6 | 📚 |
| Meme | Orange #F97316 | 😂 |
| Thread | Purple #A855F7 | 🧵 |

## 🔧 Files Modified/Added

### Backend
- ✅ `backend/ai/categorizer.go` - Claude API integration
- ✅ `backend/handlers/ai.go` - AI categorization endpoints
- ✅ `backend/database/queries.go` - Added `GetUncategorizedBookmarks`, `GetBookmarkByID`
- ✅ `backend/main.go` - Added `/api/ai/*` routes

### Frontend
- ✅ `frontend/src/pages/Settings.jsx` - Added "Auto-Categorize" button
- ✅ Existing search and filter functionality works out of the box!

### Documentation
- ✅ `backend/README_AI_SETUP.md` - Detailed AI setup guide
- ✅ `AI_FEATURES.md` - This file!

## ✅ Existing Features (Already Working!)

These features were already implemented and working:
- ✅ **Search by keyword** - `SearchBar.jsx` with debouncing
- ✅ **Filter by category** - `CategorySidebar.jsx` with event dispatch
- ✅ **Category management** - Create, edit, delete categories
- ✅ **Bookmark import** - Chrome extension integration
- ✅ **Export** - Download as JSON

## 🐛 Troubleshooting

**"ANTHROPIC_API_KEY not set"**
- Add the API key to `backend/.env`
- Restart the backend

**"No bookmarks categorized"**
- Make sure you have uncategorized bookmarks
- Check bookmarks have `tweet_text` field populated

**Search not working**
- Frontend search requires backend to be running
- Check API connection in browser console

**Categories not showing**
- Refresh the page after running auto-categorize
- Check backend logs for errors

## 📈 Next Steps

Potential enhancements:
- Background job queue for large batches (1000+ bookmarks)
- Confidence scores on AI suggestions
- Manual review before applying categories
- Category merge suggestions
- Multi-language tweet support

## 🎉 You're Ready!

1. Set up your `ANTHROPIC_API_KEY`
2. Import some bookmarks
3. Click "Auto-Categorize" in Settings
4. Enjoy organized bookmarks! 🚀

Questions? Check `backend/README_AI_SETUP.md` for detailed setup instructions.


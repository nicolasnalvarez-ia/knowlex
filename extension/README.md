# Knowlex Chrome Extension

Automatically sync your Twitter/X bookmarks to Knowlex for better organization and searchability.

## Features

- 🔄 **Auto-sync**: Automatically extracts and syncs your bookmarks
- 📊 **Progress tracking**: Visual feedback during bookmark extraction
- 🔐 **Secure**: Uses your personal auth token
- ⚡ **Fast**: Efficiently scrolls and extracts all bookmarks
- 🎯 **Smart**: Avoids duplicates and handles errors gracefully

## Installation (Development)

1. **Clone/Download** the extension folder

2. **Open Chrome Extensions**
   - Navigate to `chrome://extensions/`
   - Enable "Developer mode" (toggle in top-right)

3. **Load Extension**
   - Click "Load unpacked"
   - Select the `extension` folder

4. **Setup**
   - Click the extension icon in your toolbar
   - Enter your API URL (e.g., `http://localhost:8080`)
   - Copy your auth token from Knowlex Settings
   - Paste the token and click "Save & Start Sync"

5. **Automatic Sync**
   - Extension will automatically open your Twitter bookmarks page
   - It will scroll to load all bookmarks
   - Bookmarks will be synced to your Knowlex account
   - Check your dashboard to see them!

## Usage

### First Time Setup
1. Install the extension
2. Get your auth token from Knowlex app (Settings page)
3. Enter API URL and token in the extension popup
4. Click "Save & Start Sync"
5. Wait for automatic sync to complete

### Manual Sync
- Click the extension icon and then "Sync Now"
- **OR** use the "Refresh Bookmarks" button in the Knowlex dashboard (opens Twitter bookmarks automatically)
- Extension will open Twitter bookmarks and start syncing

## How It Works

1. **Content Script**: Runs on twitter.com/*/bookmarks pages
2. **Auto-scroll**: Automatically scrolls to load all bookmarks
3. **Extract**: Parses tweet data (text, author, media, timestamp)
4. **Sync**: Sends bookmarks to Knowlex API via background worker
5. **Complete**: Shows success message with count

## Permissions

- **storage**: Store API URL and auth token locally
- **tabs**: Open Twitter bookmarks page automatically
- **scripting**: Inject content script on bookmarks page
- **host_permissions**: Access twitter.com and x.com to extract bookmarks

## Configuration

The extension stores:
- `apiUrl`: Your Knowlex API endpoint
- `authToken`: Your personal auth token
- `autoStarted`: Flag to prevent repeated auto-starts
- `totalBookmarks`: Count of synced bookmarks
- `lastSync`: Timestamp of last sync

## Troubleshooting

### Extension not syncing
- Make sure you're logged into Twitter/X
- Check that your auth token is valid (test in popup)
- Verify API URL is correct
- Try manual sync from popup

### Missing bookmarks
- Extension only syncs visible bookmarks
- Make sure page scrolled to bottom
- Re-run sync to catch any missed bookmarks

### Connection errors
- Check that backend is running
- Verify CORS settings allow extension requests
- Test connection in popup before syncing

## Icons

Place icon files in the `icons/` directory:
- `icon16.png` - 16x16px
- `icon48.png` - 48x48px
- `icon128.png` - 128x128px

You can use any icon generator or design tool to create these.

## Development

### File Structure
```
extension/
├── manifest.json       # Extension configuration
├── popup.html          # Extension popup UI
├── popup.js            # Popup logic and settings
├── content.js          # Twitter page bookmark extraction
├── background.js       # Background worker for API calls
├── icons/              # Extension icons
└── README.md           # This file
```

### Testing
1. Make changes to files
2. Go to `chrome://extensions/`
3. Click refresh icon on Knowlex extension
4. Test functionality

## Publishing (Future)

To publish to Chrome Web Store:
1. Create proper icon files (16, 48, 128px)
2. Update manifest with final URLs
3. Create ZIP of extension folder
4. Submit to Chrome Web Store
5. Pay $5 developer fee
6. Wait for review (usually 1-3 days)

## Support

For issues or questions:
- Check the main Knowlex README
- Review browser console for errors
- Check background service worker logs
- Verify API is accessible

## License

Same as main Knowlex project (MIT)


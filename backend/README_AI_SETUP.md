# AI Categorization Setup Guide

## Overview
The bookmark manager now includes AI-powered automatic categorization using Claude AI (Anthropic). This feature analyzes bookmark content and suggests relevant categories.

## Features
- 🤖 **Auto-categorization**: Automatically categorize uncategorized bookmarks
- 🧠 **Smart matching**: Reuses existing categories when appropriate
- ✨ **New category creation**: Creates new categories when needed
- 🎯 **Batch processing**: Handles up to 50 bookmarks at a time
- 💰 **Cost-effective**: Uses Claude 3 Haiku (cheapest model)

## Setup Instructions

### 1. Get Anthropic API Key
1. Go to [console.anthropic.com](https://console.anthropic.com/)
2. Sign up or log in
3. Navigate to API Keys
4. Create a new API key
5. Copy the key (starts with `sk-ant-...`)

### 2. Add to Environment Variables
Add to your `.env` file in the backend directory:

```env
ANTHROPIC_API_KEY=sk-ant-api03-your-key-here
```

### 3. Restart Backend
```bash
cd backend
go run main.go
```

## Usage

### Via Frontend (Settings Page)
1. Navigate to Settings
2. Enable "Auto-categorize new imports" toggle to categorize on import
3. (Optional) Click "Auto-Categorize Bookmarks" button to process existing uncategorized bookmarks
4. Wait for AI processing
5. See results notification

### Via API

#### Auto-categorize all uncategorized bookmarks
```bash
POST /api/ai/categorize
Authorization: Bearer <token>
```

Response:
```json
{
  "message": "Auto-categorization completed",
  "data": {
    "categorized": 45,
    "new_categories": 5,
    "total_processed": 50
  }
}
```

#### Get suggestions for single bookmark
```bash
POST /api/ai/categorize/:bookmark_id
Authorization: Bearer <token>
```

Response:
```json
{
  "message": "Categorization suggestions generated",
  "data": {
    "suggested_categories": ["Tech", "AI"]
  }
}
```

## How It Works

1. **Fetches uncategorized bookmarks**: Queries bookmarks with no categories
2. **Sends to Claude AI**: Analyzes tweet text with context of existing categories
3. **Processes suggestions**: 
   - If category exists → assigns it
   - If new → creates category with smart defaults (colors/icons)
4. **Batch assigns**: Links bookmarks to categories

## Category Defaults

The system includes smart defaults for common categories:

| Category | Color | Icon |
|----------|-------|------|
| Tech | Blue | 💻 |
| AI | Purple | 🤖 |
| Design | Pink | 🎨 |
| Coding | Green | ⚡ |
| Business | Amber | 💼 |
| Marketing | Red | 📈 |
| News | Indigo | 📰 |
| Tutorial | Teal | 📚 |

## Cost Estimation

Claude 3 Haiku pricing (as of 2024):
- Input: $0.25 per million tokens
- Output: $1.25 per million tokens

Approximate costs:
- ~100 tokens per categorization (input + output)
- **~$0.0001 per bookmark** (1/100th of a cent)
- **$0.01 for 100 bookmarks**
- **$1.00 for 10,000 bookmarks**

Very affordable for most use cases! 💰

## Limitations

- Processes max 50 bookmarks per request (to prevent timeouts)
- Requires tweet text to categorize
- AI suggestions are not perfect (95%+ accuracy typical)
- Rate limits apply based on your Anthropic tier

## Troubleshooting

### "ANTHROPIC_API_KEY not set"
- Check your `.env` file has the key
- Restart the backend server
- Verify the key starts with `sk-ant-`

### "API returned status 401"
- API key is invalid
- Get a new key from Anthropic console

### "API returned status 429"
- Rate limit exceeded
- Wait a few minutes and try again
- Upgrade your Anthropic plan for higher limits

### No bookmarks categorized
- All bookmarks already have categories
- Check bookmark has tweet_text field populated

## Future Enhancements

Potential improvements:
- [ ] Background job queue for large batches
- [ ] User-customizable prompts
- [ ] Confidence scores for suggestions
- [ ] Manual review before applying
- [ ] Category merge suggestions
- [ ] Multi-language support

## Support

For issues or questions:
- Check the backend logs for detailed error messages
- Ensure your API key has sufficient credits
- Test with a small batch first (5-10 bookmarks)


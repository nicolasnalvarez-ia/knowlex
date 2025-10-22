# Twitter Bookmarks Frontend

React SPA for managing Twitter bookmarks with categories, search, and export functionality.

## Setup

### Prerequisites
- Node.js 18+ and npm

### Installation

```bash
npm install
```

### Environment Variables

Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

Update `VITE_API_URL` with your backend API URL.

### Development

```bash
npm run dev
```

The app will run on http://localhost:5173

### Build for Production

```bash
npm run build
```

The production build will be in the `dist` folder.

## Features

- Twitter OAuth authentication
- Search bookmarks with real-time results
- Create and manage categories with colors and icons
- Assign bookmarks to categories
- One-click import via bookmarklet
- Export bookmarks as JSON
- Dark/Light mode toggle
- Responsive design

## Tech Stack

- React 18
- Vite
- React Router
- Axios
- Zustand (state management)
- Tailwind CSS
- Lucide React (icons)
- React Hot Toast (notifications)
- Headless UI (components)

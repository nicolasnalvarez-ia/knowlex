# Cyberpunk Design System - Backup

This file contains the cyberpunk redesign for future reference. Ask the AI to reimplement this design when ready.

## Design Concept
- Deep dark backgrounds (charcoal, near-black, dark navy) with electric accent colors
- Neon highlights: cyan (#00ffff) and magenta (#ff00ff) used sparingly
- Sleek geometric UI with sharp edges and glowing borders
- Holographic/iridescent effects on interactive elements
- Monospace fonts (JetBrains Mono) mixed with clean sans-serif (Inter)
- Subtle grid overlays and scan-line effects
- Smooth, precise animations

## Files Modified

### 1. tailwind.config.js
```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        cyber: {
          black: '#0a0a0a',
          dark: '#0f1729',
          darker: '#141b2d',
          navy: '#1a2332',
          cyan: '#00ffff',
          'cyan-dark': '#00cccc',
          magenta: '#ff00ff',
          'magenta-dark': '#cc00cc',
          purple: '#8b5cf6',
          slate: '#1e293b',
          border: '#1f2937',
        }
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'Fira Code', 'Courier New', 'monospace'],
        tech: ['Inter', 'system-ui', 'sans-serif'],
      },
      boxShadow: {
        'neon-cyan': '0 0 5px theme("colors.cyber.cyan"), 0 0 20px theme("colors.cyber.cyan")',
        'neon-magenta': '0 0 5px theme("colors.cyber.magenta"), 0 0 20px theme("colors.cyber.magenta")',
        'neon-cyan-lg': '0 0 10px theme("colors.cyber.cyan"), 0 0 40px theme("colors.cyber.cyan"), 0 0 80px theme("colors.cyber.cyan")',
        'neon-magenta-lg': '0 0 10px theme("colors.cyber.magenta"), 0 0 40px theme("colors.cyber.magenta"), 0 0 80px theme("colors.cyber.magenta")',
        'cyber': '0 0 15px rgba(0, 255, 255, 0.3)',
        'cyber-hover': '0 0 20px rgba(0, 255, 255, 0.5), 0 0 40px rgba(0, 255, 255, 0.3)',
      },
      animation: {
        'pulse-slow': 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'glow': 'glow 2s ease-in-out infinite alternate',
        'scan': 'scan 4s linear infinite',
        'flicker': 'flicker 3s linear infinite',
        'slide-in': 'slideIn 0.3s ease-out',
      },
      keyframes: {
        glow: {
          '0%': { boxShadow: '0 0 5px theme("colors.cyber.cyan"), 0 0 10px theme("colors.cyber.cyan")' },
          '100%': { boxShadow: '0 0 10px theme("colors.cyber.cyan"), 0 0 20px theme("colors.cyber.cyan"), 0 0 30px theme("colors.cyber.cyan")' },
        },
        scan: {
          '0%': { transform: 'translateY(-100%)' },
          '100%': { transform: 'translateY(100%)' },
        },
        flicker: {
          '0%, 100%': { opacity: '1' },
          '41%': { opacity: '1' },
          '42%': { opacity: '0.8' },
          '43%': { opacity: '1' },
          '45%': { opacity: '0.9' },
          '46%': { opacity: '1' },
        },
        slideIn: {
          '0%': { transform: 'translateY(-10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
      backgroundImage: {
        'grid-cyber': 'linear-gradient(rgba(0, 255, 255, 0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(0, 255, 255, 0.1) 1px, transparent 1px)',
        'gradient-cyber': 'linear-gradient(135deg, #0a0a0a 0%, #0f1729 50%, #141b2d 100%)',
      },
      backgroundSize: {
        'grid': '20px 20px',
      },
    },
  },
  plugins: [],
}
```

### 2. index.css - Custom Components and Animations
See the full CSS in the committed version with:
- `.cyber-grid` - Grid overlay effect
- `.scan-line` - Animated scan line
- `.cyber-border` - Glowing border on hover
- `.holographic` - Holographic gradient effect
- `.neon-text-cyan` / `.neon-text-magenta` - Neon text effects
- `.cyber-btn` - Cyberpunk button style
- Custom scrollbar styling
- Glitch animation effects

### 3. Component Design Patterns

All components follow these patterns:
- Background: `rgba(10, 10, 10, 0.8)` with backdrop blur
- Borders: `rgba(0, 255, 255, 0.2)` cyan or `rgba(255, 0, 255, 0.2)` magenta
- Hover effects: Increase border opacity and add glow shadow
- Clip-path for angled corners: `polygon(0 0, calc(100% - 20px) 0, 100% 20px, 100% 100%, 20px 100%, 0 calc(100% - 20px))`
- Font: `font-mono` for technical elements, regular Inter for body text
- Animations: `animate-pulse`, `animate-flicker`, `scan-line` class

### Key Features to Implement:
1. **LandingPage**: Full cyberpunk hero with animated grid, neon accents, geometric feature cards
2. **DashboardLayout**: Futuristic header with status indicators, glowing borders
3. **BookmarkCard**: Holographic effects, cyber borders, scan lines
4. **SearchBar**: Neon glow on focus, cyber styling
5. **CategorySidebar**: Tech-inspired with glowing indicators
6. **Settings**: Futuristic panels with clip-path styling
7. **All Modals**: Dark background with cyber borders and animations

## To Reimplement:
Just ask: "Apply the cyberpunk design from the backup file" and the AI will restore all the changes.


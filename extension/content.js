// Content script for X bookmarks page
console.log('Knowlex: Content script loaded on X bookmarks page');

let isExtracting = false;
let extractedTweetIds = new Set();
let collectedBookmarks = {};

const urlParams = new URLSearchParams(window.location.search);
const shouldTriggerSync = urlParams.get('knowlex_sync') === '1';

// Check if we should auto-start extraction
chrome.storage.local.get(['apiUrl', 'authToken', 'autoStarted'], (result) => {
  if (result.apiUrl && result.authToken) {
    if (shouldTriggerSync) {
      console.log('Knowlex: Sync trigger detected from knowlex_sync parameter');
      chrome.storage.local.set({ autoStarted: false });
      setTimeout(() => {
        startExtraction();
      }, 2000);
    } else if (!result.autoStarted) {
      // mark as auto-started and run initial sync once
      chrome.storage.local.set({ autoStarted: true });
      setTimeout(() => {
        startExtraction();
      }, 3000);
    }
  }
});

// Listen for messages from popup or background
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'START_EXTRACTION') {
    startExtraction();
    sendResponse({ success: true });
  }
  return true;
});

async function startExtraction() {
  if (isExtracting) {
    console.log('Knowlex: Extraction already in progress');
    return;
  }
  
  isExtracting = true;
  extractedTweetIds.clear();
  collectedBookmarks = {};
  
  console.log('Knowlex: Starting bookmark extraction...');
  
  // Show overlay
  showOverlay('Initializing...', 0);
  
  try {
    // Auto-scroll to load all bookmarks
    await autoScroll();
    
    // Collect bookmarks one final time in case new ones appeared after scrolling stopped
    collectVisibleBookmarks();
    
    const bookmarks = Object.values(collectedBookmarks);
    
    console.log(`Knowlex: Extracted ${bookmarks.length} bookmarks`);
    
    if (bookmarks.length === 0) {
      showOverlay('No bookmarks found', 100, 'error');
      setTimeout(hideOverlay, 3000);
      isExtracting = false;
      return;
    }
    
    // Send to background for API sync
    showOverlay(`Syncing ${bookmarks.length} bookmarks...`, 50);
    
    chrome.runtime.sendMessage({
      type: 'SYNC_BOOKMARKS',
      bookmarks: bookmarks
    }, (response) => {
      if (response && response.success) {
        showOverlay(`✅ Successfully synced ${bookmarks.length} bookmarks!`, 100, 'success');
      } else {
        showOverlay(`❌ Sync failed: ${response?.error || 'Unknown error'}`, 100, 'error');
      }
      
      setTimeout(() => {
        hideOverlay();
        isExtracting = false;
      }, 3000);
    });
    
  } catch (error) {
    console.error('Knowlex: Extraction error:', error);
    showOverlay(`❌ Error: ${error.message}`, 100, 'error');
    setTimeout(() => {
      hideOverlay();
      isExtracting = false;
    }, 3000);
  }
}

async function autoScroll() {
  const scrollingElement = document.scrollingElement || document.documentElement;
  let scrollCount = 0;
  let stableAttempts = 0;
  let previousTweetCount = 0;

  while (scrollCount < 150 && stableAttempts < 8) {
    scrollCount++;
    const currentTweetCount = document.querySelectorAll('article[data-testid="tweet"]').length;

    scrollingElement.scrollTo({ top: scrollingElement.scrollHeight, behavior: 'smooth' });
    showOverlay(`Loading bookmarks... (scroll ${scrollCount})`, Math.min(scrollCount * 2, 45));

    await new Promise((resolve) => setTimeout(resolve, 1200));

    collectVisibleBookmarks();
    const newTweetCount = document.querySelectorAll('article[data-testid="tweet"]').length;

    if (newTweetCount <= currentTweetCount) {
      stableAttempts++;
    } else {
      stableAttempts = 0;
      previousTweetCount = newTweetCount;
    }

    // If we haven't increased the tweet count in the last few attempts, break out
    if (stableAttempts >= 8) {
      console.log('Knowlex: No more new tweets loaded after multiple attempts');
      break;
    }
  }

  // Give page a moment to finalize layout before extraction
  await new Promise((resolve) => setTimeout(resolve, 1500));
}

function collectVisibleBookmarks() {
  const articles = document.querySelectorAll('article[data-testid="tweet"]');
  
  console.log(`Knowlex: Found ${articles.length} tweet articles`);
  
  articles.forEach((article, index) => {
    try {
      // Extract tweet URL and ID
      const tweetLink = article.querySelector('a[href*="/status/"]');
      if (!tweetLink) return;
      
      const tweetUrl = tweetLink.href;
      const tweetId = tweetUrl.split('/status/')[1]?.split('?')[0];
      
      if (!tweetId || extractedTweetIds.has(tweetId)) return;
      extractedTweetIds.add(tweetId);
      
      // Extract tweet text
      const textElement = article.querySelector('[data-testid="tweetText"]');
      const tweetText = textElement ? textElement.innerText : '';
      
      // Extract author info
      const authorLink = article.querySelector('[data-testid="User-Name"] a');
      const authorUsername = authorLink ? authorLink.href.split('/').pop() : '';
      
      const authorNameElement = article.querySelector('[data-testid="User-Name"] span');
      const authorDisplayName = authorNameElement ? authorNameElement.innerText : '';
      
      // Extract media URLs
      const images = article.querySelectorAll('img[src*="twimg.com/media"]');
      const mediaUrls = Array.from(images).map(img => img.src);
      
      // Extract timestamp
      const timeElement = article.querySelector('time');
      const bookmarkedAt = timeElement ? timeElement.getAttribute('datetime') : new Date().toISOString();
      
      if (tweetId) {
        collectedBookmarks[tweetId] = {
          tweet_id: tweetId,
          tweet_text: tweetText,
          author_username: authorUsername,
          author_display_name: authorDisplayName,
          tweet_url: tweetUrl,
          media_urls: mediaUrls,
          bookmarked_at: bookmarkedAt
        };
      }
    } catch (error) {
      console.error('Knowlex: Error parsing tweet:', error);
    }
  });
}

// Overlay functions
function showOverlay(message, progress, type = 'info') {
  let overlay = document.getElementById('knowlex-overlay');
  
  if (!overlay) {
    overlay = document.createElement('div');
    overlay.id = 'knowlex-overlay';
    overlay.style.cssText = `
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(0, 0, 0, 0.8);
      z-index: 999999;
      display: flex;
      align-items: center;
      justify-content: center;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    `;
    
    overlay.innerHTML = `
      <div style="
        background: white;
        padding: 30px 40px;
        border-radius: 12px;
        text-align: center;
        max-width: 400px;
        box-shadow: 0 10px 40px rgba(0,0,0,0.3);
      ">
        <h2 style="margin: 0 0 15px 0; font-size: 24px; color: #333;">
          📚 Knowlex Sync
        </h2>
        <p id="knowlex-message" style="margin: 0 0 20px 0; color: #666; font-size: 16px;">
          ${message}
        </p>
        <div style="
          width: 100%;
          height: 8px;
          background: #e0e0e0;
          border-radius: 4px;
          overflow: hidden;
        ">
          <div id="knowlex-progress" style="
            height: 100%;
            background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
            width: ${progress}%;
            transition: width 0.3s ease;
          "></div>
        </div>
      </div>
    `;
    
    document.body.appendChild(overlay);
  } else {
    const messageEl = overlay.querySelector('#knowlex-message');
    const progressEl = overlay.querySelector('#knowlex-progress');
    if (messageEl) messageEl.textContent = message;
    if (progressEl) progressEl.style.width = `${progress}%`;
  }
  
  if (type === 'success') {
    const messageEl = overlay.querySelector('#knowlex-message');
    if (messageEl) {
      messageEl.style.color = '#28a745';
      messageEl.style.fontWeight = 'bold';
    }
  } else if (type === 'error') {
    const messageEl = overlay.querySelector('#knowlex-message');
    if (messageEl) {
      messageEl.style.color = '#dc3545';
      messageEl.style.fontWeight = 'bold';
    }
  }
}

function hideOverlay() {
  const overlay = document.getElementById('knowlex-overlay');
  if (overlay) {
    overlay.remove();
  }
}


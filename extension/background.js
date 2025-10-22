// Background service worker for Knowlex extension
console.log('Knowlex: Background service worker loaded');

// Listen for messages from content script or popup
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'SYNC_BOOKMARKS') {
    handleBookmarkSync(message.bookmarks)
      .then(result => sendResponse(result))
      .catch(error => sendResponse({ success: false, error: error.message }));
    return true; // Keep channel open for async response
  } else if (message.type === 'TEST_CONNECTION') {
    testConnection(message.apiUrl, message.token)
      .then(result => sendResponse(result))
      .catch(error => sendResponse({ success: false, error: error.message }));
    return true; // Keep channel open for async response
  }
});

// Handle bookmark synchronization
async function handleBookmarkSync(bookmarks) {
  try {
    // Get API URL and auth token from storage
    const config = await chrome.storage.local.get(['apiUrl', 'authToken']);
    
    if (!config.apiUrl || !config.authToken) {
      throw new Error('API URL or auth token not configured');
    }
    
    console.log(`Knowlex: Syncing ${bookmarks.length} bookmarks to ${config.apiUrl}`);
    
    // Send bookmarks to API
    const response = await fetch(`${config.apiUrl}/api/bookmarks/import`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${config.authToken}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ bookmarks })
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
      throw new Error(errorData.error || `HTTP ${response.status}: ${response.statusText}`);
    }
    
    const data = await response.json();
    
    console.log('Knowlex: Sync successful', data);
    
    // Notify popup of success
    try {
      chrome.runtime.sendMessage({
        type: 'SYNC_COMPLETE',
        count: bookmarks.length,
        imported: data.imported_count,
        duplicates: data.duplicate_count
      });
    } catch (e) {
      // Popup might not be open, that's ok
    }
    
    return {
      success: true,
      count: bookmarks.length,
      imported: data.imported_count,
      duplicates: data.duplicate_count
    };
    
  } catch (error) {
    console.error('Knowlex: Sync error:', error);
    
    // Notify popup of error
    try {
      chrome.runtime.sendMessage({
        type: 'SYNC_ERROR',
        error: error.message
      });
    } catch (e) {
      // Popup might not be open, that's ok
    }
    
    return {
      success: false,
      error: error.message
    };
  }
}

// Listen for extension installation
chrome.runtime.onInstalled.addListener((details) => {
  if (details.reason === 'install') {
    console.log('Knowlex: Extension installed, opening popup');
    // Open the popup or a welcome page
    chrome.action.openPopup();
  } else if (details.reason === 'update') {
    console.log('Knowlex: Extension updated');
  }
});

// Helper to test API connection
async function testConnection(apiUrl, authToken) {
  try {
    console.log('Knowlex: Testing connection to', apiUrl);
    
    const response = await fetch(`${apiUrl}/api/auth/me`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json'
      }
    });
    
    if (response.ok) {
      const user = await response.json();
      console.log('Knowlex: Connection successful', user);
      return { success: true, user };
    } else {
      const errorText = await response.text().catch(() => 'Invalid credentials');
      console.error('Knowlex: Connection failed', response.status, errorText);
      return { success: false, error: `Invalid credentials (${response.status})` };
    }
  } catch (error) {
    console.error('Knowlex: Connection error', error);
    return { success: false, error: `Connection failed: ${error.message}` };
  }
}


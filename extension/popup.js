// DOM Elements
const setupView = document.getElementById('setupView');
const syncedView = document.getElementById('syncedView');
const apiUrlInput = document.getElementById('apiUrl');
const authTokenInput = document.getElementById('authToken');
const saveBtn = document.getElementById('saveBtn');
const testBtn = document.getElementById('testBtn');
const syncBtn = document.getElementById('syncBtn');
const resetBtn = document.getElementById('resetBtn');
const status = document.getElementById('status');
const syncStatus = document.getElementById('syncStatus');
const progress = document.getElementById('progress');
const progressText = document.getElementById('progressText');
const progressFill = document.getElementById('progressFill');
const totalBookmarks = document.getElementById('totalBookmarks');
const lastSync = document.getElementById('lastSync');

// Initialize
document.addEventListener('DOMContentLoaded', async () => {
  const config = await chrome.storage.local.get(['apiUrl', 'authToken', 'totalBookmarks', 'lastSync']);
  
  if (config.apiUrl && config.authToken) {
    // Show synced view
    setupView.classList.remove('active');
    syncedView.classList.add('active');
    
    // Update stats
    totalBookmarks.textContent = config.totalBookmarks || 0;
    lastSync.textContent = config.lastSync || 'Never';
  } else {
    // Show setup view
    setupView.classList.add('active');
    syncedView.classList.remove('active');
  }
});

// Test connection
testBtn.addEventListener('click', async () => {
  const apiUrl = apiUrlInput.value.trim();
  const token = authTokenInput.value.trim();
  
  if (!apiUrl || !token) {
    showStatus(status, 'Please enter both API URL and auth token', 'error');
    return;
  }
  
  showStatus(status, 'Testing connection...', 'info');
  
  // Send message to background worker to test connection
  chrome.runtime.sendMessage({
    type: 'TEST_CONNECTION',
    apiUrl: apiUrl,
    token: token
  }, (response) => {
    if (response && response.success) {
      showStatus(status, `✅ Connected as @${response.user.username}`, 'success');
    } else {
      showStatus(status, `❌ ${response?.error || 'Connection failed'}`, 'error');
    }
  });
});

// Save and start sync
saveBtn.addEventListener('click', async () => {
  const apiUrl = apiUrlInput.value.trim();
  const token = authTokenInput.value.trim();
  
  if (!apiUrl || !token) {
    showStatus(status, 'Please enter both API URL and auth token', 'error');
    return;
  }
  
  // Save to storage
  await chrome.storage.local.set({ apiUrl, authToken: token });
  
  showStatus(status, '✅ Settings saved! Opening Twitter bookmarks...', 'success');
  
  // Wait a moment, then open Twitter bookmarks
  setTimeout(() => {
    chrome.tabs.create({ url: 'https://x.com/i/bookmarks' });
    window.close();
  }, 1500);
});

// Manual sync trigger
syncBtn.addEventListener('click', () => {
  chrome.tabs.create({ url: 'https://x.com/i/bookmarks' });
  window.close();
});

// Reset settings
resetBtn.addEventListener('click', async () => {
  if (confirm('Are you sure you want to reset all settings?')) {
    await chrome.storage.local.clear();
    setupView.classList.add('active');
    syncedView.classList.remove('active');
  }
});

// Listen for sync progress updates
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'SYNC_PROGRESS') {
    progress.classList.add('visible');
    progressText.textContent = `Syncing... ${message.current}/${message.total}`;
    const percentage = (message.current / message.total) * 100;
    progressFill.style.width = `${percentage}%`;
  } else if (message.type === 'SYNC_COMPLETE') {
    progress.classList.remove('visible');
    showStatus(syncStatus, `✅ Synced ${message.count} bookmarks successfully!`, 'success');
    totalBookmarks.textContent = message.count;
    lastSync.textContent = new Date().toLocaleTimeString();
    
    // Save stats
    chrome.storage.local.set({
      totalBookmarks: message.count,
      lastSync: new Date().toLocaleString()
    });
  } else if (message.type === 'SYNC_ERROR') {
    progress.classList.remove('visible');
    showStatus(syncStatus, `❌ Sync failed: ${message.error}`, 'error');
  }
});

// Helper function to show status messages
function showStatus(element, message, type) {
  element.textContent = message;
  element.className = `status ${type} visible`;
  
  if (type === 'success' || type === 'error') {
    setTimeout(() => {
      element.classList.remove('visible');
    }, 5000);
  }
}


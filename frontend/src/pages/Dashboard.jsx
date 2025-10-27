import { useState, useEffect } from 'react';
import { ExternalLink, RefreshCw } from 'lucide-react';
import DashboardLayout from '../components/layout/DashboardLayout';
import BookmarkGrid from '../components/bookmarks/BookmarkGrid';
import api from '../lib/api';
import toast from 'react-hot-toast';

const Dashboard = () => {
  const [bookmarks, setBookmarks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState(null);

  const handleRefreshSync = () => {
    const syncUrl = 'https://x.com/i/bookmarks?knowlex_sync=1';
    const newWindow = window.open(syncUrl, '_blank');
    if (newWindow) {
      toast.success('Opened Twitter bookmarks to start sync. Keep the tab open until it finishes.');
    } else {
      toast.error('Please allow popups to run the sync.');
    }
  };

  useEffect(() => {
    fetchBookmarks();
  }, [page, searchQuery, selectedCategory]);

  useEffect(() => {
    const handleSearch = (e) => {
      setSearchQuery(e.detail);
      setPage(1);
    };

    const handleCategorySelect = (e) => {
      setSelectedCategory(e.detail);
      setPage(1);
    };

    window.addEventListener('searchQuery', handleSearch);
    window.addEventListener('categorySelected', handleCategorySelect);

    return () => {
      window.removeEventListener('searchQuery', handleSearch);
      window.removeEventListener('categorySelected', handleCategorySelect);
    };
  }, []);

  const fetchBookmarks = async () => {
    setLoading(true);
    try {
      let url = `/bookmarks?page=${page}&page_size=20`;
      
      if (searchQuery) {
        url = `/bookmarks/search?q=${encodeURIComponent(searchQuery)}&page=${page}&page_size=20`;
      } else if (selectedCategory) {
        url = `/bookmarks?page=${page}&page_size=20&category_id=${selectedCategory.id}`;
      }

      const response = await api.get(url);
      setBookmarks(response.data.bookmarks || []);
      setTotalPages(response.data.total_pages || 1);
    } catch (error) {
      toast.error('Failed to load bookmarks');
      setBookmarks([]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <DashboardLayout>
      {loading ? (
        <div className="flex items-center justify-center py-20">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
        </div>
      ) : bookmarks.length === 0 ? (
        <div className="max-w-3xl mx-auto mt-20">
          <div className="bg-gradient-to-br from-white to-blue-50 dark:from-gray-800 dark:to-gray-900 rounded-3xl shadow-2xl p-12 text-center border border-gray-200 dark:border-gray-700">
            <div className="inline-flex items-center justify-center w-24 h-24 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full mb-6 shadow-lg">
              <span className="text-5xl">📚</span>
            </div>
            <h2 className="text-3xl font-bold mb-4 text-gray-900 dark:text-white">Start Your Collection</h2>
            <p className="text-lg text-gray-600 dark:text-gray-300 mb-8 max-w-md mx-auto">
              Install the Chrome extension to automatically sync and organize your X bookmarks
            </p>
            <div className="bg-white dark:bg-gray-800 rounded-2xl p-6 text-left mb-6 shadow-inner border border-gray-200 dark:border-gray-700">
              <p className="font-bold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                <span className="text-2xl">⚡</span>
                Quick Setup (2 minutes):
              </p>
              <ol className="space-y-3 text-gray-700 dark:text-gray-300">
                <li className="flex items-start gap-3">
                  <span className="flex items-center justify-center w-6 h-6 bg-blue-500 text-white rounded-full text-sm font-bold flex-shrink-0 mt-0.5">1</span>
                  <span>Load the Chrome extension (unpacked)</span>
                </li>
                <li className="flex items-start gap-3">
                  <span className="flex items-center justify-center w-6 h-6 bg-blue-500 text-white rounded-full text-sm font-bold flex-shrink-0 mt-0.5">2</span>
                  <span>Go to Settings and copy your auth token</span>
                </li>
                <li className="flex items-start gap-3">
                  <span className="flex items-center justify-center w-6 h-6 bg-blue-500 text-white rounded-full text-sm font-bold flex-shrink-0 mt-0.5">3</span>
                  <span>Paste token in extension and click sync</span>
                </li>
              </ol>
            </div>
            <a
              href="/settings"
              className="inline-flex items-center gap-2 px-8 py-4 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-2xl font-semibold hover:from-blue-700 hover:to-indigo-700 transition-all shadow-lg hover:shadow-xl hover:scale-105"
            >
              Get Your Token →
            </a>
          </div>
        </div>
      ) : (
        <>
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
                {searchQuery ? `Search: "${searchQuery}"` : 
                 selectedCategory ? selectedCategory.name : 
                 'All Bookmarks'}
              </h2>
              {bookmarks.length > 0 && (
                <p className="text-sm text-gray-500 mt-1">
                  Showing {bookmarks.length} bookmarks
                </p>
              )}
            </div>
            <div className="flex items-center gap-3">
              <button
                onClick={handleRefreshSync}
                className="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition-colors"
              >
                <RefreshCw className="w-4 h-4" />
                Refresh Bookmarks
              </button>
            </div>
          </div>

          <BookmarkGrid bookmarks={bookmarks} onRefresh={fetchBookmarks} />
          
          {totalPages > 1 && (
            <div className="flex items-center justify-center gap-2 mt-8">
              <button
                onClick={() => setPage(page - 1)}
                disabled={page === 1}
                className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg disabled:opacity-50"
              >
                Previous
              </button>
              <span className="text-gray-600 dark:text-gray-400">
                Page {page} of {totalPages}
              </span>
              <button
                onClick={() => setPage(page + 1)}
                disabled={page === totalPages}
                className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg disabled:opacity-50"
              >
                Next
              </button>
            </div>
          )}
        </>
      )}
    </DashboardLayout>
  );
};

export default Dashboard;

import { useState, useEffect } from 'react';
import { Download, Moon, Sun, Trash2, User, Copy, Eye, EyeOff } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import DashboardLayout from '../components/layout/DashboardLayout';
import useAuthStore from '../lib/authStore';
import api from '../lib/api';
import toast from 'react-hot-toast';

const Settings = () => {
  const [darkMode, setDarkMode] = useState(false);
  const [categories, setCategories] = useState([]);
  const [showToken, setShowToken] = useState(false);
  const { user, logout, token } = useAuthStore();
  const navigate = useNavigate();

  useEffect(() => {
    const isDark = document.documentElement.classList.contains('dark');
    setDarkMode(isDark);
    fetchCategories();
  }, []);

  const fetchCategories = async () => {
    try {
      const response = await api.get('/categories');
      setCategories(response.data || []);
    } catch (error) {
      toast.error('Failed to load categories');
    }
  };

  const toggleDarkMode = () => {
    const html = document.documentElement;
    if (darkMode) {
      html.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    } else {
      html.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    }
    setDarkMode(!darkMode);
  };

  const handleExportAll = async () => {
    try {
      const response = await api.get('/export/bookmarks');
      const blob = new Blob([JSON.stringify(response.data, null, 2)], {
        type: 'application/json',
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = 'bookmarks-export.json';
      a.click();
      toast.success('Bookmarks exported');
    } catch (error) {
      toast.error('Failed to export bookmarks');
    }
  };

  const handleExportCategory = async (categoryId, categoryName) => {
    try {
      const response = await api.get(`/export/category/${categoryId}`);
      const blob = new Blob([JSON.stringify(response.data, null, 2)], {
        type: 'application/json',
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${categoryName}-export.json`;
      a.click();
      toast.success('Category exported');
    } catch (error) {
      toast.error('Failed to export category');
    }
  };

  const handleDeleteAccount = async () => {
    if (
      window.confirm(
        'Are you sure you want to delete your account? This action cannot be undone.'
      )
    ) {
      try {
        await api.delete('/user/account');
        toast.success('Account deleted');
        logout();
        navigate('/');
      } catch (error) {
        toast.error('Failed to delete account');
      }
    }
  };

  const copyToken = () => {
    navigator.clipboard.writeText(token);
    toast.success('Token copied to clipboard!');
  };

  return (
    <DashboardLayout>
      <div className="max-w-3xl mx-auto">
        <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
          Settings
        </h2>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            Account Information
          </h3>
          <div className="flex items-center gap-4">
            <img
              src={user?.profile_image || 'https://via.placeholder.com/80'}
              alt={user?.username}
              className="w-20 h-20 rounded-full"
            />
            <div>
              <p className="text-lg font-semibold text-gray-900 dark:text-white">
                {user?.display_name}
              </p>
              <p className="text-gray-600 dark:text-gray-400">@{user?.username}</p>
              <p className="text-sm text-gray-500 mt-1">
                Member since {new Date(user?.created_at).toLocaleDateString()}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-gradient-to-br from-white to-blue-50 dark:from-gray-800 dark:to-gray-900 rounded-2xl shadow-lg p-8 mb-6 border border-gray-200 dark:border-gray-700">
          <h3 className="text-2xl font-bold text-gray-900 dark:text-white mb-3 flex items-center gap-2">
            <span className="text-2xl">🔐</span>
            Extension Token
          </h3>
          <p className="text-gray-600 dark:text-gray-400 mb-6 leading-relaxed">
            Copy this token to configure the Chrome extension for automatic X bookmark syncing.
          </p>
          <div className="flex items-center gap-3">
            <div className="flex-1 relative">
              <input
                type={showToken ? 'text' : 'password'}
                value={token || ''}
                readOnly
                className="w-full px-5 py-4 bg-white dark:bg-gray-700 border-2 border-gray-300 dark:border-gray-600 rounded-xl font-mono text-sm shadow-inner focus:border-blue-500 focus:outline-none transition-colors"
              />
            </div>
            <button
              onClick={() => setShowToken(!showToken)}
              className="p-4 border-2 border-gray-300 dark:border-gray-600 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700 hover:border-blue-500 transition-all"
              title={showToken ? 'Hide token' : 'Show token'}
            >
              {showToken ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
            </button>
            <button
              onClick={copyToken}
              className="flex items-center gap-2 px-6 py-4 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-xl font-semibold hover:from-blue-700 hover:to-indigo-700 transition-all shadow-lg hover:shadow-xl hover:scale-105"
              title="Copy token"
            >
              <Copy className="w-5 h-5" />
              Copy
            </button>
          </div>
          <div className="mt-5 p-4 bg-amber-50 dark:bg-amber-900/20 rounded-xl border border-amber-200 dark:border-amber-800">
            <p className="text-sm text-amber-900 dark:text-amber-200 leading-relaxed">
              <strong className="font-bold">⚠️ Security:</strong> Keep this token private. If compromised, logout and login again to generate a new one.
            </p>
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            Export Data
          </h3>
          <div className="space-y-3">
            <button
              onClick={handleExportAll}
              className="w-full flex items-center justify-between p-3 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700"
            >
              <span className="flex items-center gap-2 text-gray-900 dark:text-white">
                <Download className="w-5 h-5" />
                Export All Bookmarks
              </span>
            </button>

            {categories.length > 0 && (
              <div>
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">
                  Export by Category:
                </p>
                {categories.map((category) => (
                  <button
                    key={category.id}
                    onClick={() => handleExportCategory(category.id, category.name)}
                    className="w-full flex items-center justify-between p-3 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 mb-2"
                  >
                    <span className="flex items-center gap-2 text-gray-900 dark:text-white">
                      <div
                        className="w-3 h-3 rounded-full"
                        style={{ backgroundColor: category.color }}
                      />
                      {category.name}
                    </span>
                    <Download className="w-4 h-4 text-gray-500" />
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 mb-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            Appearance
          </h3>
          <button
            onClick={toggleDarkMode}
            className="w-full flex items-center justify-between p-3 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700"
          >
            <span className="flex items-center gap-2 text-gray-900 dark:text-white">
              {darkMode ? <Moon className="w-5 h-5" /> : <Sun className="w-5 h-5" />}
              {darkMode ? 'Dark Mode' : 'Light Mode'}
            </span>
            <div
              className={`w-12 h-6 rounded-full transition-colors ${
                darkMode ? 'bg-blue-500' : 'bg-gray-300'
              }`}
            >
              <div
                className={`w-5 h-5 bg-white rounded-full shadow-md transform transition-transform ${
                  darkMode ? 'translate-x-6' : 'translate-x-1'
                } mt-0.5`}
              />
            </div>
          </button>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold text-red-600 mb-4">Danger Zone</h3>
          <button
            onClick={handleDeleteAccount}
            className="w-full flex items-center justify-center gap-2 p-3 bg-red-500 text-white rounded-lg hover:bg-red-600"
          >
            <Trash2 className="w-5 h-5" />
            Delete Account
          </button>
          <p className="text-sm text-gray-500 mt-2 text-center">
            This will permanently delete your account and all your data.
          </p>
        </div>
      </div>
    </DashboardLayout>
  );
};

export default Settings;

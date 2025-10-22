import { Sparkles, Search, Folder, Download, ArrowRight } from 'lucide-react';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

const LandingPage = () => {
  const handleLogin = () => {
    window.location.href = `${API_URL}/auth/twitter`;
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-50 dark:from-gray-900 dark:via-slate-900 dark:to-gray-900">
      <div className="container mx-auto px-4 py-20">
        {/* Hero Section */}
        <div className="text-center mb-24 max-w-4xl mx-auto">
          <div className="inline-flex items-center gap-2 px-4 py-2 bg-blue-100 dark:bg-blue-900/30 rounded-full text-blue-600 dark:text-blue-400 text-sm font-medium mb-6 animate-pulse">
            <Sparkles className="w-4 h-4" />
            Your bookmarks, finally organized
          </div>
          <h1 className="text-6xl md:text-7xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-gray-900 via-blue-900 to-indigo-900 dark:from-white dark:via-blue-100 dark:to-indigo-100 mb-6 leading-tight">
            Master Your X Bookmarks
          </h1>
          <p className="text-xl md:text-2xl text-gray-600 dark:text-gray-300 mb-10 max-w-3xl mx-auto leading-relaxed">
            Never lose a great post again. Automatically sync, search, and organize your X (formerly Twitter) bookmarks with AI-powered categories.
          </p>
          <button
            onClick={handleLogin}
            className="group bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white font-bold py-5 px-10 rounded-2xl text-lg inline-flex items-center gap-3 transition-all transform hover:scale-105 hover:shadow-2xl shadow-lg"
          >
            <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
              <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/>
            </svg>
            Continue with X
            <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
          </button>
        </div>

        {/* Features Grid */}
        <div className="grid md:grid-cols-3 gap-8 max-w-6xl mx-auto mb-20">
          <div className="group bg-white dark:bg-gray-800 p-8 rounded-3xl shadow-lg hover:shadow-2xl transition-all duration-300 border border-gray-100 dark:border-gray-700 hover:border-blue-500 dark:hover:border-blue-500 hover:-translate-y-2">
            <div className="w-14 h-14 bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl flex items-center justify-center mb-5 group-hover:scale-110 transition-transform">
              <Search className="w-7 h-7 text-white" />
            </div>
            <h3 className="text-2xl font-bold text-gray-900 dark:text-white mb-3">
              Lightning Search
            </h3>
            <p className="text-gray-600 dark:text-gray-300 leading-relaxed">
              Find any post instantly with powerful full-text search across all your saved content
            </p>
          </div>

          <div className="group bg-white dark:bg-gray-800 p-8 rounded-3xl shadow-lg hover:shadow-2xl transition-all duration-300 border border-gray-100 dark:border-gray-700 hover:border-green-500 dark:hover:border-green-500 hover:-translate-y-2">
            <div className="w-14 h-14 bg-gradient-to-br from-green-500 to-emerald-600 rounded-2xl flex items-center justify-center mb-5 group-hover:scale-110 transition-transform">
              <Folder className="w-7 h-7 text-white" />
            </div>
            <h3 className="text-2xl font-bold text-gray-900 dark:text-white mb-3">
              Smart Categories
            </h3>
            <p className="text-gray-600 dark:text-gray-300 leading-relaxed">
              Organize with custom categories, colors, and icons. Your bookmarks, your way
            </p>
          </div>

          <div className="group bg-white dark:bg-gray-800 p-8 rounded-3xl shadow-lg hover:shadow-2xl transition-all duration-300 border border-gray-100 dark:border-gray-700 hover:border-purple-500 dark:hover:border-purple-500 hover:-translate-y-2">
            <div className="w-14 h-14 bg-gradient-to-br from-purple-500 to-indigo-600 rounded-2xl flex items-center justify-center mb-5 group-hover:scale-110 transition-transform">
              <Download className="w-7 h-7 text-white" />
            </div>
            <h3 className="text-2xl font-bold text-gray-900 dark:text-white mb-3">
              Export & Backup
            </h3>
            <p className="text-gray-600 dark:text-gray-300 leading-relaxed">
              Download your bookmarks as JSON anytime. Your data is always yours
            </p>
          </div>
        </div>

        {/* How It Works Section */}
        <div className="max-w-5xl mx-auto bg-gradient-to-br from-white to-blue-50 dark:from-gray-800 dark:to-gray-900 rounded-3xl p-12 shadow-xl border border-gray-100 dark:border-gray-700">
          <h2 className="text-4xl font-bold text-center mb-4 text-gray-900 dark:text-white">
            Get Started in 3 Steps
          </h2>
          <p className="text-center text-gray-600 dark:text-gray-400 mb-12 text-lg">
            Setup takes less than 2 minutes
          </p>
          <div className="grid md:grid-cols-3 gap-10">
            <div className="text-center relative">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl mb-6 text-white text-2xl font-bold shadow-lg">
                1
              </div>
              <h3 className="text-xl font-bold mb-3 text-gray-900 dark:text-white">Connect Your X Account</h3>
              <p className="text-gray-600 dark:text-gray-400 leading-relaxed">Secure OAuth login - no passwords stored</p>
            </div>
            <div className="text-center relative">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-green-500 to-emerald-600 rounded-2xl mb-6 text-white text-2xl font-bold shadow-lg">
                2
              </div>
              <h3 className="text-xl font-bold mb-3 text-gray-900 dark:text-white">Install Extension</h3>
              <p className="text-gray-600 dark:text-gray-400 leading-relaxed">One-click Chrome extension setup</p>
            </div>
            <div className="text-center relative">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-purple-500 to-indigo-600 rounded-2xl mb-6 text-white text-2xl font-bold shadow-lg">
                3
              </div>
              <h3 className="text-xl font-bold mb-3 text-gray-900 dark:text-white">Auto-Sync & Organize</h3>
              <p className="text-gray-600 dark:text-gray-400 leading-relaxed">Watch your bookmarks sync automatically</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LandingPage;

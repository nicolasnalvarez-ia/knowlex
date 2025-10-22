import { useState } from 'react';
import { MoreVertical, Trash2, Tag, ExternalLink } from 'lucide-react';
import { format } from 'date-fns';

const BookmarkCard = ({ bookmark, onDelete, onCategoryChange, onClick }) => {
  const [menuOpen, setMenuOpen] = useState(false);

  const truncateText = (text, maxLength) => {
    if (!text) return '';
    return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
  };

  const formatDate = (dateString) => {
    try {
      return format(new Date(dateString), 'MMM d, yyyy');
    } catch {
      return '';
    }
  };

  return (
    <div className="group bg-white dark:bg-gray-800 rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 p-6 cursor-pointer border border-gray-100 dark:border-gray-700 hover:border-blue-500 dark:hover:border-blue-500 hover:-translate-y-1"
         onClick={() => onClick(bookmark)}>
      <div className="flex justify-between items-start mb-3">
        <div className="flex items-center gap-2 flex-1 min-w-0">
          <span className="font-bold text-gray-900 dark:text-white truncate">
            {bookmark.author_display_name}
          </span>
          <span className="text-gray-500 dark:text-gray-400 text-sm flex-shrink-0">@{bookmark.author_username}</span>
        </div>
        <button
          onClick={(e) => {
            e.stopPropagation();
            setMenuOpen(!menuOpen);
          }}
          className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-xl transition-colors opacity-0 group-hover:opacity-100 relative flex-shrink-0"
        >
          <MoreVertical className="w-5 h-5 text-gray-600 dark:text-gray-400" />
          {menuOpen && (
            <div className="absolute right-0 mt-2 w-52 bg-white dark:bg-gray-700 rounded-2xl shadow-2xl border border-gray-200 dark:border-gray-600 z-10 overflow-hidden">
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onCategoryChange(bookmark);
                  setMenuOpen(false);
                }}
                className="w-full flex items-center gap-3 px-4 py-3 text-sm text-gray-700 dark:text-gray-200 hover:bg-blue-50 dark:hover:bg-gray-600 transition-colors"
              >
                <Tag className="w-4 h-4" />
                Manage Categories
              </button>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onDelete(bookmark.id);
                  setMenuOpen(false);
                }}
                className="w-full flex items-center gap-3 px-4 py-3 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-gray-600 transition-colors"
              >
                <Trash2 className="w-4 h-4" />
                Delete Bookmark
              </button>
            </div>
          )}
        </button>
      </div>

      <p className="text-gray-700 dark:text-gray-300 mb-4 leading-relaxed text-[15px]">
        {truncateText(bookmark.tweet_text, 200)}
      </p>

      {bookmark.media_urls && bookmark.media_urls.length > 0 && (
        <div className="mb-4 -mx-2">
          <img
            src={bookmark.media_urls[0]}
            alt="Post media"
            className="rounded-xl w-full h-56 object-cover"
          />
        </div>
      )}

      <div className="flex items-center justify-between text-sm">
        <span className="text-gray-500 dark:text-gray-400 font-medium">{formatDate(bookmark.bookmarked_at)}</span>
        {bookmark.categories && bookmark.categories.length > 0 && (
          <div className="flex gap-2 flex-wrap justify-end">
            {bookmark.categories.slice(0, 2).map((cat) => (
              <span
                key={cat.id}
                className="px-3 py-1 rounded-full text-xs font-semibold shadow-sm"
                style={{ backgroundColor: cat.color + '25', color: cat.color }}
              >
                {cat.icon} {cat.name}
              </span>
            ))}
            {bookmark.categories.length > 2 && (
              <span className="px-3 py-1 rounded-full text-xs font-semibold bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-300">
                +{bookmark.categories.length - 2}
              </span>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default BookmarkCard;

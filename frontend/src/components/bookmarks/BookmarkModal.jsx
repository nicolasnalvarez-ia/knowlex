import { X, ExternalLink } from 'lucide-react';
import { format } from 'date-fns';

const BookmarkModal = ({ bookmark, onClose }) => {
  if (!bookmark) return null;

  const formatDate = (dateString) => {
    try {
      return format(new Date(dateString), 'MMMM d, yyyy • h:mm a');
    } catch {
      return '';
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
         onClick={onClose}>
      <div className="bg-white dark:bg-gray-800 rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto"
           onClick={(e) => e.stopPropagation()}>
        <div className="sticky top-0 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 p-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
            Bookmark Details
          </h2>
          <button
            onClick={onClose}
            className="p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="p-6">
          <div className="flex items-start gap-3 mb-4">
            <div>
              <div className="flex items-center gap-2 mb-1">
                <span className="font-bold text-gray-900 dark:text-white">
                  {bookmark.author_display_name}
                </span>
                <span className="text-gray-500">@{bookmark.author_username}</span>
              </div>
              <p className="text-gray-700 dark:text-gray-300 text-lg leading-relaxed mb-4">
                {bookmark.tweet_text}
              </p>
              <div className="text-sm text-gray-500 mb-4">
                {formatDate(bookmark.bookmarked_at)}
              </div>
            </div>
          </div>

          {bookmark.media_urls && bookmark.media_urls.length > 0 && (
            <div className="grid grid-cols-1 gap-2 mb-4">
              {bookmark.media_urls.map((url, index) => (
                <img
                  key={index}
                  src={url}
                  alt={`Media ${index + 1}`}
                  className="rounded-lg w-full"
                />
              ))}
            </div>
          )}

          {bookmark.categories && bookmark.categories.length > 0 && (
            <div className="mb-4">
              <p className="text-sm text-gray-500 mb-2">Categories:</p>
              <div className="flex gap-2 flex-wrap">
                {bookmark.categories.map((cat) => (
                  <span
                    key={cat.id}
                    className="px-3 py-1 rounded-full text-sm"
                    style={{ backgroundColor: cat.color + '20', color: cat.color }}
                  >
                    {cat.name}
                  </span>
                ))}
              </div>
            </div>
          )}

          <a
            href={bookmark.tweet_url}
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-2 text-blue-500 hover:text-blue-600"
          >
            <ExternalLink className="w-4 h-4" />
            Open on Twitter
          </a>
        </div>
      </div>
    </div>
  );
};

export default BookmarkModal;

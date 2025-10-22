import { useState, useEffect } from 'react';
import BookmarkCard from './BookmarkCard';
import BookmarkModal from './BookmarkModal';
import CategoryDropdown from '../categories/CategoryDropdown';
import api from '../../lib/api';
import toast from 'react-hot-toast';

const BookmarkGrid = ({ bookmarks, onRefresh }) => {
  const [selectedBookmark, setSelectedBookmark] = useState(null);
  const [categoryDropdown, setCategoryDropdown] = useState(null);

  const handleDelete = async (bookmarkId) => {
    if (window.confirm('Delete this bookmark?')) {
      try {
        await api.delete(`/bookmarks/${bookmarkId}`);
        toast.success('Bookmark deleted');
        onRefresh();
      } catch (error) {
        toast.error('Failed to delete bookmark');
      }
    }
  };

  const handleCategoryChange = (bookmark) => {
    setCategoryDropdown(bookmark);
  };

  const handleCategorySelected = async (categoryId) => {
    try {
      await api.post(`/bookmarks/${categoryDropdown.id}/category`, {
        category_id: categoryId,
      });
      toast.success('Category assigned');
      setCategoryDropdown(null);
      onRefresh();
    } catch (error) {
      toast.error('Failed to assign category');
    }
  };

  return (
    <>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {bookmarks.map((bookmark) => (
          <BookmarkCard
            key={bookmark.id}
            bookmark={bookmark}
            onDelete={handleDelete}
            onCategoryChange={handleCategoryChange}
            onClick={setSelectedBookmark}
          />
        ))}
      </div>

      {selectedBookmark && (
        <BookmarkModal
          bookmark={selectedBookmark}
          onClose={() => setSelectedBookmark(null)}
        />
      )}

      {categoryDropdown && (
        <CategoryDropdown
          onSelect={handleCategorySelected}
          onClose={() => setCategoryDropdown(null)}
        />
      )}
    </>
  );
};

export default BookmarkGrid;

import { useState, useEffect } from 'react';
import { Plus, Folder, MoreVertical, Edit2, Trash2 } from 'lucide-react';
import api from '../../lib/api';
import CategoryModal from './CategoryModal';
import toast from 'react-hot-toast';

const CategorySidebar = () => {
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState(null);
  const [menuOpen, setMenuOpen] = useState(null);

  useEffect(() => {
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

  const handleCategoryClick = (category) => {
    setSelectedCategory(category);
    window.dispatchEvent(new CustomEvent('categorySelected', { detail: category }));
  };

  const handleEdit = (category) => {
    setEditingCategory(category);
    setIsModalOpen(true);
    setMenuOpen(null);
  };

  const handleDelete = async (category) => {
    if (window.confirm(`Delete category "${category.name}"?`)) {
      try {
        await api.delete(`/categories/${category.id}`);
        toast.success('Category deleted');
        fetchCategories();
        setMenuOpen(null);
      } catch (error) {
        toast.error('Failed to delete category');
      }
    }
  };

  const handleModalClose = () => {
    setIsModalOpen(false);
    setEditingCategory(null);
    fetchCategories();
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Categories</h2>
        <button
          onClick={() => setIsModalOpen(true)}
          className="p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded"
        >
          <Plus className="w-5 h-5 text-blue-500" />
        </button>
      </div>

      <button
        onClick={() => {
          setSelectedCategory(null);
          window.dispatchEvent(new CustomEvent('categorySelected', { detail: null }));
        }}
        className={`w-full flex items-center justify-between p-3 rounded-lg mb-2 ${
          selectedCategory === null
            ? 'bg-blue-50 dark:bg-blue-900 text-blue-600 dark:text-blue-300'
            : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300'
        }`}
      >
        <span className="flex items-center gap-2">
          <Folder className="w-4 h-4" />
          All Bookmarks
        </span>
      </button>

      <div className="space-y-1">
        {categories.map((category) => (
          <div key={category.id} className="relative">
            <button
              onClick={() => handleCategoryClick(category)}
              className={`w-full flex items-center justify-between p-3 rounded-lg ${
                selectedCategory?.id === category.id
                  ? 'bg-blue-50 dark:bg-blue-900 text-blue-600 dark:text-blue-300'
                  : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300'
              }`}
            >
              <span className="flex items-center gap-2">
                <div
                  className="w-3 h-3 rounded-full"
                  style={{ backgroundColor: category.color }}
                />
                <span className="truncate">{category.name}</span>
              </span>
              <div className="flex items-center gap-2">
                <span className="text-xs text-gray-500">{category.count || 0}</span>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setMenuOpen(menuOpen === category.id ? null : category.id);
                  }}
                  className="p-1 hover:bg-gray-200 dark:hover:bg-gray-600 rounded"
                >
                  <MoreVertical className="w-4 h-4" />
                </button>
              </div>
            </button>

            {menuOpen === category.id && (
              <div className="absolute right-0 top-12 w-32 bg-white dark:bg-gray-700 rounded-lg shadow-lg border border-gray-200 dark:border-gray-600 z-10">
                <button
                  onClick={() => handleEdit(category)}
                  className="w-full flex items-center gap-2 px-3 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-600"
                >
                  <Edit2 className="w-3 h-3" />
                  Edit
                </button>
                <button
                  onClick={() => handleDelete(category)}
                  className="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-gray-100 dark:hover:bg-gray-600"
                >
                  <Trash2 className="w-3 h-3" />
                  Delete
                </button>
              </div>
            )}
          </div>
        ))}
      </div>

      {isModalOpen && (
        <CategoryModal
          category={editingCategory}
          onClose={handleModalClose}
        />
      )}
    </div>
  );
};

export default CategorySidebar;

package database

import (
	"context"
	"fmt"
	"twitter-bookmarks-api/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func CreateUser(ctx context.Context, twitterID, username, displayName, profileImage string) (*models.User, error) {
	user := &models.User{}
	query := `
		INSERT INTO users (twitter_id, username, display_name, profile_image)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (twitter_id) 
		DO UPDATE SET username = $2, display_name = $3, profile_image = $4
		RETURNING id, twitter_id, username, display_name, profile_image, created_at
	`
	err := DB.QueryRow(ctx, query, twitterID, username, displayName, profileImage).Scan(
		&user.ID, &user.TwitterID, &user.Username, &user.DisplayName, &user.ProfileImage, &user.CreatedAt,
	)
	return user, err
}

func GetUserByTwitterID(ctx context.Context, twitterID string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, twitter_id, username, display_name, profile_image, created_at FROM users WHERE twitter_id = $1`
	err := DB.QueryRow(ctx, query, twitterID).Scan(
		&user.ID, &user.TwitterID, &user.Username, &user.DisplayName, &user.ProfileImage, &user.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, twitter_id, username, display_name, profile_image, created_at FROM users WHERE id = $1`
	err := DB.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.TwitterID, &user.Username, &user.DisplayName, &user.ProfileImage, &user.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func CreateBookmark(ctx context.Context, bookmark *models.Bookmark) error {
	query := `
		INSERT INTO bookmarks (user_id, tweet_id, tweet_text, author_username, author_display_name, tweet_url, media_urls, bookmarked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, tweet_id) DO NOTHING
		RETURNING id, created_at
	`
	return DB.QueryRow(ctx, query,
		bookmark.UserID, bookmark.TweetID, bookmark.TweetText, bookmark.AuthorUsername,
		bookmark.AuthorDisplayName, bookmark.TweetURL, bookmark.MediaURLs, bookmark.BookmarkedAt,
	).Scan(&bookmark.ID, &bookmark.CreatedAt)
}

func GetBookmarksByUserID(ctx context.Context, userID uuid.UUID, params models.PaginationParams, categoryID *uuid.UUID) (*models.BookmarksResponse, error) {
	var bookmarks []models.Bookmark
	var total int

	countQuery := `SELECT COUNT(*) FROM bookmarks WHERE user_id = $1`
	baseQuery := `
		SELECT DISTINCT b.id, b.user_id, b.tweet_id, b.tweet_text, b.author_username, 
		       b.author_display_name, b.tweet_url, b.media_urls, b.bookmarked_at, b.created_at
		FROM bookmarks b
	`

	countArgs := []interface{}{userID}
	queryArgs := []interface{}{userID}
	whereClauses := "WHERE b.user_id = $1"

	if categoryID != nil {
		baseQuery += " LEFT JOIN bookmark_categories bc ON b.id = bc.bookmark_id"
		whereClauses += " AND bc.category_id = $2"
		countQuery += " AND id IN (SELECT bookmark_id FROM bookmark_categories WHERE category_id = $2)"
		countArgs = append(countArgs, *categoryID)
		queryArgs = append(queryArgs, *categoryID)
	}

	err := DB.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, err
	}

	fullQuery := fmt.Sprintf("%s %s ORDER BY b.bookmarked_at DESC LIMIT $%d OFFSET $%d",
		baseQuery, whereClauses, len(queryArgs)+1, len(queryArgs)+2)
	queryArgs = append(queryArgs, params.PageSize, params.Offset)

	rows, err := DB.Query(ctx, fullQuery, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b models.Bookmark
		err := rows.Scan(&b.ID, &b.UserID, &b.TweetID, &b.TweetText, &b.AuthorUsername,
			&b.AuthorDisplayName, &b.TweetURL, &b.MediaURLs, &b.BookmarkedAt, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories, _ := GetCategoriesByBookmarkID(ctx, b.ID)
		b.Categories = categories
		bookmarks = append(bookmarks, b)
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize
	return &models.BookmarksResponse{
		Bookmarks:  bookmarks,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func SearchBookmarks(ctx context.Context, userID uuid.UUID, query string, params models.PaginationParams) (*models.BookmarksResponse, error) {
	var bookmarks []models.Bookmark
	var total int

	searchPattern := "%" + query + "%"
	countQuery := `
		SELECT COUNT(*) FROM bookmarks 
		WHERE user_id = $1 AND (tweet_text ILIKE $2 OR author_username ILIKE $2 OR author_display_name ILIKE $2)
	`
	err := DB.QueryRow(ctx, countQuery, userID, searchPattern).Scan(&total)
	if err != nil {
		return nil, err
	}

	searchQuery := `
		SELECT id, user_id, tweet_id, tweet_text, author_username, author_display_name, 
		       tweet_url, media_urls, bookmarked_at, created_at
		FROM bookmarks
		WHERE user_id = $1 AND (tweet_text ILIKE $2 OR author_username ILIKE $2 OR author_display_name ILIKE $2)
		ORDER BY bookmarked_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := DB.Query(ctx, searchQuery, userID, searchPattern, params.PageSize, params.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b models.Bookmark
		err := rows.Scan(&b.ID, &b.UserID, &b.TweetID, &b.TweetText, &b.AuthorUsername,
			&b.AuthorDisplayName, &b.TweetURL, &b.MediaURLs, &b.BookmarkedAt, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories, _ := GetCategoriesByBookmarkID(ctx, b.ID)
		b.Categories = categories
		bookmarks = append(bookmarks, b)
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize
	return &models.BookmarksResponse{
		Bookmarks:  bookmarks,
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
	}, nil
}

func DeleteBookmark(ctx context.Context, bookmarkID, userID uuid.UUID) error {
	query := `DELETE FROM bookmarks WHERE id = $1 AND user_id = $2`
	result, err := DB.Exec(ctx, query, bookmarkID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("bookmark not found")
	}
	return nil
}

func CreateCategory(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO categories (user_id, name, color, icon)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return DB.QueryRow(ctx, query, category.UserID, category.Name, category.Color, category.Icon).
		Scan(&category.ID, &category.CreatedAt)
}

func GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	query := `
		SELECT c.id, c.user_id, c.name, c.color, c.icon, c.created_at,
		       COALESCE(COUNT(bc.bookmark_id), 0) as count
		FROM categories c
		LEFT JOIN bookmark_categories bc ON c.id = bc.category_id
		WHERE c.user_id = $1
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`
	rows, err := DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Color, &c.Icon, &c.CreatedAt, &c.Count)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func GetCategoriesByBookmarkID(ctx context.Context, bookmarkID uuid.UUID) ([]models.Category, error) {
	query := `
		SELECT c.id, c.user_id, c.name, c.color, c.icon, c.created_at
		FROM categories c
		INNER JOIN bookmark_categories bc ON c.id = bc.category_id
		WHERE bc.bookmark_id = $1
	`
	rows, err := DB.Query(ctx, query, bookmarkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Color, &c.Icon, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func GetCategoryByID(ctx context.Context, categoryID, userID uuid.UUID) (*models.Category, error) {
	category := &models.Category{}
	query := `SELECT id, user_id, name, color, icon, created_at FROM categories WHERE id = $1 AND user_id = $2`
	err := DB.QueryRow(ctx, query, categoryID, userID).Scan(
		&category.ID, &category.UserID, &category.Name, &category.Color, &category.Icon, &category.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return category, err
}

func UpdateCategory(ctx context.Context, categoryID, userID uuid.UUID, name, color, icon string) error {
	query := `
		UPDATE categories 
		SET name = COALESCE(NULLIF($1, ''), name),
		    color = COALESCE(NULLIF($2, ''), color),
		    icon = COALESCE(NULLIF($3, ''), icon)
		WHERE id = $4 AND user_id = $5
	`
	result, err := DB.Exec(ctx, query, name, color, icon, categoryID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

func DeleteCategory(ctx context.Context, categoryID, userID uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`
	result, err := DB.Exec(ctx, query, categoryID, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

func AssignBookmarkToCategory(ctx context.Context, bookmarkID, categoryID, userID uuid.UUID) error {
	var bookmarkUserID uuid.UUID
	err := DB.QueryRow(ctx, `SELECT user_id FROM bookmarks WHERE id = $1`, bookmarkID).Scan(&bookmarkUserID)
	if err != nil {
		return fmt.Errorf("bookmark not found")
	}
	if bookmarkUserID != userID {
		return fmt.Errorf("unauthorized")
	}

	var categoryUserID uuid.UUID
	err = DB.QueryRow(ctx, `SELECT user_id FROM categories WHERE id = $1`, categoryID).Scan(&categoryUserID)
	if err != nil {
		return fmt.Errorf("category not found")
	}
	if categoryUserID != userID {
		return fmt.Errorf("unauthorized")
	}

	query := `
		INSERT INTO bookmark_categories (bookmark_id, category_id)
		VALUES ($1, $2)
		ON CONFLICT (bookmark_id, category_id) DO NOTHING
	`
	_, err = DB.Exec(ctx, query, bookmarkID, categoryID)
	return err
}

func RemoveBookmarkFromCategory(ctx context.Context, bookmarkID, categoryID, userID uuid.UUID) error {
	query := `
		DELETE FROM bookmark_categories
		WHERE bookmark_id = $1 AND category_id = $2
		AND bookmark_id IN (SELECT id FROM bookmarks WHERE user_id = $3)
	`
	_, err := DB.Exec(ctx, query, bookmarkID, categoryID, userID)
	return err
}

func GetUncategorizedBookmarks(ctx context.Context, userID uuid.UUID, limit int) ([]models.Bookmark, error) {
	query := `
		SELECT b.id, b.user_id, b.tweet_id, b.tweet_text, b.author_username, 
		       b.author_display_name, b.tweet_url, b.media_urls, b.bookmarked_at, b.created_at
		FROM bookmarks b
		LEFT JOIN bookmark_categories bc ON b.id = bc.bookmark_id
		WHERE b.user_id = $1 AND bc.id IS NULL
		ORDER BY b.created_at DESC
		LIMIT $2
	`
	
	rows, err := DB.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []models.Bookmark
	for rows.Next() {
		var b models.Bookmark
		err := rows.Scan(&b.ID, &b.UserID, &b.TweetID, &b.TweetText, &b.AuthorUsername,
			&b.AuthorDisplayName, &b.TweetURL, &b.MediaURLs, &b.BookmarkedAt, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}

func GetBookmarkByID(ctx context.Context, bookmarkID, userID uuid.UUID) (*models.Bookmark, error) {
	bookmark := &models.Bookmark{}
	query := `
		SELECT id, user_id, tweet_id, tweet_text, author_username, author_display_name, 
		       tweet_url, media_urls, bookmarked_at, created_at
		FROM bookmarks
		WHERE id = $1 AND user_id = $2
	`
	err := DB.QueryRow(ctx, query, bookmarkID, userID).Scan(
		&bookmark.ID, &bookmark.UserID, &bookmark.TweetID, &bookmark.TweetText,
		&bookmark.AuthorUsername, &bookmark.AuthorDisplayName, &bookmark.TweetURL,
		&bookmark.MediaURLs, &bookmark.BookmarkedAt, &bookmark.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("bookmark not found")
	}
	if err != nil {
		return nil, err
	}
	
	categories, _ := GetCategoriesByBookmarkID(ctx, bookmark.ID)
	bookmark.Categories = categories
	return bookmark, nil
}

func DeleteUserAndAllData(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := DB.Exec(ctx, query, userID)
	return err
}

func GetAllBookmarksByUserID(ctx context.Context, userID uuid.UUID) ([]models.Bookmark, error) {
	query := `
		SELECT id, user_id, tweet_id, tweet_text, author_username, author_display_name,
		       tweet_url, media_urls, bookmarked_at, created_at
		FROM bookmarks
		WHERE user_id = $1
		ORDER BY bookmarked_at DESC
	`
	rows, err := DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []models.Bookmark
	for rows.Next() {
		var b models.Bookmark
		err := rows.Scan(&b.ID, &b.UserID, &b.TweetID, &b.TweetText, &b.AuthorUsername,
			&b.AuthorDisplayName, &b.TweetURL, &b.MediaURLs, &b.BookmarkedAt, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}

func GetBookmarksByCategoryID(ctx context.Context, categoryID, userID uuid.UUID) ([]models.Bookmark, error) {
	query := `
		SELECT b.id, b.user_id, b.tweet_id, b.tweet_text, b.author_username, b.author_display_name,
		       b.tweet_url, b.media_urls, b.bookmarked_at, b.created_at
		FROM bookmarks b
		INNER JOIN bookmark_categories bc ON b.id = bc.bookmark_id
		WHERE bc.category_id = $1 AND b.user_id = $2
		ORDER BY b.bookmarked_at DESC
	`
	rows, err := DB.Query(ctx, query, categoryID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []models.Bookmark
	for rows.Next() {
		var b models.Bookmark
		err := rows.Scan(&b.ID, &b.UserID, &b.TweetID, &b.TweetText, &b.AuthorUsername,
			&b.AuthorDisplayName, &b.TweetURL, &b.MediaURLs, &b.BookmarkedAt, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, nil
}

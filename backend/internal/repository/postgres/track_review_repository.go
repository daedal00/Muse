package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/daedal00/muse/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type trackReviewRepository struct {
	db *database.PostgresDB
}

func NewTrackReviewRepository(db *database.PostgresDB) repository.TrackReviewRepository {
	return &trackReviewRepository{db: db}
}

func (r *trackReviewRepository) Create(ctx context.Context, review *models.TrackReview) error {
	query := `
		INSERT INTO track_reviews (id, user_id, track_id, rating, review_text, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		review.ID, review.UserID, review.TrackID, review.Rating,
		review.ReviewText, review.CreatedAt, review.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create track review: %w", err)
	}

	return nil
}

func (r *trackReviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TrackReview, error) {
	query := `
		SELECT id, user_id, track_id, rating, review_text, created_at, updated_at
		FROM track_reviews
		WHERE id = $1
	`

	review := &models.TrackReview{}
	if err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&review.ID, &review.UserID, &review.TrackID, &review.Rating,
		&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("track review not found")
		}
		return nil, fmt.Errorf("failed to get track review: %w", err)
	}

	return review, nil
}

func (r *trackReviewRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.TrackReview, error) {
	query := `
		SELECT id, user_id, track_id, rating, review_text, created_at, updated_at
		FROM track_reviews
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list track reviews by user: %w", err)
	}
	defer rows.Close()

	var reviews []*models.TrackReview
	for rows.Next() {
		review := &models.TrackReview{}
		if err := rows.Scan(
			&review.ID, &review.UserID, &review.TrackID, &review.Rating,
			&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan track review: %w", err)
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating track reviews: %w", err)
	}

	return reviews, nil
}

func (r *trackReviewRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM track_reviews WHERE user_id = $1`

	var count int
	if err := r.db.Pool.QueryRow(ctx, query, userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count track reviews by user: %w", err)
	}

	return count, nil
}

func (r *trackReviewRepository) GetByTrackID(ctx context.Context, trackID uuid.UUID, limit, offset int) ([]*models.TrackReview, error) {
	query := `
		SELECT id, user_id, track_id, rating, review_text, created_at, updated_at
		FROM track_reviews
		WHERE track_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, trackID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list track reviews by track: %w", err)
	}
	defer rows.Close()

	var reviews []*models.TrackReview
	for rows.Next() {
		review := &models.TrackReview{}
		if err := rows.Scan(
			&review.ID, &review.UserID, &review.TrackID, &review.Rating,
			&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan track review: %w", err)
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating track reviews: %w", err)
	}

	return reviews, nil
}

func (r *trackReviewRepository) CountByTrackID(ctx context.Context, trackID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM track_reviews WHERE track_id = $1`

	var count int
	if err := r.db.Pool.QueryRow(ctx, query, trackID).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count track reviews by track: %w", err)
	}

	return count, nil
}

func (r *trackReviewRepository) AverageRatingByTrackID(ctx context.Context, trackID uuid.UUID) (*float64, error) {
	query := `SELECT AVG(rating) FROM track_reviews WHERE track_id = $1`

	var avg sql.NullFloat64
	if err := r.db.Pool.QueryRow(ctx, query, trackID).Scan(&avg); err != nil {
		return nil, fmt.Errorf("failed to average rating by track: %w", err)
	}

	if !avg.Valid {
		return nil, nil
	}

	return &avg.Float64, nil
}

func (r *trackReviewRepository) GetByUserAndTrack(ctx context.Context, userID, trackID uuid.UUID) (*models.TrackReview, error) {
	query := `
		SELECT id, user_id, track_id, rating, review_text, created_at, updated_at
		FROM track_reviews
		WHERE user_id = $1 AND track_id = $2
	`

	review := &models.TrackReview{}
	if err := r.db.Pool.QueryRow(ctx, query, userID, trackID).Scan(
		&review.ID, &review.UserID, &review.TrackID, &review.Rating,
		&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("track review not found")
		}
		return nil, fmt.Errorf("failed to get track review: %w", err)
	}

	return review, nil
}

func (r *trackReviewRepository) ListTopTracksByAverageRating(ctx context.Context, limit int) ([]*models.TrackReviewSummary, error) {
	query := `
		SELECT track_id, AVG(rating) AS average_rating, COUNT(*) AS review_count
		FROM track_reviews
		GROUP BY track_id
		ORDER BY average_rating DESC, review_count DESC
		LIMIT $1
	`

	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list top tracks by average rating: %w", err)
	}
	defer rows.Close()

	var summaries []*models.TrackReviewSummary
	for rows.Next() {
		summary := &models.TrackReviewSummary{}
		if err := rows.Scan(&summary.TrackID, &summary.AverageRating, &summary.ReviewCount); err != nil {
			return nil, fmt.Errorf("failed to scan track summary: %w", err)
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating track summaries: %w", err)
	}

	return summaries, nil
}

func (r *trackReviewRepository) ListTopTracksByReviewCount(ctx context.Context, limit int) ([]*models.TrackReviewSummary, error) {
	query := `
		SELECT track_id, AVG(rating) AS average_rating, COUNT(*) AS review_count
		FROM track_reviews
		GROUP BY track_id
		ORDER BY review_count DESC, average_rating DESC
		LIMIT $1
	`

	rows, err := r.db.Pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list top tracks by review count: %w", err)
	}
	defer rows.Close()

	var summaries []*models.TrackReviewSummary
	for rows.Next() {
		summary := &models.TrackReviewSummary{}
		if err := rows.Scan(&summary.TrackID, &summary.AverageRating, &summary.ReviewCount); err != nil {
			return nil, fmt.Errorf("failed to scan track summary: %w", err)
		}
		summaries = append(summaries, summary)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating track summaries: %w", err)
	}

	return summaries, nil
}

func (r *trackReviewRepository) ListFavoritesByUserID(ctx context.Context, userID uuid.UUID, minRating, limit int) ([]*models.TrackReview, error) {
	query := `
		SELECT id, user_id, track_id, rating, review_text, created_at, updated_at
		FROM track_reviews
		WHERE user_id = $1 AND rating >= $2
		ORDER BY rating DESC, created_at DESC
		LIMIT $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, minRating, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list favorite tracks: %w", err)
	}
	defer rows.Close()

	var reviews []*models.TrackReview
	for rows.Next() {
		review := &models.TrackReview{}
		if err := rows.Scan(
			&review.ID, &review.UserID, &review.TrackID, &review.Rating,
			&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan favorite track review: %w", err)
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating favorite track reviews: %w", err)
	}

	return reviews, nil
}

func (r *trackReviewRepository) List(ctx context.Context, limit, offset int) ([]*models.TrackReview, error) {
	query := `
		SELECT id, user_id, track_id, rating, review_text, created_at, updated_at
		FROM track_reviews
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list track reviews: %w", err)
	}
	defer rows.Close()

	var reviews []*models.TrackReview
	for rows.Next() {
		review := &models.TrackReview{}
		if err := rows.Scan(
			&review.ID, &review.UserID, &review.TrackID, &review.Rating,
			&review.ReviewText, &review.CreatedAt, &review.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan track review: %w", err)
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating track reviews: %w", err)
	}

	return reviews, nil
}

func (r *trackReviewRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM track_reviews`

	var count int
	if err := r.db.Pool.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count track reviews: %w", err)
	}

	return count, nil
}

func (r *trackReviewRepository) Update(ctx context.Context, review *models.TrackReview) error {
	query := `
		UPDATE track_reviews
		SET rating = $2, review_text = $3, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, review.ID, review.Rating, review.ReviewText)
	if err != nil {
		return fmt.Errorf("failed to update track review: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("track review not found")
	}

	return nil
}

func (r *trackReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM track_reviews WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete track review: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("track review not found")
	}

	return nil
}

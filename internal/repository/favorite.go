package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	// Р—Р°РјРµРЅРё РїСѓС‚СЊ РЅР° РёРјСЏ РјРѕРґСѓР»СЏ РёР· go.mod.
	"github.com/Zuf4ik/Z_JobScout/internal/vacancy"
)

type FavoriteRepository struct {
	db *pgxpool.Pool
}

func NewFavoriteRepository(db *pgxpool.Pool) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Add(
	ctx context.Context,
	favorite vacancy.Favorite,
) error {
	query := `
		INSERT INTO favorite_vacancies (
			user_id,
			vacancy_id,
			source,
			title,
			company_name,
			url
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, source, vacancy_id)
		DO UPDATE SET
			title = EXCLUDED.title,
			company_name = EXCLUDED.company_name,
			url = EXCLUDED.url
	`

	_, err := r.db.Exec(
		ctx,
		query,
		favorite.UserID,
		favorite.VacancyID,
		favorite.Source,
		favorite.Title,
		favorite.CompanyName,
		favorite.URL,
	)
	if err != nil {
		return fmt.Errorf("РЅРµ СѓРґР°Р»РѕСЃСЊ СЃРѕС…СЂР°РЅРёС‚СЊ РІР°РєР°РЅСЃРёСЋ: %w", err)
	}

	return nil
}

func (r *FavoriteRepository) GetAll(
	ctx context.Context,
	userID int64,
) ([]vacancy.Favorite, error) {
	query := `
		SELECT
			id,
			user_id,
			vacancy_id,
			source,
			title,
			company_name,
			url,
			saved_at
		FROM favorite_vacancies
		WHERE user_id = $1
		ORDER BY saved_at DESC, id DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("РѕС€РёР±РєР° РїРѕР»СѓС‡РµРЅРёСЏ РёР·Р±СЂР°РЅРЅРѕРіРѕ: %w", err)
	}
	defer rows.Close()

	favorites := make([]vacancy.Favorite, 0)

	for rows.Next() {
		var favorite vacancy.Favorite

		err := rows.Scan(
			&favorite.ID,
			&favorite.UserID,
			&favorite.VacancyID,
			&favorite.Source,
			&favorite.Title,
			&favorite.CompanyName,
			&favorite.URL,
			&favorite.SavedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("РѕС€РёР±РєР° С‡С‚РµРЅРёСЏ РІР°РєР°РЅСЃРёРё: %w", err)
		}

		favorites = append(favorites, favorite)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("РѕС€РёР±РєР° С‡С‚РµРЅРёСЏ СЂРµР·СѓР»СЊС‚Р°С‚РѕРІ: %w", err)
	}

	return favorites, nil
}

func (r *FavoriteRepository) Delete(
	ctx context.Context,
	userID int64,
	source string,
	vacancyID string,
) error {
	query := `
		DELETE FROM favorite_vacancies
		WHERE user_id = $1
		  AND source = $2
		  AND vacancy_id = $3
	`

	_, err := r.db.Exec(ctx, query, userID, source, vacancyID)
	if err != nil {
		return fmt.Errorf("РЅРµ СѓРґР°Р»РѕСЃСЊ СѓРґР°Р»РёС‚СЊ РІР°РєР°РЅСЃРёСЋ: %w", err)
	}

	return nil
}

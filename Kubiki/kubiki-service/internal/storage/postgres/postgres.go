package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"kubiki_Service/internal/domain/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Khranilishche struct {
	db *sql.DB
}

func NovoeKhranilishche(dsn string) (*Khranilishche, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("ne udalos otkrit bazu: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("net svyazi s bazoy: %w", err)
	}

	return &Khranilishche{db: db}, nil
}

func (k *Khranilishche) SohranitBrosok(ctx context.Context, polzovatelID int64, brosok models.Brosok) error {
	rezultatyJSON, err := json.Marshal(brosok.Rezultaty)
	if err != nil {
		return fmt.Errorf("ne udalos kodirovat rezultaty: %w", err)
	}

	zapros := `INSERT INTO istoriya_broskov (polzovatel_id, kolichestvo, grani, modifikator, rezultaty, itog) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = k.db.ExecContext(ctx, zapros, polzovatelID, brosok.Kolichestvo, brosok.Grani, brosok.Modifikator, rezultatyJSON, brosok.Itog)
	if err != nil {
		return fmt.Errorf("ne udalos sohranit brosok: %w", err)
	}

	return nil
}

func (k *Khranilishche) IstoriyaBroskov(ctx context.Context, polzovatelID int64, predel int32) ([]models.Brosok, error) {
	if predel <= 0 {
		predel = 20
	}

	zapros := `SELECT kolichestvo, grani, modifikator, rezultaty, itog, sozdan_v FROM istoriya_broskov WHERE polzovatel_id = $1 ORDER BY id DESC LIMIT $2`

	stroki, err := k.db.QueryContext(ctx, zapros, polzovatelID, predel)
	if err != nil {
		return nil, fmt.Errorf("oshibka polucheniya istorii broskov: %w", err)
	}
	defer stroki.Close()

	var broski []models.Brosok

	for stroki.Next() {
		var brosok models.Brosok
		var rezultatyJSON []byte

		err := stroki.Scan(&brosok.Kolichestvo, &brosok.Grani, &brosok.Modifikator, &rezultatyJSON, &brosok.Itog, &brosok.SozdanV)
		if err != nil {
			return nil, fmt.Errorf("oshibka chteniya broska: %w", err)
		}

		json.Unmarshal(rezultatyJSON, &brosok.Rezultaty)

		broski = append(broski, brosok)
	}

	return broski, stroki.Err()
}

func (k *Khranilishche) Stop() error {
	return k.db.Close()
}
package kubiki

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"

	"kubiki_Service/internal/domain/models"
)

type Kubiki struct {
	log           *slog.Logger
	khranilishche Khranilishche
}

type Khranilishche interface {
	SohranitBrosok(ctx context.Context, polzovatelID int64, brosok models.Brosok) error
	IstoriyaBroskov(ctx context.Context, polzovatelID int64, predel int32) ([]models.Brosok, error)
}

func New(log *slog.Logger, khranilishche Khranilishche) *Kubiki {
	return &Kubiki{
		log:           log,
		khranilishche: khranilishche,
	}
}

func (k *Kubiki) Brosit(ctx context.Context, polzovatelID int64, kolichestvo int32, grani int32, modifikator int32) (models.Brosok, error) {
	const op = "kubiki.Brosit"

	log := k.log.With(slog.String("op", op))

	if kolichestvo <= 0 {
		kolichestvo = 1
	}

	if grani < 2 {
		grani = 6
	}

	if kolichestvo > 100 {
		kolichestvo = 100
	}

	rezultaty := make([]int32, kolichestvo)
	var summa int32

	for i := int32(0); i < kolichestvo; i++ {
		rezultat := int32(rand.Intn(int(grani)) + 1)
		rezultaty[i] = rezultat
		summa += rezultat
	}

	brosok := models.Brosok{
		Kolichestvo: kolichestvo,
		Grani:       grani,
		Modifikator: modifikator,
		Rezultaty:   rezultaty,
		Summa:       summa,
		Itog:        summa + modifikator,
	}

	log.Info("brosok kubikov", slog.Int("grani", int(grani)), slog.Int("itog", int(brosok.Itog)))

	if polzovatelID > 0 {
		if err := k.khranilishche.SohranitBrosok(ctx, polzovatelID, brosok); err != nil {
			return brosok, fmt.Errorf("%s: %w", op, err)
		}
	}

	return brosok, nil
}

func (k *Kubiki) Istoriya(ctx context.Context, polzovatelID int64, predel int32) ([]models.Brosok, error) {
	const op = "kubiki.Istoriya"

	broski, err := k.khranilishche.IstoriyaBroskov(ctx, polzovatelID, predel)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return broski, nil
}
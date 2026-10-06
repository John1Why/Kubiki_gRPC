package kubiki

import (
	"context"

	kubiki "github.com/John1Why/proto_nri_kubiki/gen/go/kubiki"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"kubiki_Service/internal/domain/models"
)

type Kubiki interface {
	Brosit(ctx context.Context, polzovatelID int64, kolichestvo int32, grani int32, modifikator int32) (models.Brosok, error)
	Istoriya(ctx context.Context, polzovatelID int64, predel int32) ([]models.Brosok, error)
}

type serverAPI struct {
	kubiki.UnimplementedKubikiServer
	servis Kubiki
}

func Register(gRPC *grpc.Server, servis Kubiki) {
	kubiki.RegisterKubikiServer(gRPC, &serverAPI{servis: servis})
}

func (s *serverAPI) Brosit(ctx context.Context, req *kubiki.Brosok_Zapros) (*kubiki.Brosok_Otvet, error) {
	if req.GetGrani() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Grani is required")
	}

	brosok, err := s.servis.Brosit(ctx, req.GetPolzovatel_ID(), req.GetKolichestvo(), req.GetGrani(), req.GetModifikator())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &kubiki.Brosok_Otvet{
		Rezultaty:   brosok.Rezultaty,
		Summa:       brosok.Summa,
		Modifikator: brosok.Modifikator,
		Itog:        brosok.Itog,
	}, nil
}

func (s *serverAPI) Istoriya(ctx context.Context, req *kubiki.Istoriya_Zapros) (*kubiki.Istoriya_Otvet, error) {
	if req.GetPolzovatel_ID() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Polzovatel_ID is required")
	}

	broski, err := s.servis.Istoriya(ctx, req.GetPolzovatel_ID(), req.GetPredel())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	otvet := &kubiki.Istoriya_Otvet{}

	for _, brosok := range broski {
		otvet.Zapisi = append(otvet.Zapisi, &kubiki.Zapis_Istorii{
			Kolichestvo: brosok.Kolichestvo,
			Grani:       brosok.Grani,
			Modifikator: brosok.Modifikator,
			Rezultaty:   brosok.Rezultaty,
			Itog:        brosok.Itog,
			Sozdan_V:    brosok.SozdanV,
		})
	}

	return otvet, nil
}
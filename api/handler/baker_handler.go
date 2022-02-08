package handler

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pancake.maker/api/gen/api"
)

func init() {
	// seedを初期化
	rand.Seed(time.Now().UnixNano())
}

type BakerHandler struct {
	report *report
}

type report struct {
	sync.Mutex // クリティカルセクションをロックできる
	data       map[api.Pancake_Menu]int
}

func NewBakerHandler() *BakerHandler {
	return &BakerHandler{
		report: &report{
			data: make(map[api.Pancake_Menu]int),
		},
	}
}

func (h *BakerHandler) Bake(
	ctx context.Context,
	req *api.BakeRequest,
) (*api.BakeResponse, error) {

	// バリデーション
	if req.Menu == api.Pancake_UNKNOWN || req.Menu > api.Pancake_BANANA {
		return nil, status.Errorf(codes.InvalidArgument, "パンケーキを選んでください")
	}

	// パンケーキを焼いて数を記録
	now := time.Now()
	h.report.Lock()
	h.report.data[req.Menu] = h.report.data[req.Menu] + 1
	h.report.Unlock()

	// レスポンスを作って返す
	return &api.BakeResponse{
		Pancake: &api.Pancake{
			Menu:           req.Menu,
			ChefName:       "lamp",
			TechnicalScore: rand.Float32(),
			CreateTime: &timestamp.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		},
	}, nil
}

func (h *BakerHandler) Report(
	ctx context.Context,
	req *api.BakeRequest,
) (*api.ReportResponse, error) {
	counts := make([]*api.Report_BakeCount, 0)

	// レポートを作る
	h.report.Lock()
	for k, v := range h.report.data {
		counts = append(counts, &api.Report_BakeCount{
			Menu:  k,
			Count: int32(v),
		})
	}
	h.report.Unlock()

	return &api.ReportResponse{
		Report: &api.Report{
			BakeCounts: counts,
		},
	}, nil
}

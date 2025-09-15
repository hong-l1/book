package service

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"github.com/hong-l1/project/webook/internal/domain"
	"github.com/hong-l1/project/webook/internal/repository"
	"math"
	"time"
)

type RankingService interface {
	TOPN(ctx context.Context) error
}
type BatchRankingService struct {
	artSvc   ArticleService
	intrSvc  InteractiveService
	repo     repository.RankingRepository
	capicity int
	caculate func(timer time.Time, likeCnt int64) float64
	batchCnt int
}

func NewBatchRankingService(artSvc ArticleService, intrSvc InteractiveService, repo repository.RankingRepository) RankingService {
	return &BatchRankingService{
		artSvc:   artSvc,
		intrSvc:  intrSvc,
		capicity: 100,
		repo:     repo,
		batchCnt: 100,
		caculate: func(timer time.Time, likeCnt int64) float64 {
			duration := time.Since(timer).Seconds()
			return float64(likeCnt) / math.Pow(duration, 1.5)
		}}
}

func (s *BatchRankingService) TOPN(ctx context.Context) error {
	arts, err := s.tOPN(ctx)
	if err != nil {
		return err
	}
	err = s.repo.ReplaceTopN(ctx, arts)
	return err
}
func (s *BatchRankingService) tOPN(ctx context.Context) ([]domain.Article, error) {
	offset := 0
	start := time.Now()
	type Score struct {
		art   domain.Article
		score float64
	}
	temp := queue.NewPriorityQueue[Score](s.capicity, func(src Score, dst Score) int {
		if src.score > dst.score {
			return 1
		} else if src.score < dst.score {
			return -1
		} else {
			return 0
		}
	})
	for {
		arts, err := s.artSvc.ListPub(ctx, start, offset, s.capicity)
		if err != nil {
			return nil, err
		}
		ids := slice.Map[domain.Article, int64](arts, func(idx int, src domain.Article) int64 {
			return src.Id
		})
		intrs, err := s.intrSvc.Getbyids(ctx, "article", ids)
		if err != nil {
			return nil, err
		}
		for _, art := range arts {
			inter := intrs[art.Id]
			score := s.caculate(art.Ctime, inter.LikeCnt)
			ele := Score{
				score: score,
				art:   art,
			}
			err = temp.Enqueue(ele)
			if errors.Is(err, queue.ErrOutOfCapacity) {
				minEle, _ := temp.Dequeue()
				if score > minEle.score {
					_ = temp.Enqueue(ele)
				} else {
					_ = temp.Enqueue(minEle)
				}
			}
		}
		if len(arts) < s.capicity || start.Sub(arts[len(arts)-1].Utime).Hours() > 7*24 {
			break
		}
		offset += s.batchCnt
	}
	res := make([]domain.Article, temp.Len())
	for i := len(res) - 1; i >= 0; i-- {
		ele, _ := temp.Dequeue()
		res[i] = ele.art
	}
	return res, nil
}

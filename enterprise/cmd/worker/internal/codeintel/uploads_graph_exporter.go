package codeintel

import (
	"context"

	"github.com/sourcegraph/sourcegraph/cmd/worker/job"
	"github.com/sourcegraph/sourcegraph/enterprise/cmd/worker/shared/init/codeintel"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/ranking"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/internal/observation"
)

type rankingJob struct{}

func NewRankingFileReferenceCounter() job.Job {
	return &rankingJob{}
}

func (j *rankingJob) Description() string {
	return ""
}

func (j *rankingJob) Config() []env.Config {
	return []env.Config{
		ranking.ConfigInst,
	}
}

func (j *rankingJob) Routines(_ context.Context, observationCtx *observation.Context) ([]goroutine.BackgroundRoutine, error) {
	services, err := codeintel.InitServices(observationCtx)
	if err != nil {
		return nil, err
	}

	jobs := ranking.NewSymbolExporter(observationCtx, services.RankingService)
	jobs = append(jobs, ranking.NewMapper(observationCtx, services.RankingService)...)
	jobs = append(jobs, ranking.NewReducer(observationCtx, services.RankingService)...)
	return jobs, nil
}

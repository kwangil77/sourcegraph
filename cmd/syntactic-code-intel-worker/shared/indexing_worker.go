package shared

import (
	"context"
	"time"

	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/internal/codeintel/syntactic_indexing/job_store"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/internal/workerutil"
	"github.com/sourcegraph/sourcegraph/internal/workerutil/dbworker"
)

func NewIndexingWorker(ctx context.Context,
	observationCtx *observation.Context,
	store job_store.SyntacticIndexingJobStore,
	config IndexingWorkerConfig) *workerutil.Worker[*job_store.SyntacticIndexingJob] {

	name := "syntactic_code_intel_indexing_worker"

	return dbworker.NewWorker[*job_store.SyntacticIndexingJob](ctx, store.DBWorkerStore(), &indexingHandler{}, workerutil.WorkerOptions{
		Name:                 name,
		Interval:             config.PollInterval,
		HeartbeatInterval:    10 * time.Second,
		Metrics:              workerutil.NewMetrics(observationCtx, name),
		NumHandlers:          config.Concurrency,
		MaximumRuntimePerJob: config.MaximumRuntimePerJob,
	})

}

type indexingHandler struct{}

var _ workerutil.Handler[*job_store.SyntacticIndexingJob] = &indexingHandler{}

func (i indexingHandler) Handle(ctx context.Context, logger log.Logger, record *job_store.SyntacticIndexingJob) error {
	logger.Info("Stub indexing worker handling record",
		log.Int("id", record.ID),
		log.String("repository name", record.RepositoryName),
		log.String("commit", record.Commit))
	return nil
}

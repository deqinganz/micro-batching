package preprocess

import (
	. "github.com/deqinganz/batching-api/api"
	"github.com/stretchr/testify/assert"
	"micro-batching/internal/service/preprocess/processors"
	"testing"
)

type RemoveTwoJobsProcessor struct {
}

func (d *RemoveTwoJobsProcessor) Process(jobs []Job) []Job {
	return jobs[2:]
}

func TestJobProcess(t *testing.T) {
	jobPreprocessing := NewJobProcess()

	jobPreprocessing.Use("A", &processors.DummyProcessor{})
	jobPreprocessing.Use("A", &RemoveTwoJobsProcessor{})
	jobPreprocessing.Use("B", &processors.DummyProcessor{})

	jobs := jobPreprocessing.Process([]Job{
		{Type: "A"}, {Type: "B"}, {Type: "A"}, {Type: "A"},
	})

	assert.ElementsMatch(t, []Job{{Type: "A"}, {Type: "B"}}, jobs)
}

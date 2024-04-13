package processors

import (
	. "github.com/deqinganz/batching-api/api"
)

type DummyProcessor struct {
}

func (m *DummyProcessor) Process(jobs []Job) []Job {
	return jobs
}

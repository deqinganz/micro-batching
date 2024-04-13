package processors

import (
	. "github.com/deqinganz/batching-api/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDummyProcessor(t *testing.T) {
	processor := &DummyProcessor{}

	jobs := processor.Process([]Job{{}, {}})

	assert.Equal(t, []Job{{}, {}}, jobs)
}

package batching

import (
	"github.com/deqinganz/batch-processor"
	. "github.com/deqinganz/batching-api/api"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"log"
	"micro-batching/internal/config"
	"micro-batching/internal/service"
	"micro-batching/internal/service/preprocess"
	"micro-batching/internal/service/preprocess/processors"
	"time"
)

type Batching struct {
	config         config.RunConfig
	queue          service.Queue
	batchProcessor batchprocessor.BatchProcessor
	scheduler      gocron.Scheduler
	preProcess     *preprocess.JobProcess
}

func NewBatching() (Batching, error) {
	c, err := config.ReadConfig("config.json")
	if err != nil {
		return Batching{}, err
	}

	return Batching{
		config: c,
		queue:  service.Queue{},
	}, nil
}

func NewBatchingWithConfig(c config.RunConfig) Batching {
	return Batching{
		config: c,
		queue:  service.Queue{},
	}
}

// Start create a new scheduler and start it
func (b *Batching) Start() {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal("Failed to create new scheduler")
		return
	}

	b.scheduler = s
	_, err = b.scheduler.NewJob(
		gocron.DurationJob(
			time.Duration(b.config.Frequency)*time.Second,
		),
		gocron.NewTask(
			func() {
				b.Post()
			},
		),
	)

	if err != nil {
		log.Fatal("Failed to create new job")
		return
	}
	b.scheduler.Start()
}

// Restart stops the scheduler and starts it again
func (b *Batching) Restart() {
	err := b.scheduler.Shutdown()
	if err != nil {
		log.Fatal("Failed to shutdown scheduler")
		return
	}
	b.Start()
}

// Take creates a new job and adds it to the queue
func (b *Batching) Take(jobRequest JobRequest) Job {
	job := Job{
		Id:     uuid.New(),
		Status: QUEUED,
		Type:   jobRequest.Type,
		Name:   jobRequest.Name,
		Params: Job_Params(jobRequest.Params),
	}

	b.queue.Enqueue(job)

	return job
}

// JobInfo returns the job by the given id
func (b *Batching) JobInfo(id uuid.UUID) (Job, error) {
	j, err := b.queue.Find(id)
	if err != nil {
		return Job{}, err
	}

	return j, nil
}

func (b *Batching) GetFrequency() BatchFrequency {
	return BatchFrequency{
		Frequency: b.config.Frequency,
	}
}

func (b *Batching) GetBatchSize() BatchSize {
	return BatchSize{
		BatchSize: b.config.BatchSize,
	}
}

func (b *Batching) SetFrequency(frequency BatchFrequency) {
	b.config.Frequency = frequency.Frequency
}

func (b *Batching) SetBatchSize(batchSize BatchSize) {
	b.config.BatchSize = batchSize.BatchSize
}

func (b *Batching) SetPreProcess(on bool) {
	if on {
		if b.preProcess != nil {
			log.Print("preprocessing is already on")
		} else {
			b.preProcess = preprocess.NewJobProcess()
			b.preProcess.Use(UPDATEUSERINFO, &processors.DummyProcessor{})
			b.preProcess.Use(BALANCEUPDATE, &processors.BalanceUpdate{})
			log.Print("Added preprocessing")
		}
	} else {
		b.preProcess = nil
		log.Print("Removed preprocessing")
	}
}

func (b *Batching) Post() {
	if b.queue.Size() == 0 {
		log.Print("No jobs to process")
		return
	}

	if b.preProcess != nil {
		b.queue.EnqueueJobs(
			b.preProcess.Process(b.queue.Dequeue(b.queue.Size())))
	}

	jobs := b.queue.Dequeue(b.config.BatchSize)

	b.batchProcessor.Process(jobs)
}

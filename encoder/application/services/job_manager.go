package services

import (
	"encoding/json"
	"enconder/application/repositories"
	"enconder/domain"
	"enconder/framework/queue"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db *gorm.DB
	Domain domain.Job
	MessageChannel chan amqp.Delivery
	JobResultChannel chan JobWorkerResult
	RabbitMQ *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error string `json:"error"`
}

func NewJobManager(db *gorm.DB, rabbiMQ *queue.RabbitMQ, jobReturnChannel chan JobWorkerResult, messageChannel chan amqp.Delivery) *JobManager {
	return &JobManager{
		Db: db,
		Domain: domain.Job{},
		MessageChannel: messageChannel,
		JobResultChannel: jobReturnChannel,
		RabbitMQ: rabbiMQ,
	}
}

func (j *JobManager) Start(ch *amqp.Channel) {
	videoService := NewVideoService()
	videoService.VideoRepository = repositories.VideoRepositoryDb{Db: j.Db}

	jobService := JobService{
		JobRepository: repositories.JobRepositoryDb{Db: j.Db},
		VideoService: videoService,
	}

	concurrence, err := strconv.Atoi(os.Getenv("CONCURRANCY_WORKERS"))
	if err != nil {
		log.Fatal("error loading var: CONCURRANCY_WORKERS")
	}

	for qtdProcess := 0; qtdProcess < concurrence; qtdProcess++ {
		go JobWorker(j.MessageChannel, j.JobResultChannel, jobService, j.Domain, qtdProcess)
	}

	for jobResult := range j.JobResultChannel {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}

}

func (j *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing Job: #{jobRsult.Job.ID}")
	} else {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing message: #{jobRsult.Error}")
	}

	errMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error: jobResult.Error.Error(),
	}

	jobJson, err := json.Marshal(errMsg)
	if err != nil {
		return err
	}

	err = j.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) notifySuccess(jobResult JobWorkerResult, ch *amqp.Channel) error {
	jobJson, err := json.Marshal(jobResult.Job)
	if err != nil {
		return err
	}

	err = j.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)
	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) notify(jobJson []byte) error {
	err := j.RabbitMQ.Notify(
		string(jobJson),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)

	if err != nil {
		return err
	}

	return nil
}
package repositories_test

import (
	"enconder/application/repositories"
	"enconder/domain"
	"enconder/framework/database"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestJobRepositoryDbInsert(t *testing.T) {
	db := database.NewDbTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db:db}
	repo.Insert(video)

	job, err := domain.NewJob("output _path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db:db}
	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestJobRepositoryDbUpdate(t *testing.T) {
		db := database.NewDbTest()
		defer db.Close()

		video := domain.NewVideo()
		video.ID = uuid.NewV4().String()
		video.FilePath = "path"
		video.CreatedAt = time.Now()
		
		repo := repositories.VideoRepositoryDb{Db:db}
		repo.Insert(video)

		job, err := domain.NewJob("output", "Pending", video)
		require.Nil(t, err)

		repoJob := repositories.JobRepositoryDb{Db:db}
		repoJob.Insert(job)

		job.Status = "Converting"
		
		repoJob.Update(job)

		jb, err := repoJob.Find(job.ID)
		require.Nil(t, err)
		require.Equal(t, jb.Status, job.Status)
}
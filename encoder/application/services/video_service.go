package services

import (
	"context"
	"enconder/application/repositories"
	"enconder/domain"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
)

type VideoService struct {
	Video *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (service *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bucket := client.Bucket(bucketName)
	object := bucket.Object(service.Video.FilePath)

	reader, err := object.NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil
	}

	file, err := os.Create(os.Getenv("localStoragePath") + "/" + service.Video.ID + ".mp4")
	if err != nil {
		return nil
	}

	_, err = file.Write(body)
	if err != nil {
		return nil
	}

	defer file.Close()

	log.Printf("video %v has been stored", service.Video.ID)

	return nil
}

func (service *VideoService) Frament() error {
	err := os.Mkdir(os.Getenv("localStoragePath")+"/"+service.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("localStoragePath") + "/" + service.Video.ID + ".mp4"
	target := os.Getenv("localStoragePath") + "/" + service.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)
	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("======> Output: %s\n", string(out))
	}
}
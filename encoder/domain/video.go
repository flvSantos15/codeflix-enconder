package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Video struct {
	ID         string `valid:"uuid"`
	ResourceID string `valid:"notnull"`
	FilePath   string `valid:"notnull"`
	CreatedAt    time.Time `valid:"-"`
}

// Esse codigo vai rodar antes de todo mundo
func init() {
	// essa parte vai exigir que tudo seja true
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewVideo() *Video {
	return &Video{}
}

func (video *Video) Validate() error {
	_, err := govalidator.ValidateStruct(video)
	
	if err != nil {
		return err
	}

	return nil
}

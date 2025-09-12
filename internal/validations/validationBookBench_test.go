package validations

import (
	"testing"
	"time"

	"github.com/Talos-hub/BooksRestApi/internal/models"
)

func BenchmarkValidate_validBook(b *testing.B) {
	general := models.GeneralBook{
		ID:              1,
		Title:           "Some",
		Author:          "Some",
		Genre:           "some",
		PublicationDate: time.Now(),
	}

	for i := 0; i < b.N; i++ {
		Validate(general)
	}
}

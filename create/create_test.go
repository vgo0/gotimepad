package create

import (
	"testing"

	"github.com/vgo0/gotimepad/db"
	"github.com/vgo0/gotimepad/models"
)

func TestExec(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		count int64
		size int
		args args
	}{
		{"Test Default Create", 500, 2000, args{ args: []string{"-d", "file::memory:?cache=shared"}}},
		{"Test Page Count Create", 1000, 2000, args{ args: []string{"-d", "file::memory:?cache=shared", "-p", "1000"}}},
		{"Test Page Size Create", 500, 1000, args{ args: []string{"-d", "file::memory:?cache=shared", "-s", "1000"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Exec(tt.args.args)
			var count int64
			db.DB.Model(&models.OTPage{}).Count(&count)
			
			if count != tt.count {
				t.Errorf("got %d pages, want %d", count, tt.count)
			}

			var page models.OTPage
			db.DB.First(&page)

			if len(page.Page) != tt.size {
				t.Errorf("got %d page length, want %d", count, tt.count)
			}

			// get rid of current database
			imdb, _ := db.DB.DB()
			imdb.Close();
		})
	}
}

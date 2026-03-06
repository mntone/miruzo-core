package variant

import (
	"testing"

	"github.com/mntone/miruzo-core/miruzo/internal/config"
)

func TestMediaURLBuilder_Build(t *testing.T) {
	testCases := []struct {
		name     string
		basePath string
		rel      string
		want     string
	}{
		{
			name:     "keeps normalized base path",
			basePath: "/media/",
			rel:      "l1w320/xyz/a.webp",
			want:     "/media/l1w320/xyz/a.webp",
		},
		{
			name:     "adds leading and trailing slash to base path",
			basePath: "media",
			rel:      "l1w320/xyz/a.webp",
			want:     "/media/l1w320/xyz/a.webp",
		},
		// {
		// 	name:     "trims leading slash from relative path",
		// 	basePath: "/media/",
		// 	rel:      "/l1w320/xyz/a.webp",
		// 	want:     "/media/l1w320/xyz/a.webp",
		// },
		{
			name:     "handles empty base path",
			basePath: "",
			rel:      "l1w320/xyz/a.webp",
			want:     "/l1w320/xyz/a.webp",
		},
		{
			name:     "keeps base path when relative path is empty",
			basePath: "/media/",
			rel:      "",
			want:     "/media/",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			builder := NewMediaURLBuilder(config.MediaPublicConfig{
				BasePath: testCase.basePath,
			})

			got := builder.Build(testCase.rel)

			if got != testCase.want {
				t.Fatalf("Build() = %q, want %q", got, testCase.want)
			}
		})
	}
}

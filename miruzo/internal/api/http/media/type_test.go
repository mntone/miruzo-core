package media_test

import (
	"testing"

	httpmedia "github.com/mntone/miruzo-core/miruzo/internal/api/http/media"
)

func TestMediaTypeString(t *testing.T) {
	got := httpmedia.MediaType{
		Type:    "application",
		SubType: "json",
	}.String()

	if got != "application/json" {
		t.Fatalf("String() = %q, want %q", got, "application/json")
	}
}

func TestParseMediaType(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   httpmedia.MediaType
		wantOK bool
	}{
		{
			name:   "basic",
			input:  "application/json",
			want:   httpmedia.JSON,
			wantOK: true,
		},
		{
			name:   "trim and lowercase",
			input:  "  Application/JSON  ",
			want:   httpmedia.JSON,
			wantOK: true,
		},
		{
			name:   "with params",
			input:  "application/json; q=0.85",
			want:   httpmedia.JSON,
			wantOK: true,
		},
		{
			name:   "wildcard",
			input:  "*/*",
			want:   httpmedia.MediaType{Type: "*", SubType: "*"},
			wantOK: true,
		},
		{
			name:   "type wildcard",
			input:  "application/*",
			want:   httpmedia.MediaType{Type: "application", SubType: "*"},
			wantOK: true,
		},
		{
			name:   "restricted chars",
			input:  "application/vnd.example+json",
			want:   httpmedia.MediaType{Type: "application", SubType: "vnd.example+json"},
			wantOK: true,
		},
		{
			name:   "empty",
			input:  "",
			wantOK: false,
		},
		{
			name:   "no slash",
			input:  "application",
			wantOK: false,
		},
		{
			name:   "missing type",
			input:  "/json",
			wantOK: false,
		},
		{
			name:   "missing subtype",
			input:  "application/",
			wantOK: false,
		},
		{
			name:   "space before slash",
			input:  "application /json",
			wantOK: false,
		},
		{
			name:   "space after slash",
			input:  "application/ protobuf",
			wantOK: false,
		},
		{
			name:   "only params",
			input:  ";q=1",
			wantOK: false,
		},
		{
			name:   "invalid type first char",
			input:  "-application/json",
			wantOK: false,
		},
		{
			name:   "invalid subtype first char",
			input:  "application/-json",
			wantOK: false,
		},
		{
			name:   "invalid type char",
			input:  "app@lication/json",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := httpmedia.ParseMediaType(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("ParseMediaType() ok = %v, want %v", ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}
			if got != tt.want {
				t.Fatalf("ParseMediaType() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestParseMediaTypes(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		got := httpmedia.ParseMediaTypes(nil)
		if got != nil {
			t.Fatalf("ParseMediaTypes(nil) = %#v, want nil", got)
		}
	})

	t.Run("filters invalid entries", func(t *testing.T) {
		got := httpmedia.ParseMediaTypes([]string{
			"application/json",
			"invalid",
			"text/plain; charset=utf-8",
		})
		want := []httpmedia.MediaType{
			{Type: "application", SubType: "json"},
			{Type: "text", SubType: "plain"},
		}

		if len(got) != len(want) {
			t.Fatalf("len(ParseMediaTypes()) = %d, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("ParseMediaTypes()[%d] = %#v, want %#v", i, got[i], want[i])
			}
		}
	})
}

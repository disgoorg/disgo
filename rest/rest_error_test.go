package rest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/disgoorg/json/v2"
)

func TestError_Error(t *testing.T) {
	data := []struct {
		Name     string
		TestFile string
		Expected string
	}{
		{
			Name:     "array error",
			TestFile: "array_error.json",
			Expected: "50035: Invalid Form Body\nactivities -> 0 -> platform -> BASE_TYPE_CHOICES: Value must be one of ('desktop', 'android', 'ios').\nactivities -> 0 -> type -> BASE_TYPE_CHOICES: Value must be one of (0, 1, 2, 3, 4, 5).",
		},
		{
			Name:     "object error",
			TestFile: "object_error.json",
			Expected: "50035: Invalid Form Body\naccess_token -> BASE_TYPE_REQUIRED: This field is required",
		},
		{
			Name:     "request error",
			TestFile: "request_error.json",
			Expected: "50035: Invalid Form Body\n -> APPLICATION_COMMAND_TOO_LARGE: Command exceeds maximum size (8000)",
		},
		{
			Name:     "other error",
			TestFile: "other_error.json",
			Expected: "50035: Invalid Form Body\ndata -> components -> BASE_TYPE_MAX_LENGTH: Must be 5 or fewer in length.",
		},
	}

	for _, d := range data {
		t.Run(d.Name, func(t *testing.T) {
			raw, err := os.ReadFile(filepath.Join("testdata", d.TestFile))
			if err != nil {
				t.Fatal(err)
			}
			var restErr *Error
			if err = json.Unmarshal(raw, &restErr); err != nil {
				t.Fatal(err)
			}

			actual := restErr.Error()
			if actual != d.Expected {
				t.Errorf("got %q, want %q", actual, d.Expected)
			}
		})
	}
}

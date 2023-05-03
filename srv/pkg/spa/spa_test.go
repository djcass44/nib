package spa

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

//go:embed testdata/test.txt
var textData string

//go:embed testdata/index.html
var htmlData string

func TestFileSystem_Open(t *testing.T) {
	fs := NewFileSystem(http.Dir("./testdata"))
	ts := httptest.NewServer(http.FileServer(fs))
	defer ts.Close()

	var cases = []struct {
		path string
		code int
		body string
	}{
		{
			"/",
			http.StatusOK,
			htmlData,
		},
		{
			"/test.txt",
			http.StatusOK,
			textData,
		},
		{
			"/notfound",
			http.StatusOK,
			htmlData,
		},
	}
	for _, tt := range cases {
		t.Run(tt.path, func(t *testing.T) {
			resp, err := ts.Client().Get(ts.URL + tt.path)
			assert.NoError(t, err)
			defer resp.Body.Close()
			assert.EqualValues(t, tt.code, resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.body, string(body))
		})
	}
}

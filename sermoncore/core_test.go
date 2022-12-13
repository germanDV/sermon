package sermoncore

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"gitlab.com/germandv/sermon/expect"
	"gitlab.com/germandv/sermon/internal/httpclient"
)

func TestIn(t *testing.T) {
	codes := []StatusCode{{200}, {301}, {415}, {500}}

	t.Run("ReturnsFalseWhenCodeIsNotInSlice", func(t *testing.T) {
		t.Parallel()
		found := in(codes, 400)
		expect.Equal(t, found, false)
	})

	t.Run("ReturnsTrueWhenCodeIsInSlice", func(t *testing.T) {
		t.Parallel()
		found := in(codes, 200)
		expect.Equal(t, found, true)
	})

	t.Run("ReturnsFalseWhenSliceIsEmpty", func(t *testing.T) {
		t.Parallel()
		found := in([]StatusCode{}, 204)
		expect.Equal(t, found, false)
	})
}

func TestHealthRequest(t *testing.T) {
	url, _ := url.Parse("http://localhost:4000/health_check")
	service := &Service{
		Name:     "localhost",
		Endpoint: Endpoint{URL: url},
		Codes:    []StatusCode{{Code: 200}},
		Timeout:  Timeout{Duration: 5 * time.Second},
	}

	t.Run("NoErrorWhenResponseHasExpectedStatusCode", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		mockClient := httpclient.NewMock(ts.Client(), ts.URL)
		err := service.Health(mockClient)
		expect.NoError(t, err)
	})

	t.Run("ErrorWhenResponseHasUnexpectedStatusCode", func(t *testing.T) {
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))
		defer ts.Close()

		mockClient := httpclient.NewMock(ts.Client(), ts.URL)
		err := service.Health(mockClient)
		expect.Contains(t, err.Error(), "Got status 502, want one of [{200}]")
	})
}

package acme

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBoulderTestSRVPresent(t *testing.T) {
	var requestReceived bool

	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true

		if got, want := r.Method, "POST"; got != want {
			t.Errorf("Expected method to be '%s' but got '%s'", want, got)
		}
		if got, want := r.URL.Path, "/set-txt"; got != want {
			t.Errorf("Expected path to be '%s' but got '%s'", want, got)
		}
		if got, want := r.Header.Get("Content-Type"), "application/json"; got != want {
			t.Errorf("Expected Content-Type to be '%s' but got '%s'", want, got)
		}

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Error reading request body: %v", err)
		}
		if got, want := string(reqBody), `{"host":"_acme-challenge.example.com.","value":"w6uP8Tcg6K2QR905Rms8iXTlksL6OD1KOWBxTK7wxPI"}`; got != want {
			t.Errorf("Expected body data to be: `%s` but got `%s`", want, got)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer mock.Close()
	dnsURL := mock.URL

	boulderProv, err := NewDNSProviderBoulderTestSRV(dnsURL)
	if boulderProv == nil {
		t.Fatal("Expected non-nil Boulder DNS Test Server provider, but was nil")
	}
	if err != nil {
		t.Fatalf("Expected no error creating provider, but got: %v", err)
	}

	err = boulderProv.Present("example.com", "", "foobar")
	if err != nil {
		t.Fatalf("Expected no error creating TXT record, but got: %v", err)
	}
	if !requestReceived {
		t.Error("Expected request to be received by mock backend, but it wasn't")
	}
}

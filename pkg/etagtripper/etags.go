package etagtripper

import "net/http"

type contextKey string

func (c contextKey) String() string {
	return "mylib " + string(c)
}

// ContextETag is the context to pass etags to the transport
var (
	ContextETag = contextKey("etag")
)

// Custom transport to chain into the HTTPClient to gather statistics.
type ETagTransport struct {
	Next *http.Transport
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to provide stats and handle error rates.
func (t *ETagTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if etag, ok := req.Context().Value(ContextETag).(string); ok {
		req.Header.Set("if-none-match", etag)
	}

	// Run the request.
	return t.Next.RoundTrip(req)
}

package listener

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/garyburd/go-oauth/oauth"
)

// filterURL is Twitter Filter Streaming API endpoint
const filterURL = "https://stream.twitter.com/1.1/statuses/filter.json"

// stream initiates streaming connection and starts receiving in an infinite loop.
func (s *httpStreamer) stream(c chan int) {
	defer func() {
		c <- 1
	}()

	reader, err := s.open()
	if err != nil {
		log.Printf("ERROR opening stream for %q: %v", s.app.Name, err)
		return
	}
	defer reader.Close()

	log.Printf("CONNECTED to %q stream", s.app.Name)
	s.loop(reader)
}

// open creates a new HTTP streaming connection to the Filter endpoint.
// The caller is responsible for closing the stream.
func (s *httpStreamer) open() (io.ReadCloser, error) {
	params := url.Values{"follow": {strings.Join(s.users, ",")}}
	req, err := http.NewRequest("POST", filterURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	// TODO: extract req signing and move this into auth.go
	cl := &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  s.app.TwConsumerKey,
			Secret: s.app.TwConsumerSecret,
		},
	}
	creds := &oauth.Credentials{
		Token:  s.app.TwAccessToken,
		Secret: s.app.TwTokenSecret,
	}
	auth := cl.AuthorizationHeader(creds, req.Method, req.URL, params)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", auth)

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return nil, respErr
	}
	// TODO: handle transient errors 420 & 503 with exponential backoff.
	// https://dev.twitter.com/docs/streaming-apis/connecting#HTTP_Error_Codes
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return resp.Body, nil
}

// loop reads from provided reader stream and logs received data indefinitely.
func (s *httpStreamer) loop(reader io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, readErr := reader.Read(buf)
		if n > 0 {
			log.Printf("Got %d bytes: \n%s\n", n, string(buf[:n]))
		}

		if readErr != nil {
			log.Printf("Connection error: %v", readErr)
			break
		}
	}
}
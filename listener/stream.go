package listener

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/garyburd/go-oauth/oauth"
)

// filterURL is Twitter Filter Streaming API endpoint
const filterURL = "https://stream.twitter.com/1.1/statuses/filter.json"

var (
	// endOfTweet is a tweet delimiter in the stream.
	endOfTweet = []byte{13, 10}
	lenEOT     = len(endOfTweet)
)

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
	// params := url.Values{"track": []string{"Twitter"}}
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
	var (
		buf     = make([]byte, 1024)
		stream  bytes.Buffer
		lastIdx = 0
	)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			stream.Write(buf[:n])
			body := stream.Bytes()
			if i := bytes.Index(body[lastIdx:], endOfTweet); i >= 0 {
				i += lastIdx
				go s.digest(unmarshalTweet(body[:i+lenEOT-1]))
				stream.Truncate(0)
				stream.Write(body[i+lenEOT:])
				lastIdx = 0
			} else {
				lastIdx = len(body) - lenEOT
				if lastIdx < 0 {
					lastIdx = 0
				}
			}
		}

		if err != nil {
			log.Printf("%q stream: %v", s.app.Name, err)
			break
		}
	}
}

// digest pushes the tweet to a processing queue or silently ignores it
// if nil or does not belong to the application's users.
func (s *httpStreamer) digest(tweet *Tweet) {
	// TODO: favorites might be useful too.
	if tweet == nil || !isInList(s.users, tweet.User.ID) {
		return
	}

	log.Printf("%v", tweet)
}

// unmarshalTweet parses jsonTweet bytes into Tweet struct.
func unmarshalTweet(jsonTweet []byte) *Tweet {
	if len(jsonTweet) < 3 {
		return nil
	}

	tweet := &Tweet{}
	if err := json.Unmarshal(jsonTweet, tweet); err != nil {
		log.Printf("ERROR parsing tweet: %v", err)
		return nil
	}

	return tweet
}

// isInList returns true if list contains elem using binary search.
// list is assumed to be already sorted in ascending order.
func isInList(list []string, elem string) bool {
	i := sort.SearchStrings(list, elem)
	return i < len(list) && list[i] == elem
}

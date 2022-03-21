package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const searchRecentTweetsURL = "https://api.twitter.com/2/tweets/search/recent"

type client struct {
	bearerToken string
	client      *http.Client
}

type Tweet struct {
	ID                string `json:"id"`
	Text              string `json:"text"`
	AuthorID          string `json:"author_id,omitempty"`
	ConversationID    string `json:"conversation_id,omitempty"`
	CreatedAt         string `json:"created_at"`
	InReplyToUserID   string `json:"in_reply_to_user_id,omitempty"`
	Lang              string `json:"lang,omitempty"`
	PossiblySensitive bool   `json:"possibly_sensitive,omitempty"`
	ReplySettings     string `json:"reply_settings,omitempty"`
	Source            string `json:"source,omitempty"`
}

type HTTPError struct {
	APIName string
	Status  string
	URL     string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%s: %s %s", e.APIName, e.Status, e.URL)
}

type SearchTweetsResponse struct {
	Tweets []*Tweet `json:"data"`
	Title  string   `json:"title,omitempty"`
	Detail string   `json:"detail,omitempty"`
	Type   string   `json:"type,omitempty"`
}

func (c *client) SearchRecentTweets(ctx context.Context, tweet string) (*SearchTweetsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchRecentTweetsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("search recent tweets new request with context: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.bearerToken))
	params := req.URL.Query()
	params.Add("query", tweet)
	req.URL.RawQuery = params.Encode()

	fmt.Println(req.URL.String())
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search recent tweets do: %w", err)
	}

	defer resp.Body.Close()

	var str SearchTweetsResponse
	if err := json.NewDecoder(resp.Body).Decode(&str); err != nil {
		return nil, fmt.Errorf("search recent tweets: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return &str, &HTTPError{
			APIName: "search recent tweets",
			Status:  resp.Status,
			URL:     req.URL.String(),
		}
	}

	return &str, nil
}

func newClient(bearerToken string) *client {
	return &client{
		bearerToken: bearerToken,
		client:      http.DefaultClient,
	}
}

func main() {
	http.HandleFunc("/search", searchTweet)
	http.ListenAndServe(":8080", nil)
}

func searchTweet(w http.ResponseWriter, r *http.Request) {
	client := newClient("")

	res, err := client.SearchRecentTweets(context.Background(), "ビットコイン")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(res); err != nil {
		panic(err)
	}

	_, err = fmt.Fprint(w, buf.String())
	if err != nil {
		return
	}
}

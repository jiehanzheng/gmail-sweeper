package api

import (
	"fmt"
	"log"

	"google.golang.org/api/gmail/v1"
)

func GetAllMessages(srv *gmail.Service, lc *gmail.UsersMessagesListCall) ([]*gmail.Message, error) {
	// https://github.com/google/google-api-go-client/blob/a69f0f19d246419bb931b0ac8f4f8d3f3e6d4feb/examples/gmail.go
	var msgs []*gmail.Message
	pageToken := ""
	for {
		if pageToken != "" {
			lc.PageToken(pageToken)
		}
		r, err := lc.Do()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve messages: %v", err)
		}

		log.Printf("Processing %v messages...\n", len(r.Messages))
		for _, m := range r.Messages {
			msg, err := srv.Users.Messages.Get(ME, m.Id).Do()
			if err != nil {
				return nil, fmt.Errorf("Unable to retrieve message %v: %v", m.Id, err)
			}

			msgs = append(msgs, msg)
		}

		if r.NextPageToken == "" {
			break
		}
		pageToken = r.NextPageToken
	}

	return msgs, nil
}

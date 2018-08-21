package main // import "github.com/jiehanzheng/gmail-sweeper"

import (
	"fmt"
	"log"

	"github.com/jiehanzheng/gmail-sweeper/gmail/api"
	"github.com/jiehanzheng/gmail-sweeper/sweep_policy"
	"google.golang.org/api/gmail/v1"
)

func main() {
	srv := api.GetAPIClient()

	// Fetch
	listCall := srv.Users.Messages.List(api.ME)
	listCall = sweep_policy.Fetch(listCall)
	messages, err := api.GetAllMessages(srv, listCall)
	if err != nil {
		log.Fatal(err)
	}

	// Group
	groupedMsgs := make(map[sweep_policy.GroupId][]*gmail.Message)
	for _, msg := range messages {
		grp := sweep_policy.Group(msg)
		groupedMsgs[grp] = append(groupedMsgs[grp], msg)
	}

	fmt.Println(groupedMsgs)

	// Sweep
	for gId, msgs := range groupedMsgs {
		log.Printf("Sweeping group `%v`...", gId)
		sweep_policy.ExecuteGroupPolicy(sweep_policy.Sweep(gId), srv, msgs, false)
	}
}

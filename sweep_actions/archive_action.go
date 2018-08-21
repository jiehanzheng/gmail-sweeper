package sweep_actions

import (
	"log"

	"github.com/jiehanzheng/gmail-sweeper/gmail/api"
	"google.golang.org/api/gmail/v1"
)

type ArchiveDoer struct{}

func (d ArchiveDoer) Do(srv *gmail.Service, ids []api.MessageId) {
	var idStrs []string
	for _, id := range ids {
		idStrs = append(idStrs, string(id))
	}

	err := srv.Users.Messages.BatchModify(api.ME, &gmail.BatchModifyMessagesRequest{
		Ids:            idStrs,
		RemoveLabelIds: []string{"INBOX"}}).Do()
	if err != nil {
		log.Fatalf("ArchiveDoer.Do: BatchModify failed: %v", err)
	}
}

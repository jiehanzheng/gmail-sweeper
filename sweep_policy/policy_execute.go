package sweep_policy

import (
	"log"
	"reflect"

	"github.com/jiehanzheng/gmail-sweeper/gmail/api"
	"google.golang.org/api/gmail/v1"
)

func ExecuteGroupPolicy(gp GroupPolicy, srv *gmail.Service, ms []*gmail.Message, dry bool) {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	// Do action on messages that we do no not retain
	// Retain the first (gp.retain) messages
	var ids []api.MessageId
	for _, m := range ms[min(len(ms), gp.retain):] {
		ids = append(ids, api.MessageId(m.Id))
	}

	log.Printf("ExecuteGroupPolicy: will perform %v on %v", reflect.TypeOf(gp.action_doer), ids)

	if !dry {
		for i := 0; i < len(ids); i += api.BATCH_SIZE {
			batch := ids[i:min(i+api.BATCH_SIZE, len(ids))]
			gp.action_doer.Do(srv, batch)
		}
	} else {
		log.Printf("ExecuteGroupPolicy: dry run")
	}
}

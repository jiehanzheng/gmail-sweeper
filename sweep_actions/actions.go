package sweep_actions

import (
	"github.com/jiehanzheng/gmail-sweeper/gmail/api"
	"google.golang.org/api/gmail/v1"
)

type ActionDoer interface {
	Do(*gmail.Service, []api.MessageId)
}

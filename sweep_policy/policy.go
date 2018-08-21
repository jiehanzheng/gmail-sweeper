package sweep_policy

import (
	"log"
	"net/mail"
	"regexp"

	"github.com/jiehanzheng/gmail-sweeper/gmail/api"
	"github.com/jiehanzheng/gmail-sweeper/sweep_actions"
	"google.golang.org/api/gmail/v1"
)

type GroupId string

type GroupPolicy struct {
	retain      int
	action_doer sweep_actions.ActionDoer
}

func Fetch(lc *gmail.UsersMessagesListCall) *gmail.UsersMessagesListCall {
	return lc.LabelIds("INBOX").Q("label:newsletters")

	// ProTip: Use the following to quickly target a test email
	//   .Q("subject:test newer_than:1d")
}

func Group(m *gmail.Message) GroupId {
	from, err := api.GetMessageHeader(m.Payload.Headers, "From")
	if err != nil {
		log.Fatalf("Group: unable to find From header in message: %v", err)
	}

	parsedFrom, err := mail.ParseAddress(from)
	if err != nil {
		log.Fatalf("Group: unable to parse From header [%v]: %v", from, err)
	}

	parsedFromAddr := parsedFrom.Address

	switch parsedFromAddr {
	default:
		return GroupId(parsedFromAddr)
	case "nytdirect@nytimes.com":
		// If NYT, find the name of newsletter
		msgText, err := api.ExtractMessageText(m.Payload)
		if err != nil {
			log.Fatalf("Group: NYT case: %s", err)
		}
		re := regexp.MustCompile("You received this message because you signed up for NYTimes.com's (.{1,128}) newsletter.")
		matches := re.FindStringSubmatch(msgText)
		nytNewsletterName := "UNKNOWN"
		if matches != nil {
			nytNewsletterName = matches[1]
		}
		return GroupId(parsedFromAddr + "~" + nytNewsletterName)
	}
}

func Sweep(g GroupId) GroupPolicy {
	return GroupPolicy{retain: 1, action_doer: sweep_actions.ArchiveDoer{}}
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"
	"gopkg.in/yaml.v3"
)

var app = struct {
	ID     string
	Token  string
	Config string
}{
	ID:     os.Getenv("SLACK_APP_ID"),
	Token:  os.Getenv("SLACK_TOKEN"),
	Config: os.Getenv("SLACK_APP_URL"),
}

const (
	Limit = 100
)

var tpl = template.Must(template.New("slack").Parse(`---
{{ .FrontMatter -}}
---
# Title

<!-- raw
` + "```json" + `
{{ .Raw }}` + "```" + `
`))

type FrontMatter struct {
	InternalID uuid.UUID `json:"uid" yaml:"uid"`
	ExternalID string    `json:"xid" yaml:"xid"`
	Tags       []string  `json:"tags" yaml:"tags"`
	URL        string    `json:"url" yaml:"url"`
}

func (fm FrontMatter) String() string {
	buf := bytes.NewBuffer(nil)
	enc := yaml.NewEncoder(buf)
	enc.SetIndent(2)
	unsafe.Ignore(enc.Encode(fm))
	return buf.String()
}

type Markdown struct {
	FrontMatter
	raw interface{}
}

func (md Markdown) Raw() string {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	unsafe.Ignore(enc.Encode(md.raw))
	return buf.String()
}

func (md Markdown) DumpTo(file io.WriteCloser) error {
	defer safe.Close(file, unsafe.Ignore)
	return tpl.Execute(file, md)
}

func main() {
	var cursor string

	api := slack.New(app.Token)
	ctx := context.Background()

	user, err := api.AuthTestContext(ctx)
	if err != nil {
		panic(err)
	}

	var target *slack.Channel
	cursor = ""
	for {
		channels, next, err := api.GetConversationsContext(ctx, &slack.GetConversationsParameters{
			Cursor:          cursor,
			ExcludeArchived: true,
			Limit:           Limit,
			Types:           []string{slack.TYPE_IM},
		})
		if err != nil {
			panic(err)
		}

		for _, channel := range channels {
			if channel.NumMembers == 0 && channel.User == user.UserID {
				channel := channel
				target = &channel
				break
			}
		}
		cursor = next
		if cursor == "" {
			break
		}
	}

	if target == nil {
		panic("there is no saved messages chat")
	}

	var messages []slack.Message
	cursor = ""
	for {
		resp, err := api.GetConversationHistoryContext(ctx, &slack.GetConversationHistoryParameters{
			ChannelID: target.ID,
			Cursor:    cursor,
			Limit:     Limit,
		})
		if err != nil {
			panic(err)
		}
		if !resp.Ok {
			panic(errors.New(resp.Error))
		}

		messages = append(messages, resp.Messages...)
		cursor = resp.ResponseMetaData.NextCursor
		if cursor == "" {
			break
		}
	}

	if len(messages) == 0 {
		panic("there are no saved messages")
	}

	tags := []string{"slack", "message"}
	for _, message := range messages {
		f, err := os.Create(fmt.Sprintf("stream/slack/slack-%s.md", normalize(message.ClientMsgID)))
		if err != nil {
			panic(err)
		}

		md := Markdown{
			FrontMatter: FrontMatter{
				InternalID: uuid.New(),
				ExternalID: normalize(message.ClientMsgID),
				Tags:       tags,
				URL:        message.Permalink,
			},
			raw: message,
		}
		if err := md.DumpTo(f); err != nil {
			panic(err)
		}
	}
}

func normalize(in string) string {
	return strings.ToLower(in)
}

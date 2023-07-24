//go:build telegram

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/google/uuid"
	tdlib "github.com/zelenin/go-tdlib/client"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"
	"gopkg.in/yaml.v3"
)

var app = struct {
	ID     int32
	Hash   string
	Title  string
	Short  string
	Config string
}{
	ID:     int32(unsafe.Return(strconv.Atoi(os.Getenv("TELEGRAM_APP_ID"))).(int)),
	Hash:   os.Getenv("TELEGRAM_APP_HASH"),
	Title:  "Sparkle Service",
	Short:  "Sparkle",
	Config: "https://my.telegram.org/apps",
}

const (
	Fatal int32 = iota
	Error
	Warning
	Info
	Debug
	Verbose

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
	ExternalID int64     `json:"xid" yaml:"xid"`
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
	ctx := context.Background()

	if _, err := tdlib.SetLogVerbosityLevel(&tdlib.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: Error,
	}); err != nil {
		panic(err)
	}
	tdlib.WithExtraGenerator(uuid.NewString)

	auth := tdlib.NewAppAuthorizer()
	auth.TdlibParameters <- &tdlib.SetTdlibParametersRequest{
		ApiId:   app.ID,
		ApiHash: app.Hash,

		DatabaseDirectory: filepath.Join(".obsidian", "db"),
		FilesDirectory:    filepath.Join(".obsidian", "files"),

		EnableStorageOptimizer: true,
		UseChatInfoDatabase:    true,
		UseFileDatabase:        true,
		UseMessageDatabase:     true,

		// in the future
		ApplicationVersion:    "v0.0.1",
		DatabaseEncryptionKey: nil,
		DeviceModel:           "Service",
		SystemLanguageCode:    "en",
	}
	go login(auth.State, auth.PhoneNumber, auth.Code, auth.Password)

	client := tdlib.NewClient()
	if err := tdlib.Authorize(ctx, client, auth); err != nil {
		panic(err)
	}

	user, err := client.GetMe(ctx)
	if err != nil {
		panic(err)
	}

	var messages []*tdlib.Message
	cursor := int64(0)
	for {
		resp, err := client.GetChatHistory(ctx, &tdlib.GetChatHistoryRequest{
			ChatId:        user.Id,
			FromMessageId: cursor,
			Limit:         Limit,
			OnlyLocal:     false,
		})
		if err != nil {
			panic(err)
		}
		if resp.TotalCount == 0 {
			break
		}

		messages = append(messages, resp.Messages...)
		cursor = resp.Messages[len(resp.Messages)-1].Id
	}

	if len(messages) == 0 {
		panic("there are no saved messages")
	}

	tags := []string{"telegram", "message"}
	for _, message := range messages {
		f, err := os.Create(fmt.Sprintf("stream/telegram/tg-%d.md", message.Id))
		if err != nil {
			panic(err)
		}

		link, err := client.GetMessageLink(ctx, &tdlib.GetMessageLinkRequest{
			ChatId:    message.ChatId,
			MessageId: message.Id,
		})
		if err != nil {
			var valid tdlib.ResponseError
			if !errors.As(err, &valid) {
				panic(err)
			}
			link = &tdlib.MessageLink{Link: ""}
		}

		md := Markdown{
			FrontMatter: FrontMatter{
				InternalID: uuid.New(),
				ExternalID: message.Id,
				Tags:       append(tags, classifier(message.Content.MessageContentType())),
				URL:        link.Link,
			},
			raw: message,
		}
		if err := md.DumpTo(f); err != nil {
			panic(err)
		}
	}
}

func login(status <-chan tdlib.AuthorizationState, phone, code, pass chan<- string) {
	for state := range status {
		switch state.AuthorizationStateType() {
		case tdlib.TypeAuthorizationStateWaitPhoneNumber:
			fmt.Println("Enter phone number:")
			var value string
			_, _ = fmt.Scanln(&value)
			phone <- value

		case tdlib.TypeAuthorizationStateWaitCode:
			fmt.Println("Enter code:")
			var value string
			_, _ = fmt.Scanln(&value)
			code <- value

		case tdlib.TypeAuthorizationStateWaitPassword:
			fmt.Println("Enter password:")
			var value string
			_, _ = fmt.Scanln(&value)
			pass <- value

		case tdlib.TypeAuthorizationStateReady:
			return
		}
	}
}

func classifier(t string) string {
	return strings.ToLower(strings.TrimPrefix(t, "message"))
}

// if it returns error 400 Chat not found
func debug(user int64, client interface {
	GetChats(req *tdlib.GetChatsRequest) (*tdlib.Chats, error)
}) {
	chats, err := client.GetChats(&tdlib.GetChatsRequest{
		Limit: 100,
	})
	if err != nil {
		panic(err)
	}
	for _, chat := range chats.ChatIds {
		fmt.Println("chat id:", chat)
	}
	fmt.Println("user id:", user)
}

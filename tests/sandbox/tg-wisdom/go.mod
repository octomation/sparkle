module sparkle/sandbox/tg-wisdom

go 1.21.3

require (
	github.com/google/uuid v1.3.1
	github.com/zelenin/go-tdlib v0.7.0
	go.octolab.org v0.12.2
	golang.org/x/crypto v0.14.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/zelenin/go-tdlib => github.com/withsparkle/go-tdlib v0.7.1-0.20231011160900-c4aee4541527

require (
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/term v0.13.0 // indirect
)

module go.octolab.org/ecosystem/sparkle

go 1.21

toolchain go1.21.3

require (
	connectrpc.com/connect v1.11.1
	github.com/fatih/color v1.15.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.0
	github.com/slack-go/slack v0.12.3
	github.com/spf13/afero v1.10.0
	github.com/spf13/cobra v1.7.0
	github.com/stretchr/testify v1.8.4
	github.com/zelenin/go-tdlib v0.7.0
	go.octolab.org v0.12.2
	go.octolab.org/toolkit/cli v0.6.3
	go.octolab.org/toolkit/config v0.0.4
	golang.org/x/crypto v0.14.0
	golang.org/x/net v0.17.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/zelenin/go-tdlib => github.com/withsparkle/go-tdlib v0.7.1-0.20231011160900-c4aee4541527

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/sagikazarmark/locafero v0.3.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.17.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/term v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

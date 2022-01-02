module github.com/nuvolaris/nuvolaris-cli/nuv

go 1.17

replace github.com/apache/openwhisk-cli/wski18n => ../openwhisk-cli/wski18n

replace github.com/apache/openwhisk-cli/commands => ../openwhisk-cli/commands

replace github.com/go-task/task/cmd/task => ../task/cmd/task

require (
	github.com/alecthomas/kong v0.2.22
	github.com/apache/openwhisk-cli/commands v0.0.0-00010101000000-000000000000
	github.com/apache/openwhisk-cli/wski18n v0.0.0-00010101000000-000000000000
	github.com/apache/openwhisk-client-go v0.0.0-20211007130743-38709899040b
	github.com/go-task/task/cmd/task v0.0.0-00010101000000-000000000000
	github.com/nicksnyder/go-i18n v1.10.1
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/alessio/shellescape v1.4.1 // indirect
	github.com/apache/openwhisk-wskdeploy v0.0.0-20211214002128-60983d9412cc // indirect
	github.com/cloudfoundry/jibber_jabber v0.0.0-20151120183258-bcc4c8345a21 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/evanphx/json-patch/v5 v5.2.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/go-task/task/v3 v3.9.2 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20210113012101-fb4e108d2519 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-zglob v0.0.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/radovskyb/watcher v1.0.7 // indirect
	github.com/spf13/cobra v1.3.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	golang.org/x/term v0.0.0-20210916214954-140adaaadfaf // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/apimachinery v0.20.2 // indirect
	mvdan.cc/sh/v3 v3.4.2 // indirect
	sigs.k8s.io/kind v0.11.1 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

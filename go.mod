module github.com/leezer3379/sms-sender

go 1.14

require (
	github.com/garyburd/redigo v1.6.0
	github.com/toolkits/pkg v1.1.1
	go.uber.org/automaxprocs v1.3.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
replace (
	github.com/leezer3379/sms-sender/config => ./config
	github.com/leezer3379/sms-sender/cron => ./corn
	github.com/leezer3379/sms-sender/redisc => ./redisc
	github.com/leezer3379/sms-sender/corp => ./corp
)

package config

import (
	"fmt"
	"os"
	//    "time"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type topic struct {
	Enabled          bool   `yaml:"enabled"`
	ChatTopicNetCard int    `yaml:"chatTopicNetCard"`
	RootSleepTime    uint64 `yaml:"rootSleepTime"`
	ChatSleepTime    uint64 `yaml:"chatSleepTime"`
}
type bot struct {
	GuildID       string `yaml:"GuildID"`
	RootChannelID string `yaml:"RootChannelID"`
	ChatChannelID string `yaml:"ChatChannelID"`
	ChannelTopic  topic  `yaml:"ChannelTopic"`
	Token         string `yaml:"Token"`
}
type bdsEz struct {
	CrashStart     bool   `yaml:"CrashStart"`
	ZeroCrashStart bool   `yaml:"ZeroCrashStart"`
	StartPath      string `yaml:"StartPath"`
	CronBackup     string `yaml:"CronBackup"`
}
type Config struct {
	LogFile bool  `yaml:"logFile"`
	Bot     bot   `yaml:"bot"`
	Bds     bdsEz `yaml:"bdsEz"`
}

const (
	FileName = "strbot.yaml"
	Data     = `logFile: false
bdsEz:
  StartPath: "bedrock_server_mod.exe"
  CrashStart: false
  ZeroCrashStart: false
  CronBackup: "* * * * * *"
bot:
  GuildID: ""
  RootChannelID: ""
  ChatChannelID: ""
  ChannelTopic:
    enabled: false
    chatTopicNetCard: 0
#   Channel Topic Update Time, Second
    rootSleepTime: 120
    chatSleepTime: 60
  Token: ""`
)

func YamlConfig() (*Config, bool) {
	yamlFile, err := ioutil.ReadFile(FileName)
	if err != nil {
		fmt.Println("GetYamlFileErr", err)
		f, err := os.Create(FileName)
		defer f.Close()
		if err != nil {
			fmt.Println("CreateErr", err.Error())
			conf := Config{}
			return &conf, false
		} else {
			_, err = f.Write([]byte(Data))
			if err != nil {
				fmt.Println("WriteErr", err)
				conf := Config{}
				return &conf, false
			} else {
				fmt.Println("The configuration file is successfully created, please edit and restart.")
			}
		}
	}
	conf := Config{}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println("config.Yaml.UnmarshalError: %v", err)
	}
	if confNoNil(&conf) {
		return &conf, false
	}
	return &conf, true
}

const ConfNil = "Please fill in the configuration file before starting."

func confNoNil(c *Config) bool {
	if c.Bot.GuildID == "" || c.Bot.RootChannelID == "" {
		return false
	}
	if c.Bot.ChatChannelID == "" || c.Bot.Token == "" {
		return false
	}
	if c.Bds.StartPath == "" {
		return false
	}
	return true
}

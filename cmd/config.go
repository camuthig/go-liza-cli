package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	UserUUID     string                 `json:"user_uuid" mapstructure:"user_uuid"`
	Username     string                 `json:"username" mapstructure:"username"`
	Token        string                 `json:"token" mapstructure:"token"`
	Repositories map[string]*Repository `json:"repositories" mapstructure:"repositories"`
}

type Repository struct {
	Name         string               `json:"name" mapstructure:"name"`
	PullRequests map[int]*PullRequest `json:"pull_requests" mapstructure:"pull_requests"`
}

type PullRequest struct {
	ID          int                 `json:"id" mapstructure:"id"`
	Title       string              `json:"title" mapstructure:"title"`
	LastRead    time.Time           `json:"last_read" mapstructure:"last_read"`
	LastUpdated time.Time           `json:"last_updated" mapstructure:"last_updated"`
	Author      User                `json:"author" mapstructure:"author"`
	Links       PullRequestLinks    `json:"links" mapstructure:"links"`
	Updates     []PullRequestUpdate `json:"updates" mapstructure:"updates"`
}

func (p PullRequest) UnreadUpdates() int {
	if p.LastRead.After(p.LastUpdated) {
		return 0
	}

	return len(p.Updates)
}

type PullRequestLinks struct {
	HTML Link `json:"html" mapstructure:"html"`
}

type ActivityType string

const (
	Update   ActivityType = "update"
	Approval              = "approval"
	Comment               = "comment"
)

type PullRequestUpdate struct {
	Date         time.Time    `json:"date" mapstructure:"date"`
	ActivityType ActivityType `json:"activity_type" mapstructure:"activity_type"`
	Author       User         `json:"author" mapstructure:"author"`
}

type User struct {
	DisplayName string    `json:"display_name" mapstructure:"display_name"`
	UUID        string    `json:"uuid" mapstructure:"uuid"`
	Links       UserLinks `json:"links" mapstructure:"links"`
}

type UserLinks struct {
	Avatar Link `json:"avatar" mapstructure:"avatar"`
}

type Link struct {
	Href string `json:"href" mapstructure:"href"`
}

func toTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}

func ParseConfig() Config {
	var c Config
	err := viper.Unmarshal(&c, func(config *mapstructure.DecoderConfig) {
		config.DecodeHook = toTimeHookFunc()
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return c
}

func (c *Config) Write() {
	s, err := json.MarshalIndent(c, "", "  ")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(viper.ConfigFileUsed(), s, 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return
}

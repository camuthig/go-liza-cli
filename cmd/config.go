package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
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
	ID                 int                 `json:"id" mapstructure:"id"`
	Title              string              `json:"title" mapstructure:"title"`
	ReadAt             time.Time           `json:"read_at" mapstructure:"read_at"`
	PreviouslyReadAt   time.Time           `json:"previously_read_at" mapstructure:"previously_read_at"`
	LastUpdated        time.Time           `json:"last_updated" mapstructure:"last_updated"`
	Author             User                `json:"author" mapstructure:"author"`
	Links              PullRequestLinks    `json:"links" mapstructure:"links"`
	UnreadUpdatesCount int                 `json:"unread_updates_count" mapstructure:"unread_updates_count"`
	Updates            []PullRequestUpdate `json:"updates" mapstructure:"updates"`
}

func (p *PullRequest) MarkRead() {
	p.UnreadUpdatesCount = 0
	p.PreviouslyReadAt = p.ReadAt
	p.ReadAt = time.Now()
}

func (p *PullRequest) MarkUnread(c *Config) {
	p.ReadAt = p.PreviouslyReadAt
	p.UnreadUpdatesCount = p.CountUnread(c)
}

func (p PullRequest) CountUnread(c *Config) int {
	unread := 0
	for _, u := range p.Updates {
		if u.Date.Before(p.ReadAt) {
			break
		}

		p.Updates = append(p.Updates, u)

		if u.IsNewToUserUUID(c.UserUUID, p.ReadAt) {
			unread++
		}
	}

	return unread
}

type PullRequestLinks struct {
	HTML Link `json:"html" mapstructure:"html"`
}

type ActivityType string

const (
	Update   ActivityType = "update"
	Approval ActivityType = "approval"
	Comment  ActivityType = "comment"
)

type PullRequestUpdate struct {
	Date         time.Time    `json:"date" mapstructure:"date"`
	ActivityType ActivityType `json:"activity_type" mapstructure:"activity_type"`
	Author       User         `json:"author" mapstructure:"author"`
}

func (pru PullRequestUpdate) IsNewToUserUUID(uuid string, t time.Time) bool {
	if pru.Date.Before(t) {
		return false
	}

	if pru.Author.UUID == uuid {
		return false
	}

	return true
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

type PullRequestWithRepository struct {
	*PullRequest
	Repository *Repository
}

func (c *Config) AllPullRequests(repo *string) []PullRequestWithRepository {
	// Zip up all PRs
	prs := make([]PullRequestWithRepository, 0)
	if repo != nil {
		r, found := c.Repositories[*repo]

		if !found {
			fmt.Printf("You are not watching the repository %s\n", *repo)
			os.Exit(1)
		}

		for _, pr := range r.PullRequests {
			prs = append(prs, PullRequestWithRepository{Repository: r, PullRequest: pr})
		}
	} else {
		for _, r := range c.Repositories {
			for _, pr := range r.PullRequests {
				prs = append(prs, PullRequestWithRepository{Repository: r, PullRequest: pr})
			}
		}
	}

	// Sort them by unread updates, repository name (asc), then ID (desc)
	sort.Slice(prs, func(i, j int) bool {
		l := prs[i]
		r := prs[j]

		if l.UnreadUpdatesCount == 0 && r.UnreadUpdatesCount > 0 {
			return false
		}

		if l.UnreadUpdatesCount > 0 && r.UnreadUpdatesCount == 0 {
			return true
		}

		switch strings.Compare(l.Repository.Name, r.Repository.Name) {
		case -1:
			return true
		case 1:
			return false
		}

		return l.ID < r.ID
	})

	return prs
}

func (c *Config) ChunkPullRequests(size int) [][]PullRequestWithRepository {
	prs := c.AllPullRequests(nil)

	// Chunk them out
	var chunked [][]PullRequestWithRepository
	for i := 0; i < len(prs); i += size {
		end := i + size

		if end > len(prs) {
			end = len(prs)
		}

		chunked = append(chunked, prs[i:end])
	}

	return chunked
}

func (c *Config) Write() {
	s, err := json.MarshalIndent(c, "", "  ")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(viper.ConfigFileUsed(), s, 0644); err != nil {
		fmt.Println("Unable to write config file")
		fmt.Println(err)
		os.Exit(1)
	}

	return
}

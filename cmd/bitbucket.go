package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/mitchellh/mapstructure"
)

func GetWatchedPullRequests(config Config, repo string) []*PullRequest {
	b := bitbucket.NewBasicAuth(config.Username, config.Token)
	b.Pagelen = 25

	parts := strings.SplitN(repo, "/", 2)
	opts := &bitbucket.PullRequestsOptions{
		Owner:    parts[0],
		RepoSlug: parts[1],
		Query:    fmt.Sprintf(`state="OPEN" AND (author.uuid="%s" OR reviewers.uuid="%s")`, config.UserUUID, config.UserUUID),
		States:   []string{"OPEN"},
	}

	resp, err := b.Repositories.PullRequests.Gets(opts)

	if err != nil {
		fmt.Printf("Error reading pull requests for repository %s\n", repo)
		fmt.Println(err)
		os.Exit(1)
	}

	m := resp.(map[string]interface{})

	var prs []*PullRequest
	err = mapstructure.Decode(m["values"], &prs)

	if err != nil {
		fmt.Printf("Error reading pull requests for repository %s\n", repo)
		fmt.Println(err)
		os.Exit(1)
	}

	return prs
}

func execute(username string, password string, method string, urlStr string, params map[string]string, text string) (map[string]interface{}, error) {
	body := strings.NewReader(text)

	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if text != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	req.SetBasicAuth(username, password)
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var decoded map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&decoded)

	return decoded, err
}

func parseUpdate(data map[string]interface{}) PullRequestUpdate {
	supportedTypes := map[string]bool{"approval": true, "comment": true, "update": true}

	var aData map[string]interface{}
	var aType string
	for t, d := range data {
		found := supportedTypes[t]

		if found {
			aData = d.(map[string]interface{})
			aType = t
		}
	}

	if aData == nil {
		fmt.Println("Unable to parse pull request activity")
		os.Exit(1)
	}

	dateKeys := map[string]string{
		"approval": "date",
		"comment":  "created_on",
		"update":   "date",
	}

	userKeys := map[string]string{
		"approval": "user",
		"comment":  "user",
		"update":   "author",
	}

	dValue := aData[dateKeys[aType]]
	aDate, err := time.Parse(time.RFC3339, dValue.(string))
	if err != nil {
		fmt.Printf("Unable to parse date %s\n", dValue.(string))
		os.Exit(1)
	}

	var aUser User
	uValue := aData[userKeys[aType]]
	err = mapstructure.Decode(uValue, &aUser)
	if err != nil {
		fmt.Println("Unable to parse user data")
		os.Exit(1)
	}

	return PullRequestUpdate{
		Date:         aDate,
		Author:       aUser,
		ActivityType: ActivityType(aType),
	}
}

func GetPullRequestActivity(config Config, repo string, ID int, cursor string) (page []PullRequestUpdate, next string) {
	var resp map[string]interface{}
	var err error
	var urlStr string
	params := make(map[string]string)
	if cursor != "" {
		urlStr = cursor
	} else {
		urlStr = fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/pullrequests/%d/activity", repo, ID)
		params["pagelen"] = "25"
	}

	resp, err = execute(config.Username, config.Token, "GET", urlStr, params, "")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	vs := resp["values"].([]interface{})
	for _, v := range vs {
		page = append(page, parseUpdate(v.(map[string]interface{})))
	}

	n, hasNext := resp["next"]
	if !hasNext {
		next = ""
	} else {
		next = n.(string)
	}

	return page, next
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kenzo0107/backlog"
)

var (
	baseUrl = os.Getenv("BACKLOG_BASE_URL")
	token   = os.Getenv("BACKLOG_TOKEN")
)

func main() {
	c := backlog.New(token, baseUrl)

	user, err := c.GetUserMySelf()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("user ID: %d, Name %s\n", user.ID, *user.Name)

	input := &backlog.GetUserMySelfRecentrlyViewedIssuesOptions{
		Order:  backlog.OrderAsc,
		Offset: backlog.Int(0),
		Count:  backlog.Int(10),
	}
	issues, err := c.GetUserMySelfRecentrlyViewedIssues(input)

	if err != nil {
		log.Fatal(err)
		return
	}

	var key string
	for _, i := range issues {
		fmt.Printf("id: %d, issue key: %s, summary: %s\n",
			*i.Issue.ID, *i.Issue.IssueKey, *i.Issue.Summary)
		key = *i.Issue.IssueKey
	}

	issue, err := c.GetIssue(key)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%s\n", *issue.Assignee.Name)
	fmt.Printf("%s\n", *issue.Category[0].Name)
	fmt.Printf("%s\n", *issue.Description)
	fmt.Printf("%+v\n", issue.CustomFields)

	icInput := &backlog.CreateIssueCommentInput{
		Content: backlog.String("api でコメント追加"),
	}
	ic, err := c.CreateIssueComment(key, icInput)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("%+v\n", *ic.CreatedUser)
}

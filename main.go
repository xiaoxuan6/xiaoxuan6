package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic(errors.New("token invalid"))
	}

	ctx := context.Background()
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	cli := github.NewClient(client)

	repo, err := listStar(ctx, cli, "xiaoxuan6")
	if err != nil {
		panic(err)
	}

	buf := new(strings.Builder)
	for _, val := range repo {
		buf.WriteString(fmt.Sprintf("<a href='https://github.com/%s' target='_blank'>%s</a> - %s<br/>\n", val.FullName, val.FullName, val.CreatedAt.Format("2006-01-02")))
	}

	putContent(buf.String())

}

type Repo struct {
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
}

func listStar(ctx context.Context, client *github.Client, username string) (sta []*Repo, error error) {

	res, _, err := client.Activity.ListStarred(ctx, username, &github.ActivityListStarredOptions{
		Sort:      "created",
		Direction: "",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 5,
		},
	})

	if err != nil {
		return nil, err
	}

	for _, val := range res {
		sta = append(sta, &Repo{
			FullName:  *val.Repository.FullName,
			CreatedAt: val.StarredAt.Time,
		})
	}

	return sta, nil
}

func putContent(s string) {

	by, _ := ioutil.ReadFile("./README.md")
	content := string(by)

	//str := regexp.MustCompile(`<!-- Star starts -->(?s).*<!-- Star ends -->`).FindAllStringSubmatch(strings.TrimSpace(content), -1)
	//fmt.Printf("%v", str[0][0])

	rel := fmt.Sprintf("<!-- Star starts -->\n%s\n<!-- Star ends -->", s)
	str := regexp.MustCompile(`<!-- Star starts -->(?s).*<!-- Star ends -->`).ReplaceAllString(content, rel)
	_ = ioutil.WriteFile("./README.md", []byte(str), os.ModePerm)
}

package main

import (
	"fmt"
	"github.com/geoffreywiseman/gh-actions-usage/client"
	"github.com/geoffreywiseman/gh-actions-usage/format"
	"os"
)

var gh client.Client

func main() {
	fmt.Println("GitHub Actions Usage")
	fmt.Println()

	gh = client.New()
	if len(os.Args) <= 1 {
		tryDisplayCurrentRepo()
	} else {
		tryDisplayAllSpecified(os.Args)
	}
}

func tryDisplayCurrentRepo() {
	repo, err := gh.GetCurrentRepository()
	if repo == nil {
		fmt.Printf("No repository specified; %s\n", err)
		return
	}
	displayRepoUsage(repo)
}

func tryDisplayAllSpecified(names []string) {
	// TODO: Get Repos by Owner
	// TODO: Display Each

	repoName := os.Args[1]
	repo, err := gh.GetRepository(repoName)
	if err != nil {
		panic(err)
	}
	if repo == nil {
		fmt.Printf("Cannot find repo: %s\n", repoName)
		return
	}

	displayRepoUsage(repo)
}

func displayRepoUsage(repo *client.Repository) {
	workflows, err := gh.GetWorkflows(*repo)
	if err != nil {
		panic(err)
	}

	if len(workflows) == 0 {
		fmt.Printf("%s (0 workflows)\n\n", repo.FullName)
		return
	}

	var lines []string = make([]string, 0, len(workflows))
	var repoTotal uint
	for _, flow := range workflows {
		usage, err := gh.GetWorkflowUsage(*repo, flow)
		if err != nil {
			panic(err)
		}
		repoTotal += usage.TotalMs()
		line := fmt.Sprintf("- %s (%s, %s, %s)", flow.Name, flow.Path, flow.State, format.Humanize(usage.TotalMs()))
		lines = append(lines, line)
	}

	fmt.Printf("%s (%d workflows; %s): \n", repo.FullName, len(workflows), format.Humanize(repoTotal))
	for _, line := range lines {
		fmt.Println(line)
	}
}

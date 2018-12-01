package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
	"time"
)

func main() {

	var owner string
	flag.StringVar(&owner, "owner", "", "The github Owner")

	var repos string
	flag.StringVar(&repos, "repo", "", "The github Repository")

	flag.Parse()

	reposArray := strings.Split(repos, ",")

	var user string
	fmt.Println("Github user: ")
	fmt.Scanln(&user)

	fmt.Println("Github password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	if err != nil {
		fmt.Println("Password can not be empty")
		os.Exit(3)
	}

	fmt.Println()

	fmt.Println("Starting analyse")

	var dateStart string
	fmt.Println("Date start: ")
	fmt.Scanln(&dateStart)
	dateStart += "T00:00:00Z"
	dateStartParsed, err := time.Parse(time.RFC3339, dateStart)

	if err != nil {
		fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
		os.Exit(3)
	}

	var dateEnd string
	fmt.Println("Date end: ")
	fmt.Scanln(&dateEnd)
	dateEnd += "T23:59:59Z"
	dateEndParsed, err := time.Parse(time.RFC3339, dateEnd)

	if err != nil {
		fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
		os.Exit(3)
	}

	var githubUser string
	fmt.Println("Github user to calculate: ")
	fmt.Scanln(&githubUser)

	pullRequestsCount := 0
	pullReviewsCount := 0

	fmt.Println()
	fmt.Println("Calculating ")

	for _, repo := range reposArray {
		fmt.Println()
		fmt.Println(repo)
		pullRequests := getPullRequests("pullRequests", user, password, owner, repo, 1)
		for f := 2; len(pullRequests) > 0; f++ {
			for _, pullRequest := range pullRequests {
				update, err := time.Parse(time.RFC3339, pullRequest.Created_at)
				if err != nil {
					fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
					os.Exit(3)
				}
				if dateStartParsed.After(update) || dateEndParsed.Before(update) {
					continue
				}
				pullRequestsCount++

				reviews := getPullReviews("pullReviews", user, password, owner, repo, pullRequest.Number)

				for _, review := range reviews {
					if review.User.Login == githubUser {
						pullReviewsCount++
					}
				}
				fmt.Print(".")
			}
			pullRequests = getPullRequests("pullRequests", user, password, owner, repo, f)
		}
	}

	fmt.Println()
	fmt.Println()
	reviewsPercentual := (100 / pullRequestsCount) * pullReviewsCount
	fmt.Printf("User %s did %d reviews about %d pull requests (%d%%) on interval between %s and %s", githubUser, pullReviewsCount, pullRequestsCount, reviewsPercentual, dateStart, dateEnd)

	fmt.Println()
	fmt.Println("Researchs ends for", owner, reposArray)
}

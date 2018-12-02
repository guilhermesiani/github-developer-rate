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

	var dateStart string
	flag.StringVar(&dateStart, "dateStart", "", "Date start of period")

	var dateEnd string
	flag.StringVar(&dateEnd, "dateEnd", "", "Date end of period")

	var githubUser string
	flag.StringVar(&githubUser, "githubUser", "", "Github User to compare")

	flag.Parse()

	if "" == owner {
		fmt.Println("You must set `owner` param")
		os.Exit(3)
	}

	if "" == repos {
		fmt.Println("You must set `repo` param")
		os.Exit(3)
	}

	if "" == dateStart {
		fmt.Println("You must set `dateStart` param")
		os.Exit(3)
	}

	if "" == dateEnd {
		fmt.Println("You must set `dateEnd` param")
		os.Exit(3)
	}

	if "" == githubUser {
		fmt.Println("You must set `githubUser` param")
		os.Exit(3)
	}

	dateStart += "T00:00:00Z"
	dateStartParsed, err := time.Parse(time.RFC3339, dateStart)
	dateEnd += "T23:59:59Z"
	dateEndParsed, err := time.Parse(time.RFC3339, dateEnd)

	reposArray := strings.Split(repos, ",")

	fmt.Println()
	fmt.Println("You must log in to proceed")

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

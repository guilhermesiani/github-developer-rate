package main

import (
    "fmt"
    "flag"
    "time"
    "syscall"
    "os"
    "golang.org/x/crypto/ssh/terminal"
)

func main() {

    var owner string
    flag.StringVar(&owner, "owner", "", "The github Owner")

    var repo string
    flag.StringVar(&repo, "repo", "", "The github Repository")

    flag.Parse()

    var user string;
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

    var dateStart string;
    fmt.Println("Date start: ")
    fmt.Scanln(&dateStart)
    dateStart += "T00:00:00Z"
    dateStartParsed, err := time.Parse(time.RFC3339, dateStart)

    if err != nil {
        fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
        os.Exit(3)
    }

    var dateEnd string;
    fmt.Println("Date end: ")
    fmt.Scanln(&dateEnd)
    dateEnd += "T23:59:59Z"
    dateEndParsed, err := time.Parse(time.RFC3339, dateEnd)

    if err != nil {
        fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
        os.Exit(3)
    }

    var githubUser string;
    fmt.Println("Github user to calculate: ")
    fmt.Scanln(&githubUser)

    pullRequestsCount := 0
    pullReviewsCount := 0

    fmt.Println()
    fmt.Println("Calculating ")

    keys := getPullRequests("pullRequests", user, password, owner, repo, 1)
    for f := 2; len(keys) > 0; f++ {
        for i := 0; i < len(keys); i++ {
            update, err := time.Parse(time.RFC3339, keys[i].Created_at)
            if err != nil {
                fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
                os.Exit(3)
            }
            if dateStartParsed.After(update) || dateEndParsed.Before(update) {
                continue
            }
            pullRequestsCount++

            reviewKeys := getPullReviews("pullReviews", user, password, owner, repo, keys[i].Number)

            for j := 0; j < len(reviewKeys); j++ {
                if reviewKeys[j].User.Login == githubUser {
                    pullReviewsCount++
                }
            }
            fmt.Print(".")
        }
        keys = getPullRequests("pullRequests", user, password, owner, repo, f)
    }

    fmt.Println()
    fmt.Println()
    fmt.Printf("User %s did %d reviews about %d pull requests on interval between %s and %s", githubUser, pullReviewsCount, pullRequestsCount, dateStart, dateEnd)

    fmt.Println()
    fmt.Println("Researchs ends for", owner, repo)
}
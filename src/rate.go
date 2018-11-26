package main

import (
    "fmt"
    "flag"
    "time"
    "syscall"
    "os"
    "encoding/json"
    "golang.org/x/crypto/ssh/terminal"
    "net/http"
    "io/ioutil"
)

type PullRequest struct {
    Id int
    Url string
    Number int
    Title string
    Created_at string
    Updated_at string
}

type PullRequests struct {
    Collection []PullRequest
}

type PullReview struct {
    Id int
    User  struct {
        Id int
        Login string
    }
}

type PullReviews struct {
    Collection []PullReview
}

func getPullRequests(dataType string, user string, password string, owner string, repo string) []PullRequest {

    endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)
    request, err := http.NewRequest("GET", endpoint, nil)
    request.SetBasicAuth(user, password)
    client := &http.Client{}
    response, err := client.Do(request)

    if err != nil {
        fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
        os.Exit(3)
    }

    data, _ := ioutil.ReadAll(response.Body)
    jsonResponse := string(data)

    keysBody := []byte(jsonResponse)
    keys := make([]PullRequest, 0)
    json.Unmarshal(keysBody, &keys)

    return keys
}

func getPullReviews(dataType string, user string, password string, owner string, repo string, pullRequestNumber int) []PullReview {

    endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/reviews", owner, repo, pullRequestNumber)
    request, err := http.NewRequest("GET", endpoint, nil)
    request.SetBasicAuth(user, password)
    client := &http.Client{}
    response, err := client.Do(request)

    if err != nil {
        fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
        os.Exit(3)
    }

    data, _ := ioutil.ReadAll(response.Body)
    jsonResponse := string(data)

    keysBody := []byte(jsonResponse)
    keys := make([]PullReview, 0)
    json.Unmarshal(keysBody, &keys)

    return keys
}

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

    keys := getPullRequests("pullRequests", user, password, owner, repo)

    var dateStart string;
    fmt.Println("Date start: ")
    fmt.Scanln(&dateStart)
    dateStartParsed, err := time.Parse(time.RFC3339, dateStart)

    if err != nil {
        fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
        os.Exit(3)
    }

    var dateEnd string;
    fmt.Println("Date end: ")
    fmt.Scanln(&dateEnd)
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

        reviewKeys := getPullReviews("pullRequests", user, password, owner, repo, keys[i].Number)

        for j := 0; j < len(reviewKeys); j++ {
            if reviewKeys[j].User.Login == githubUser {
                pullReviewsCount++
            }
        }
        fmt.Print(".")
    }

    fmt.Println()
    fmt.Println()
    fmt.Printf("User %s did %d reviews about %d pull requests on interval between %s and %s", githubUser, pullReviewsCount, pullRequestsCount, dateStart, dateEnd)

    fmt.Println()
    fmt.Println("Researchs ends for", owner, repo)
}

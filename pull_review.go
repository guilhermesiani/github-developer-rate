package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"
)

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

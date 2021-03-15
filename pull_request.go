package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type PullRequest struct {
	Id         int
	Url        string
	Number     int
	Title      string
	Created_at string
	Updated_at string
}

type PullRequests struct {
	Collection []PullRequest
}

func getPullRequests(dataType string, user string, password string, owner string, repo string, page int) []PullRequest {

	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?page=%d", owner, repo, page)
	request, err := http.NewRequest("GET", endpoint, nil)
	//request.SetBasicAuth(user, password)
	request.Header.Add("Authorization", "token "+password)
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

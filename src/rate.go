package main

import (
    "fmt"
    "flag"
    "net/http"
    "io/ioutil"
)

func main() {

    var owner string
    flag.StringVar(&owner, "owner", "", "The github Owner")

    var repo string
    flag.StringVar(&repo, "repo", "", "The github Repository")

    flag.Parse()

    fmt.Println("Starting analyse")

    endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)
    fmt.Println(endpoint)
    response, err := http.Get(endpoint)

    if err != nil {
        fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
    }

    fmt.Println("Researchs ends", owner, repo)
}

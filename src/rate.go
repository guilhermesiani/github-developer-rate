package main

import (
    "fmt"
    "flag"
    "syscall"
    "os"

    "golang.org/x/crypto/ssh/terminal"
    "net/http"
    "io/ioutil"
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

    if err != nil {
        fmt.Println("Password can not be empty")
        os.Exit(3)
    }

    fmt.Println()

    fmt.Println("Starting analyse")

    endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", owner, repo)
    request, err := http.NewRequest("GET", endpoint, nil)
    request.SetBasicAuth(user, string(bytePassword))
    client := &http.Client{}
    response, err := client.Do(request)

    if err != nil {
        fmt.Println(fmt.Sprintf("The http request failed with error %s\n", err))
    } else {
        data, _ := ioutil.ReadAll(response.Body)
        fmt.Println(string(data))
    }

    fmt.Println("Researchs ends", owner, repo)
}

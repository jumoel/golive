package main

import (
  "net/http"
  "log"
  "io/ioutil"
  "encoding/json"
  "regexp"
  "fmt"
  "os/exec"
)

// A map of repos => (a map of branches => (a list of actions))
// Todo: Idiomatic way of doing this?
type Config map[string]map[string][]string


type HookMsg struct {
  CanonicalUrl string `json:"canon_url"`

  Commits []struct {
    Branch string
  }

  Repository struct {
    AbsoluteUrl string `json:"absolute_url"`
  }
}

type Commit struct {
  Repository string
  Branch string
}

func main() {
  config_raw, _ := ioutil.ReadFile("golive.json")
  var config Config;
  json.Unmarshal(config_raw, &config)

  msgs := make(chan HookMsg, 100)
  commits := make(chan Commit, 100)
  jobs := make(chan string, 100)

  go hookWrangler(msgs, commits)
  go commitWrangler(commits, jobs, config)
  go jobRunner(jobs)

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    payload := r.FormValue("payload")

    var hookmsg HookMsg
    if err := json.Unmarshal([]byte(payload), &hookmsg); err != nil {
      http.Error(w, "Could not decode json", 500)
			log.Print(err)
			return
    }

    msgs <- hookmsg
  })

  log.Fatal(http.ListenAndServe(":8080", nil))
}

func hookWrangler(msgs <-chan HookMsg, results chan<- Commit) {
  reg := regexp.MustCompile("https?://")

  for msg := range msgs {
    baserepo := reg.ReplaceAllString(msg.CanonicalUrl, "")
    repository := baserepo + msg.Repository.AbsoluteUrl

    for _, commit := range msg.Commits {
      results <- Commit{ repository, commit.Branch}
    }
  }
}

func commitWrangler(commits <-chan Commit, results chan<- string, config Config) {
  for commit := range commits {
    if commit.Branch == "" || commit.Repository == "" {
      continue
    }

    if branches, ok := config[commit.Repository]; ok {
      if actions, ok := branches[commit.Branch]; ok {
        for _, action := range actions {
          results <- action
        }
      }
    }
  }
}

func jobRunner(jobs <-chan string) {
  for job := range jobs {
    command := exec.Command("bash", "-c", job)
    go command.Run()
  }
}

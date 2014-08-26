package main

import (
  "net/http"
  "log"
  "io/ioutil"
  "encoding/json"
  "os/exec"
  "flag"
  "strconv"
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

type Job struct {
  Commit Commit
  Action string
}

var listenPort = flag.Int("port", 8080, "portnumber to listen on")
var configFile = flag.String("config", "golive.json", "the configfile to read")
var verbose = flag.Bool("v", false, "print more output")

func main() {
	flag.Parse()

  config_raw, err := ioutil.ReadFile(*configFile)
  if err != nil {
    log.Fatal(err)
  }

  var config Config;
  json.Unmarshal(config_raw, &config)

  msgs := make(chan HookMsg, 100)
  commits := make(chan Commit, 100)
  jobs := make(chan Job, 100)

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

    if *verbose {
      log.Print("Received commit: ", hookmsg)
    }

    msgs <- hookmsg
  })

  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(*listenPort), nil))
}

func hookWrangler(msgs <-chan HookMsg, results chan<- Commit) {
  for msg := range msgs {
    repository := msg.CanonicalUrl + msg.Repository.AbsoluteUrl

    for _, commit := range msg.Commits {
      results <- Commit{ repository, commit.Branch}
    }
  }
}

func commitWrangler(commits <-chan Commit, results chan<- Job, config Config) {
  for commit := range commits {
    if commit.Branch == "" || commit.Repository == "" {
      continue
    }

    if branches, ok := config[commit.Repository]; ok {
      if actions, ok := branches[commit.Branch]; ok {
        for _, action := range actions {
          results <- Job{commit, action}
        }
      }
    }
  }
}

func jobRunner(jobs <-chan Job) {
  for job := range jobs {
    if *verbose {
      log.Print("Running job: ", job)
    }

    command := exec.Command("bash", "-c", job.Action)
    go command.Run()
  }
}

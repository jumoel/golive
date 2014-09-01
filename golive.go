package main

import (
  "net/http"
  "log"
  "io/ioutil"
  "encoding/json"
  "os/exec"
  "flag"
  "strconv"
  "crypto/md5"
  "text/template"
  "bytes"
  "gopkg.in/fsnotify.v1"
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

var jobTemplates = make(map[[16]byte]template.Template)
var config Config

func main() {
	flag.Parse()

  parseConfig(*configFile)
  watchConfig(*configFile)

  msgs := make(chan HookMsg, 100)
  commits := make(chan Commit, 100)
  jobs := make(chan Job, 100)
  actions := make(chan string, 100)

  go hookWrangler(msgs, commits)
  go commitWrangler(commits, jobs, config)
  go jobWrangler(jobs, actions)
  go actionRunner(actions)

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    payload := r.FormValue("payload")

    var hookmsg HookMsg
    if err := json.Unmarshal([]byte(payload), &hookmsg); err != nil {
      http.Error(w, "Could not decode json", 500)
      log.Fatal(err)
    }

    if *verbose {
      log.Print("Received commit: ", hookmsg)
    }

    msgs <- hookmsg
  })

  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(*listenPort), nil))
}

func watchConfig(configFile string) {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    log.Fatal(err)
  }

  log.Print("Watching config file: ", configFile)

  defer watcher.Close()

	done := make(chan bool)

  go func() {
    for {
      select {
        case event := <-watcher.Events:
          if event.Name == "" {
            continue
          }

          if *verbose {
            log.Println(event.String())
          }

          if event.Op & fsnotify.Write == fsnotify.Write ||
             event.Op & fsnotify.Chmod == fsnotify.Chmod {
            log.Print("Relading config file ", configFile)

            parseConfig(configFile)
          }

        case err := <-watcher.Errors:
          if err == nil {
            continue
          }

          done <- true
          log.Fatal("Error when watching config file: ", err)
      }
    }
  }()

  err = watcher.Add(configFile)
  if err != nil {
    log.Fatal(err)
  }

  <- done
}

func parseConfig(configFile string) {
    config_raw, err := ioutil.ReadFile(configFile)
    if err != nil {
      log.Fatal(err)
    }

    var newConfig Config
    json.Unmarshal(config_raw, &newConfig)

    config = newConfig

    if *verbose {
      log.Print("Loaded config: ", config)
    }
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
          results <- Job{commit, string(action)}
        }
      }
    }
  }
}

func jobWrangler(jobs <-chan Job, actions chan<- string) {
  for job := range jobs {
    if *verbose {
      log.Print("Running job: ", job)
    }

    hash := md5.Sum([]byte(job.Action))

    if _, ok := jobTemplates[hash]; !ok {
      t, err := template.New(string(hash[:])).Parse(job.Action)
      if err != nil {
        log.Fatal("Could not compile template: ", job.Action, " - ", err)
      }

      jobTemplates[hash] = *t
    }

    t := jobTemplates[hash]

    var buff bytes.Buffer
    (&t).Execute(&buff, job.Commit)
    s := buff.String()

    actions <- s
  }
}

func actionRunner(actions <-chan string) {
  for action := range actions {
    command := exec.Command("bash", "-c", action)

    command.Run()
  }
}

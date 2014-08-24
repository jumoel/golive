package main

import (
  "fmt"
  "net/http"
  "log"
  "io/ioutil"
  "encoding/json"
)

// A map of repos => (a map of branches => (a list of actions))
// Todo: Idiomatic way of doing this?
type Config map[string]map[string][]string


type HookMsg struct {
  CanonicalUrl string "canon_url"

  Commits []struct {
    Branch string
  }

  Repository struct {
    AbsoluteUrl string "absolute_url"
  }
}

type Commit struct {
  Repository string
  Branch string
}
/*
{
    "canon_url": "https://bitbucket.org",
    "commits": [
        {
            "author": "marcus",
            "branch": "master",
            "files": [
                {
                    "file": "somefile.py",
                    "type": "modified"
                }
            ],
            "message": "Added some more things to somefile.py\n",
            "node": "620ade18607a",
            "parents": [
                "702c70160afc"
            ],
            "raw_author": "Marcus Bertrand <marcus@somedomain.com>",
            "raw_node": "620ade18607ac42d872b568bb92acaa9a28620e9",
            "revision": null,
            "size": -1,
            "timestamp": "2012-05-30 05:58:56",
            "utctimestamp": "2012-05-30 03:58:56+00:00"
        }
    ],
    "repository": {
        "absolute_url": "/marcus/project-x/",
        "fork": false,
        "is_private": true,
        "name": "Project X",
        "owner": "marcus",
        "scm": "git",
        "slug": "project-x",
        "website": "https://atlassian.com/"
    },
    "user": "marcus"
}
*/

func main() {

  /*
  - Modtag POST fra Bitbucket, unmarshal JSON til HookMsg

  - Ryd op i HookMsg og konstruer (evt. flere) Commit objekter

  - For hver Commit, hvis config[Commit.Repository][Commit.Branch] findes,
    lav map[string]int af actions -> antal gange den er blevet bedt om at blive
    kørt

    For hver action, kør den
  */
  config_raw, _ := ioutil.ReadFile("golive.json")
  var config Config;
  json.Unmarshal(config_raw, &config)

  http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	   fmt.Fprintf(w, "Hello, %q, %v, %q", r.URL.Path, config, config_raw)
  })

  log.Fatal(http.ListenAndServe(":8080", nil))
}

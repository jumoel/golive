# Golive

A minimal continuous *someting* system. The code is very much in flux, so don't
expect things to work the same *N* commits from now!

It (by default) listens on port 8080 for a POST webhook request from hosted git
services (only Bitbucket at the moment).

If the pushed repository and branch is present in the config file `golive.json`,
all the commands listed run, as-is. Wildcards are also available as fallback matches
for the branch names. They will be run if no other branch matches. They are denoted with `*`
as the branch name.

If `golive` is run with `-bootstrap`, it will run all non-wildcard actions a single time on startup.
This can be used to ensure that all `golive`-managed sites are present when provisioning new servers.

Hopelessly insecure, but very minimal and fully configurable with text files.
In fact, there is only `golive.json`. :-)

See `golive.example.json` for an example configuration.

`golive` automatically re-reads the configuration file if it is changed.

## Usage

    $ golive -help
    Usage of ./golive:
      -config="golive.json": the configfile to read
      -port=8080: portnumber to listen on
      -v=false: print more output
      -bootstrap=false: Run all non-wildcard actions on startup

## Example of use

With `golive.json` containing:

    {
      "https://bitbucket.org/username/repo/": {
        "master": [
          "echo 'Commit from: {{.Repository}}{{.Branch}}' >> test.txt"
        ],
        "*": [
          "echo 'Commit from: {{.Repository}}{{.Branch}} WILDCARD' >> test.txt"
        ]
      }
    }


 1. Fill in Bitbucket repository URL and branch in `golive.example.json`
 2. Run `golive --config=golive.json --port=8080 -v` at `servername.tld`.
 3. Set up POST hook in a Bitbucket repository to point to `servername.tld:8080`
 4. Commit a change to the `master` branch repository and push it to Bitbucket
 5. Commit a change to the `test` branch repository and push it to Bitbucket
 6. Watch, as `test.txt` in the folder you ran `golive` from contains
    `Commit from: <repository>/master` as well as `Commit from: <repository>/test WILDCARD`
    


### Real world example

I've used `golive` to run [Ansible](http://www.ansible.com) playbooks that
deploy websites to servers when their repositories have been updated.

`golive` could easily run a test suite or other shenanigans as well.

## Dependencies

`bash`, so `bash -c $job` can run.

## RFC

This is my first Golang project - I'd very much like to know how it can be
improved to be more idiomatic, so I welcome all and any critique/comments.

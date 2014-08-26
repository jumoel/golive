# Golive

A minimal continious *someting* system.

It (by default) listens on port 8080 for a POST webhook request from hosted git
services (only Bitbucket at the moment).

If the pushed repository and branch is present in the config file `golive.json`,
all the commands listed run, as-is.

Hopelessly insecure, but very minimal and fully configurable with text files.
In fact, there is only `golive.json`. :-)

See `golive.example.json` for an example configuration. If you ran `golive` and
received a POST webhook request for my test repo for this project, a `test.txt`
file would be created in the same directory with the contents `foo`.

## Usage

    $ golive -help
    Usage of golive:
      -port=8080: portnumber to listen on

## RFC

This is my first Golang project - I'd very much like to know how it can be
improved to be more idiomatic, so I welcome all and any critique/comments.

## Example of use

We run [Ansible](http://www.ansible.com) playbooks that deploy the pushed
repository/branch to our servers.

## Dependencies

`bash`, so `bash -c $job` can run.

# Golive

A minimal continuous *someting* system. The code is very much in flux, so don't
expect things to work the same *N* commits from now!

It (by default) listens on port 8080 for a POST webhook request from hosted git
services (only Bitbucket at the moment).

If the pushed repository and branch is present in the config file `golive.json`,
all the commands listed run, as-is.

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

## Example of use

With the supplied `golive.example.json`:

 1. Fill in Bitbucket repository URL and branch in `golive.example.json`
 2. Run `golive --config=golive.example.json -v` at `servername.tld`.
 3. Set up POST hook in a Bitbucket repository to point to `servername.tld`
 4. Commit a change to the repository and push it to Bitbucket
 5. Watch, as `test.txt` in the folder you ran `golive` from contains
    `Commit from: <yourrepository>/<yourbranch>`

### Real world example

I've used `golive` to run [Ansible](http://www.ansible.com) playbooks that
deploy websites to servers when their repositories have been updated.

`golive` could easily run a test suite or other shenanigans as well.

## Dependencies

`bash`, so `bash -c $job` can run.

## RFC

This is my first Golang project - I'd very much like to know how it can be
improved to be more idiomatic, so I welcome all and any critique/comments.

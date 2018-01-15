# dcr_seed_recovery
Allows decred wallet seed recovery when a couple of words are incorrect
## Updating

#### Windows

Install a newer MSI

#### Linux/BSD/MacOSX/POSIX - Build from Source

## Requirements

[Go](http://golang.org) 1.7 or newer.

## Updating

#### Linux/BSD/MacOSX/POSIX - Build from Source

- **Dep**

  Dep is used to manage project dependencies and provide reproducible builds.
  To install:

  `go get -u github.com/golang/dep/cmd/dep`

Unfortunately, the use of `dep` prevents a handy tool such as `go get` from
automatically downloading, building, and installing the source in a single
command.  Instead, the latest project and dependency sources must be first
obtained manually with `git` and `dep`, and then `go` is used to build and
install the project.

**Getting the source**:

For a first time installation, the project and dependency sources can be
obtained manually with `git` and `dep` (create directories as needed):

```
git clone https://github.com/annawaw/dcr_seed_recovery $GOPATH/src/github.com/annawaw/dcr_seed_recovery
cd $GOPATH/src/github.com/annawaw/dcr_seed_recovery
dep ensure
go build
```

## Examples

** One invalid word **
```
$ ./dcr_seed_recovery -backupSeed  "crusade graduate swelter maritime brickyard Atlantic slingshot aggregate brickyard handiwork spellbind unicorn select yesteryear sugar Chicago Mohawk belowground bluebird adviser tumor torpedo bison headwaters deckhand Jamaica sawdust yesteryear chatter article tapeworm unicorn tapeworm" -addr DsaAWrhMCr4ASMSifnMAbR99LigfRUg18uG
         7782 / 8448 [=======================================================================>------]  92 %
Found wallet mnemonic:
crusade graduate swelter maritime brickyard Atlantic slingshot aggregate brickyard handiwork spellbind unicorn select yesteryear sugar Chicago Mohawk belowground bluebird adviser tumor torpedo bison headwaters deckhand Jamaica sawdust yesteryear chatter article framework unicorn tapeworm
```

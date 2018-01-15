# dcr_seed_recovery
Allows recovery of a decred wallet seed in a scenario when only a couple of words are incorrect

#### Features ####
dcr_seed_recovery will attempt to recover the valid decred wallet seed, assuming that following are provided:
1. An approximate 33 word backup seed. Preferably with only up to 1 or 2 words incorrect. 

2. By default only single invalid word is assumed. It is possible to try increase the number of words tired, but the time required grows exponentially. 2 words are ok and should take less than hour, 3 words may require days.

3. A sample public address from the wallet that needs revocery. Ideally the earlier generated the better.
   The tool generates all possible valid wallet seeds and checks if either external or internal address matches the address provided.
   By default it checks only first 128 addresses in the wallet. If it is known that the provided public address was generated earlier or later, it can be configured to limit the number of addresses to a smaller or larger value.

4. By default checks only first (default, with id 0) account in the wallet. It's possible to increase number of accounts checked.

5. By default it only checks BIP44 wallet accounts, with coin type 42. It's possible to also check legacy accounts, with coin type 20. 

#### Linux/BSD/MacOSX/POSIX - Build from Source

## Requirements

[Go](http://golang.org) 1.7 or newer.

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


## Usage ##
```
$ ./dcr_seed_recovery --help
Usage of ./dcr_seed_recovery:
  -accountLimit int
    	max number of wallet accounts to check (optional) (default 1)
  -addr string
    	a public address that belongs to the wallet (required)
  -addrLimit int
    	max number of addresses in a wallet to check (optional) (default 128)
  -allowLegacy
    	allow legacy coin type (optional)
  -backupSeed string
    	backup seed to recover/fix (required)
  -depth int
    	max number of invalid words, recommended <=2 (optional) (default 1)
```

## Examples

#### One invalid word
```
$ ./dcr_seed_recovery -backupSeed  "crusade graduate swelter maritime brickyard Atlantic slingshot aggregate brickyard handiwork spellbind unicorn select yesteryear sugar Chicago Mohawk belowground bluebird adviser tumor torpedo bison headwaters deckhand Jamaica sawdust yesteryear chatter article tapeworm unicorn tapeworm" -addr DsaAWrhMCr4ASMSifnMAbR99LigfRUg18uG
         7782 / 8448 [=======================================================================>------]  92 %
Found wallet mnemonic:
crusade graduate swelter maritime brickyard Atlantic slingshot aggregate brickyard handiwork spellbind unicorn select yesteryear sugar Chicago Mohawk belowground bluebird adviser tumor torpedo bison headwaters deckhand Jamaica sawdust yesteryear chatter article framework unicorn tapeworm
```

#### Two invalid words
```
$ ././dcr_seed_recovery -backupSeed  "crusade graduate swelter maritime brickyard Atlantic slingshot aggregate brickyard handiwork spellbind unicorn select yesteryear sugar Chicago Mohawk belowground bluebird adviser tumor torpedo bison headwaters deckhand Jamaica sawdust yesteryear tapeworm article tapeworm unicorn tapeworm" -addr DsaAWrhMCr4ASMSifnMAbR99LigfRUg18uG -depth 2
 33997383 / 34603008 [============================================================================>-]  98 %
Found wallet mnemonic:
crusade graduate swelter maritime brickyard Atlantic slingshot aggregate brickyard handiwork spellbind unicorn select yesteryear sugar Chicago Mohawk belowground bluebird adviser tumor torpedo bison headwaters deckhand Jamaica sawdust yesteryear chatter article framework unicorn tapeworm
```

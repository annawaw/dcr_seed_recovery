package main

import (
	"fmt"
	"crypto/sha256"
	"flag"
	"strings"
	"github.com/decred/dcrwallet/wallet/udb"
    "github.com/decred/dcrwallet/walletseed"
    "github.com/decred/dcrd/chaincfg"
    "github.com/decred/dcrd/hdkeychain"
	"github.com/decred/dcrwallet/pgpwordlist"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type SearchOptions struct {
	addr string
	addrLimit int
	accountLimit int
	allowLegacy bool
}

func main() {
	backupSeedPtr := flag.String("backupSeed", "", "backup seed to recover/fix (required)")
	addrPtr := flag.String("addr", "", "a public address that belongs to the wallet (required)")
	depthPtr := flag.Int("depth", 1, "max number of invalid words, recommended <=2 (optional)")
	addrLimitPtr := flag.Int("addrLimit", 128, "max number of addresses in a wallet to check (optional)")
	accountLimitPtr := flag.Int("accountLimit", 1, "max number of wallet accounts to check (optional)")
	allowLegacyPtr := flag.Bool("allowLegacy", false, "allow legacy coin type (optional)")

	flag.Parse()

	if (*backupSeedPtr == "") {
		fmt.Println("Please specify backupSeed")
		return
	}

	if (*addrPtr == "") {
		fmt.Println("Please specify addr")
		return
	}

	if (*depthPtr < 0) {
		fmt.Println("depth must be non negative number")
		return
	}

    if (*addrLimitPtr <= 0) {
    	fmt.Println("addrLimit must be positive number")
		return
    }

    if (*accountLimitPtr <= 0) {
    	fmt.Println("accountLimit must be positive number")
		return
    }


	words := strings.Split(strings.TrimSpace(*backupSeedPtr), " ")
    decoded, err := pgpwordlist.DecodeMnemonics(words)
	if err != nil {
		fmt.Println(err)
		return
	}

	opt := SearchOptions{
		addr: *addrPtr,
		addrLimit: *addrLimitPtr,
		accountLimit: *accountLimitPtr,
		allowLegacy: *allowLegacyPtr,
	}

    progress := mpb.New()

	totalIteration := countIterations(decoded, *depthPtr)
	bar := createProgressBar(progress, totalIteration)

	seed := findSeed(decoded, 0, *depthPtr, &opt, bar)
	
	//bar.SetTotal(bar.Current(), true)
	bar.Complete()
	
	progress.Stop()

	if (seed != nil) {
		mnemonic := walletseed.EncodeMnemonic(seed)
		
		fmt.Printf("Found wallet mnemonic:\n")
		fmt.Printf("%s\n", mnemonic)
	} else {
		fmt.Printf("Wallet mnemonic not found\n")
	}
}

func countIterations(data[] byte, depth int) int64 {
	 return ncr(len(data), depth)*pow(256, depth)
}

func createProgressBar(progress *mpb.Progress, total int64) *mpb.Bar {
	return progress.AddBar(total,
		mpb.PrependDecorators(
			decor.CountersNoUnit("%d / %d", 20, 0),
		),
		// Appending decorators
		mpb.AppendDecorators(
			// Percentage decorator with minWidth and no extra config
			decor.Percentage(5, 0),
		),
	)
}

func findSeed(data []byte, offset int, remainingDepth int, opt *SearchOptions, bar *mpb.Bar) []byte {
	if remainingDepth > 0 {
		for idx := offset; idx < len(data); idx++ {
			b := data[idx]
			for x := 0; x < 256; x++ {
				data[idx] = byte(x)

				result := findSeed(data, idx + 1, remainingDepth - 1, opt, bar)
				if (result != nil) {
					data[idx] = b	
					return result 	
				}
			}
			
			data[idx] = b
		}

		return nil
	} else {
		if validateSeed(data) {
			seed := data[:len(data)-1]

			if checkWallet(seed, opt) {
				var result = make([]byte, len(data) - 1)
				copy(result, data)

				return result
			}
		}


		bar.Increment()
		return nil
	}
}



func validateSeed(input []byte) bool {
	return checksumByte(input[:len(input)-1]) == input[len(input)-1]
}

func checkWallet(seed []byte, opt *SearchOptions) bool {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		fmt.Println(err)
		return false
	}
	
	legacyCoinType, slip0044CoinType := udb.CoinTypes(&chaincfg.MainNetParams)

	slip0044CoinTypeKey, err := deriveCoinTypeKey(masterKey, slip0044CoinType)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if checkCoinTypeKey(slip0044CoinTypeKey, opt) {
		return true
	}

    if (opt.allowLegacy) {
		legacyCoinTypeKey, err := deriveCoinTypeKey(masterKey, legacyCoinType)
		if err != nil {
			fmt.Println(err)
			return false
		}
		
		if checkCoinTypeKey(legacyCoinTypeKey, opt) {
			return true
		}
	}

	return false
}

func checkCoinTypeKey(coinTypeKey *hdkeychain.ExtendedKey, opt *SearchOptions) bool {
	for account := 0; account < opt.accountLimit; account++ {
		accountKey, err := deriveAccountKey(coinTypeKey, uint32(account))
		if err != nil {
			fmt.Println(err)
			return false
		}

		if checkAccountKey(accountKey, opt) {
			return true
		}
	}

	return false
}

func checkAccountKey(accountKey *hdkeychain.ExtendedKey, opt *SearchOptions) bool {
	for kind := 0; kind < 2; kind++ {
		kindKey, err := accountKey.Child(uint32(kind))
		if err != nil {
			fmt.Println(err)
			return false
		}

		if checkAddressGenKey(kindKey, opt) {
			return true
		}
	}

	return false
}

func checkAddressGenKey(addressGenKey *hdkeychain.ExtendedKey, opt *SearchOptions) bool {
	for idx := 0; idx < opt.addrLimit; idx++ {
		addrKey, err := addressGenKey.Child(uint32(idx))
		if err != nil {
			fmt.Println(err)
			return false
		}

		addr, err := addrKey.Address(&chaincfg.MainNetParams)
		if err != nil {
			fmt.Println(err)
			return false
		}

		if addr.String() == opt.addr {
			return true
		}
	}

	return false
}

// checksumByte returns the checksum byte used at the end of the seed mnemonic
// encoding.  The "checksum" is the first byte of the double SHA256.
func checksumByte(data []byte) byte {
	intermediateHash := sha256.Sum256(data)
	return sha256.Sum256(intermediateHash[:])[0]
}

// deriveCoinTypeKey derives the cointype key which can be used to derive the
// extended key for an account according to the hierarchy described by BIP0044
// given the coin type key.
//
// In particular this is the hierarchical deterministic extended key path:
// m/44'/<coin type>'
func deriveCoinTypeKey(masterNode *hdkeychain.ExtendedKey,
	coinType uint32) (*hdkeychain.ExtendedKey, error) {
	
	// The hierarchy described by BIP0043 is:
	//  m/<purpose>'/*
	// This is further extended by BIP0044 to:
	//  m/44'/<coin type>'/<account>'/<branch>/<address index>
	//
	// The branch is 0 for external addresses and 1 for internal addresses.

	// Derive the purpose key as a child of the master node.
	purpose, err := masterNode.Child(44 + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}

	// Derive the coin type key as a child of the purpose key.
	coinTypeKey, err := purpose.Child(coinType + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}

	return coinTypeKey, nil
}

// deriveAccountKey derives the extended key for an account according to the
// hierarchy described by BIP0044 given the master node.
//
// In particular this is the hierarchical deterministic extended key path:
//   m/44'/<coin type>'/<account>'
func deriveAccountKey(coinTypeKey *hdkeychain.ExtendedKey,
	account uint32) (*hdkeychain.ExtendedKey, error) {
	
	// Derive the account key as a child of the coin type key.
	return coinTypeKey.Child(account + hdkeychain.HardenedKeyStart)
}

func ncr(n int, k int) int64 {
	var r = int64(1)
	kk := min(k, n - k)
	for i := 0; i < kk; i++ {
		r = r * int64(n - i)
		r = r / int64(i + 1)
	}

	return r
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func pow(a int, e int) int64 {
	var r = int64(1)
	for i := 0; i < e; i++ {
		r = r * int64(a)
	}

	return r
}

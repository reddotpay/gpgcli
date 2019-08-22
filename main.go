package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/reddotpay/gogpg"
)

// HELP displays the help text message
const HELP = `
encrypt or decrypt using gnu/pg or gpg

gpgcli help
gpgcli encrypt <file> --public <public-file-gpg> [--output <outfile>]
gpgcli decrypt <file> --secret <secret-file-gpg> --passphrase <passphrase> [--output <outfile>]

Commands:
encrypt - encrypts the file
decrypt - decrypts the file
help    - displays this help

Options:
--public     - gpg public key
--secret     - gpg secret key
--passphrase - gpg secret passphrase
--output     - dumps the output to this file
`

func main() {
	var (
		command string
		err     error
		buf     []byte

		args = parseArgs()

		output     = flagAssign(args, "--output", "")
		public     = flagAssign(args, "--public", "")
		secret     = flagAssign(args, "--secret", "")
		passphrase = flagAssign(args, "--passphrase", "")
	)

	command, err = getCommand()
	handleErr(err)

	switch command {
	case "help":
	case "-h":
		fmt.Print(HELP)
		return

	case "encrypt":
		if buf, err = encrypt(public); nil != err {
			throwErr(err.Error())
		}
		break

	case "decrypt":
		if buf, err = decrypt(secret, passphrase); nil != err {
			throwErr(err.Error())
		}
		break

	default:
		throwErr(fmt.Sprintf("command `%s` not supported", command))
	}

	if "" != output {
		ioutil.WriteFile(output, buf, 0644)
	} else {
		fmt.Printf("%s", buf)
	}
}

// parseArgs parses the arguments and maps with the -- parameter
func parseArgs() map[string]string {
	var (
		flagPattern = regexp.MustCompile(`^\-\-[\-a-z0-9]+$`)
		params      = map[string]string{}
		activeFlag  = ""
	)

	for _, f := range os.Args {
		if "" != activeFlag {
			params[activeFlag] = f
			activeFlag = ""
			continue
		}

		if flagPattern.MatchString(f) {
			activeFlag = f
			continue
		}
	}

	return params
}

// flagAssign searches the flag map and assigns the default value what the key is not found
func flagAssign(flag map[string]string, key, defaultVal string) string {
	if _, ok := flag[key]; ok {
		return flag[key]
	}

	return defaultVal
}

// throwErr displays the [m]essage and kills the program
func throwErr(m string) {
	fmt.Printf("error: %s. use help for help\n", m)
	os.Exit(1)
}

// handleErr validates an error and uses `throwErr`
func handleErr(err error) {
	if nil == err {
		return
	}

	throwErr(err.Error())
}

// getCommand fetches the first argument from the CLI
func getCommand() (string, error) {
	if 1 >= len(os.Args) {
		return "", errors.New("missing command")
	}

	return os.Args[1], nil
}

// loadTargetFile checks the second argument and assumes it's the target file
func loadTargetFile() ([]byte, error) {
	if 2 >= len(os.Args) {
		return nil, errors.New("missing file")
	}

	return ioutil.ReadFile(os.Args[2])
}

// encrypt encrypts the file with public key
func encrypt(key string) ([]byte, error) {
	var (
		f   []byte
		kf  *os.File
		err error
	)

	if f, err = loadTargetFile(); nil != err {
		return nil, err
	}

	if "" == key {
		return nil, errors.New("missing `--public` argument while using encrypt")
	}

	if kf, err = os.Open(key); nil != err {
		return nil, err
	}

	return gogpg.Encrypt(kf, f)
}

// decrypt decrypts the file with private key and passphrase
func decrypt(key, passphrase string) ([]byte, error) {
	var (
		f   []byte
		kf  *os.File
		err error
	)

	if f, err = loadTargetFile(); nil != err {
		return nil, err
	}

	if "" == key {
		return nil, errors.New("missing `--secret` argument while using decrypt")
	}

	if "" == passphrase {
		return nil, errors.New("missing `--passphrase` argument while using decrypt")
	}

	if kf, err = os.Open(key); nil != err {
		return nil, err
	}

	return gogpg.Decrypt(kf, passphrase, f)
}

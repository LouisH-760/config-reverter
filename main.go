package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func fileExists(path string) bool {
	if fileInfo, err := os.Stat(path); err != nil {
		return false
	} else {
		return !fileInfo.IsDir()
	}
}

func fileHash(path string) string {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Could not open %s (%s)\n", path, err.Error())
		return ""
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Printf("Could not hash %s (%s)\n", path, err.Error())
		return ""
	}
	return string(hash.Sum(nil))
}

func restore(source string, replace string, backup string) {
	if err := os.Rename(replace, source); err != nil {
		fmt.Printf("Could not restore the config file %s (%s).\n", replace, err.Error())
		return
	}
	if err := os.Rename(backup, replace); err != nil {
		fmt.Printf("Could not restore the config file %s (%s).\n", replace, err.Error())
		return
	}
}

func runcmd(cmd []string) ([]byte, error) {
	if len(cmd) > 1 {
		return exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	} else {
		return exec.Command(cmd[0]).Output()
	}
}

func main() {
	sleeptime := 60 * time.Second

	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Print("Usage: reverter newconfig oldconfig command\n")
		return
	}

	source := args[0]              // new file
	replace := args[1]             // file to replace
	cmd := strings.Fields(args[2]) // command to run once the file is replaced

	backup := fmt.Sprintf("%s.reverterbak", replace)

	if !fileExists(source) {
		fmt.Printf("%s does not exist or is a directory.\n", source)
		return
	}

	if !fileExists(replace) {
		fmt.Printf("%s does not exist or is a directory.\n", replace)
		return
	}

	if fileExists(backup) {
		fmt.Printf("A backup already exists (%s). Please delete it before continuing.\n", backup)
		return
	}

	sourceHash := fileHash(source)
	replaceHash := fileHash(replace)

	if sourceHash == replaceHash {
		fmt.Printf("Both files have the same hash, no need to do anything!\n")
		return
	}

	if err := os.Rename(replace, backup); err != nil { // take a backup of the old file
		fmt.Printf("Could not rename old config file %s (%s).\n", replace, err.Error())
		return
	}

	if err := os.Rename(source, replace); err != nil { // move the new file in place of the old config
		fmt.Printf("Could not rename new config file %s (%s).\n", replace, err.Error())
		return
	}

	output, err := runcmd(cmd)

	if err != nil {
		fmt.Printf("Command failed: %s", err.Error())
		restore(source, replace, backup)
		return
	}

	fmt.Printf("---Command succeeded with following output:---\n%s\n---If not interrupted, the config will be restored in %s.---", output, sleeptime)
	time.Sleep(sleeptime)
	fmt.Printf("Not interrupted, restoring %s", replace)
	restore(source, replace, backup)
	runcmd(cmd)
}

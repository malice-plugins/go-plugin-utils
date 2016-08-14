package utils

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/parnurzeal/gorequest"
)

func Getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

func Assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GetSHA256 calculates a file's sha256sum
func GetSHA256(name string) string {

	dat, err := ioutil.ReadFile(name)
	Assert(err)

	h256 := sha256.New()
	_, err = h256.Write(dat)
	Assert(err)

	return fmt.Sprintf("%x", h256.Sum(nil))
}

// RunCommand runs cmd on file
func RunCommand(cmd string, args ...string) string {

	cmdOut, err := exec.Command(cmd, args...).Output()
	if len(cmdOut) == 0 {
		Assert(err)
	}

	return string(cmdOut)
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(resp.Status)
}

// RemoveDuplicates removes duplicate items from a list
func RemoveDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

package gopsearch

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	regList               []*regexp.Regexp
	extensionWhiteListPtr *[]string
	extensionBlackListPtr *[]string
	onlyFilesPtr          *bool
)

// RunSearchCmd : Run the Search command
func RunSearchCmd(patternList []string, pathList []string, extensionWhiteList []string, extensionBlackList []string, onlyFiles bool) {
	extensionWhiteListPtr = &extensionWhiteList
	extensionBlackListPtr = &extensionBlackList
	onlyFilesPtr = &onlyFiles

	// Compile all patterns
	for _, expr := range patternList {
		regExp, err := regexp.Compile(expr)
		if err != nil {
			fmt.Printf("[X] Error with expression : %s\n", expr)
			break
		}
		regList = append(regList, regExp)
	}

	// Walk from each given location in order to found files
	for _, path := range pathList {
		err := filepath.Walk(path, findInPath)
		if err != nil {
			fmt.Printf("Error during walk in location : %s\n", path)
		}
	}
}

func findInPath(path string, info os.FileInfo, err error) error {
	// Apply white list option. If extension file is blacklist then do to check the file
	if len(*extensionWhiteListPtr) > 0 {
		found := false
		for _, extension := range *extensionWhiteListPtr {
			if strings.HasSuffix(info.Name(), "."+extension) {
				found = true
			}
		}

		if found == false {
			return nil
		}

	} else {
		// Apply black list option. If extension file is blacklist then do to check the file
		for _, extension := range *extensionBlackListPtr {
			if strings.HasSuffix(info.Name(), "."+extension) {
				return nil
			}
		}
	}

	for _, re := range regList {
		res := re.MatchString(info.Name())
		if res == true {
			if info.IsDir() && *onlyFilesPtr == false {
				fmt.Printf("[+] [D] %s\n", path)
			} else {
				fmt.Printf("[+] [F] %s\n", path)
			}
		}
	}

	// return errors.New("Could not find in the path.")
	return nil
}
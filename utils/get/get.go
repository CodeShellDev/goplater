package get

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/codeshelldev/goplater/utils/fsutils"
)

func Remote(key string) string {
	_, err := url.Parse(key)

	if err != nil {
		fmt.Println("error parsing url:", err.Error())
		return "invalid url: " + key
	}

	response, err := http.Get(key)
	if err != nil {
		fmt.Println("error requesting remote:", err.Error())
		return "remote failed: " + key
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading body:", err.Error())
		return "body malformed: " + key
	}

	return string(body)
}

func Local(key, context string) string {
	contextAbs, _ := filepath.Abs(context)

	fullPath := fsutils.Relative(contextAbs, key)
	
	fullPath, err := filepath.Abs(fullPath)

	if err == nil {
		data, err := os.ReadFile(fullPath)
		
		if err != nil {
			fmt.Println("error reading file:", err.Error())
			return "file not found: " + fullPath
		}

		return string(data)
	} else {
		fmt.Println("error reading file:", err.Error())
		return "invalid path: " + fullPath 
	}
}
package fetchutils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/codeshelldev/goplater/utils/fsutils"
)

func Remote(urlStr string) string {
	_, err := url.Parse(urlStr)

	if err != nil {
		fmt.Println("error parsing url:", err.Error())
		return "invalid url: " + urlStr
	}

	response, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("error requesting remote:", err.Error())
		return "remote failed: " + urlStr
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading body:", err.Error())
		return "body malformed: " + urlStr
	}

	return string(body)
}

func Local(path, context string) string {
	contextAbs, _ := filepath.Abs(context)

	fullPath := fsutils.Relative(contextAbs, path)
	
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
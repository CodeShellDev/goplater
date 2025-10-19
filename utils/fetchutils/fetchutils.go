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

func Remote(path string) string {
	_, err := url.Parse(path)

	if err != nil {
		fmt.Println("error parsing url:", err.Error())
		return "invalid url: " + path
	}

	response, err := http.Get(path)
	if err != nil {
		fmt.Println("error requesting remote:", err.Error())
		return "remote failed: " + path
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading body:", err.Error())
		return "body malformed: " + path
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
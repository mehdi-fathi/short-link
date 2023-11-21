package url

import (
	"log"
	"net/http"
	"strings"
)

func CheckURL(url string) bool {
	url = strings.TrimSpace(url)

	resp, err := http.Get(url)

	if err != nil {
		return false // Treat any error as a failed check
	}

	log.Println(resp.StatusCode)

	defer resp.Body.Close()

	// Return true if the status code is not 404, false otherwise
	return resp.StatusCode == http.StatusOK
}

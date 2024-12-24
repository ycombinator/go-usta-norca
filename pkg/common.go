package pkg

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strconv"
)

const baseUrl = "https://leagues.ustanorcal.com"

func getDoc(u string) (*goquery.Document, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve organization details from url [%s]: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unable to retrieve organization details from url [%s]; got status code [%d]", u, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to parse organization details: %w", err)
	}

	return doc, nil
}

func parseIDFromUrl(u string) (int, error) {
	up, err := url.Parse(u)
	if err != nil {
		return 0, fmt.Errorf("unable to parse URL [%s]: %w", u, err)
	}
	idStr := up.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("unable to convert ID [%s] to integer: %w", idStr, err)
	}

	return id, nil
}

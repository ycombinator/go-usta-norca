package pkg

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type OrganizationDetails struct {
	Name string
}

type OrganizationTeam struct {
	Status    string
	Name      string
	ID        int
	Area      string
	Captain   string
	StartDate time.Time
}

const organizationURL = "/organization.asp?id=%d"

func GetOrganizationDetails(id int) (*OrganizationDetails, error) {
	doc, err := getOrganizationDoc(id)
	if err != nil {
		return nil, err
	}

	sel := doc.Find("table tbody tr td font b")
	if sel == nil {
		return nil, fmt.Errorf("unable to parse organization name: %w", err)
	}

	o := new(OrganizationDetails)
	caser := cases.Title(language.English)
	o.Name = caser.String(sel.First().Text())
	return o, nil
}

func GetOrganizationTeams(id int) ([]OrganizationTeam, error) {
	doc, err := getOrganizationDoc(id)
	if err != nil {
		return nil, err
	}

	sel := doc.Find("table tbody tr td table tbody tr")
	if sel == nil {
		return nil, fmt.Errorf("unable to parse organization teams: %w", err)
	}

	orgTeams := make([]OrganizationTeam, 0)
	headerSeen := false
	sel.Each(func(i int, row *goquery.Selection) {
		cells := row.Children()
		if cells.Length() < 6 {
			return
		}
		if !headerSeen {
			headerSeen = true
			return
		}

		status := cells.First().Text()

		cells = cells.Next()
		name := cells.First().Text()

		u := cells.First().Get(0).FirstChild.Attr[0].Val
		up, err := url.Parse(u)
		if err != nil {
			return
		}
		idStr := up.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}

		cells = cells.Next()
		area := cells.First().Text()

		cells = cells.Next()
		captain := cells.First().Text()

		cells = cells.Next()
		cells = cells.Next()
		startDateStr := cells.First().Text()
		startDate, err := time.ParseInLocation("01/02/2006", startDateStr, time.Local)
		if err != nil {
			return
		}

		orgTeam := OrganizationTeam{
			Status:    status,
			Name:      name,
			ID:        id,
			Area:      area,
			Captain:   captain,
			StartDate: startDate,
		}
		orgTeams = append(orgTeams, orgTeam)
	})

	return orgTeams, nil
}

func getOrganizationDoc(id int) (*goquery.Document, error) {
	u := baseUrl + fmt.Sprintf(organizationURL, id)

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

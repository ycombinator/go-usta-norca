package pkg

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const teamURL = "/teaminfo.asp?id=%d"

// All N at HH:MM AMPM
// M/N at HH:MM AMPM
var matchTimesRegex = regexp.MustCompile(`^(All\s+\d|\d/\d)\s+at\s+(\d\d?):(\d\d)\s+([AP]M)`)

type TeamMatch struct {
	Week         int
	Date         time.Time
	OpponentName string
	OpponentID   int
	Location     string
	IsScheduled  bool
}

func GetTeamMatches(id int) ([]TeamMatch, error) {
	doc, err := getTeamDoc(id)
	if err != nil {
		return nil, err
	}

	sel := doc.Find("table tbody tr td table tbody tr")
	if sel == nil {
		return nil, fmt.Errorf("unable to parse team matches: %w", err)
	}

	teamMatches := make([]TeamMatch, 0)
	headerSeen := false
	sel.Each(func(i int, row *goquery.Selection) {
		cells := row.Children()
		if cells.Length() < 10 {
			return
		}
		if !headerSeen {
			headerSeen = true
			return
		}

		// Status; skip because it's not accurate
		_ = cells.First().Text()

		// Week number
		cells = cells.Next()
		weekCell := cells.First().Text()
		weekCell = html.UnescapeString(weekCell)
		week, err := parseWeek(weekCell)
		if err != nil {
			return
		}

		// Match date
		matchDateStr := cells.First().Text()
		matchDate, err := time.ParseInLocation("01/02/2006", matchDateStr, time.Local)

		// Match day (of the week); skip because it can be derived
		cells = cells.Next()

		// Match notes (includes time, if scheduled)
		matchNotes := cells.First().Text()
		matchNotes = html.UnescapeString(matchNotes)
		matchDate, isScheduled, err := parseMatchTime(matchDate, matchNotes)
		if err != nil {
			return
		}

		// Opponent team name
		cells = cells.Next()
		opponentName := cells.First().Text()

		// Opponent team ID
		u := cells.First().Get(0).FirstChild.Attr[0].Val
		opponentID, err := parseIDFromUrl(u)
		if err != nil {
			return
		}

		// Location (home or away)
		cells = cells.Next()
		location := cells.First().Text()

		teamMatch := TeamMatch{
			Week:         week,
			Date:         matchDate,
			OpponentName: opponentName,
			OpponentID:   opponentID,
			Location:     location,
			IsScheduled:  isScheduled,
		}
		teamMatches = append(teamMatches, teamMatch)
	})

	return teamMatches, nil
}

func getTeamDoc(id int) (*goquery.Document, error) {
	u := baseUrl + fmt.Sprintf(teamURL, id)
	return getDoc(u)
}

func parseWeek(weekCell string) (int, error) {
	weekCell = strings.TrimSpace(weekCell)
	parts := strings.Split(weekCell, " ")
	weekStr := parts[0]
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		return 0, fmt.Errorf("unable to parse week [%s] as integer: %w", weekStr, err)
	}

	return week, nil
}

func parseMatchTime(matchDate time.Time, matchNotes string) (time.Time, bool, error) {
	matchNotes = strings.TrimSpace(matchNotes)

	// Match notes are empty --> match is not scheduled yet.
	if matchNotes == "" {
		return matchDate, false, nil
	}

	matches := matchTimesRegex.FindStringSubmatch(matchNotes)
	if len(matches) >= 5 {
		hourStr := matches[2]
		minutesStr := matches[3]
		amPmStr := matches[4]

		hour, err := strconv.Atoi(hourStr)
		if err != nil {
			return matchDate, false, fmt.Errorf("unable to parse hour [%s] from match start time as integer: %w", hourStr, err)
		}

		minutes, err := strconv.Atoi(minutesStr)
		if err != nil {
			return matchDate, false, fmt.Errorf("unable to parse minutes [%s] from match start time as integer: %w", minutesStr, err)
		}

		if amPmStr == "PM" {
			hour += 12
		}

		timeDuration := time.Duration(hour)*time.Hour + time.Duration(minutes)*time.Minute
		matchDate = matchDate.Add(timeDuration)
	}

	return matchDate, true, nil
}

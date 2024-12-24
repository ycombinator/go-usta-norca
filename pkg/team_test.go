package pkg

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetTeamMatches(t *testing.T) {
	teamMatches, err := GetTeamMatches(104288)
	require.NoError(t, err)
	require.Len(t, teamMatches, 10)

	firstMatch := teamMatches[0]
	expectedMatch := TeamMatch{
		Week:         1,
		Date:         time.Date(2025, 1, 8, 18, 30, 0, 0, time.Local),
		OpponentName: "SILVER CREEK 40AM4.0A",
		OpponentID:   104945,
		Location:     "Home",
		IsScheduled:  true,
	}
	require.Equal(t, expectedMatch, firstMatch)
}

func TestParseWeek(t *testing.T) {
	cases := map[string]int{
		" 3 (1/20-1/26) ": 3,
		"5 (2/3-2/9)":     5,
		"10 (3/10-3/16)":  10,
	}

	for weekCell, expected := range cases {
		t.Run(weekCell, func(t *testing.T) {
			actual, err := parseWeek(weekCell)
			require.NoError(t, err)
			require.Equal(t, expected, actual)
		})
	}
}

func TestParseMatchTime(t *testing.T) {
	matchDate := time.Date(2024, 12, 24, 0, 0, 0, 0, time.Local)
	cases := map[string]struct {
		matchNotes  string
		matchDate   time.Time
		isScheduled bool
		error       string
	}{
		"empty":      {"", matchDate, false, ""},
		"whitespace": {" ", matchDate, false, ""},
		"all_4":      {"All 4 at  7:30 PM Water on courts", matchDate.Add(19*time.Hour + 30*time.Minute), true, ""},
		"split_3_1":  {"3/1 at  6:30 PM and 7:45 PM - D3 at 7:45 PM. Please bring refillable water bottles. ", matchDate.Add(18*time.Hour + 30*time.Minute), true, ""},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			actualMatchDate, isScheduled, err := parseMatchTime(matchDate, test.matchNotes)
			if test.error == "" {
				require.NoError(t, err)
			} else {
				require.Equal(t, test.error, err.Error())
			}

			require.Equal(t, test.matchDate, actualMatchDate)
			require.Equal(t, test.isScheduled, isScheduled)
		})
	}
}

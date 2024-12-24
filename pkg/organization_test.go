package pkg

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetOrganizationDetails(t *testing.T) {
	org, err := GetOrganizationDetails(225)
	require.NoError(t, err)
	require.Equal(t, "Almaden Swim And Racquet Club", org.Name)
}

func TestGetOrganizationTeams(t *testing.T) {
	orgTeams, err := GetOrganizationTeams(225)
	require.NoError(t, err)
	require.Len(t, orgTeams, 60)

	firstTeam := orgTeams[0]
	expectedTeam := OrganizationTeam{
		Status:          "Register",
		Name:            "ALMADEN SR 40AM3.5A-DT",
		ID:              105115,
		Area:            "Peninsula - Lower",
		Captain:         "Caouette, Cory",
		SeasonStartDate: time.Date(2025, 01, 06, 0, 0, 0, 0, time.Local),
	}
	require.Equal(t, expectedTeam, firstTeam)
}

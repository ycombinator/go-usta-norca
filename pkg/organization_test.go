package pkg

import (
	"github.com/stretchr/testify/require"
	"testing"
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
	require.Equal(t, "Register", firstTeam.Status)
	require.Equal(t, "ALMADEN SR 40AM3.5A-DT", firstTeam.Name)
	require.Equal(t, "Peninsula - Lower", firstTeam.Area)
	require.Equal(t, "Caouette, Cory", firstTeam.Captain)
	require.Equal(t, "01/06/2025", firstTeam.StartDate.Format("01/02/2006"))
}

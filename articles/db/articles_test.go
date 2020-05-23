package db

import (
	"testing"

	"github.com/ellemouton/snell/db"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	dbc := db.ConnectForTesting(t)

	_, err := Create(dbc, "title", "the dash", 1, "blah blah blah.....")
	require.NoError(t, err)
}

func TestLookupInfo(t *testing.T) {
	dbc := db.ConnectForTesting(t)

	id, err := Create(dbc, "title", "the dash", 1, "blah blah blah.....")
	require.NoError(t, err)

	info, err := LookupInfo(dbc, id)
	require.NoError(t, err)

	require.Equal(t, info.Name, "title")
	require.Equal(t, info.Description, "the dash")
	require.Equal(t, info.Price, int64(1))
}

func TestLookupContent(t *testing.T) {
	dbc := db.ConnectForTesting(t)

	id, err := Create(dbc, "title", "the dash", 1, "blah blah blah")
	require.NoError(t, err)

	info, err := LookupInfo(dbc, id)
	require.NoError(t, err)

	content, err := LookupContent(dbc, info.ContentID)
	require.NoError(t, err)
	require.Equal(t, content.Text, "blah blah blah")
}

func TestListAllInfos(t *testing.T) {
	dbc := db.ConnectForTesting(t)

	_, err := Create(dbc, "title", "the dash", 1, "blah blah blah")
	require.NoError(t, err)

	_, err = Create(dbc, "title", "the dash", 1, "blah blah blah")
	require.NoError(t, err)

	_, err = Create(dbc, "title", "the dash", 1, "blah blah blah")
	require.NoError(t, err)

	infos, err := ListAllInfo(dbc)
	require.NoError(t, err)
	require.Equal(t, len(infos), 3)
}

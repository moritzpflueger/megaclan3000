package database

import (
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/google/go-cmp/cmp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pinpox/megaclan3000/internal/steamclient"
)

var (
	db       *DataStorage
	fixtures *testfixtures.Loader
)

func prepareDB() {

	var err error
	if db, err = NewDataStorage("../../test/database/test.db"); err != nil {
		panic(err)
	}

	if fixtures, err = testfixtures.New(
		testfixtures.Database(db.db),
		testfixtures.Dialect("sqlite"),
		testfixtures.Directory(
			"../../test/database/fixtures",
		),
	); err != nil {
		panic(err)
	}

	if err = fixtures.Load(); err != nil {
		panic(err)
	}
}

func TestDataStorage_GetPlayerInfoBySteamID(t *testing.T) {

	summary, err := db.GetPlayerSummary("all_columns")
	if err != nil {
		panic(err)
	}
	stats, err := db.GetUserStatsForGame("all_columns")
	if err != nil {
		panic(err)
	}
	recent, err := db.GetRecentlyPlayedGames("all_columns")
	if err != nil {
		panic(err)
	}
	history, err := db.GetPlayerHistory("all_columns")
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		steamID string
		want    steamclient.PlayerInfo
		wantErr bool
	}{
		{
			name:    "Retrieve PlayerInfo from fixtures (ID: all_columns)",
			steamID: "all_columns",
			want: steamclient.PlayerInfo{
				PlayerSummary:       summary,
				UserStatsForGame:    stats,
				RecentlyPlayedGames: recent,
				PlayerHistory:       history,
			},
			wantErr: false,
		},
		{
			name:    "Retrieve PlayerInfo from fixtures (ID: no_exist)",
			steamID: "all_columns",
			want: steamclient.PlayerInfo{
				PlayerSummary:       steamclient.PlayerSummary{},
				UserStatsForGame:    steamclient.UserStatsForGame{},
				RecentlyPlayedGames: steamclient.RecentlyPlayedGames{},
				PlayerHistory:       steamclient.PlayerHistory{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepareDB()
			got, err := db.GetPlayerInfoBySteamID(tt.steamID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataStorage.GetPlayerInfoBySteamID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("DataStorage.GetPlayerInfoBySteamID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

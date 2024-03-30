package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumYoutubeToSpotify(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  utils.FixturesMap
	}{
		{
			name:         "when youtube album link given and spotify album found",
			input:        "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/sf",
			expectedText: "https://open.spotify.com/album/1HrMmB5useeZ0F5lHrMvl0",
			fixturesMap: utils.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read("youtube/get_album_radiohead_amnesiac.json"),
				},
				// TODO: fix adapter query
				SpotifySearchAlbums: map[string][]byte{
					"artist:Radiohead Amnesiac (2001) album:": fixture.Read("spotify/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:         "when youtube album link given and spotify album not found",
			input:        "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/sf",
			expectedText: "Album not found in Spotify",
			fixturesMap: utils.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read("youtube/get_album_radiohead_amnesiac.json"),
				},
				SpotifySearchAlbums: map[string][]byte{},
			},
		},
		{
			name:         "when youtube album not found",
			input:        "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/sf",
			expectedText: "",
			fixturesMap: utils.FixturesMap{
				YoutubeAlbums:       map[string][]byte{},
				SpotifySearchAlbums: map[string][]byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixturesMap.Merge(&tt.fixturesMap)
			defer fixturesMap.Reset()
			defer senderMock.Reset()

			callback := telebot.Callback{
				Sender: user,
				Data:   tt.input,
			}

			vs.CallbackHandler(&callback)

			if tt.expectedText == "" {
				require.Nil(t, senderMock.Response)
			} else {
				require.NotNil(t, senderMock.Response)
				require.Equal(t, user, senderMock.Response.To)
				require.Equal(t, tt.expectedText, senderMock.Response.Text)
			}
		})
	}
}

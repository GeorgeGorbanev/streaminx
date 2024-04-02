package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumAppleToYandex(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  fixture.FixturesMap
	}{
		{
			name:         "when apple album link given and yandex album found",
			input:        "cnval/ap/1097864180/ya",
			expectedText: "https://music.yandex.com/album/3389008",
			fixturesMap: fixture.FixturesMap{
				AppleAlbums: map[string][]byte{
					"1097864180": fixture.Read("apple/get_album_radiohead_amnesiac.json"),
				},
				YandexSearchAlbums: map[string][]byte{
					"radiohead – amnesiac": fixture.Read("yandex/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:         "when apple album link given and yandex album not found",
			input:        "cnval/ap/1097864180/ya",
			expectedText: "Album not found in Yandex",
			fixturesMap: fixture.FixturesMap{
				AppleAlbums: map[string][]byte{
					"1097864180": fixture.Read("apple/get_album_radiohead_amnesiac.json"),
				},
				YandexSearchAlbums: map[string][]byte{},
			},
		},
		{
			name:         "when apple album not found",
			input:        "cnval/ap/1097864180/ya",
			expectedText: "",
			fixturesMap: fixture.FixturesMap{
				AppleAlbums:        map[string][]byte{},
				YandexSearchAlbums: map[string][]byte{},
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

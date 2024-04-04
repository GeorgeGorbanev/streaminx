package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackYoutubeToApple(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  fixture.FixturesMap
	}{
		{
			name:         "when youtube track link given and apple track found",
			input:        "cnvtr/yt/hbe3CQamF8k/ap",
			expectedText: "https://music.apple.com/us/album/angel/724466069?i=724466660",
			fixturesMap: fixture.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"hbe3CQamF8k": fixture.Read("youtube/get_track_massive_attack_angel.json"),
				},
				AppleSearchTracks: map[string][]byte{
					"Massive Attack Angel": fixture.Read("apple/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when youtube track link given and apple track not found",
			input:        "cnvtr/yt/hbe3CQamF8k/ap",
			expectedText: "Track not found in Apple",
			fixturesMap: fixture.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"hbe3CQamF8k": fixture.Read("youtube/get_track_massive_attack_angel.json"),
				},
				AppleSearchTracks: map[string][]byte{},
			},
		},
		{
			name:         "when youtube track not found",
			input:        "cnvtr/yt/hbe3CQamF8k/ap",
			expectedText: "",
			fixturesMap: fixture.FixturesMap{
				YoutubeTracks:     map[string][]byte{},
				AppleSearchTracks: map[string][]byte{},
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

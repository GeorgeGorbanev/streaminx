package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackYoutubeToApple(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when youtube track link given and apple track found",
			input: "cnvtr/yt/hbe3CQamF8k/ap",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.apple.com/us/album/angel/724466069?i=724466660",
				},
			},
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
			name:  "when youtube autogenerated track link given and apple track found",
			input: "cnvtr/yt/5PgdZDXg0z0/ap",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.apple.com/us/album/space-oddity/697650603?i=697651126",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"5PgdZDXg0z0": fixture.Read("youtube/get_track_david_bowie_space_oddity.json"),
				},
				AppleSearchTracks: map[string][]byte{
					"David Bowie Space Oddity": fixture.Read("apple/search_track_david_bowie_space_oddity.json"),
				},
			},
		},
		{
			name:  "when youtube track link given and apple track not found",
			input: "cnvtr/yt/hbe3CQamF8k/ap",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Track not found in Apple",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"hbe3CQamF8k": fixture.Read("youtube/get_track_massive_attack_angel.json"),
				},
				AppleSearchTracks: map[string][]byte{},
			},
		},
		{
			name:             "when youtube track not found",
			input:            "cnvtr/yt/hbe3CQamF8k/ap",
			expectedMessages: []*telegram.Message{},
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

			require.Equal(t, tt.expectedMessages, senderMock.AllSent)
		})
	}
}

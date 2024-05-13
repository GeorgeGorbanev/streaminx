package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/templates"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackAppleToYandex(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when apple track link given and yandex track found",
			input: "cnvtr/ap/us-724466660/ya",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.yandex.com/album/35627/track/354093",
				},
				{
					To:   user,
					Text: templates.SpecifyRegion,
					ReplyMarkup: &telebot.ReplyMarkup{
						InlineKeyboard: [][]telebot.InlineButton{
							{
								{
									Text: "🇧🇾 Belarus",
									Data: "regtr/354093/by",
								},
							},
							{
								{
									Text: "🇰🇿 Kazakhstan",
									Data: "regtr/354093/kz",
								},
							},
							{
								{
									Text: "🇷🇺 Russia",
									Data: "regtr/354093/ru",
								},
							},
							{
								{
									Text: "🇺🇿 Uzbekistan",
									Data: "regtr/354093/uz",
								},
							},
						},
					},
				},
			},
			fixturesMap: fixture.FixturesMap{
				AppleTracks: map[string][]byte{
					"us-724466660": fixture.Read("apple/get_track_massive_attack_angel.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"massive attack – angel": fixture.Read("yandex/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:  "when apple track link given and yandex track not found",
			input: "cnvtr/ap/us-724466660/ya",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Track not found in Yandex",
				},
			},
			fixturesMap: fixture.FixturesMap{
				AppleTracks: map[string][]byte{
					"us-724466660": fixture.Read("apple/get_track_massive_attack_angel.json"),
				},
				YandexSearchTracks: map[string][]byte{},
			},
		},
		{
			name:             "when apple track not found",
			input:            "cnvtr/ap/us-724466660/ya",
			expectedMessages: []*telegram.Message{},
			fixturesMap: fixture.FixturesMap{
				AppleTracks:        map[string][]byte{},
				YandexSearchTracks: map[string][]byte{},
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

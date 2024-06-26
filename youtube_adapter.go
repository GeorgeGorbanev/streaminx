package streaminx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/GeorgeGorbanev/streaminx/internal/youtube"
)

type YoutubeAdapter struct {
	client youtube.Client
}

var nonTitleContentRe = regexp.MustCompile(`\s*\[.*?\]|\s*\{.*?\}|\s*\(.*?\)`)

func newYoutubeAdapter(client youtube.Client) *YoutubeAdapter {
	return &YoutubeAdapter{
		client: client,
	}
}

func (a *YoutubeAdapter) DetectTrackID(trackURL string) (string, error) {
	if matches := youtube.VideoRe.FindStringSubmatch(trackURL); len(matches) > 1 {
		return matches[1], nil
	}
	return "", IDNotFoundError
}

func (a *YoutubeAdapter) DetectAlbumID(albumURL string) (string, error) {
	if matches := youtube.PlaylistRe.FindStringSubmatch(albumURL); len(matches) > 1 {
		return matches[1], nil
	}
	return "", IDNotFoundError
}

func (a *YoutubeAdapter) GetTrack(id string) (*Track, error) {
	video, err := a.client.GetVideo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get video from youtube: %w", err)
	}
	if video == nil {
		return nil, nil
	}
	return a.adaptTrack(video), nil
}

func (a *YoutubeAdapter) SearchTrack(artistName, trackName string) (*Track, error) {
	query := fmt.Sprintf("%s – %s", artistName, trackName)
	video, err := a.client.SearchVideo(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search video on youtube: %w", err)
	}
	if video == nil {
		return nil, nil
	}

	return a.adaptTrack(video), nil
}

func (a *YoutubeAdapter) GetAlbum(id string) (*Album, error) {
	album, err := a.client.GetPlaylist(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist from youtube: %w", err)
	}
	if album == nil {
		return nil, nil
	}
	return a.adaptAlbum(album)
}

func (a *YoutubeAdapter) SearchAlbum(artistName, albumName string) (*Album, error) {
	query := fmt.Sprintf("%s – %s", artistName, albumName)
	album, err := a.client.SearchPlaylist(query)
	if err != nil {
		return nil, fmt.Errorf("failed to search playlist on youtube: %w", err)
	}
	if album == nil {
		return nil, nil
	}

	return a.adaptAlbum(album)
}

func (a *YoutubeAdapter) adaptTrack(video *youtube.Video) *Track {
	trackTitle := a.extractTrackTitle(video)
	artist, track := a.cleanAndSplitTitle(trackTitle)

	return &Track{
		ID:       video.ID,
		Title:    track,
		Artist:   artist,
		URL:      video.URL(),
		Provider: Youtube,
	}
}

func (a *YoutubeAdapter) extractTrackTitle(video *youtube.Video) string {
	if video.IsAutogenerated() {
		return fmt.Sprintf("%s - %s", video.Artist(), video.Title)
	}
	return video.Title
}

func (a *YoutubeAdapter) adaptAlbum(playlist *youtube.Playlist) (*Album, error) {
	albumTitle, err := a.extractAlbumTitle(playlist)
	if err != nil {
		return nil, fmt.Errorf("failed to extract album title: %w", err)
	}

	artist, album := a.cleanAndSplitTitle(albumTitle)

	return &Album{
		ID:       playlist.ID,
		Title:    album,
		Artist:   artist,
		URL:      playlist.URL(),
		Provider: Youtube,
	}, nil
}

func (a *YoutubeAdapter) extractAlbumTitle(playlist *youtube.Playlist) (string, error) {
	if playlist.IsAutogenerated() {
		videos, err := a.client.GetPlaylistItems(playlist.ID)
		if err != nil {
			return "", fmt.Errorf("failed to get playlist items from youtube: %w", err)
		}
		if len(videos) == 0 {
			return playlist.Title, nil
		}

		v := videos[0]
		if !v.IsAutogenerated() {
			return playlist.Title, nil
		}

		return fmt.Sprintf("%s - %s", v.Artist(), playlist.Album()), nil
	}
	return playlist.Title, nil
}

func (a *YoutubeAdapter) cleanAndSplitTitle(title string) (artist, entity string) {
	cleanTitle := nonTitleContentRe.ReplaceAllString(title, "")

	separators := []string{" - ", " – ", " — ", "|"}
	for _, sep := range separators {
		if strings.Contains(cleanTitle, sep) {
			parts := strings.Split(cleanTitle, sep)
			return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		}
	}

	words := strings.Fields(cleanTitle)
	if len(words) > 1 {
		return strings.TrimSpace(words[0]), strings.TrimSpace(strings.Join(words[1:], " "))
	}

	return cleanTitle, ""
}

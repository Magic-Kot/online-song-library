package musicinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SongDetail struct {
	ReleaseData string
	Text        string
	Link        string
}

type MusicInfo struct {
	url string
}

func NewMusicInfo(url string) *MusicInfo {
	return &MusicInfo{
		url: url,
	}
}

func (m *MusicInfo) Info(group string, song string) (SongDetail, error) {
	var yr SongDetail

	req, err := http.NewRequest(http.MethodGet, m.url, nil)
	if err != nil {
		return yr, err
	}

	req.URL.RawQuery = url.Values{
		"group": {group},
		"song":  {song},
	}.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return yr, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return yr, err
	}

	err = resp.Body.Close()
	if err != nil {
		return yr, err
	}

	if err = json.Unmarshal(body, &yr); err != nil {
		return yr, err
	}

	if resp.StatusCode != http.StatusOK {
		return yr, errors.New(fmt.Sprint("Response status: ", resp.Status))
	}

	return yr, nil
}

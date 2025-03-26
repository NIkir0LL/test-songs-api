package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"songs-api/models"

	"github.com/sirupsen/logrus"
)

func GetSongDetail(apiURL, group, song string) (*models.SongDetail, error) {
	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)
	resp, err := http.Get(apiURL + "/info?" + params.Encode())
	if err != nil {
		logrus.Errorf("Ошибка при запросе к внешнему API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Внешний API вернул код %d", resp.StatusCode)
		return nil, fmt.Errorf("external API returned status: %d", resp.StatusCode)
	}

	var detail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		logrus.Errorf("Ошибка при парсинге ответа внешнего API: %v", err)
		return nil, err
	}
	logrus.Debugf("Получены данные из внешнего API для %s - %s", group, song)
	return &detail, nil
}

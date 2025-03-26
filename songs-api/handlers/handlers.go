package handlers

import (
	"net/http"
	"songs-api/api"
	"songs-api/config"
	"songs-api/models"
	"songs-api/storage"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Получение списка песен
// @Description Получение списка песен с фильтрацией и пагинацией
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Фильтр по группе"
// @Param song query string false "Фильтр по названию песни"
// @Param releaseDate query string false "Фильтр по дате релиза"
// @Param text query string false "Фильтр по тексту"
// @Param link query string false "Фильтр по ссылке"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит записей" default(10)
// @Success 200 {array} models.Song
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs [get]
func GetSongs(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters := map[string]string{
			"group_name":   c.Query("group"),
			"song":         c.Query("song"),
			"release_date": c.Query("releaseDate"),
			"text":         c.Query("text"),
			"link":         c.Query("link"),
		}
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		songs, err := storage.GetSongs(filters, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список песен"})
			return
		}
		c.JSON(http.StatusOK, songs)
	}
}

// @Summary Добавление новой песни
// @Description Добавление новой песни с запросом к внешнему API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.SongInput true "Данные песни"
// @Success 201 {object} models.Song
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs [post]
func AddSong(cfg *config.Config, storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.SongInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		detail, err := api.GetSongDetail(cfg.APIURL, input.Group, input.Song)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить данные из внешнего API"})
			return
		}

		song := &models.Song{
			Group:       input.Group,
			Song:        input.Song,
			ReleaseDate: detail.ReleaseDate,
			Text:        detail.Text,
			Link:        detail.Link,
		}

		if err := storage.CreateSong(song); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить песню"})
			return
		}

		c.JSON(http.StatusCreated, song)
	}
}

// @Summary Получение текста песни с пагинацией по куплетам
// @Description Возвращает текст песни, разбитый на куплеты
// @Tags songs
// @Produce json
// @Param id path int true "ID песни"
// @Param verse query int false "Номер куплета" default(1)
// @Param limit query int false "Количество куплетов" default(1)
// @Success 200 {object} map[string][]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id} [get]
func GetSongText(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		verse, _ := strconv.Atoi(c.DefaultQuery("verse", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1"))

		song, err := storage.GetSongByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
			return
		}

		verses := strings.Split(song.Text, "\n\n")
		start := verse - 1
		if start >= len(verses) {
			c.JSON(http.StatusOK, gin.H{"verses": []string{}})
			return
		}

		end := start + limit
		if end > len(verses) {
			end = len(verses)
		}

		c.JSON(http.StatusOK, gin.H{"verses": verses[start:end]})
	}
}

// @Summary Обновление данных песни
// @Description Обновляет информацию о песне
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.Song true "Новые данные песни"
// @Success 200 {object} models.Song
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id} [put]
func UpdateSong(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var song models.Song
		if err := c.ShouldBindJSON(&song); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, err := storage.GetSongByID(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
			return
		}

		song.ID = id
		if err := storage.UpdateSong(id, &song); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить песню"})
			return
		}

		c.JSON(http.StatusOK, song)
	}
}

// @Summary Удаление песни
// @Description Удаляет песню по ID
// @Tags songs
// @Produce json
// @Param id path int true "ID песни"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id} [delete]
func DeleteSong(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if _, err := storage.GetSongByID(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
			return
		}

		if err := storage.DeleteSong(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить песню"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

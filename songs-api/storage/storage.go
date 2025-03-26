package storage

import (
	"database/sql"
	"fmt"
	"songs-api/config"
	"songs-api/models"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) CreateSong(song *models.Song) error {
	query := `INSERT INTO songs (group_name, song, release_date, text, link) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := s.db.QueryRow(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link).Scan(&song.ID)
	if err != nil {
		logrus.Errorf("Ошибка при создании песни: %v", err)
		return err
	}
	logrus.Infof("Песня %s - %s успешно добавлена с ID %d", song.Group, song.Song, song.ID)
	return nil
}

func (s *Storage) GetSongs(filters map[string]string, page, limit int) ([]models.Song, error) {
	query := "SELECT id, group_name, song, release_date, text, link FROM songs WHERE 1=1"
	var args []interface{}
	argIndex := 1

	for key, value := range filters {
		if value != "" {
			query += fmt.Sprintf(" AND %s ILIKE $%d", key, argIndex)
			args = append(args, "%"+value+"%")
			argIndex++
		}
	}

	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		logrus.Errorf("Ошибка при получении списка песен: %v", err)
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			logrus.Errorf("Ошибка при сканировании песни: %v", err)
			return nil, err
		}
		songs = append(songs, song)
	}
	logrus.Debugf("Получено %d песен", len(songs))
	return songs, nil
}

func (s *Storage) GetSongByID(id int) (*models.Song, error) {
	query := "SELECT id, group_name, song, release_date, text, link FROM songs WHERE id = $1"
	var song models.Song
	err := s.db.QueryRow(query, id).Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		logrus.Errorf("Ошибка при получении песни с ID %d: %v", id, err)
		return nil, err
	}
	logrus.Debugf("Получена песня с ID %d", id)
	return &song, nil
}

func (s *Storage) UpdateSong(id int, song *models.Song) error {
	query := `UPDATE songs SET group_name = $1, song = $2, release_date = $3, text = $4, link = $5 
	          WHERE id = $6`
	_, err := s.db.Exec(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, id)
	if err != nil {
		logrus.Errorf("Ошибка при обновлении песни с ID %d: %v", id, err)
		return err
	}
	logrus.Infof("Песня с ID %d успешно обновлена", id)
	return nil
}

func (s *Storage) DeleteSong(id int) error {
	query := "DELETE FROM songs WHERE id = $1"
	_, err := s.db.Exec(query, id)
	if err != nil {
		logrus.Errorf("Ошибка при удалении песни с ID %d: %v", id, err)
		return err
	}
	logrus.Infof("Песня с ID %d успешно удалена", id)
	return nil
}

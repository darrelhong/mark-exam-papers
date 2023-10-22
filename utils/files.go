package utils

import (
	"regexp"
	"strings"
)

func SanitiseFilename(fileName string) string {
	// trim any leading or trailing whitespace
	sanitised := strings.Trim(fileName, "")

	// replace any non-alphanumeric characters with underscores
	reg := regexp.MustCompile("[^a-zA-Z0-9.]+")
	sanitised = reg.ReplaceAllString(sanitised, "_")

	return sanitised
}

type File struct {
	Id       int64
	Filename string
}

func GetAllFiles() ([]File, error) {
	rows, err := DB.Query("SELECT * FROM files;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File

	for rows.Next() {

		var file File
		err = rows.Scan(&file.Id, &file.Filename)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func GetFilename(id int64) (string, error) {
	var filename string
	err := DB.QueryRow("SELECT filename FROM files WHERE id = ?", id).Scan(&filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func DeleteFile(id int64) error {
	_, err := DB.Exec("DELETE FROM files WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

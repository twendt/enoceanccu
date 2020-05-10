package repository

import (
	"encoding/json"
	"enoceanccu/pkg/entities"
	"os"
)

type Repository interface {
	GetDevices() (map[string]entities.Device, error)
}

type JSONRepository struct {
	filename string
}

func New(filename string) JSONRepository {
	return JSONRepository{filename: filename}
}

func (j JSONRepository) GetDevices() (map[string]entities.Device, error) {
	devices := map[string]entities.Device{}

	f, err := os.Open(j.filename)
	if err != nil {
		return devices, err
	}

	err = json.NewDecoder(f).Decode(&devices)
	if err != nil {
		return devices, err
	}

	return devices, nil
}

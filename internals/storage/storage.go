package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type Storager interface {
	Save(user []User) error
	Load() ([]User, error)
}

type FileStorage struct {
	filepath string
}

func NewStorage(filepath string) FileStorage {
	return FileStorage{filepath: filepath}
}

func (s FileStorage) Save(user []User) error {
	file, err := os.Create(s.filepath)

	if err != nil {
		return fmt.Errorf("could not create file : %s %w ", s.filepath, err)
	}

	defer file.Close()

	b, err := json.MarshalIndent(user, " ", "		")
	if err != nil {
		return fmt.Errorf("could not marshal json %w", err)
	}

	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("Could not write file %w", err)
	}

	return nil
}

func (s FileStorage) Load() ([]User, error) {
	b, err := os.ReadFile(s.filepath)

	if err != nil {
		if os.IsNotExist(err){
			return []User{}, nil
		}
		return []User{}, fmt.Errorf("could not read file : %s %w ", s.filepath, err)
	}

	var data []User

	if err = json.Unmarshal(b, &data); err != nil {
		return []User{}, fmt.Errorf("could not unmarshal file : %s %w ", s.filepath, err)
	}

	return data, nil
}

//how to connect without using load function
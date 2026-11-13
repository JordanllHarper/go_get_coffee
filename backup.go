package main

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"log"
	"os"
)

type jsonFileBackup[T any] struct {
	path string
}

func (ig *jsonFileBackup[T]) Backup(items []T) {
	log.Println("Backup - running")
	json, err := json.Marshal(items)
	if err != nil {
		log.Println("Backup - failed to serialize json:", err)
		return
	}
	// TODO: this could go wrong between the steps so we would want to create FIRST and then clear old
	if _, err := os.Create(ig.path); err != nil {
		log.Println("Backup - failed to clear backup:", err)
		return
	}

	if err := os.WriteFile(ig.path, json, 0644); err != nil {
		log.Println("Backup - failed to clear backup:", err)
		return
	}
	log.Println("Backup - successful!")
}

func (ig *jsonFileBackup[T]) Fetch() ([]T, error) {
	log.Println("Fetch - Fetching from", ig.path)
	f, err := os.Open(ig.path)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Fetch - No file found at %s. Creating new empty file", ig.path)
		f, err = os.Create(ig.path)
		f.Close()
		return make([]T, 0), nil
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching backup from path %s: %w", ig.path, err)
	}
	defer f.Close()

	var backup []T
	if err := json.UnmarshalRead(f, &backup); err != nil {
		return nil, fmt.Errorf("error unmarshalling backup from path %s: %w", ig.path, err)
	}

	log.Println("Fetch - Successful!")

	return backup, nil
}

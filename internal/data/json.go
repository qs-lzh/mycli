package data

import (
	"encoding/json"
	"errors"
	"os"
)

const TasksFile string = "tasks.json"
const MetaFile string = "meta.json"

type Meta struct {
	NextID int `json:"next_id"`
}

func LoadMeta() (*Meta, error) {
	b, err := os.ReadFile(MetaFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Meta{NextID: 1}, nil
		}
		return nil, err
	}
	var m Meta
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func SaveMeta(m *Meta) error {
	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(MetaFile, b, 0644)
}

func AppendJSONL(v any) error {
	f, err := os.OpenFile(TasksFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)

	return enc.Encode(v)
}

func RewriteJSON(tasks []Task) error {
	_, err := os.Create(TasksFile)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		if err := AppendJSONL(t); err != nil {
			return err
		}
	}
	return nil
}

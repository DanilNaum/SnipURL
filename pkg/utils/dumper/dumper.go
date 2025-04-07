package dumper

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
)

//go:generate moq -out logger_moq_test.go . logger
type logger interface {
	Fatalf(format string, v ...any)
}

type dumper struct {
	file   *os.File
	logger logger
}

type URLRecord struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewDumper(path string, log logger) (*dumper, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &dumper{
		file:   file,
		logger: log,
	}, nil
}

func (d *dumper) Add(record *URLRecord) error {
	data, err := json.Marshal(record)

	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = d.file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (d *dumper) ReadAll() (chan URLRecord, error) {
	c := make(chan URLRecord, 10)
	go func() {
		defer close(c)
		writer := bufio.NewReader(d.file)
		for data, err := writer.ReadBytes('\n'); !errors.Is(err, io.EOF); data, err = writer.ReadBytes('\n') {
			if err != nil {
				d.logger.Fatalf("error read data from file: %s", err)
			}

			var record URLRecord
			err := json.Unmarshal(data, &record)
			if err != nil {
				d.logger.Fatalf("error unmarshal data: %s", err)
			}
			c <- record
		}

	}()
	return c, nil
}

func (d *dumper) Close() error {
	return d.file.Close()
}

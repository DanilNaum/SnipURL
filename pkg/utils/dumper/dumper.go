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

type Record struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewDumper(file *os.File, log logger) *dumper {
	return &dumper{
		file:   file,
		logger: log,
	}
}

func (d *dumper) Add(record *Record) error {
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

func (d *dumper) ReadAll() (chan Record, error) {
	c := make(chan Record, 10)
	go func() {
		defer close(c)
		writer := bufio.NewReader(d.file)
		for data, err := writer.ReadBytes('\n'); !errors.Is(err, io.EOF); data, err = writer.ReadBytes('\n') {
			if err != nil {
				d.logger.Fatalf("error read data from file: %s", err)
			}

			var record Record
			err := json.Unmarshal(data, &record)
			if err != nil {
				d.logger.Fatalf("error unmarshal data: %s", err)
			}
			c <- record
		}

	}()
	return c, nil
}

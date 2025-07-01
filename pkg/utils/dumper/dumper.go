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

// URLRecord represents a mapping between a unique identifier, a shortened URL, and its original URL.
// It is used for storing and serializing URL shortening records with JSON tags for marshaling.
type URLRecord struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// NewDumper creates a new dumper with the specified file path and logger.
// It opens the file in append, read-write, and create modes with 0666 permissions.
// Returns a pointer to the dumper and an error if file opening fails.
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

// Add writes a URLRecord to the file as a JSON-encoded line.
// It marshals the record to JSON, appends a newline, and writes to the file.
// Returns an error if JSON marshaling or file writing fails.
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

// ReadAll reads all URLRecords from the file asynchronously.
// It returns a channel of URLRecords and an error.
// Each record is read line by line, unmarshaled from JSON, and sent to the channel.
// Logs fatal errors for read or unmarshal failures.
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

// Close closes the underlying file.
func (d *dumper) Close() error {
	return d.file.Close()
}

package dumper

import (
	"os"
	"strconv"
	"testing"
)

func BenchmarkDumper_Add(b *testing.B) {
	d, err := NewDumper("testfile.json", &loggerMock{})
	if err != nil {
		b.Fatalf("failed to create dumper: %v", err)
	}
	defer d.Close()

	record := &URLRecord{
		UUID:        1,
		ShortURL:    "http://short.url/1",
		OriginalURL: "http://original.url/1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := d.Add(record); err != nil {
			b.Fatalf("failed to add record: %v", err)
		}
	}

	// Удаление файла после теста
	if err := os.Remove("testfile.json"); err != nil {
		b.Fatalf("failed to remove test file: %v", err)
	}
}

func BenchmarkDumper_ReadAll(b *testing.B) {

	d, err := NewDumper("testfile.json", &loggerMock{
		FatalfFunc: func(format string, v ...any) {},
	})
	if err != nil {
		b.Fatalf("failed to create dumper: %v", err)
	}
	defer d.Close()

	// Заполнение файла тестовыми данными
	for i := range 1000 {
		record := &URLRecord{
			UUID:        i,
			ShortURL:    "http://short.url/" + strconv.Itoa(i),
			OriginalURL: "http://original.url/" + strconv.Itoa(i),
		}
		if err := d.Add(record); err != nil {
			b.Fatalf("failed to add record: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, err := d.ReadAll()

		if err != nil {
			b.Fatalf("failed to read all records: %v", err)
		}
		go func() {
			for {
				_, ok := <-c
				if !ok {
					return
				}
			}
		}()
	}
	b.StopTimer()

	// Удаление файла после теста
	if err := os.Remove("testfile.json"); err != nil {
		b.Fatalf("failed to remove test file: %v", err)
	}
}

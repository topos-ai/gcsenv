package gcsenv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
)

func setenv(reader io.Reader) error {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '='
	csvReader.Comment = '#'
	csvReader.FieldsPerRecord = 2
	csvReader.LazyQuotes = true
	csvReader.ReuseRecord = true
	csvReader.TrimLeadingSpace = true
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if len(row) != 2 {
			return fmt.Errorf("gcsenv: invalid row length %d", len(row))
		}

		os.Setenv(row[0], row[1])
	}

	return nil
}

func Setenv(ctx context.Context, bucket, object string) (err error) {
	var client *storage.Client
	client, err = storage.NewClient(ctx)
	if err != nil {
		return
	}

	defer func() {
		if cerr := client.Close(); cerr != nil {
			err = cerr
		}
	}()

	var reader io.ReadCloser
	reader, err = client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return
	}

	defer func() {
		if cerr := reader.Close(); cerr != nil {
			err = cerr
		}
	}()

	err = setenv(reader)
	return
}

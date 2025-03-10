package csv

import (
	"bufio"
	"os"
	"strings"
)

type CSVReader struct {
	FilePath string
}

func (c *CSVReader) ReadLines(urlChan chan<- string) error {
	file, err := os.Open(c.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Skip header line
	scanner.Scan()

	// Read the file line by line and send URLs to the channel
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			urlChan <- line
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

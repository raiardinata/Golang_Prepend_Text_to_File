package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	Filename string
	Contents []string
}

func main() {
	f, err := os.OpenFile("./test.csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	fileStats, err := f.Stat()
	if err != nil {
		fmt.Println(err.Error())
	}

	precision := strconv.Itoa(NumDecPlaces(float64(fileStats.Size()) / 1000))

	size := fmt.Sprintf("%."+precision+"f", float64(fileStats.Size())/1000)
	f.Close()

	NewRecord("./test.csv").Prepend(fmt.Sprintf("\"File Size\" : \""+size+"KB\",%02d", 0))
}

// read length of precision
func NumDecPlaces(v float64) int {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	i := strings.IndexByte(s, '.')
	if i > -1 {
		return len(s) - i - 1
	}
	return 0
}

func NewRecord(filename string) *Record {
	return &Record{
		Filename: filename,
		Contents: make([]string, 0),
	}
}

func (r *Record) readLines() error {
	if _, err := os.Stat(r.Filename); err != nil {
		return nil
	}

	f, err := os.OpenFile(r.Filename, os.O_RDONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if tmp := scanner.Text(); len(tmp) != 0 {
			r.Contents = append(r.Contents, tmp)
		}
	}

	return nil
}

func (r *Record) Prepend(content string) error {
	err := r.readLines()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(r.Filename, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)

	// prepend text before recreate it again
	writer.WriteString(fmt.Sprintf("%s\n", strings.TrimSuffix(content, "00")))
	for _, line := range r.Contents {
		_, err := writer.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

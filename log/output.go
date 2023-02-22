package log

import (
	"bufio"
	"fmt"
	"os"
)

type output struct {
	dir    string
	file   *os.File
	writer *bufio.Writer
}

func newOutput(dir, key string) (*output, error) {
	o := &output{
		dir: dir,
	}

	file, err := o.createFile(key)
	if err != nil {
		return nil, err
	}

	o.file = file
	o.writer = bufio.NewWriter(file)

	return o, nil
}

func (o *output) Writeln(p []byte) error {
	_, err := o.writer.Write(p)
	if err != nil {
		return err
	}

	return o.writer.WriteByte('\n')
}

func (o *output) Flush() error {
	return o.writer.Flush()
}

func (o *output) ResetFile(key string) error {
	err := o.writer.Flush()
	if err != nil {
		return err
	}

	newFile, err := o.createFile(key)
	if err != nil {
		return fmt.Errorf("failed to create new file: %w", err)
	}

	fileName := o.file.Name()
	err = o.file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to close file %s, err: %v", fileName, err)
	}

	o.file = newFile
	o.writer.Reset(o.file)

	return nil
}

func (o *output) createFile(key string) (*os.File, error) {
	fileName := "log." + key
	filePath := o.dir + "/" + fileName
	return os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

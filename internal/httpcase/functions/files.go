package functions

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Files struct {
}

func (f *Files) AppendFile(filename string, text string) error {
	return f.writeFile(filename, text, os.O_WRONLY|os.O_CREATE|os.O_APPEND)
}

func (f *Files) WriteFile(filename string, text string) error {
	return f.writeFile(filename, text, os.O_WRONLY|os.O_CREATE)
}

func (f *Files) writeFile(filename string, text string, flag int) error {
	file, err := os.OpenFile(filename, flag, 0666)
	if err != nil {
		return err
	}

	data := []byte(text)
	n, err := file.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := file.Close(); err == nil {
		err = err1
	}
	return err
}

func (f *Files) ReadFile(filename string) (string, error) {
	text := ""
	file, err := os.Open(filename)
	if err != nil {
		return text, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text = fmt.Sprintf("%s%s\n", text, scanner.Text())
	}
	return strings.Trim(text, "\n"), nil
}

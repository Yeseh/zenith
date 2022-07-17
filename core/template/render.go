package template

import (
	"bytes"
	"html/template"
	"os"
)

func renderTemplate(tmpl string, data interface{}) (*bytes.Buffer, error) {
	t, err := template.New("main").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	if err := t.Execute(buffer, data); err != nil {
		return nil, err
	}

	return buffer, nil
}

func Render(tmpl string, data interface{}) (string, error) {
	tBuf, err := renderTemplate(tmpl, data)
	if err != nil {
		return "", err
	}

	return tBuf.String(), nil
}

func RenderToFile(tmpl string, data interface{}, outputPath string) error {
	tBuf, err := renderTemplate(tmpl, data)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, tBuf.Bytes(), 0666); err != nil {
		return err
	}

	return nil
}

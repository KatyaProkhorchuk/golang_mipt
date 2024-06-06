//go:build !solution

package ciletters

// package main
import (
	"embed"
	"strings"
	"text/template"
)

func LastLog(job string) []string {
	//  берем последние 10 записей
	result := strings.Split(job, "\n")
	if len(result) > 9 {
		result = result[len(result)-10 : len(result)]
	}
	return result

}

//go:embed index.txt
var content embed.FS

func MakeLetter(n *Notification) (string, error) {
	data, _ := content.ReadFile("index.txt")
	funcMap := template.FuncMap{
		"lastLog": LastLog,
	}
	t, err := template.New("letter").Funcs(funcMap).Parse(string(data))
	if err != nil {
		return "", err
	}
	var res strings.Builder
	if err = t.Execute(&res, n); err != nil {
		return "", err
	}
	return res.String(), nil
}

//go:build !solution

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	flag "github.com/spf13/pflag"

	"gitlab.com/slon/shad-go/gitfame/configs"
	"gitlab.com/slon/shad-go/gitfame/internal"
	"gitlab.com/slon/shad-go/gitfame/internal/git"
)

var (
	flagRepository   = flag.StringP("repository", "r", "./", "путь до Git репозитория; по умолчанию текущая директория")
	flagExtensions   = flag.StringSlice("extensions", []string{}, "список расширений, сужающий список файлов в расчёте")
	flagRevision     = flag.StringP("revision", "v", "HEAD", "указатель на коммит")
	flagOrderBy      = flag.StringP("order-by", "o", "lines", "ключ сортировки результатов")
	flagUseCommitter = flag.BoolP("use-committer", "u", false, "булев флаг, заменяющий в расчётах автора (дефолт) на коммиттера")
	flagFormat       = flag.StringP("format", "f", "tabular", "формат вывода; один из `tabular` (дефолт), `csv`, `json`, `json-lines`")
	flagLanguages    = flag.StringSlice("languages", []string{}, "список языков (программирования, разметки и др.), сужающий список файлов в расчёте; множество ограничений разделяется запятыми, например `'go,markdown'`")
	flagExclude      = flag.StringSlice("exclude", []string{}, "набор паттернов, исключающих файлы из расчёта")
	flagRestrictTo   = flag.StringSlice("restrict-to", []string{}, "набор Glob паттернов, исключающий все файлы, не удовлетворяющие ни одному из паттернов набора")
)

func fillFlagInfo(values *map[string]struct{}, flagData []string) {
	for _, value := range flagData {
		(*values)[value] = struct{}{}
	}
}

func checkIncorrectFlag(flagFormat, flagOrderBy string) {
	if !(flagFormat == "tabular" || flagFormat == "csv" || flagFormat == "json" || flagFormat == "json-lines") {
		panic("format error")
	}
	if !(flagOrderBy == "lines" || flagOrderBy == "commits" || flagOrderBy == "files") {
		panic("order-by error")
	}
}

func languageData(flagLanguages []string, languages []internal.Language, extensions *map[string]struct{}) {
	for _, lang := range flagLanguages {
		for _, language := range languages {
			if strings.EqualFold(strings.ToUpper(language.Name), strings.ToUpper(lang)) {
				for _, extension := range language.Extensions {
					(*extensions)[extension] = struct{}{}
				}
			}
		}
	}
}

func sortedAuthor(authors *[]string, statistics map[string]*internal.Statisctics) {
	sort.SliceStable((*authors), func(i, j int) bool {
		stat1 := statistics[(*authors)[i]]
		stat2 := statistics[(*authors)[j]]
		if *flagOrderBy == "lines" {
			if stat1.Lines != stat2.Lines {
				return stat1.Lines > stat2.Lines
			} else {
				if stat1.Commits != stat2.Commits {
					return stat1.Commits > stat2.Commits
				} else if stat1.Files != stat2.Files {
					return stat1.Files > stat2.Files
				}
			}
		} else if *flagOrderBy == "commits" {
			if stat1.Commits != stat2.Commits {
				return stat1.Commits > stat2.Commits
			} else if stat1.Lines != stat2.Lines {
				return stat1.Lines > stat2.Lines
			} else if stat1.Files != stat2.Files {
				return stat1.Files > stat2.Files
			}
		} else if *flagOrderBy == "files" {
			if stat1.Files != stat2.Files {
				return stat1.Files > stat2.Files
			} else if stat1.Lines != stat2.Lines {
				return stat1.Lines > stat2.Lines
			} else if stat1.Commits != stat2.Commits {
				return stat1.Commits > stat2.Commits
			}
		}
		return (*authors)[i] < (*authors)[j]

	})
}

func makeStdOut(stdoutput *[][]string, authors []string, statistics map[string]*internal.Statisctics) {
	(*stdoutput) = append((*stdoutput), []string{"Name", "Lines", "Commits", "Files"})
	for _, author := range authors {
		stat := statistics[author]
		lines := strconv.Itoa(stat.Lines)
		commits := strconv.Itoa(stat.Commits)
		files := strconv.Itoa(stat.Files)
		(*stdoutput) = append((*stdoutput), []string{strings.TrimSpace(author), strings.TrimSpace(lines), strings.TrimSpace(commits), strings.TrimSpace(files)})
	}
}

func printTabular(stdoutput [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, data := range stdoutput {
		fmt.Fprintln(w, data[0]+"\t"+data[1]+"\t"+data[2]+"\t"+data[3])
	}
	w.Flush()
}

func printCSV(stdoutput [][]string) {
	w := csv.NewWriter(os.Stdout)
	for _, data := range stdoutput {
		if err := w.Write(data); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func printJSON(authors []string, statistics map[string]*internal.Statisctics) {
	var jsonOut []interface{}
	for _, author := range authors {
		stat := statistics[author]
		lines := stat.Lines
		commits := stat.Commits
		files := stat.Files
		jsonOut = append(jsonOut, map[string]interface{}{
			"name":    author,
			"lines":   lines,
			"commits": commits,
			"files":   files,
		})
	}
	jsonData, err := json.Marshal(jsonOut)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}

func printJSONLines(authors []string, statistics map[string]*internal.Statisctics) {
	for _, author := range authors {
		stat := statistics[author]
		lines := stat.Lines
		commits := stat.Commits
		files := stat.Files
		jsonData, err := json.Marshal(map[string]interface{}{
			"name":    author,
			"lines":   lines,
			"commits": commits,
			"files":   files,
		})
		if err != nil {
			fmt.Println("Ошибка при преобразовании в JSON:", err)
			return
		}
		fmt.Println(string(jsonData))
	}
}

func main() {
	flag.Parse()
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка при получении текущего рабочего каталога:", err)
		return
	}
	languages := configs.GetLanguages(currentDir + "/../../configs/language_extensions.json")
	useCommitter := *flagUseCommitter
	extensions := make(map[string]struct{})
	fillFlagInfo(&extensions, *flagExtensions)
	exclude := make(map[string]struct{})
	fillFlagInfo(&exclude, *flagExclude)
	restrictTo := make(map[string]struct{})
	fillFlagInfo(&restrictTo, *flagRestrictTo)
	languageData(*flagLanguages, languages, &extensions)
	checkIncorrectFlag(*flagFormat, *flagOrderBy)

	statistics := git.GetStatistics(*flagRepository, *flagRevision, extensions, exclude, restrictTo, useCommitter)
	authors := make([]string, 0)
	for author := range statistics {
		authors = append(authors, author)
	}
	sortedAuthor(&authors, statistics)

	stdoutput := make([][]string, 0)
	makeStdOut(&stdoutput, authors, statistics)
	if *flagFormat == "tabular" {
		printTabular(stdoutput)
	} else if *flagFormat == "csv" {
		printCSV(stdoutput)
	} else if *flagFormat == "json" {
		printJSON(authors, statistics)
	} else if *flagFormat == "json-lines" {
		printJSONLines(authors, statistics)
	}
}

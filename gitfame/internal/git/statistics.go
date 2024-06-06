package git

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/internal"
	"gitlab.com/slon/shad-go/gitfame/internal/command"
)

func checkExclude(exclude map[string]struct{}, files, repository string) bool {
	for pattern := range exclude {
		matched, _ := filepath.Match(pattern, files)
		if matched {
			return false
		}
	}
	path := strings.Split(repository, "/")
	for pattern := range exclude {
		patterns := strings.Split(pattern, "/")
		n := len(path)
		checker := false
		j := 0
		cnt := 0
		for i := 0; i < n; i++ {
			if patterns[j] == path[i] {
				j++
				cnt++
				checker = true

			} else if checker && cnt == len(patterns)-1 {
				return false
			} else if checker && patterns[j] != path[i] {
				checker = false
				j = 0
			}
		}
		if checker && cnt == len(patterns)-1 {
			return false
		}
	}
	return true
}
func checkRestrictTo(restrictTo map[string]struct{}, files, repository string) bool {
	if len(restrictTo) == 0 {
		return true
	}
	return !checkExclude(restrictTo, files, repository)
}
func checkExtensions(extensions map[string]struct{}, files string) bool {
	if len(extensions) == 0 {
		return true
	}
	ext := strings.Split(files, ".")
	if _, ok := extensions["."+ext[len(ext)-1]]; ok {
		return true
	}
	return false
}
func parsePorcelain(statistics map[string]*internal.Statisctics, commit, repository string, author *map[string]string, useCommitter bool) {
	// хотим распарсить коммиты
	// надо узнать структуру коммита где число строк
	commits := strings.Split(commit, "\n")
	changedFiles := make(map[string]int)
	authorsFile := make(map[string]struct{})
	for i, cmt := range commits {
		info := strings.Split(cmt, " ")
		var authorName string
		if info[0] == "author" {
			infoAuthor := strings.Join(info[1:], " ")
			prev := strings.Split(commits[i-1], " ")
			if !useCommitter {
				authorName = infoAuthor
			} else {
				authorName = command.GetGitCommitter(prev[0], repository)
				authorName = strings.Split(authorName, "\n")[0]
			}
			_, authorExists := authorsFile[authorName]
			if !authorExists {
				authorsFile[authorName] = struct{}{}
			}
			_, commitExists := (*author)[prev[0]]
			if !commitExists {

				(*author)[prev[0]] = authorName

				if len(prev) >= 3 {
					countLine, err := strconv.Atoi(prev[3])
					if err == nil {
						addAuthor(statistics, authorName, countLine, authorExists, commitExists)
					} else {
						fmt.Println("Error converting line count:", err)
						return
					}
					if _, ok := changedFiles[prev[0]]; !ok {
						changedFiles[prev[0]] = 0
					}
					changedFiles[prev[0]] += 1
				}
			} else {
				if len(prev) >= 3 {
					addAuthor(statistics, authorName, 0, authorExists, commitExists)
					if _, ok := changedFiles[prev[0]]; !ok {
						changedFiles[prev[0]] = 0
					}
					changedFiles[prev[0]] += 1
				}
			}

		}

		infoAuthor, ok := (*author)[info[0]]
		if len(info) == 4 && ok {
			statistic, err1 := statistics[infoAuthor]
			countLine, err := strconv.Atoi(info[3])
			if err1 && err == nil {
				(*statistic).Lines += countLine
			} else if !err1 && err == nil {
				statistics[infoAuthor] = &internal.Statisctics{
					Lines:   countLine,
					Commits: 1,
					Files:   0,
				}
			}
		}
	}

}
func addAuthor(statistics map[string]*internal.Statisctics, author string, countLine int, authorExists, err bool) {
	if _, ok := statistics[author]; !ok {
		statistics[author] = &internal.Statisctics{
			Lines:   0,
			Commits: 0,
			Files:   0,
		}
	}
	statistic := statistics[author]
	(*statistic).Lines += countLine
	if !err {
		(*statistic).Commits++
	}
	if !authorExists {
		(*statistic).Files++
	}
}

func TreeGit(files, filesAndDirectories, repository, revision string, exclude, restrictTo, extensions map[string]struct{}, author *map[string]string, statistics map[string]*internal.Statisctics, useCommitter bool) {
	filesSplit := strings.Split(files, "\n")
	filesAndDirectoriesSplit := strings.Split(filesAndDirectories, "\n")
	for i, file := range filesAndDirectoriesSplit {
		fileInfo := strings.Split(file, " ")
		if len(fileInfo) >= 2 {
			if fileInfo[1] == "blob" && checkExtensions(extensions, filesSplit[i]) && checkExclude(exclude, filesSplit[i], repository) && checkRestrictTo(restrictTo, filesSplit[i], repository) {
				//  если расширение подходит
				commit := command.GetCommit(revision, filesSplit[i], repository)
				if len(commit) != 0 {
					parsePorcelain(statistics, commit, repository, author, useCommitter)
				} else {
					// если коммитов не было
					name := command.GetGitEmptyFile(revision, filesSplit[i], repository)
					name = strings.Split(name, "\n")[0]
					name = name[1 : len(name)-1]
					part := strings.Split(name, " ")
					hash := part[len(part)-1]
					name = strings.Join(part[:len(part)-1], " ")
					_, err := (*author)[hash]
					if !err {
						(*author)[hash] = name
					}

					addAuthor(statistics, name, 0, false, err)
				}
			} else if fileInfo[1] == "tree" {
				path := filepath.Join(repository, filesSplit[i])
				dirs := command.GetGitFilesAndDirs(revision, path)
				fileNames := command.GetGitFiles(revision, path)
				TreeGit(fileNames, dirs, path, revision, exclude, restrictTo, extensions, author, statistics, useCommitter)

			}
		}

	}
}

func GetStatistics(repository, revision string, extensions, exclude, restrictTo map[string]struct{}, useCommitter bool) map[string]*internal.Statisctics {
	filesAndDirectories := command.GetGitFilesAndDirs(revision, repository)
	files := command.GetGitFiles(revision, repository)
	statistics := make(map[string]*internal.Statisctics)
	author := make(map[string]string)
	TreeGit(files, filesAndDirectories, repository, revision, exclude, restrictTo, extensions, &author, statistics, useCommitter)

	return statistics
}

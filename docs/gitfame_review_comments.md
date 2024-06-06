# Gitfame review comments

В этом файле собрано несколько основных общих замечаний к решениям.

NIT пункты (от фразы nitpick) — это мелкие придирки, которые могут быть нерелевантны в конкретной реализации.

## Парсинг флагов

### Парсинг аргументов массивов

В решениях часто встречается
```golang
flagExclude = flag.String("exclude", "", "Globs to exclude")
exclude := strings.Split(*flagExclude, ",")
```

Для подобных аргументов, передаваемых через запятую можно использовать библиотеку `pflag` — альтернативу стандартному `flag`:  https://pkg.go.dev/github.com/spf13/pflag#StringSlice

```golang
import "github.com/spf13/pflag"
flagExclude = pflag.StringSlice("exclude", nil, "Globs to exclude")
```

`flagExclude` сразу будет слайсом.

### NIT Имена флагов.

Обычно имена переменных начинаются с подстроки, описывающих их природу,
например `ErrNotFound` для ошибки, `flagVerbose` для флага

т.е. вместо
```golang
exclude = flag.String("exclude", "", "Globs to exclude")
restrictTo = flag.String("restrict-to", "", "Globs to restrict")
```

канонично было бы написать
```golang
flagExclude = flag.String("exclude", "", "Globs to exclude")
flagRestrictTo = flag.String("restrict-to", "", "Globs to restrict")
```


## NIT Импорты

В go обычно импорты сгруппированы в три группы
```
import (
    stdlib

    external deps

    current module deps
)
```
Между группами должна быть пустая строка.

Внутри каждой группы импорты должны быть отсортированы. Это проверяет fmt линтер.

Вот так выглядят плохие импорты
```golang
import (
	"flag"
	"log"

	"gitlab.com/slon/shad-go/gitfame/configs"
	"gitlab.com/slon/shad-go/gitfame/internal"

	"github.com/spf13/pflag"

	"strings"
)
```

Вот так хорошие:
```golang
import (
	"flag"
	"log"
	"strings"

	"github.com/spf13/pflag"

	"gitlab.com/slon/shad-go/gitfame/configs"
	"gitlab.com/slon/shad-go/gitfame/internal"
)
```

## Структура проекта

Для такого маленького проекта, как наш было бы абсолютно нормально все файлы, включая `main.go`, парсинг блейма и вывод результата положить в корень репозитория.

Однако, мы специально предложили познакомится со [стандартным лэйаутом go проектов](https://github.com/golang-standards/project-layout) и разложить код в разные директории.

### Не нужно класть реализацию непосредственно в internal или pkg

В месте вызова использование пакета будет выглядеть как-то так:
```golang
internal.ListGitFiles(*flagRepository, *flagRevision)
```

Тут слово `internal` не добавляет никакого смысла.

Вместо этого стоит положить код, например, в подпакет `internal/git`, чтобы вызов выглядел так
```golang
git.ListFiles(*flagRepository, *flagRevision)
```

## JSON encoding

Никогда нельзя собирать JSON с помощью формат строки!

Вместо этого нужно использовать [encoding/json](https://pkg.go.dev/encoding/json).

Неправильно:
```golang
js := fmt.Sprintf("{\"name\":\"%s\",\"lines\":%d,\"commits\":%d,\"files\":%d}", author, lines[author], commits[author], files[author])
```

Правильно:
```golang
import "encoding/json"

js, err := json.Marshal(info)
```
где `info` — это структура или `map[string]interface{}` со всеми данными.

В чём здесь проблема?

Во-первых, может быть сложно правильно экранировать спецсимволы в аргументах, те же кавычки и фигурные скобки.

Во-вторых, собирание JSON'а через формат строку во многих сеттингах подвержено уязвимости инъекции.

Представьте, что вы пишете приложение, которое получает автора `X` от пользователя и кладёт json в sql базу данных.
```
insert `{"author":"X"}` into authors;
```
где json `{"author":"X"}` вы собираете с помощью `fmt.Sprintf`.

И пользователь передаст вам X, равный `"} into authors; drop authors; "`.
тогда он сможет прочитать/затереть приватные данные других пользователей, потому что выполнится такой запрос
```
insert {"author":""} into authors; drop authors; ""}` into authors;
```
Для полного примера тут не хватает деталей (нужно правильно закрыть инъекцию), но, идея должна быть понятна.

В нашем случае у нас нет никаких БД, но чтобы об этом всём не думать, нужно везде использовать стандартный сериализатор.

## Naming

### Локальные переменные должны быть с маленькой буквы

Неправильно:
```golang
Authors := make([]string)
```

Правильно:
```golang
authors := make([]string)
```

### Аббревиатуры должны писаться в single-case

В том числе и конкатенации аббревиатур.

https://github.com/golang/go/wiki/CodeReviewComments#initialisms


Неправильно:
```
JsonLines
Sha1
printCsv
XMLHttpRequest
```

Правильно:
```
JSONLines
SHA1
printCSV
XMLHTTPRequest
```

### Имена файлов должны быть в snake case

Неправильно:
```
listFiles.go
```

Правильно:
```
list_files.go
```

## Работа с language_extensions.json

Некоторые решения читают файл с экстеншионами по относительному пути.

Неправильно:
```golang
mappingFile, err := os.Open("../../configs/language_extensions.json")
```

В таком случае вы не сможете распространять свою утилиту.
Она будет работать только если рядом лежит json файл.

Утилита должна работать вне зависимости от директории, в которой она была запущена.

Правильно:
```
//go:embed language_extensions.json
var file []byte
```

Нужно вкомпилить все зависимости в утилиту, например, с помощью embed.

## Работа с `os.Exec`

Нельзя делать `os.Chdir`.

После работы утилиты пользователь ожидает, что он останется в той же директории, в которой запускал утилиту.

Неправильно:
```golang
err := os.Chdir(repository)
cmd := exec.Command("git", "blame", "--porcelain", revision, "--", file)
```

Правильно:
```golang
cmd := exec.Command("git", "blame", "--porcelain", revision, "--", file)
cmd.Dir = repository
```

## Стиль

### NIT Используйте общие var и const декларации

Для однородных переменных и констант нет необходимости писать `var` перед каждой строкой

Вместо
```golang
var flagRepo = flag.String("repository", ".", "repo")
var flagRev = flag.String("revision", "HEAD", "revision")
```

Можно написать
```golang
var (
	flagRepo = flag.String("repository", ".", "repo")
	flagRev = flag.String("revision", "HEAD", "revision")
)
```

## Вывод результатов

Вместо
```golang
os.Stdout.Write(bytes)
fmt.Println()
```

или

```golang
os.Stdout.Write(fmt.Sprintf("%s\n"), string(bytes))
```

Используйте более читаемое однострочное
```golang
fmt.Println(string(bytes))
```

## Работа с горутинами

Часто встречается подобная конструкция с бесконтрольной параллелизацией:

```golang
ch := make(chan struct{}, len(files))
for _, file := range files {
	go blame(file, ch)
}
```

На большом репозитории одновременно будет запущено неопределённое количество горутин и подпроцессов.

Во-первых, каждая горутина требует сколько-то килобайт на стек и потенциально может закончиться память.

Во-вторых, в OS есть ограничение на количество процессов + каждый процесс потребляет сколько-то системных ресурсов (ram, cpu) и суммарное потребление может оказаться неопределённо большим.

Вместо этого стоит явно ограничить количество одновременно запущенных горутин-воркеров, обрабатывающих blame, константой с небольшим дефолтом (8), либо дополнительно можно вынести степень параллелизации за флаг.

Прочитайте [пост в go блоге](https://go.dev/blog/pipelines) про типичные pipeline паттерны.

Можно сделать одного producer'а workload'а, который будет писать файлы для обработки в канал фиксированного размера, а `n` worker'ов будут читать из этого общего канала и обрабатывать файлы.

## User experience

### Не нужно паниковать на невалидном input'е

Пользователь утилиты не очень хочет видеть stacktrace, если он неправильно указал путь до репозитория.

Вместо
```golang
files, err := git.ListFiles(*repository, *revision)
if err != nil {
	panic(err)
}
```

Стоит написать человеческую ошибку:
```golang
files, err := git.ListFiles(*repository, *revision)
if err != nil {
	_, _ = fmt.Fprintf(os.Stderr, "File listing failed: %s", err)
	os.Exit(1)
}
```

### NIT Не печатать ненужную информацию в логах

При использовании стандартного пакета log в log message добавится время.

Возможно, эта информация пользователю не очень нужна.

Вместо
```golang
files, err := git.ListFiles(*repository, *revision)
if err != nil {
	log.Fatalf("File listing failed: %s", err)
}
```

который напечатает что-то вроде:
```
2023/03/12 21:26:15 File listing failed: path not found
```

Можно написать
```golang
files, err := git.ListFiles(*repository, *revision)
if err != nil {
	_, _ = fmt.Fprintf(os.Stderr, "File listing failed: %s", err)
	os.Exit(1)
}
```

Либо можно сконфигурировать log так, чтобы он не печатал лишнего.

Если вы печатает прогресс по мере обработки файлов, то время в логе может быть уместно. Ещё лучше прогрессбар рисовать.

### Проверять невалидные значения формат флагов в начале работы, а не после

На большом репозитории программа может работать долго и после этого увидеть опечатку в `--format cvs` вдвойне обидно.

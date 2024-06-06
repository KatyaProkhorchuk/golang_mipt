package internal

type Language struct {
	Name            string
	ApplicationArea string `json:"type"`
	Extensions      []string
}

type Statisctics struct {
	Lines   int
	Commits int
	Files   int
}

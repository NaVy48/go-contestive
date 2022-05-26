package entity

// Problem db model
type Problem struct {
	Entity
	AuthorID    int64
	ExternalURL string
	Name        string
	Revisions   []ProblemRevision
}

func (p Problem) ActiveRevision() (ProblemRevision, bool) {
	for _, v := range p.Revisions {
		if !v.Outdated {
			return v, true
		}
	}

	return ProblemRevision{}, false
}

// Problem db model
type ProblemRevision struct {
	Entity
	AuthorID       int64
	ProblemID      int64
	Revision       int
	Title          string
	MemoryLimit    int
	TimeLimit      int
	StatementHtml  string
	StatementPdf   []byte
	PackageArchive []byte
	Outdated       bool
}

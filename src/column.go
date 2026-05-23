package src


// структура колонки
type Column struct {
	ID      int
	Title   string
	Tasks   map[int]struct{}
	BoardID int
}



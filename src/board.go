package src

import (
	"errors"
)

// структура доски
type Board struct {
	ID      int
	Title   string
	Columns map[int]struct{}
}

// принадлежат ли переданные колонки конкретной доске  
func (b Board) validateBelongBoard(fromColumnID, toColumnID int) error {
	if _, ok := b.Columns[fromColumnID]; !ok {
		return errors.New("колонка не принадлежит этой доске")
	}

	if _, ok := b.Columns[toColumnID]; !ok {
		return errors.New("колонка не принадлежит этой доске")
	}

	return nil
}


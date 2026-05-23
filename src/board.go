package src

import (
	"errors"
)

// структура доски
type Board struct {
	ID      int
	Title   string
	Columns map[int]*Column
}


// метод получения колонки
func (b *Board) GetColumn(columnID int) (*Column, error) {
	if _, ok := b.Columns[columnID]; !ok {
		err := errors.New("такой колонки не существует")
		return nil, err
	}
	return b.Columns[columnID], nil
}

// метод удаления колонки
func (b *Board) RemoveColumn(columnID int) error {
	if _, ok := b.Columns[columnID]; !ok {
		return errors.New("такой колонки не существует")
	}
	delete(b.Columns, columnID)
	return nil
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

// валидация имени колонки  
func (b Board) validateColumnTitle(title string, excludeID int) error {
    for id, col := range b.Columns {
        if excludeID >= 0 && id == excludeID {
            continue
        }
        if col.Title == title {
            return errors.New("колонка с таким именем уже есть в доске")
        }
    }
    if len(title) < 1 {
        return errors.New("имя колонки не может быть пустым")
    }
    return nil
}

// получение (наибольшего) id среди колонок  
func (b Board) getLastColumnID() int {
	id := 0
	for k, _ := range b.Columns {
		if k > id {
			id = k
		}
	}
	return id
}


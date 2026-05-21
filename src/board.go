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

// хранилище досок для удобства
var BoardMap map[int]*Board = make(map[int]*Board)

// метод создания новой дсоки
func NewBoard(title string) (*Board, error) {
	if err := validateBoardTitle(title); err != nil {
		return nil, err
	}

	//если в хранилище есть доски
	if len(BoardMap) > 0 {
		boardPointer := &Board{
			ID:      getLastBoardID() + 1,
			Title:   title,
			Columns: make(map[int]*Column),
		}
		//добавляем доску в хранилище по ее же id
		BoardMap[boardPointer.ID] = boardPointer
		return boardPointer, nil
	} else {
		//если хранилище пустое, у первой доски id = 0
		boardPointer := &Board{
			ID:      0,
			Title:   title,
			Columns: make(map[int]*Column),
		}
		//добавляем доску в хранилище по ее же id
		BoardMap[boardPointer.ID] = boardPointer
		return boardPointer, nil
	}
}

// метод создания и добавления новой колонки в доску
func (b *Board) AddColumn(title string) (*Column, error) {
	if err := b.validateColumnTitle(title); err != nil {
		return nil, err
	}

	if len(b.Columns) > 0 {
		columnPointer := &Column{
			ID:      b.getLastColumnID() + 1,
			Title:   title,
			Tasks:   make(map[int]*Task),
			BoardID: b.ID,
		}
		b.Columns[columnPointer.ID] = columnPointer
		return columnPointer, nil
	} else {
		columnPointer := &Column{
			ID:      0,
			Title:   title,
			Tasks:   make(map[int]*Task),
			BoardID: b.ID,
		}
		b.Columns[columnPointer.ID] = columnPointer
		return columnPointer, nil
	}
}

// метод получения доски
func GetBoard(boardID int) (*Board, error) {
	if _, ok := BoardMap[boardID]; !ok {
		err := errors.New("такой доски не существует")
		return nil, err
	}
	return BoardMap[boardID], nil
}

// метод обновления имени доски
func (b *Board) UpdateBoardTitle(title string) (*Board, error) {
	if err := validateBoardTitle(title); err != nil {
		return nil, err
	}

	b.Title = title
	return b, nil
}

// метод удаления доски
func (b *Board) RemoveBoard() error {
	if _, ok := BoardMap[b.ID]; !ok {
		return errors.New("такой доски не существует")
	}
	delete(BoardMap, b.ID)
	return nil
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

// метод перемещения задачи из колонки в колонку
func (b *Board) MoveTask(taskID, fromColumnID, toColumnID int) error {
	if err := b.validateBelongBoard(fromColumnID, toColumnID); err != nil {
		return errors.New("передана неверная колонка")
	}

	necessaryTask, ok := b.Columns[fromColumnID].Tasks[taskID]
	if !ok {
		return errors.New("задача не найдена")
	}
	b.Columns[fromColumnID].RemoveTask(taskID)
	lastTaskId := b.Columns[toColumnID].getLastTaskID()
	necessaryTask.ID = lastTaskId + 1
	necessaryTask.ColumnID = toColumnID
	b.Columns[toColumnID].Tasks[lastTaskId+1] = necessaryTask
	return nil
}

// принадлежат ли переданные колонки кокретной доске
func (b Board) validateBelongBoard(fromColumnID, toColumnID int) error {
	if _, ok := b.Columns[fromColumnID]; !ok {
		return errors.New("колонка не принадлежит этой доске")
	}

	if _, ok := b.Columns[toColumnID]; !ok {
		return errors.New("колонка не принадлежит этой доске")
	}

	return nil
}

// валидация имени доски
func validateBoardTitle(title string) error {
	for _, v := range BoardMap {
		if v.Title == title {
			return errors.New("доска с таким именем уже существует")
		}
	}

	if len(title) < 1 {
		return errors.New("имя доски не может быть пустым")
	} else {
		return nil
	}
}

// валидация имени колонки
func (b Board) validateColumnTitle(title string) error {
	for _, v := range b.Columns {
		if v.Title == title {
			return errors.New("колонка с таким именем уже есть в доске")
		}
	}

	if len(title) < 1 {
		return errors.New("имя колонки не может быть пустым")
	} else {
		return nil
	}
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

// получение последнего (наибольшего) id среди досок
func getLastBoardID() int {
	id := 0
	for k, _ := range BoardMap {
		if k > id {
			id = k
		}
	}
	return id
}

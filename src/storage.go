package src

import (
	"errors"
	"time"
)


type Storage struct {
	Boards map[int]*Board
	Columns map[int]*Column
	Tasks map[int]*Task
}

func NewStorage() *Storage {
	return &Storage{
		Boards: make(map[int]*Board),
		Columns: make(map[int]*Column),
		Tasks: make(map[int]*Task),
	}
}


// метод создания новой дсоки
func (s *Storage) NewBoard(title string) (*Board, error) {
	if err := s.validateBoardTitle(title, -1); err != nil {
		return nil, err
	}


	boardID := s.getLastBoardID() + 1
	boardPointer := &Board{
			ID:      boardID,
			Title:   title,
			Columns: make(map[int]struct{}),
	}
	s.Boards[boardID] = boardPointer
	return boardPointer, nil
}

// метод получения доски
func (s *Storage) GetBoard(boardID int) (*Board, error) {
	if _, ok := s.Boards[boardID]; !ok {
		err := errors.New("такой доски не существует")
		return nil, err
	}
	return s.Boards[boardID], nil
}

// метод обновления имени доски
func (s *Storage) UpdateBoardTitle(boardID int, title string) (*Board, error) {
	board, ok := s.Boards[boardID]
	if !ok {
		return nil, errors.New("такой доски не существует")
	}

	if err := s.validateBoardTitle(title, boardID); err != nil {
		return nil, err
	}

	board.Title = title
	return board, nil
}

// метод удаления доски
func (s *Storage) RemoveBoard(boardID int) error {
	if _, ok := s.Boards[boardID]; !ok {
		return errors.New("такой доски не существует")
	}

	board := s.Boards[boardID]
	for colID := range board.Columns {
		if col, ok := s.Columns[colID]; ok {
			for taskID := range col.Tasks {
				delete(s.Tasks, taskID)
			}
		}
		delete(s.Columns, colID)
		delete(board.Columns, colID)
	}
	delete(s.Boards, boardID)
	return nil
}

func (s *Storage) AddColumn(boardID int, title string) (*Column, error) {
	board, ok := s.Boards[boardID]
	if !ok {
		return nil, errors.New("такой доски не существует")
	}
	
	if err := s.validateColumnTitle(board, title, -1); err != nil {
		return nil, err
	}

	colID := s.getLastColumnID() + 1
	column := &Column{
		ID:      colID,
		Title:   title,
		Tasks:   make(map[int]struct{}),
		BoardID: board.ID,
	}
	board.Columns[colID] = struct{}{}
	s.Columns[colID] = column
	return column, nil
}

// метод получения колонки
func (s *Storage) GetColumn(boardID, columnID int) (*Column, error) {
	board, ok := s.Boards[boardID]
	if !ok {
		return nil, errors.New("такой доски не существует")
	}

	column, err := s.validateColumnExist(board, columnID)
	if err != nil {
		return nil, err
	}
	
	return column, nil
}

// метод удаления колонки
func (s *Storage) RemoveColumn(boardID, columnID int) error {
    board, ok := s.Boards[boardID]
    if !ok {
        return errors.New("такой доски не существует")
    }

    column, err := s.validateColumnExist(board, columnID)
	if err != nil {
		return err
	}

    for taskID := range column.Tasks {
        delete(s.Tasks, taskID)
    }

    delete(s.Columns, columnID)
    delete(board.Columns, columnID)
    return nil
}

// метод обновления имени колонки
func (s *Storage) UpdateColumnTitle(boardID, columnID int, title string) (*Column, error) {
	board, ok := s.Boards[boardID]
	if !ok {
		return nil, errors.New("такой доски не существует")
	}

	column, err := s.validateColumnExist(board, columnID)
	if err != nil {
		return nil, err
	}
	
	if err := s.validateColumnTitle(board, title, columnID); err != nil {
		return nil, err
	}

	column.Title = title
	return column, nil
}

// метод создания и добавления новой задачи в колонку
func (s *Storage) AddTask(boardID, columnID int, title, description string) (*Task, error) {
	board, ok := s.Boards[boardID]
	if !ok {
		return nil, errors.New("такой доски не существует")
	}

	column, err := s.validateColumnExist(board, columnID)
	if err != nil {
		return nil, err
	}
	
	if err := s.validateTaskTitle(column, title, -1); err != nil {
		return nil, err
	}

	taskID := s.getLastTaskID() + 1
	task := &Task{
		ID:          taskID,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		ColumnID:    column.ID,
	}
	column.Tasks[taskID] = struct{}{}
	s.Tasks[taskID] = task
	return task, nil
}

// метод получения задачи
func (s *Storage) GetTask(boardID, columnID, taskID int) (*Task, error) {
	board, ok := s.Boards[boardID]
	if !ok {
		return nil, errors.New("такой доски не существует")
	}

	column, err := s.validateColumnExist(board, columnID)
	if err != nil {
		return nil, err
	}
	
	task, err := s.validateTaskExist(column, taskID)
	if err != nil {
		return nil, err
	}
	
	return task, nil
}

// метод перемещения задачи из колонки в колонку
func (s *Storage) MoveTask(boardID, taskID, fromColumnID, toColumnID int) error {
	if fromColumnID == toColumnID {
		return nil
	}
	
	board, ok := s.Boards[boardID]
	if !ok {
		return errors.New("такой доски не существует")
	}

	fromCol, err := s.validateColumnExist(board, fromColumnID)
	if err != nil {
		return err
	}
	
	toCol, err := s.validateColumnExist(board, toColumnID)
	if err != nil {
		return err
	}

	_, err = s.validateTaskExist(fromCol, taskID)
	if err != nil {
		return err
	}
	
	delete(fromCol.Tasks, taskID)
	toCol.Tasks[taskID] = struct{}{}
	s.Tasks[taskID].ColumnID = toColumnID
	return nil
}

// метод удаления задачи
func (s *Storage) RemoveTask(boardID, columnID, taskID int) error {
	board, ok := s.Boards[boardID]
	if !ok {
		return errors.New("такой доски не существует")
	}

	column, err := s.validateColumnExist(board, columnID)
	if err != nil {
		return err
	}
	
	if _, err = s.validateTaskExist(column, taskID); err != nil {
		return err
	}
	
	delete(column.Tasks, taskID)
	delete(s.Tasks, taskID)
	return nil
}

//метод обновления имени задачи
func (s *Storage) UpdateTaskTitle(boardID, columnID, taskID int, title string) (*Task, error) {
	board, ok := s.Boards[boardID]
	if !ok {
		return nil, errors.New("такой доски не существует")
	}

	column, err := s.validateColumnExist(board, columnID)
	if err != nil {
		return nil, err
	}
	
	task, err := s.validateTaskExist(column, taskID)
	if err != nil {
		return nil, err
	}
	
	if err := s.validateTaskTitle(column, title, taskID); err != nil {
		return nil, err
	}

	task.Title = title
	return task, nil
}
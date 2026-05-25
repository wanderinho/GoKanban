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


	var id int
	if len(s.Boards) > 0 {
		id = s.getLastBoardID() + 1
	} else {
		id = 0
	}
	boardPointer := &Board{
			ID:      id,
			Title:   title,
			Columns: make(map[int]struct{}),
	}
	s.Boards[id] = boardPointer
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
		delete(s.Columns, colID)
		if col, ok := s.Columns[colID]; ok {
			for taskID := range col.Tasks {
				delete(s.Tasks, taskID)
			}
		}
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

	var colID int
	if len(s.Columns) > 0 {
		colID = s.getLastColumnID() + 1
	} else {
		colID = 0
	}
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

	if _, ok := s.Columns[columnID]; !ok {
		return nil, errors.New("такой колонки не существует во глобальном хранилище")
	}

	_, ok = board.Columns[columnID]
	if !ok {
		return nil, errors.New("такой колонки не существует в доске")
	}
	
	return s.Columns[columnID], nil
}

// метод удаления колонки
func (s *Storage) RemoveColumn(boardID, columnID int) error {
    board, ok := s.Boards[boardID]
    if !ok {
        return errors.New("такой доски не существует")
    }

    if _, ok := board.Columns[columnID]; !ok {
        return errors.New("такой колонки не существует в доске")
    }

    column, ok := s.Columns[columnID]
    if !ok {
        return errors.New("такой колонки не существует в глобальном хранилище")
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

	if _, ok = board.Columns[columnID]; !ok {
		return nil, errors.New("такой колонки не существует в доске")
	}

	column, ok := s.Columns[columnID]
	if !ok {
		return nil, errors.New("такой колонки не существует во глобальном хранилище")
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

	if _, ok = board.Columns[columnID]; !ok {
		return nil, errors.New("такой колонки не существует в доске")
	}

	column, ok := s.Columns[columnID]
	if !ok {
		return nil, errors.New("такой колонки не существует во глобальном хранилище")
	}
	
	if err := s.validateTaskTitle(column, title); err != nil {
		return nil, err
	}

	var taskID int
	if len(s.Tasks) > 0 {
		taskID = s.getLastTaskID() + 1
	} else {
		taskID = 0
	}
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

	if _, ok = board.Columns[columnID]; !ok {
		return nil, errors.New("такой колонки не существует в доске")
	}

	column, ok := s.Columns[columnID]
	if !ok {
		return nil, errors.New("такой колонки не существует в глобальном хранилище")
	}
	
	if _, ok := column.Tasks[taskID]; !ok {
		return nil, errors.New("такой задачи не существует в колонке")
	}

	if _, ok := s.Tasks[taskID]; !ok {
		return nil, errors.New("такой задачи не существует во глобальном хранилище")
	}
	
	return s.Tasks[taskID], nil
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
	
	if err := board.validateBelongBoard(fromColumnID, toColumnID); err != nil {
		return errors.New("передана неверная колонка")
	}

	fromCol, ok := s.Columns[fromColumnID]
	if !ok {
		return errors.New("колонка-источник не найдена в глобальном хранилище")
	}
	toCol, ok := s.Columns[toColumnID]
	if !ok {
		return errors.New("колонка-назначение не найдена в глобальном хранилище")
	}

	if _, ok := fromCol.Tasks[taskID]; !ok {
		return errors.New("такой задачи не существует в колонке")
	}

	if _, ok := s.Tasks[taskID]; !ok {
		return errors.New("такой задачи не существует во глобальном хранилище")
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

	_, ok = board.Columns[columnID]
	if !ok {
		return errors.New("такой колонки не существует в доске")
	}

	column, ok := s.Columns[columnID]
	if !ok {
		return errors.New("такой колонки не существует в глобальном хранилище")
	}
	
	if _, ok := column.Tasks[taskID]; !ok {
		return errors.New("такой задачи не существует в колонке")
	}

	if _, ok := s.Tasks[taskID]; !ok {
		return errors.New("такой задачи не существует во глобальном хранилище")
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

	if _, ok = board.Columns[columnID]; !ok {
		return nil, errors.New("такой колонки не существует в доске")
	}

	column, ok := s.Columns[columnID]
	if !ok {
		return nil, errors.New("такой колонки не существует в глобальном хранилище")
	}
	
	if _, ok := column.Tasks[taskID]; !ok {
		return nil, errors.New("такой задачи не существует в колонке")
	}

	task, ok := s.Tasks[taskID]
	if !ok {
		return nil, errors.New("такой задачи не существует во глобальном хранилище")
	}
	
	if err := s.validateTaskTitle(column, title); err != nil {
		return nil, err
	}

	task.Title = title
	return task, nil
}

// валидация имени доски
func (s *Storage) validateBoardTitle(title string, excludeID int) error {
    for id, b := range s.Boards {
        if excludeID >= 0 && id == excludeID {
            continue
        }
        if b.Title == title {
            return errors.New("доска с таким именем уже существует")
        }
    }
    if len(title) < 1 {
        return errors.New("имя доски не может быть пустым")
    }
    return nil
}

// валидация имени колонки  
func (s *Storage) validateColumnTitle(board *Board, title string, excludeID int) error {
    for colID, _ := range board.Columns {
    	column, ok := s.Columns[colID]
     	if !ok {
      		return errors.New("внутренняя ошибка: колонка из доски не найдена в глобальном хранилище")
      	}
        if excludeID >= 0 && colID == excludeID {
            continue
        }
        if column.Title == title {
            return errors.New("колонка с таким именем уже есть в доске")
        }
    }
    if len(title) < 1 {
        return errors.New("имя колонки не может быть пустым")
    }
    return nil
}

// валидация имени задачи
func (s *Storage) validateTaskTitle(column *Column, title string) error {
    for taskID, _ := range column.Tasks {
        task, ok := s.Tasks[taskID]
        if !ok {
            return errors.New("внутренняя ошибка: задача из колонки не найдена в глобальном хранилище")
        }
        if task.Title == title {
            return errors.New("задача с таким именем уже есть в колонке")
        }
    }
    if len(title) < 1 {
        return errors.New("имя задачи не может быть пустым")
    }
    return nil
}

// получение последнего (наибольшего) id среди досок
func (s *Storage) getLastBoardID() int {
	id := 0
	for k, _ := range s.Boards {
		if k > id {
			id = k
		}
	}
	return id
}

// получение (наибольшего) id среди колонок  
func (s *Storage) getLastColumnID() int {
	id := 0
	for k, _ := range s.Columns {
		if k > id {
			id = k
		}
	}
	return id
}

// получение последнего (наибольшего) id среди задач
func (s *Storage) getLastTaskID() int {
	id := 0
	for k, _ := range s.Tasks {
		if k > id {
			id = k
		}
	}
	return id
}
package src

import (
	"errors"
	"time"
)


type Storage struct {
	Boards map[int]*Board
	Tasks map[int]*Task
}

// метод создания новой дсоки
func (s *Storage) NewBoard(title string) (*Board, error) {
	if err := s.validateBoardTitle(title, -1); err != nil {
		return nil, err
	}

	//если в хранилище есть доски
	if len(s.Boards) > 0 {
		boardPointer := &Board{
			ID:      s.getLastBoardID() + 1,
			Title:   title,
			Columns: make(map[int]*Column),
		}
		//добавляем доску в хранилище по ее же id
		s.Boards[boardPointer.ID] = boardPointer
		return boardPointer, nil
	} else {
		//если хранилище пустое, у первой доски id = 0
		boardPointer := &Board{
			ID:      0,
			Title:   title,
			Columns: make(map[int]*Column),
		}
		//добавляем доску в хранилище по ее же id
		s.Boards[boardPointer.ID] = boardPointer
		return boardPointer, nil
	}
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
	delete(s.Boards, boardID)
	return nil
}

func (s *Storage) AddColumn(boardID int, title string) (*Column, error) {
	board, ok := s.Boards[boardID]
	if !ok {
		return nil, errors.New("такой доски не существует")
	}
	
	if err := board.validateColumnTitle(title, -1); err != nil {
		return nil, err
	}

	if len(board.Columns) > 0 {
		columnPointer := &Column{
			ID:      board.getLastColumnID() + 1,
			Title:   title,
			Tasks:   make(map[int]struct{}),
			BoardID: board.ID,
		}
		board.Columns[columnPointer.ID] = columnPointer
		return columnPointer, nil
	} else {
		columnPointer := &Column{
			ID:      0,
			Title:   title,
			Tasks:   make(map[int]struct{}),
			BoardID: board.ID,
		}
		board.Columns[columnPointer.ID] = columnPointer
		return columnPointer, nil
	}
}

// метод обновления имени колонки
func (s *Storage) UpdateColumnTitle(boardID, columnID int, title string) (*Column, error) {
	board, ok := s.Boards[boardID]
	if !ok {
		return nil, errors.New("такой доски не существует")
	}

	column, ok := board.Columns[columnID]
	if !ok {
		return nil, errors.New("такой колонки не существует")
	}
	
	if err := board.validateColumnTitle(title, columnID); err != nil {
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

	column, ok := board.Columns[columnID]
	if !ok {
		return nil, errors.New("такой колонки не существует")
	}
	
	if err := s.validateTaskTitle(column, title); err != nil {
		return nil, err
	}

	if len(s.Tasks) > 0 {
		taskPointer := &Task{
			ID:          s.getLastTaskID() + 1,
			Title:       title,
			Description: description,
			CreatedAt:   time.Now(),
			ColumnID:    column.ID,
		}
		column.Tasks[taskPointer.ID] = struct{}{}
		s.Tasks[taskPointer.ID] = taskPointer
		return taskPointer, nil
	} else {
		taskPointer := &Task{
			ID:          0,
			Title:       title,
			Description: description,
			CreatedAt:   time.Now(),
			ColumnID:    column.ID,
		}
		column.Tasks[taskPointer.ID] = struct{}{}
		s.Tasks[taskPointer.ID] = taskPointer
		return taskPointer, nil
	}
}

// метод получения задачи
func (s *Storage) GetTask(boardID, columnID, taskID int) (*Task, error) {
	board, ok := s.Boards[boardID]
	if !ok {
		return nil, errors.New("такой доски не существует")
	}

	column, ok := board.Columns[columnID]
	if !ok {
		return nil, errors.New("такой колонки не существует")
	}
	
	if _, ok := column.Tasks[taskID]; !ok {
		return nil, errors.New("такой задачи не существует")
	}

	if _, ok := s.Tasks[taskID]; !ok {
		return nil, errors.New("такой задачи не существует")
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

	if _, ok := board.Columns[fromColumnID].Tasks[taskID]; !ok {
		return errors.New("такой задачи не существует")
	}

	if _, ok := s.Tasks[taskID]; !ok {
		return errors.New("такой задачи не существует")
	}
	delete(board.Columns[fromColumnID].Tasks, taskID)
	board.Columns[toColumnID].Tasks[taskID] = struct{}{}
	s.Tasks[taskID].ColumnID = toColumnID
	return nil
}

// метод удаления задачи
func (s *Storage) RemoveTask(boardID, columnID, taskID int) error {
	board, ok := s.Boards[boardID]
	if !ok {
		return errors.New("такой доски не существует")
	}

	column, ok := board.Columns[columnID]
	if !ok {
		return errors.New("такой колонки не существует")
	}
	
	if _, ok := column.Tasks[taskID]; !ok {
		return errors.New("такой задачи нет в этой колонке")
	}

	if _, ok := s.Tasks[taskID]; !ok {
		return errors.New("такой задачи не существует")
	}
	
	delete(column.Tasks, taskID)
	delete(s.Tasks, taskID)
	return nil
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

// валидация имени задачи
func (s *Storage) validateTaskTitle(column *Column, title string) error {
    for k, _ := range column.Tasks {
        task, ok := s.Tasks[k]
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
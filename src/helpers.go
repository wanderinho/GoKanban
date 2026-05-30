package src

import "errors"


func (s *Storage) validateColumnExist(board *Board, columnID int) (*Column, error) {
	column, ok := s.Columns[columnID]
	if !ok {
		return nil, errors.New("такой колонки не существует в глобальном хранилище")
	}

	_, ok = board.Columns[columnID]
	if !ok {
		return nil, errors.New("такой колонки не существует в доске")
	}

	return column, nil
}


func (s *Storage) validateTaskExist(column *Column, taskID int) (*Task, error) {
	if _, ok := column.Tasks[taskID]; !ok {
		return nil, errors.New("такой задачи не существует в колонке")
	}

	task, ok := s.Tasks[taskID]
	if !ok {
		return nil, errors.New("такой задачи не существует во глобальном хранилище")
	}

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
func (s *Storage) validateTaskTitle(column *Column, title string, excludeID int) error {
    for taskID := range column.Tasks {
        task, ok := s.Tasks[taskID]
        if !ok {
            return errors.New("внутренняя ошибка: задача из колонки не найдена в глобальном хранилище")
        }
        if excludeID >= 0 && taskID == excludeID {
            continue
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
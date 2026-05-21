package src

import (
	"errors"
	"time"
)

// структура колонки
type Column struct {
	ID      int
	Title   string
	Tasks   map[int]*Task
	BoardID int
}

// метод создания и добавления новой задачи в колонку
func (c *Column) AddTask(title, description string) (*Task, error) {
	if err := c.validateTaskTitle(title); err != nil {
		return nil, err
	}

	if len(c.Tasks) > 0 {
		taskPointer := &Task{
			ID:          c.getLastTaskID() + 1,
			Title:       title,
			Description: description,
			CreatedAt:   time.Now(),
			ColumnID:    c.ID,
		}
		c.Tasks[taskPointer.ID] = taskPointer
		return taskPointer, nil
	} else {
		taskPointer := &Task{
			ID:          0,
			Title:       title,
			Description: description,
			CreatedAt:   time.Now(),
			ColumnID:    c.ID,
		}
		c.Tasks[taskPointer.ID] = taskPointer
		return taskPointer, nil
	}
}

// метод получения задачи
func (c *Column) GetTask(taskID int) (*Task, error) {
	if _, ok := c.Tasks[taskID]; !ok {
		err := errors.New("такой задачи не существует")
		return nil, err
	}
	return c.Tasks[taskID], nil
}

// метод удаления задачи
func (c *Column) RemoveTask(taskID int) error {
	if _, ok := c.Tasks[taskID]; !ok {
		return errors.New("такой задачи не существует")
	}
	delete(c.Tasks, taskID)
	return nil
}

// метод обновления имени задачи
func (c *Column) UpdateColumnTitle(title string) (*Column, error) {
	if err := BoardMap[c.BoardID].validateColumnTitle(title); err != nil {
		return nil, err
	}

	c.Title = title
	return c, nil
}

// валидация имени задачи
func (c Column) validateTaskTitle(title string) error {
	for _, v := range c.Tasks {
		if v.Title == title {
			return errors.New("задача с таким именем уже есть в колонке")
		}
	}

	if len(title) < 1 {
		return errors.New("имя задачи не может быть пустым")
	} else {
		return nil
	}
}

// получение последнего (наибольшего) id среди задач
func (c Column) getLastTaskID() int {
	id := 0
	for k, _ := range c.Tasks {
		if k > id {
			id = k
		}
	}
	return id
}

package src

import (
	"time"
	"errors"
)

//структура колонки
type Column struct {
	ID      int
	title   string
	tasks   map[int]*Task
	boardID int
}

//метод создания и добавления новой задачи в колонку
func (c *Column) AddTask(title, description string) (*Task, error) {
	if err := c.validateTaskTitle(title); err != nil {
		return nil, err
	}
	
	if len(c.tasks) > 0 {
		taskPointer := &Task{
			ID: c.getLastTaskID() + 1, 
			title: title,
			description: description,
			createdAt: time.Now(),
			columnID: c.ID,
		}
		c.tasks[taskPointer.ID] = taskPointer
		return taskPointer, nil
	} else {
		taskPointer := &Task{
			ID: 0, 
			title: title,
			description: description,
			createdAt: time.Now(),
			columnID: c.ID,
		}
		c.tasks[taskPointer.ID] = taskPointer
		return taskPointer, nil
	}
}

//метод получения задачи
func (c *Column) GetTask(taskID int) (*Task, error) {
	if _, ok := c.tasks[taskID]; !ok {
		err := errors.New("такой задачи не существует")
      	return nil, err
    }
    return c.tasks[taskID], nil
}

//метод удаления задачи
func (c *Column) RemoveTask(taskID int) error {
	if _, ok := c.tasks[taskID]; !ok {
      	return errors.New("такой задачи не существует")
    }
    delete(c.tasks, taskID)
    return nil
}

//метод обновления имени задачи
func (c *Column) UpdateColumnTitle(title string) (*Column, error) {
	if err := boardMap[c.boardID].validateColumnTitle(title); err != nil {
		return nil, err
	}

	c.title = title
	return c, nil
}

//валидация имени задачи
func (c Column) validateTaskTitle(title string) error {
	for _, v := range c.tasks {
		if v.title == title {
			return errors.New("задача с таким именем уже есть в колонке")
		}
	}
	
	if len(title) < 1 {
		return errors.New("имя задачи не может быть пустым")
	} else {
		return nil
	}
}

//получение последнего (наибольшего) id среди задач
func (c Column) getLastTaskID() int {
	id := 0
	for k, _ := range c.tasks {
		if k > id {
			id = k
		}
	}
	return id
}





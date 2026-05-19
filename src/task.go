package src

import (
	"time"
)


//структура таски
type Task struct {
	ID int
	title       string
	description string
	createdAt time.Time
	columnID int
}

//метод обновления имени задачи
/*func (t *Task) UpdateTaskTitle(title string) (*Task, error) {
	без ссылки на родителя невозможно реализовать валидацию заголовка
	она будет реализована после подключения бд в проект
} */

//метод обновления описания задачи
func (t *Task) UpdateTaskDescription(description string) *Task {
	t.description = description
	return t
}







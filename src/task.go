package src

import (
	"time"
)

// структура таски
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ColumnID    int       `json:"column_id"`
}

//метод обновления имени задачи
/*func (t *Task) UpdateTaskTitle(title string) (*Task, error) {
	без ссылки на родителя невозможно реализовать валидацию заголовка
	она будет реализована после подключения бд в проект
} */

// метод обновления описания задачи
func (t *Task) UpdateTaskDescription(description string) *Task {
	t.Description = description
	return t
}

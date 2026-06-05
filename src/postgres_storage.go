package src

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	db *pgxpool.Pool 
}

type StorageInterface interface {
	NewBoard(ctx context.Context, title string) (*Board, error)
    GetBoard(ctx context.Context, boardID int) (*Board, error)
    RemoveBoard(ctx context.Context, boardID int) error
    UpdateBoardTitle(ctx context.Context, boardID int, title string) (*Board, error)
    AddColumn(ctx context.Context, boardID int, title string) (*Column, error)
    GetColumn(ctx context.Context, boardID, columnID int) (*Column, error)
    RemoveColumn(ctx context.Context, boardID, columnID int) error
    UpdateColumnTitle(ctx context.Context, boardID, columnID int, title string) (*Column, error)
    AddTask(ctx context.Context, boardID, columnID int, title, description string) (*Task, error)
    GetTask(ctx context.Context, boardID, columnID, taskID int) (*Task, error)
    RemoveTask(ctx context.Context, boardID, columnID, taskID int) error
    MoveTask(ctx context.Context, boardID, taskID, fromColumnID, toColumnID int) error
    UpdateTaskTitle(ctx context.Context, boardID, columnID, taskID int, title string) (*Task, error)
    UpdateTaskDescription(ctx context.Context, boardID, columnID, taskID int, description string) (*Task, error)
    GetColumnsByBoard(ctx context.Context, boardID int) ([]*Column, error)
    GetTasksByColumn(ctx context.Context, columnID int) ([]*Task, error)
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}
	return &PostgresStorage{db: pool}, nil
}

func (ps *PostgresStorage) Close() {
	ps.db.Close()
}


func (ps *PostgresStorage) NewBoard(ctx context.Context, title string) (*Board, error) {
	if err := validateTitle(title, "доски"); err != nil {
		return nil, err
	}
	var id int
	err := ps.db.QueryRow(ctx, "INSERT INTO boards (title) VALUES ($1) RETURNING id", title).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании доски: %w", err)
	}
	return &Board{
		ID:    id,
		Title: title,
	}, nil
}

func (ps *PostgresStorage) GetBoard(ctx context.Context, boardID int) (*Board, error) {
	var id int
	var title string
	err := ps.db.QueryRow(ctx, "SELECT id, title FROM boards WHERE id = $1", boardID).Scan(&id, &title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("доска не найдена: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка при получении доски: %w", err)
	}
	return &Board{
		ID:    id,
		Title: title,
	}, nil
}

func (ps *PostgresStorage) UpdateBoardTitle(ctx context.Context, boardID int, title string) (*Board, error) {
	if err := validateTitle(title, "доски"); err != nil {
		return nil, err
	}
	var id int
	err := ps.db.QueryRow(ctx, "UPDATE boards SET title = $1 WHERE id = $2 RETURNING id", title, boardID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("доска не найдена: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка при обновлении доски: %w", err)
	}
	return &Board{
		ID:    id,
		Title: title,
	}, nil
}

func (ps *PostgresStorage) RemoveBoard(ctx context.Context, boardID int) error {
	var id int
	err := ps.db.QueryRow(ctx, "DELETE FROM boards WHERE id = $1 RETURNING id", boardID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("доска не найдена: %w", ErrNotFound)
		}
		return fmt.Errorf("ошибка при удалении доски: %w", err)
	}
	return nil
}

func (ps *PostgresStorage) AddColumn(ctx context.Context, boardID int, title string) (*Column, error) {
	if err := validateTitle(title, "колонки"); err != nil {
		return nil, err
	}
	if err := ps.boardExists(ctx, boardID); err != nil {
		return nil, err
	}
	var count int
	err := ps.db.QueryRow(ctx, "SELECT COUNT(*) FROM columns WHERE board_id = $1 AND title = $2", boardID, title).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке колонки: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("колонка с таким именем уже существует: %w", ErrAlreadyExists)
	}
	var colID int
	err = ps.db.QueryRow(ctx, "INSERT INTO columns (title, board_id) VALUES ($1, $2) RETURNING id", title, boardID).Scan(&colID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании колонки: %w", err)
	}
	return &Column{
		ID:      colID,
		Title:   title,
		BoardID: boardID,
	}, nil
}

func (ps *PostgresStorage) GetColumn(ctx context.Context, boardID, columnID int) (*Column, error) {
	if err := ps.boardExists(ctx, boardID); err != nil {
		return nil, err
	}
	var id, parentID int
	var title string
	err := ps.db.QueryRow(ctx, "SELECT id, title, board_id FROM columns WHERE id = $1 AND board_id = $2", columnID, boardID).Scan(&id, &title, &parentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("колонка не найдена: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка при поиске колонки: %w", err)
	}
	return &Column{
		ID:      id,
		Title:   title,
		BoardID: parentID,
	}, nil
}

func (ps *PostgresStorage) UpdateColumnTitle(ctx context.Context, boardID, columnID int, title string) (*Column, error) {
	if err := validateTitle(title, "колонки"); err != nil {
		return nil, err
	}
	if err := ps.boardExists(ctx, boardID); err != nil {
		return nil, err
	}
	var id, parentID int
	err := ps.db.QueryRow(ctx, "UPDATE columns SET title = $1 WHERE id = $2 AND board_id = $3 RETURNING id, board_id", title, columnID, boardID).Scan(&id, &parentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("колонка не найдена: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка при обновлении имени колонки: %w", err)
	}
	return &Column{
		ID:      id,
		Title:   title,
		BoardID: parentID,
	}, nil
}

func (ps *PostgresStorage) RemoveColumn(ctx context.Context, boardID, columnID int) error {
	if err := ps.boardExists(ctx, boardID); err != nil {
		return err
	}
	var id int
	err := ps.db.QueryRow(ctx, "DELETE FROM columns WHERE id = $1 AND board_id = $2 RETURNING id", columnID, boardID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("колонка не найдена: %w", ErrNotFound)
		}
		return fmt.Errorf("ошибка при удалении колонки: %w", err)
	}
	return nil
}

func (ps *PostgresStorage) AddTask(ctx context.Context, boardID, columnID int, title, description string) (*Task, error) {
	if err := validateTitle(title, "задачи"); err != nil {
		return nil, err
	}
	if err := ps.boardExists(ctx, boardID); err != nil {
		return nil, err
	}
	if err := ps.columnExists(ctx, boardID, columnID); err != nil {
		return nil, err
	}
	var count int
	err := ps.db.QueryRow(ctx, "SELECT COUNT(*) FROM tasks WHERE column_id = $1 AND title = $2", columnID, title).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке задачи: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("задача с таким именем уже существует в колонке: %w", ErrAlreadyExists)
	}
	var id int
	var createdAt time.Time
	err = ps.db.QueryRow(ctx, "INSERT INTO tasks (title, description, column_id) VALUES ($1, $2, $3) RETURNING id, created_at", title, description, columnID).Scan(&id, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании задачи: %w", err)
	}
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
		ColumnID:    columnID,
	}, nil
}

func (ps *PostgresStorage) GetTask(ctx context.Context, boardID, columnID, taskID int) (*Task, error) {
	if err := ps.boardExists(ctx, boardID); err != nil {
		return nil, err
	}
	if err := ps.columnExists(ctx, boardID, columnID); err != nil {
		return nil, err
	}
	var id, colID int
	var title, description string
	var createdAt time.Time
	err := ps.db.QueryRow(ctx, "SELECT id, title, description, created_at, column_id FROM tasks WHERE id = $1 AND column_id = $2", taskID, columnID).Scan(&id, &title, &description, &createdAt, &colID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("задача не найдена: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка при получении задачи: %w", err)
	}
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
		ColumnID:    colID,
	}, nil
}

func (ps *PostgresStorage) UpdateTaskTitle(ctx context.Context, boardID, columnID, taskID int, title string) (*Task, error) {
	if err := validateTitle(title, "задачи"); err != nil {
		return nil, err
	}
	if err := ps.boardExists(ctx, boardID); err != nil {
		return nil, err
	}
	if err := ps.columnExists(ctx, boardID, columnID); err != nil {
		return nil, err
	}
	var count int
	err := ps.db.QueryRow(ctx, "SELECT COUNT(*) FROM tasks WHERE column_id = $1 AND title = $2 AND id != $3", columnID, title, taskID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("ошибка при проверке задачи: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("задача с таким именем уже существует в колонке: %w", ErrAlreadyExists)
	}
	var id, colID int
	var description string
	var createdAt time.Time
	err = ps.db.QueryRow(ctx, "UPDATE tasks SET title = $1 WHERE id = $2 AND column_id = $3 RETURNING id, description, created_at, column_id", title, taskID, columnID).Scan(&id, &description, &createdAt, &colID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("задача не найдена: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка при обновлении имени задачи: %w", err)
	}
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
		ColumnID:    colID,
	}, nil
}

func (ps *PostgresStorage) UpdateTaskDescription(ctx context.Context, boardID, columnID, taskID int, description string) (*Task, error) {
	if err := ps.boardExists(ctx, boardID); err != nil {
		return nil, err
	}
	if err := ps.columnExists(ctx, boardID, columnID); err != nil {
		return nil, err
	}
	var id, colID int
	var title string
	var createdAt time.Time
	err := ps.db.QueryRow(ctx, "UPDATE tasks SET description = $1 WHERE id = $2 AND column_id = $3 RETURNING id, title, created_at, column_id", description, taskID, columnID).Scan(&id, &title, &createdAt, &colID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("задача не найдена: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("ошибка при обновлении описания задачи: %w", err)
	}
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
		ColumnID:    colID,
	}, nil
}

func (ps *PostgresStorage) RemoveTask(ctx context.Context, boardID, columnID, taskID int) error {
	if err := ps.boardExists(ctx, boardID); err != nil {
		return err
	}
	if err := ps.columnExists(ctx, boardID, columnID); err != nil {
		return err
	}
	var id int
	err := ps.db.QueryRow(ctx, "DELETE FROM tasks WHERE id = $1 AND column_id = $2 RETURNING id", taskID, columnID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("задача не найдена: %w", ErrNotFound)
		}
		return fmt.Errorf("ошибка при удалении задачи: %w", err)
	}
	return nil
}

func (ps *PostgresStorage) MoveTask(ctx context.Context, boardID, taskID, fromColumnID, toColumnID int) error {
	if fromColumnID == toColumnID {
		return nil
	}
	if err := ps.boardExists(ctx, boardID); err != nil {
		return err
	}
	if err := ps.columnExists(ctx, boardID, fromColumnID); err != nil {
		return fmt.Errorf("исходная колонка не найдена: %w", ErrNotFound)
	}
	if err := ps.columnExists(ctx, boardID, toColumnID); err != nil {
		return fmt.Errorf("целевая колонка не найдена: %w", ErrNotFound)
	}
	var id int
	err := ps.db.QueryRow(ctx, "UPDATE tasks SET column_id = $1 WHERE id = $2 AND column_id = $3 RETURNING id", toColumnID, taskID, fromColumnID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("задача не найдена в исходной колонке: %w", ErrNotFound)
		}
		return fmt.Errorf("ошибка при перемещении задачи: %w", err)
	}
	return nil
}

func (ps *PostgresStorage) GetColumnsByBoard(ctx context.Context, boardID int) ([]*Column, error) {
	if err := ps.boardExists(ctx, boardID); err != nil {
		return nil, err
	}
	rows, err := ps.db.Query(ctx, "SELECT id, title, board_id FROM columns WHERE board_id = $1", boardID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении колонок: %w", err)
	}
	defer rows.Close()
	columns := make([]*Column, 0)
	for rows.Next() {
		var id, boardID int
		var title string
		if err := rows.Scan(&id, &title, &boardID); err != nil {
			return nil, fmt.Errorf("ошибка при чтении колонки: %w", err)
		}
		columns = append(columns, &Column{
			ID:      id,
			Title:   title,
			BoardID: boardID,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе колонок: %w", err)
	}
	return columns, nil
}

func (ps *PostgresStorage) GetTasksByColumn(ctx context.Context, columnID int) ([]*Task, error) {
	if err := ps.columnExistsByID(ctx, columnID); err != nil {
		return nil, err
	}
	rows, err := ps.db.Query(ctx, "SELECT id, title, description, created_at, column_id FROM tasks WHERE column_id = $1", columnID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении задач: %w", err)
	}
	defer rows.Close()
	tasks := make([]*Task, 0)
	for rows.Next() {
		var id, colID int
		var title, description string
		var createdAt time.Time
		if err := rows.Scan(&id, &title, &description, &createdAt, &colID); err != nil {
			return nil, fmt.Errorf("ошибка при чтении задачи: %w", err)
		}
		tasks = append(tasks, &Task{
			ID:          id,
			Title:       title,
			Description: description,
			CreatedAt:   createdAt,
			ColumnID:    colID,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при переборе задач: %w", err)
	}
	return tasks, nil
}

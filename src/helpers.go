package src

import (
	"context"
	"fmt"
)

func validateTitle(title, entityName string) error {
	if len(title) < 1 {
		return fmt.Errorf("имя %s не может быть пустым: %w", entityName, ErrInvalidInput)
	}
	if len(title) > 75 {
		return fmt.Errorf("имя %s не может быть длиннее 75 символов: %w", entityName, ErrInvalidInput)
	}
	return nil
}

func (ps *PostgresStorage) boardExists(ctx context.Context, boardID int) error {
	var exists bool
	err := ps.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM boards WHERE id = $1)", boardID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка при проверке доски: %w", err)
	}
	if !exists {
		return fmt.Errorf("доска не найдена: %w", ErrNotFound)
	}
	return nil
}

func (ps *PostgresStorage) columnExists(ctx context.Context, boardID, columnID int) error {
	var exists bool
	err := ps.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM columns WHERE id = $1 AND board_id = $2)", columnID, boardID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка при проверке колонки: %w", err)
	}
	if !exists {
		return fmt.Errorf("колонка не найдена: %w", ErrNotFound)
	}
	return nil
}

func (ps *PostgresStorage) columnExistsByID(ctx context.Context, columnID int) error {
	var exists bool
	err := ps.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM columns WHERE id = $1)", columnID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка при проверке колонки: %w", err)
	}
	if !exists {
		return fmt.Errorf("колонка не найдена: %w", ErrNotFound)
	}
	return nil
}
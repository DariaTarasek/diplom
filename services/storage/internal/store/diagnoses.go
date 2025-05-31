package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DariaTarasek/diplom/services/storage/internal/model"
	"github.com/Masterminds/squirrel"
)

// GetICDCodes Получение списка всех кодов МКБ
func (s *Store) GetICDCodes(ctx context.Context) ([]model.ICD, error) {
	query, args, err := s.builder.
		Select("*").
		From("icd-codes").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка кодов МКБ: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var icd []model.ICD
	err = s.db.SelectContext(dbCtx, &icd, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка кодов МКБ: %w", err)
	}

	return icd, nil
}

func (s *Store) GetICDCodeByID(ctx context.Context, id model.ICDCodeID) (model.ICD, error) {
	query, args, err := s.builder.
		Select("*").
		From("icd-codes").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.ICD{}, fmt.Errorf("не удалось сформировать запрос для получения кода МКБ по id: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var icd model.ICD
	err = s.db.GetContext(dbCtx, &icd, query, args...)
	if err != nil {
		return model.ICD{}, fmt.Errorf("не удалось выполнить запрос для получения кода МКБ по id: %w", err)
	}

	return icd, nil
}

func (s *Store) GetICDCodeByCode(ctx context.Context, code string) (model.ICD, error) {
	query, args, err := s.builder.
		Select("*").
		From("icd-codes").
		Where(squirrel.Eq{"code": code}).
		ToSql()
	if err != nil {
		return model.ICD{}, fmt.Errorf("не удалось сформировать запрос для получения кода МКБ по коду: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var icd model.ICD
	err = s.db.GetContext(dbCtx, &icd, query, args...)
	if err != nil {
		return model.ICD{}, fmt.Errorf("не удалось выполнить запрос для получения кода МКБ по коду: %w", err)
	}

	return icd, nil
}

func (s *Store) GetDiagnosesByVisitID(ctx context.Context, visitID model.VisitID) ([]model.Diagnose, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointment_diagnoses").
		Where(squirrel.Eq{"visit_id": visitID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("не удалось сформировать запрос для получения диагнозов по id приема: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var diagnoses []model.Diagnose
	err = s.db.SelectContext(dbCtx, &diagnoses, query, args...)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос для получения диагнозов по id приема: %w", err)
	}

	return diagnoses, nil
}

func (s *Store) GetDiagnoseByID(ctx context.Context, id int) (model.Diagnose, error) {
	query, args, err := s.builder.
		Select("*").
		From("appointment_diagnoses").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Diagnose{}, fmt.Errorf("не удалось сформировать запрос для получения диагноза по id: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var diagnose model.Diagnose
	err = s.db.GetContext(dbCtx, &diagnose, query, args...)
	if err != nil {
		return model.Diagnose{}, fmt.Errorf("не удалось выполнить запрос для получения диагноза по id: %w", err)
	}

	return diagnose, nil
}

func (s *Store) AddDiagnoses(ctx context.Context, diagnoses []model.Diagnose) error {
	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для добавления диагнозов: %w", err)
	}
	defer tx.Rollback()

	builder := s.builder.
		Insert("appointment_diagnoses").
		Columns("visit_id", "icd_code_id", "diagnosis_note")

	for _, d := range diagnoses {
		builder = builder.Values(d.VisitID, d.ICDCodeID, d.DiagnosisNote)
	}
	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для добавления диагнозов: %w", err)
	}

	_, err = tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для добавления диагнозов: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для добавления диагнозов: %w", err)
	}

	return nil
}

func (s *Store) DeleteDiagnoseByID(ctx context.Context, id model.Diagnose) error {
	query, args, err := s.builder.
		Delete("appointment_diagnoses").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("не удалось сформировать запрос для удаления диагноза: %w", err)
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := s.db.BeginTxx(dbCtx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("не удалось начать транзакцию для удаления диагноза: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(dbCtx, query, args...)
	if err != nil {
		return fmt.Errorf("не удалось выполнить запрос для удаления диагноза: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось получить количество удаленных строк при удалении диагноза: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("диагноз не был")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("не удалось зафиксировать транзакцию для удаления устройства: %w", err)
	}

	return nil
}

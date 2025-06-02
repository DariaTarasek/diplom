package store

// GetICDCodes Получение списка всех кодов МКБ
//func (s *Store) GetICDCodes(ctx context.Context) ([]model.ICDCode, error) {
//	query, args, err := s.builder.
//		Select("*").
//		From("icd_codes").
//		ToSql()
//	if err != nil {
//		return nil, fmt.Errorf("не удалось сформировать запрос для получения списка кодов МКБ: %w", err)
//	}
//
//	dbCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
//	defer cancel()
//
//	var codes []model.ICDCode
//	err = s.db.SelectContext(dbCtx, &codes, query, args...)
//	if err != nil {
//		return nil, fmt.Errorf("не удалось выполнить запрос для получения списка кодов МКБ: %w", err)
//	}
//
//	return codes, nil
//}

package service

import (
	"context"
	"fmt"
	storagepb "github.com/DariaTarasek/diplom/services/admin/proto/storage"
)

func (s *AdminService) DeleteUser(ctx context.Context, id int) error {
	_, err := s.StorageClient.Client.DeleteUser(ctx, &storagepb.DeleteRequest{Id: int32(id)})
	if err != nil {
		return fmt.Errorf("не удалось удалить пользователя: %w", err)
	}
	return nil
}

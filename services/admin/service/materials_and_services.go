package service

import (
	"context"
	"fmt"
	"github.com/DariaTarasek/diplom/services/admin/model"
	storagepb "github.com/DariaTarasek/diplom/services/admin/proto/storage"
	"github.com/DariaTarasek/diplom/services/admin/sharederrors"
)

func (s *AdminService) AddMaterial(ctx context.Context, material model.Material) error {
	price := material.Price
	if price < 0 || price > 1000000 {
		return fmt.Errorf("некорректное значение цены материала: %w", sharederrors.ErrInvalidValue)
	}
	_, err := s.StorageClient.Client.AddMaterial(ctx, &storagepb.AddMaterialRequest{
		Name:  material.Name,
		Price: int32(material.Price),
	})
	if err != nil {
		return fmt.Errorf("не удалось добавить материал: %w", err)
	}
	return nil
}

func (s *AdminService) AddService(ctx context.Context, service model.Service) error {
	price := service.Price
	if price < 0 || price > 1000000 {
		return fmt.Errorf("некорректное значение цены услуги: %w", sharederrors.ErrInvalidValue)
	}
	_, err := s.StorageClient.Client.AddService(ctx, &storagepb.AddServiceRequest{
		Name:  service.Name,
		Price: int32(service.Price),
		Type:  int32(service.Category),
	})
	if err != nil {
		return fmt.Errorf("не удалось добавить услугу: %w", err)
	}
	return nil
}

func (s *AdminService) UpdateMaterial(ctx context.Context, material model.Material) error {
	price := material.Price
	if price < 0 || price > 1000000 {
		return fmt.Errorf("некорректное значение цены материала: %w", sharederrors.ErrInvalidValue)
	}
	_, err := s.StorageClient.Client.UpdateMaterial(ctx, &storagepb.UpdateMaterialRequest{
		Id:    int32(material.ID),
		Name:  material.Name,
		Price: int32(price),
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить материал: %w", err)
	}
	return nil
}

func (s *AdminService) UpdateService(ctx context.Context, service model.Service) error {
	price := service.Price
	if price < 0 || price > 1000000 {
		return fmt.Errorf("некорректное значение цены услуги: %w", sharederrors.ErrInvalidValue)
	}
	_, err := s.StorageClient.Client.UpdateService(ctx, &storagepb.UpdateServiceRequest{
		Id:    int32(service.ID),
		Name:  service.Name,
		Price: int32(service.Price),
		Type:  int32(service.Category),
	})
	if err != nil {
		return fmt.Errorf("не удалось обновить услугу: %w", err)
	}
	return nil
}

func (s *AdminService) DeleteMaterial(ctx context.Context, id int) error {
	_, err := s.StorageClient.Client.DeleteMaterial(ctx, &storagepb.DeleteRequest{Id: int32(id)})
	if err != nil {
		return fmt.Errorf("не удалось удалить материал %w", err)
	}
	return nil
}

func (s *AdminService) DeleteService(ctx context.Context, id int) error {
	_, err := s.StorageClient.Client.DeleteService(ctx, &storagepb.DeleteRequest{Id: int32(id)})
	if err != nil {
		return fmt.Errorf("не удалось удалить услугу %w", err)
	}
	return nil
}

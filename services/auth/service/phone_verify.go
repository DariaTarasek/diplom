package service

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/DariaTarasek/diplom/services/auth/sharederrors"
//	"github.com/DariaTarasek/diplom/services/auth/utils"
//	"github.com/redis/go-redis/v9"
//	"time"
//)
//
//const maxAttempts = 5
//
//func (s *AuthService) RequestCode(ctx context.Context, phone string) error {
//	code := utils.GenerateCode()
//
//	codeKey := fmt.Sprintf("verif:code:%s", phone)
//	sentKey := fmt.Sprintf("verif:sent:%s", phone)
//	attemptsKey := fmt.Sprintf("verif:attempts:%s", phone)
//
//	if s.RedisClient.Exists(ctx, sentKey).Val() > 0 {
//		return fmt.Errorf("код уже был отправлен, повторная отправка возможна через 60 секунд: %w", sharederrors.ErrRateLimited)
//	}
//
//	pipe := s.RedisClient.TxPipeline()
//	pipe.Set(ctx, codeKey, code, 5*time.Minute)
//	pipe.Set(ctx, attemptsKey, 0, 5*time.Minute)
//	pipe.Set(ctx, sentKey, 1, 60*time.Second)
//	_, err := pipe.Exec(ctx)
//	if err != nil {
//		return fmt.Errorf("не удалось записать код подтверждения: %w", err)
//	}
//
//	fmt.Printf("Код подтверждения для номера %s: %s", phone, code)
//	//msg := fmt.Sprintf("Код подтверждения: %s", code)
//	//phoneToSMS, err := strconv.Atoi(phone)
//	//if err != nil {
//	//	return fmt.Errorf("не удалось преобразовать номер телефона для отправки СМС: %w", err)
//	//}
//	//smsResponse, err := s.SMSClient.Client.SendSms(phoneToSMS, msg)
//	//fmt.Println(smsResponse)
//	//if err != nil {
//	//	return fmt.Errorf("не удалось отправить СМС с кодом подтверждения: %w", err)
//	//}
//	return nil
//}
//
//func (s *AuthService) VerifyCode(ctx context.Context, phone string, input string) error {
//	codeKey := fmt.Sprintf("verif:code:%s", phone)
//	attemptsKey := fmt.Sprintf("verif:attempts:%s", phone)
//
//	attempts, err := s.RedisClient.Get(ctx, attemptsKey).Int()
//	if err != nil {
//		return fmt.Errorf("не удалось получить количество попыток ввода кода: %w", err)
//	}
//	if attempts == maxAttempts {
//		return fmt.Errorf("слишком большое количество попыток. Отправьте код заново: %w", sharederrors.ErrTooManyAttempts)
//	}
//
//	code, err := s.RedisClient.Get(ctx, codeKey).Result()
//	if err != nil {
//		if errors.Is(err, redis.Nil) {
//			return fmt.Errorf("код не найден: %w", sharederrors.ErrCodeInvalid)
//		}
//		return err
//	}
//
//	if code != input {
//		s.RedisClient.Incr(ctx, attemptsKey)
//		return fmt.Errorf("неверный код: %w", sharederrors.ErrCodeInvalid)
//	}
//
//	s.RedisClient.Del(ctx, codeKey)
//	s.RedisClient.Del(ctx, attemptsKey)
//	return nil
//}

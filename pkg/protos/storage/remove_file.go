package storage

import (
	"context"

	"github.com/go-playground/validator/v10"

	pb "git.astralnalog.ru/utils/protos/storage/proto"
)

type RemoveFileInput struct {
	Id       string `validate:"required_without=CustomId"` // ID файла. Обязательный, если не указан CustomId
	CustomId string `validate:"omitempty"`                 // Уникальный ID в произвольном формате, задаваемый создателем файла (Необязательный) (макс. 128 символов)
}

// RemoveFile - Удаление файла по ID
func (s *Conn) RemoveFile(in RemoveFileInput) (err error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()

	return s.RemoveFileWithContext(ctx, in)
}

// RemoveFileWithContext удаляет файла по ID c заданным контекстом
func (s *Conn) RemoveFileWithContext(ctx context.Context, in RemoveFileInput) (err error) {
	defer func() {
		if err != nil {
			err = errWrapper(err)
		}
	}()

	if err = validator.New().Struct(in); err != nil {
		return
	}

	doneChan := make(chan bool, 1)
	errChan := make(chan error, 1)

	go s.removeFile(ctx, &in, doneChan, errChan)

	select {
	case <-ctx.Done():
		err = errTimeout(methodRemoveFile, in)
	case <-doneChan:
		return
	case v := <-errChan:
		err = v
	}

	return
}

func (s *Conn) removeFile(ctx context.Context, in *RemoveFileInput, doneChan chan<- bool, errChan chan<- error) {
	if _, err := s.client.RemoveFile(ctx, &pb.RequestRemoveFile{Id: in.Id, CustomId: in.CustomId, ServiceName: s.serviceName}); err != nil {
		errChan <- err
		return
	}

	doneChan <- true
}

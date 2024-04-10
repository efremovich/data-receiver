package storage

import (
	"context"

	"github.com/go-playground/validator/v10"

	pb "github.com/efremovich/data-receiver/pkg/protos/storage/proto"
)

type SetFileAttrsInput struct {
	Id       string     `validate:"required_without=CustomId"` // ID файла. Обязательный, если не указан CustomId
	CustomId string     `validate:"-"`                         // Уникальный ID в произвольном формате, задаваемый создателем файла (Необязательный) (макс. 128 символов)
	Attrs    *FileAttrs `validate:"required"`                  // Редактируемые атрибуты файла
}

// SetFileAttrs - Обновление атрибутов файла
func (s *Conn) SetFileAttrs(in SetFileAttrsInput) (err error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	return s.SetFileAttrsWithContext(ctx, in)
}

// SetFileAttrsWithContext обновляет атрибуты файла с заданным контекстом
func (s *Conn) SetFileAttrsWithContext(ctx context.Context, in SetFileAttrsInput) (err error) {
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

	go s.setFileAttrs(ctx, &in, doneChan, errChan)

	select {
	case <-ctx.Done():
		err = errTimeout(methodSetFileAttrs, in)
	case <-doneChan:
		return
	case v := <-errChan:
		err = v
	}

	return
}

func (s *Conn) setFileAttrs(ctx context.Context, in *SetFileAttrsInput, doneChan chan<- bool, errChan chan<- error) {
	attrs := new(pb.FileAttrs)
	attrs.TTL = in.Attrs.TTL
	attrs.Type = in.Attrs.Type
	attrs.Filename = in.Attrs.Filename
	attrs.SubType = in.Attrs.SubType
	attrs.Readonly = in.Attrs.Readonly
	attrs.Protected = in.Attrs.Protected

	if _, err := s.client.SetFileAttrs(ctx,
		&pb.RequestSetFileAttrs{Id: in.Id, CustomId: in.CustomId, ServiceName: s.serviceName, Attrs: attrs}); err != nil {
		errChan <- err
		return
	}

	doneChan <- true
}

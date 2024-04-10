package storage

import (
	"context"

	"github.com/go-playground/validator/v10"

	pb "github.com/efremovich/data-receiver/pkg/protos/storage/proto"
)

type GetFileAttrsInput struct {
	Id       string `validate:"required_without=CustomId"` // ID файла. Обязательный, если не указан CustomId
	CustomId string `validate:"omitempty"`                 // Уникальный ID в произвольном формате, задаваемый создателем файла (Необязательный) (макс. 128 символов)
}

type GetFileAttrsOutput struct {
	Id           string                   // ID файла
	Attrs        *AllFileAttrs            // Атрибуты файла
	ServiceAttrs map[string]*ServiceAttrs // Сервисные атрибуты файла
}

// GetFileAttrs - Получение атрибутов файла по ID
func (s *Conn) GetFileAttrs(in GetFileAttrsInput) (out *GetFileAttrsOutput, err error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	return s.GetFileAttrsWithContext(ctx, in)
}

// GetFileAttrsWithContext - Получение атрибутов файла по ID с заданным контекстом
func (s *Conn) GetFileAttrsWithContext(ctx context.Context, in GetFileAttrsInput) (out *GetFileAttrsOutput, err error) {
	defer func() {
		if err != nil {
			err = errWrapper(err)
		}
	}()

	if err = validator.New().Struct(in); err != nil {
		return
	}
	out = new(GetFileAttrsOutput)
	resultChan := make(chan *GetFileAttrsOutput, 1)
	errChan := make(chan error, 1)

	go s.getFileAttrs(ctx, &in, resultChan, errChan)

	select {
	case <-ctx.Done():
		err = errTimeout(methodGetFileAttrs, in)
	case v := <-resultChan:
		out = v
	case v := <-errChan:
		err = v
	}

	return
}

func (s *Conn) getFileAttrs(ctx context.Context, in *GetFileAttrsInput, resultChan chan<- *GetFileAttrsOutput, errChan chan<- error) {
	var resp *pb.ResponseGetFileAttrs
	var err error
	if resp, err = s.client.GetFileAttrs(ctx, &pb.RequestGetFileAttrs{Id: in.Id, CustomId: in.CustomId}); err != nil {
		errChan <- err
		return
	}

	result := new(GetFileAttrsOutput)
	result.Id = resp.Id
	result.Attrs = new(AllFileAttrs)
	result.Attrs.Created = resp.Attrs.Created
	result.Attrs.Expires = resp.Attrs.Expires
	result.Attrs.Creator = resp.Attrs.Creator
	result.Attrs.CustomId = resp.Attrs.CustomId
	result.Attrs.Filename = resp.Attrs.Filename
	result.Attrs.Size = resp.Attrs.Size
	result.Attrs.StorageType = resp.Attrs.StorageType
	result.Attrs.Type = resp.Attrs.Type
	result.Attrs.SubType = resp.Attrs.SubType
	result.Attrs.Readonly = resp.Attrs.Readonly
	result.Attrs.Protected = resp.Attrs.Protected
	result.ServiceAttrs = map[string]*ServiceAttrs{}
	for k, v := range resp.ServiceAttrs {
		result.ServiceAttrs[k] = &ServiceAttrs{
			TTL: v.TTL,
		}
	}

	resultChan <- result
}

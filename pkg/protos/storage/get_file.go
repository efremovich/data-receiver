package storage

import (
	"bytes"
	"context"
	"io"

	"github.com/go-playground/validator/v10"

	pb "git.astralnalog.ru/utils/protos/storage/proto"
)

type GetFileInput struct {
	Id       string `validate:"required_without=CustomId"` // ID файла. Обязательный, если не указан CustomId
	CustomId string `validate:"omitempty"`                 // Уникальный ID в произвольном формате, задаваемый создателем файла (Необязательный) (макс. 128 символов)
}

type GetFileOutput struct {
	File *bytes.Buffer // Содержимое файла
}

// GetFile - Получение содержимого файла по id
func (s *Conn) GetFile(in GetFileInput) (out *GetFileOutput, err error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	return s.GetFileWithContext(ctx, in)
}

// GetFileWithContext - получение содержимого файла по ID с заданным контекстом
func (s *Conn) GetFileWithContext(ctx context.Context, in GetFileInput) (out *GetFileOutput, err error) {
	defer func() {
		if err != nil {
			err = errWrapper(err)
		}
	}()

	if err = validator.New().Struct(in); err != nil {
		return
	}
	out = new(GetFileOutput)

	resultChan := make(chan *GetFileOutput, 1)
	errChan := make(chan error, 1)

	go s.getFile(ctx, &in, resultChan, errChan)

	select {
	case <-ctx.Done():
		err = errTimeout(methodGetFile, in)
	case v := <-resultChan:
		out = v
	case v := <-errChan:
		err = v
	}

	return
}

func (s *Conn) getFile(ctx context.Context, in *GetFileInput, resultChan chan<- *GetFileOutput, errChan chan<- error) {
	stream, err := s.client.GetFile(ctx, &pb.RequestGetFile{Id: in.Id, CustomId: in.CustomId})
	if err != nil {
		errChan <- err
		return
	}

	file := new(bytes.Buffer)
	for {
		if ctxDead(ctx) {
			return
		}

		var resp *pb.ResponseGetFile
		resp, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				errChan <- err
				return
			}
		}
		file.Write(resp.File)
	}

	resultChan <- &GetFileOutput{File: file}
	return
}

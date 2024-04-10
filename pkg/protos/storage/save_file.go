package storage

import (
	"context"
	"io"

	"github.com/go-playground/validator/v10"

	pb "github.com/efremovich/data-receiver/pkg/protos/storage/proto"
)

type SaveFileInput struct {
	Rewrite  bool       `validate:"-"`        // Перезаписывать ли файл, если такой ID уже имеется в хранилище
	Id       string     `validate:"-"`        // ID файла для сохранения в формате UUIDv4
	CustomId string     `validate:"-"`        // Уникальный ID в произвольном формате, задаваемый создателем файла (Необязательный) (макс. 128 символов)
	File     io.Reader  `validate:"required"` // Содержимое файла
	Attrs    *FileAttrs `validate:"-"`        // Редактируемые атрибуты файла
}

type SaveFileOutput struct {
	Id string // ID файла
}

// SaveFile - Загрузка нового файла или перезапись существующего
func (s *Conn) SaveFile(in SaveFileInput) (out *SaveFileOutput, err error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()
	return s.SaveFileWithContext(ctx, in)
}

// SaveFileWithContext загружает новый файл или перезаписывает существующий с заданным контекстом
func (s *Conn) SaveFileWithContext(ctx context.Context, in SaveFileInput) (out *SaveFileOutput, err error) {
	defer func() {
		if err != nil {
			err = errWrapper(err)
		}
	}()

	if in.Attrs == nil {
		in.Attrs = new(FileAttrs)
	}
	if err = validator.New().Struct(in); err != nil {
		return
	}
	out = new(SaveFileOutput)

	resultChan := make(chan *SaveFileOutput, 1)
	errChan := make(chan error, 1)

	go s.saveFile(ctx, &in, resultChan, errChan)

	select {
	case <-ctx.Done():
		err = errTimeout(methodSaveFile, in)
	case v := <-resultChan:
		out = v
	case v := <-errChan:
		err = v
	}

	return
}

func (s *Conn) saveFile(ctx context.Context, in *SaveFileInput, resultChan chan<- *SaveFileOutput, errChan chan<- error) {
	stream, err := s.client.SaveFile(ctx)
	if err != nil {
		errChan <- err
		return
	}

	chunkBuf := make([]byte, 1024)
	firstMsg := true
	for {
		if ctxDead(ctx) {
			return
		}

		var n int
		n, err = in.File.Read(chunkBuf)
		if err != nil {
			if err == io.EOF {
				break
			}

			errChan <- err
			return
		}
		if firstMsg {
			attrs := new(pb.FileAttrs)
			if in.Attrs != nil {
				attrs.TTL = in.Attrs.TTL
				attrs.Type = in.Attrs.Type
				attrs.Filename = in.Attrs.Filename
				attrs.SubType = in.Attrs.SubType
				attrs.Readonly = in.Attrs.Readonly
				attrs.Protected = in.Attrs.Protected
			}

			if err = stream.Send(&pb.RequestSaveFile{
				Id:          in.Id,
				CustomId:    in.CustomId,
				ServiceName: s.serviceName,
				Rewrite:     in.Rewrite,
				Attrs:       attrs,
				File:        chunkBuf[:n]},
			); err != nil {
				errChan <- err
				return
			}
			firstMsg = false
			continue
		}

		if err = stream.Send(&pb.RequestSaveFile{File: chunkBuf[:n]}); err != nil {
			// если ошибка сгенерирована не клиентом, а сервером, то прерываем стрим
			// и получаем ошибку сервера в методе CloseAndRecv()
			if err == io.EOF {
				break
			}

			errChan <- err
			return
		}
	}

	var resp *pb.ResponseSaveFile
	if resp, err = stream.CloseAndRecv(); err != nil {
		errChan <- err
		return
	}

	resultChan <- &SaveFileOutput{Id: resp.Id}
}

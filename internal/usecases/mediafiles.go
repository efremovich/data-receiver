package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setMediaFile(ctx context.Context, card *entity.Card) ([]*entity.MediaFile, error) {
	mediaFiles := []*entity.MediaFile{}
	for _, elem := range card.MediaFile {
		mediaFiles, err := s.mediafilerepo.SelectByCardID(ctx, card.ID, elem.Link)
		if errors.Is(err, ErrObjectNotFound) {
			mediaFile, err := s.mediafilerepo.Insert(ctx, entity.MediaFile{
				Link:   elem.Link,
				TypeID: elem.TypeID,
				CardID: card.ID,
			})
			if err != nil {
				return nil, wrapErr(fmt.Errorf("Ошибка при получении данных: %w", err))
			}
			mediaFiles = append(mediaFiles, mediaFile)
		}
	}

	return mediaFiles, nil
}

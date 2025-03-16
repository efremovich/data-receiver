package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setMediaFile(ctx context.Context, card *entity.Card) ([]*entity.MediaFile, error) {
	var mediaFile *entity.MediaFile

	mediaFiles := []*entity.MediaFile{}

	for _, elem := range card.MediaFile {
		mfList, err := s.mediafilerepo.SelectByCardID(ctx, card.ID, elem.Link)
		if errors.Is(err, ErrObjectNotFound) {
			mediaFile, err = s.mediafilerepo.Insert(ctx, entity.MediaFile{
				Link:   elem.Link,
				TypeID: elem.TypeID,
				CardID: card.ID,
			})
			if err != nil {
				return nil, wrapErr(fmt.Errorf("ошибка при получении данных: %w", err))
			}

			mediaFiles = append(mediaFiles, mediaFile)
		} else {
			mediaFiles = append(mediaFiles, mfList...)
		}
	}

	return mediaFiles, nil
}

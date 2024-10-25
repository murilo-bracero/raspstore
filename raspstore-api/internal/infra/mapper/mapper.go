package mapper

import (
	"fmt"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
)

func MapFilePageResponse(page int, size int, filesPage *entity.FilePage, host string) *model.FilePageResponse {
	nextUrl := buildNextUrl(filesPage, host, page, size)

	return &model.FilePageResponse{
		Page:          page,
		Size:          size,
		TotalElements: filesPage.Count,
		Next:          nextUrl,
		Content:       mapFilePageContents(filesPage.Content),
	}
}

func mapFilePageContents(entities []*entity.File) []*model.FileContent {
	res := make([]*model.FileContent, len(entities))

	for i, e := range entities {
		res[i] = mapFilePageContentParser(e)
	}

	return res
}

func mapFilePageContentParser(entity *entity.File) *model.FileContent {
	return &model.FileContent{
		FileId:    entity.FileId,
		Filename:  entity.Filename,
		Size:      entity.Size,
		Owner:     entity.Owner,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
	}
}

func buildNextUrl(filesPage *entity.FilePage, host string, page int, size int) (nextUrl string) {
	if len(filesPage.Content) == size {
		nextUrl = fmt.Sprintf("%s/file-service/v1/files?page=%d&size=%d", host, page+1, size)
	}

	return
}

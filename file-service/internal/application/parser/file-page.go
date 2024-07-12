package parser

import (
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
)

func FilePageResponseParser(page int, size int, filesPage *entity.FilePage, nextUrl string) *model.FilePageResponse {
	return &model.FilePageResponse{
		Page:          page,
		Size:          size,
		TotalElements: filesPage.Count,
		Next:          nextUrl,
		Content:       filesPage.Content,
	}
}

package converter

import (
	v1 "github.com/murilo-bracero/raspstore/file-info-service/api/v1"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
)

func ToFilePageResponse(page int, size int, filesPage *model.FilePage, nextUrl string) *v1.FilePageResponse {
	return &v1.FilePageResponse{
		Page:          page,
		Size:          size,
		TotalElements: filesPage.Count,
		Next:          nextUrl,
		Content:       filesPage.Content,
	}
}

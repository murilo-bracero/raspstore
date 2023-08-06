package model

import v1 "github.com/murilo-bracero/raspstore/idp/api/v1"

type UserPage struct {
	Content []*User `bson:"content"`
	Count   int     `bson:"count"`
}

func (up *UserPage) ToPageRepresentation(page int, size int, nextUrl string) v1.PageRepresentation {
	content := make([]*v1.UserRepresentation, len(up.Content))
	for i, usr := range up.Content {
		content[i] = usr.ToUserRepresentation()
	}

	return v1.PageRepresentation{
		Page:          page,
		Size:          size,
		TotalElements: up.Count,
		Next:          nextUrl,
		Content:       content,
	}
}

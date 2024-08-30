package http

import (
	"errors"
)

type RequestQuery struct {
	Page    int    `query:"page" url:"page"`
	PerPage int    `query:"per_page" url:"per_page"`
	Q       string `query:"q" url:"q"`
	Sort    string `query:"sort" url:"sort"`
	Filter  string `query:"filter" url:"filter"`
}

type Pagination struct {
	HasMorePage    bool   `json:"has_more_page"`
	Count          int64  `json:"count"`
	Total          int64  `json:"total"`
	PerPage        int64  `json:"per_page"`
	CurrentPage    int64  `json:"current_page"`
	LastPage       int64  `json:"last_page"`
	PrevPageUrl    string `json:"prev_page_url,omitempty"`
	CurrentPageUrl string `json:"current_page_url,omitempty"`
	NextPageUrl    string `json:"next_page_url,omitempty"`
}

func Paginate(count int64, total int64, query *RequestQuery) Pagination {
	lastPage := total / int64(query.PerPage)
	return Pagination{
		HasMorePage: int64(query.Page) < lastPage,
		Count:       count,
		Total:       total,
		PerPage:     int64(query.PerPage),
		CurrentPage: int64(query.Page),
		LastPage:    lastPage,
	}
}

var (
	reqHeader HeaderParams
)

type HeaderParams struct {
	TokenRaw string
	OutletId string
	Setting  string
}

//// SetReqHeader create global params. it's called by middleware
//func SetReqHeader(ctx *fiber.Ctx, db configs.SqlDB) error {
//	token := ctx.Locals("jwtKey").(*jwt.Token)
//
//	reqHeader = HeaderParams{
//		TokenRaw: token.Raw,
//	}
//
//	return nil
//}
//
//func SetReqHeaderRaw(token string, outletId string) {
//	reqHeader = HeaderParams{
//		TokenRaw: token,
//		OutletId: outletId,
//	}
//}

func GetReqHeader() (*HeaderParams, error) {
	if reqHeader.TokenRaw == "" {
		return nil, errors.New("authorization header must be present")
	}
	return &reqHeader, nil
}

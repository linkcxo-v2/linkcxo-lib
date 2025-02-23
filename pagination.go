package linkcxo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type GetPaginatedRequest struct {
	Sort  PaginationSortRequest `json:"sort"`
	Size  int                   `json:"size"`
	MaxID string                `json:"maxId"`
	MinID string                `json:"minId"`
}
type PaginationSortRequest struct {
	Field    string
	Order    int
	Original string
}

type PaginationMetadata struct {
	Size              int    `json:"size"`
	MinID             string `json:"minId"`
	MaxID             string `json:"maxId"`
	NextResultURL     string `json:"nextResultURL"`
	PreviousResultURL string `json:"previousResultURL"`
}

type PaginationUrlFunc func(req interface{}) string
type PaginationRequestFunc func(c echo.Context, req GetPaginatedRequest) interface{}
type PaginationConfig struct {
	UrlFunc     PaginationUrlFunc
	RequestFunc PaginationRequestFunc
}
type Pagination struct {
	req        GetPaginatedRequest
	pagination PaginationConfig
	parser     PaginationParser
	Request    interface{}
}
type PaginationBuilder struct {
	pagination PaginationConfig
	context    echo.Context
}

func NewPaginationBuilder() *PaginationBuilder {
	return &PaginationBuilder{}
}
func (pb *PaginationBuilder) WithContext(c echo.Context) *PaginationBuilder {
	pb.context = c
	return pb
}
func (pb *PaginationBuilder) WithPagination(f PaginationConfig) *PaginationBuilder {
	pb.pagination = f
	return pb
}
func (pb *PaginationBuilder) Build() *Pagination {
	pParser := PaginationParser{}
	req := pParser.ParseRequest(pb.context)
	req.Sort = pParser.ParseSort(pb.context)
	newReq := pb.pagination.RequestFunc(pb.context, req)
	return &Pagination{
		req:        req,
		Request:    newReq,
		pagination: pb.pagination,
		parser:     pParser,
	}

}

type PaginationData struct {
	Len   int
	MinID string
	MaxID string
}

func (p *Pagination) Metadata(meta PaginationData) *PaginationMetadata {
	req := p.req
	metadata := PaginationMetadata{}
	if meta.Len > 0 {
		metadata.Size = req.Size
		metadata.MinID = meta.MinID
		metadata.MaxID = meta.MaxID
		url := p.pagination.UrlFunc(p.Request)
		if meta.Len > 0 {
			metadata.NextResultURL = url + p.parser.BuildSortAndFilterQuery(metadata, true, req)
		}
		if req.MaxID != "" {
			metadata.PreviousResultURL = url + p.parser.BuildSortAndFilterQuery(metadata, false, req)
		}

	} else {
		metadata.Size = req.Size
	}
	return &metadata
}

type PaginationParser struct {
}

func (pu PaginationParser) BuildSortAndFilterQuery(metadata PaginationMetadata, isNext bool,
	req GetPaginatedRequest) string {
	url := ""

	if req.Sort.Original != "" {
		url = url + fmt.Sprintf("&sort=%s", req.Sort.Original)
	}
	if req.Sort.Order == -1 {
		if isNext {
			url = url + fmt.Sprintf("&maxId=%s", metadata.MaxID)
		} else {
			url = url + fmt.Sprintf("&minId=%s", metadata.MinID)
		}
	} else if req.Sort.Order == 1 {
		if !isNext {
			url = url + fmt.Sprintf("&maxId=%s", metadata.MaxID)
		} else {
			url = url + fmt.Sprintf("&minId=%s", metadata.MinID)
		}
	}

	return url
}

func (pu PaginationParser) getDefaultSort(c echo.Context) PaginationSortRequest {
	s := c.QueryParam("sort")
	return PaginationSortRequest{
		Field:    "_id",
		Order:    -1,
		Original: s,
	}
}
func (pu PaginationParser) ParseSort(c echo.Context) PaginationSortRequest {

	sort := pu.getDefaultSort(c)
	s := sort.Original
	if strings.TrimSpace(s) == "" {
		return sort
	}
	if strings.HasPrefix(s, "-") {
		if s[1:] != "" {
			sort.Field = s[1:]
		}
		sort.Order = -1
	} else {
		sort.Order = 1
	}
	return sort
}
func (pu PaginationParser) ParseRequest(c echo.Context) GetPaginatedRequest {
	r := GetPaginatedRequest{}
	var err error
	r.MinID = c.QueryParam("minId")
	r.MaxID = c.QueryParam("maxId")
	size := c.QueryParam("size")
	if size != "" {
		r.Size, err = strconv.Atoi(size)
		if err != nil {
			r.Size = 10
		}
		if r.Size > 1000 {
			r.Size = 1000
		}
	} else {
		r.Size = 10
	}

	return r
}

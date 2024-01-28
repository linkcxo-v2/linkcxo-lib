package linkcxo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type GetPaginatedRequest struct {
	Sort  SortRequest `json:"sort"`
	Size  int         `json:"size"`
	MaxID string      `json:"maxId"`
	MinID string      `json:"minId"`
}
type SortRequest struct {
	Field    string
	Order    int
	Original string
}

type Metadata struct {
	Size              int    `json:"size"`
	MinID             string `json:"minId"`
	MaxID             string `json:"maxId"`
	NextResultURL     string `json:"nextResultURL"`
	PreviousResultURL string `json:"previousResultURL"`
}

type PaginationUrlFunc func(req interface{}) string
type PaginationRequestFunc func(req GetPaginatedRequest) interface{}
type PaginationConfig struct {
	UrlFunc     PaginationUrlFunc
	RequestFunc PaginationRequestFunc
}
type Pagination struct {
	req        GetPaginatedRequest
	pagination PaginationConfig
	parser     paginationParser
	reqImp     interface{}
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
	pParser := paginationParser{}
	req := pParser.parseRequest(pb.context)
	req.Sort = pParser.parseSort(pb.context)
	newReq := pb.pagination.RequestFunc(req)
	return &Pagination{
		req:        req,
		reqImp:     newReq,
		pagination: pb.pagination,
		parser:     pParser,
	}

}

type PaginationData struct {
	Data  []interface{}
	MinID string
	MaxID string
}

func (p *Pagination) Metadata(meta PaginationData) *Metadata {
	req := p.req
	metadata := Metadata{}
	if len(meta.Data) > 0 {
		metadata.Size = req.Size
		metadata.MinID = meta.MinID
		metadata.MaxID = meta.MaxID
		url := p.pagination.UrlFunc(p.reqImp)
		if len(meta.Data) >= req.Size {
			metadata.NextResultURL = url + p.parser.buildSortAndFilterQuery(metadata, true, req)
		}
		if req.MaxID != "" {
			metadata.PreviousResultURL = url + p.parser.buildSortAndFilterQuery(metadata, false, req)
		}

	} else {
		metadata.Size = req.Size
	}
	return &metadata
}

type paginationParser struct {
}

func (pu paginationParser) buildSortAndFilterQuery(metadata Metadata, isNext bool, req GetPaginatedRequest) string {
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

func (pu paginationParser) getDefaultSort(c echo.Context) SortRequest {
	s := c.QueryParam("sort")
	return SortRequest{
		Field:    "_id",
		Order:    -1,
		Original: s,
	}
}
func (pu paginationParser) parseSort(c echo.Context) SortRequest {

	sort := pu.getDefaultSort(c)
	s := sort.Original
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
func (pu paginationParser) parseRequest(c echo.Context) GetPaginatedRequest {
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
		if r.Size > 20 {
			r.Size = 20
		}
	} else {
		r.Size = 10
	}

	return r
}

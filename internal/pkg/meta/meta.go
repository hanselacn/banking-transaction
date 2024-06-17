package meta

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ErrInvalidMetadata is an error when metadata is invalid.
// This error usually returned by the implementation of Filter interface.
var ErrInvalidMetadata = errors.New("invalid metadata")

// Metadata represents a metadata for HTTP API.
type Metadata struct {
	Pagination
	Filtering
	*DateRange     `json:"date_range,omitempty"`
	*DateTimeRange `json:"date_time_range,omitempty"`
	ErrorCode      string `json:"error_code,omitempty"`
}

type MetadataLimit struct {
	Limit     int `json:"limit"`
	TotalData int `json:"total_data"`
}

type MetadataErrorCode struct {
	Code string `json:"code,omitempty"`
}

type MetadataExperienceLog struct {
	Pagination
	Filtering
	*DateRange     `json:"date_range,omitempty"`
	*DateTimeRange `json:"date_time_range,omitempty"`
	Filter         FilterExperience `json:"filter,omitempty"`
}

type FilterExperience struct {
	Activity    string      `json:"activity,omitempty"`
	Application interface{} `json:"application,omitempty"`
}

// MetadataFromURL gets metadata from the given request url.
func MetadataFromURL(u url.Values) Metadata {
	return Metadata{
		Pagination: PaginationFromURL(u),
		Filtering:  FilterFromURL(u),
	}
}

// ParseMetadataFromURL gets metadata from the given request url.
func ParseMetadataFromURL(u url.Values) Metadata {
	return Metadata{
		Pagination: PaginationFromURL(u),
		Filtering:  ParseFilterFromURL(u),
		DateRange:  ParseDateRangeFromURL(u),
	}
}

// ParsingMetadataFromURL gets metadata from the given request url - default ordering desc and field created_at.
func ParsingMetadataFromURL(u url.Values) Metadata {
	return Metadata{
		Pagination: PaginationFromURL(u),
		Filtering:  ParsingFilterFromURL(u),
		DateRange:  ParseDateRangeFromURL(u),
	}
}

// DefaultPerPage is a default value for per_page query params.
const DefaultPerPage = 10

// Pagination is a meta data for pagination.
type Pagination struct {
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
	Total   int `json:"total"`
}

// PaginationFromURL gets pagination meta from request URL.
func PaginationFromURL(u url.Values) Pagination {
	p := Pagination{
		PerPage: DefaultPerPage,
		Page:    1,
	}

	pps := u.Get("per_page")
	if v, err := strconv.Atoi(pps); err == nil {
		if v <= 0 {
			v = DefaultPerPage
		}

		if v > 25 {
			v = 25
		}

		p.PerPage = v
	}

	ps := u.Get("page")
	if v, err := strconv.Atoi(ps); err == nil {
		if v < 1 {
			v = 1
		}

		p.Page = v
	}

	return p
}

// SortXXX are default values for order_type query params.
const (
	SortAscending  = "asc"
	SortDescending = "desc"
)

// Filtering represents a filterable fields.
type Filtering struct {
	OrderBy   string                   `json:"order_by,omitempty"`
	OrderType string                   `json:"order_type,omitempty"`
	Search    string                   `json:"search,omitempty"`
	Keyword   string                   `json:"keyword,omitempty"`
	SearchBy  string                   `json:"search_by,omitempty"`
	FilterBy  string                   `json:"filter_by,omitempty"`
	Filter    string                   `json:"filter,omitempty"`
	SearchOpt []map[string]interface{} `json:"search_opt,omitempty"`
}

// FilterFromURL gets filter values from query params.
func FilterFromURL(u url.Values) Filtering {
	f := Filtering{
		OrderBy:   "created_at",
		OrderType: SortAscending,
	}

	ob := u.Get("order_by")
	if len(ob) != 0 {
		f.OrderBy = strings.ToLower(strings.ToLower(ob))
	}

	ot := u.Get("order_type")
	if len(ot) != 0 {
		ot = strings.TrimSpace(strings.ToLower(ot))
		if ot == SortDescending {
			f.OrderType = SortDescending
		}
	}

	search := strings.TrimSpace(u.Get("search"))
	if len(search) == 0 {
		search = strings.TrimSpace(u.Get("keyword"))
	}

	if len(search) != 0 {
		f.Search = search
	}

	filter := strings.TrimSpace(u.Get("filter"))
	if len(filter) == 0 {
		filter = strings.TrimSpace(u.Get("keyword"))
	}

	if len(filter) != 0 {
		f.Filter = filter
	}

	searchBy := strings.TrimSpace(u.Get("search_by"))
	if len(search) == 0 {
		search = strings.TrimSpace(u.Get("search"))
	}

	if len(search) != 0 {
		f.SearchBy = searchBy
	}

	filterBy := strings.TrimSpace(u.Get("filter_by"))
	if len(filter) == 0 {
		filter = strings.TrimSpace(u.Get("filter"))
	}

	if len(filter) != 0 {
		f.FilterBy = filterBy
	}

	return f
}

// ParseFilterFromURL gets filter values from query params.
func ParseFilterFromURL(u url.Values) Filtering {
	// regex check Unicode
	regexUnicode, _ := regexp.Compile(`\\u|\\U`)

	f := Filtering{
		OrderBy:   "updated_at",
		OrderType: SortDescending,
	}

	ob := u.Get("order_by")
	if len(ob) != 0 {
		f.OrderBy = strings.ToLower(strings.ToLower(ob))
	}

	ot := u.Get("order_type")
	if len(ot) != 0 {
		ot = strings.TrimSpace(strings.ToLower(ot))
		if ot == SortAscending {
			f.OrderType = SortAscending
		}
	}
	search := strings.TrimSpace(u.Get("search"))
	matchsearch := regexUnicode.MatchString(search)

	if matchsearch {
		search = ""
	}

	if len(search) == 0 {
		search = strings.TrimSpace(u.Get("keyword"))
	}

	if len(search) != 0 {
		f.Search = search
	}

	searchBy := strings.TrimSpace(u.Get("search_by"))
	if len(search) == 0 {
		search = strings.TrimSpace(u.Get("search"))
	}

	if len(search) != 0 {
		f.SearchBy = searchBy
	}

	filter := strings.TrimSpace(u.Get("filter"))
	if len(filter) == 0 {
		filter = strings.TrimSpace(u.Get("keyword"))
	}

	if len(filter) != 0 {
		f.Filter = filter
	}
	filterBy := strings.TrimSpace(u.Get("filter_by"))
	if len(filter) == 0 {
		filter = strings.TrimSpace(u.Get("filter"))
	}

	if len(filter) != 0 {
		f.FilterBy = filterBy
	}
	return f
}

type DateRange struct {
	Field string    `json:"field"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func DateRangeFromURL(u url.Values, field string, startQuery, endQuery string) (*DateRange, error) {
	ts := u.Get(startQuery)
	te := u.Get(endQuery)
	if len(ts) == 0 || len(te) == 0 {
		return nil, nil
	}

	dr := DateRange{
		Field: "created_at",
		Start: time.Time{},
		End:   time.Time{},
	}

	if v := u.Get(field); len(v) != 0 {
		dr.Field = strings.TrimSpace(strings.ToLower(v))
	}

	t, err := time.Parse("2006-01-02", ts)
	if err != nil {
		return nil, ErrInvalidMetadata
	}

	dr.Start = t

	t, err = time.Parse("2006-01-02", te)
	if err != nil {
		return nil, ErrInvalidMetadata
	}

	dr.End = t

	return &dr, nil
}

func TimeDateRangeFromURL(u url.Values, field string, startQuery, endQuery string) (*DateRange, error) {
	ts := u.Get(startQuery)
	te := u.Get(endQuery)
	if len(ts) == 0 || len(te) == 0 {
		return nil, nil
	}

	dr := DateRange{
		Field: "created_at",
		Start: time.Time{},
		End:   time.Time{},
	}

	if v := u.Get(field); len(v) != 0 {
		dr.Field = strings.TrimSpace(strings.ToLower(v))
	}

	t, err := time.Parse("2006-01-02 15:04", ts)
	if err != nil {
		return nil, fmt.Errorf("start_date invalid format")
	}

	dr.Start = t

	t, err = time.Parse("2006-01-02 15:04", te)
	if err != nil {
		return nil, fmt.Errorf("end_date invalid format")
	}

	dr.End = t

	if dr.Start.After(dr.End) {
		return nil, ErrInvalidMetadata
	}

	return &dr, nil
}

func ParsingDateRangeFromURL(u url.Values, field string, startQuery, endQuery string) (*DateRange, error) {
	ts := u.Get(startQuery)
	te := u.Get(endQuery)
	if len(ts) == 0 || len(te) == 0 {
		return nil, nil
	}

	dr := DateRange{
		Field: "updated_at",
		Start: time.Time{},
		End:   time.Time{},
	}

	if v := u.Get(field); len(v) != 0 {
		dr.Field = strings.TrimSpace(strings.ToLower(v))
	}

	t, err := time.Parse("2006-01-02", ts)
	if err != nil {
		return nil, ErrInvalidMetadata
	}

	dr.Start = t

	t, err = time.Parse("2006-01-02", te)
	if err != nil {
		return nil, ErrInvalidMetadata
	}

	dr.End = t

	return &dr, nil
}

// Filter knows how to validate filterable fields.
// This Filter usually implemented by Repository.
type Filter interface {
	// Sortable returns true if a given field is allowed for sorting.
	Sortable(field string) bool
}

type FilterFull interface {
	// Sortable returns true if a given field is allowed for sorting.
	Sortable(field string) bool

	// Searchable return true if a given field is allowed for search
	Searchable(field string) bool
}

// ParseDateRangeFromURL
func ParseDateRangeFromURL(u url.Values) *DateRange {

	ts := u.Get("start")
	te := u.Get("end")
	if len(ts) == 0 || len(te) == 0 {
		return nil
	}

	dr := DateRange{
		Field: "created_at",
		Start: time.Time{},
		End:   time.Time{},
	}

	if v := u.Get("field"); len(v) != 0 {
		dr.Field = strings.TrimSpace(strings.ToLower(v))
	}

	t, err := time.Parse("2006-01-02", ts)
	if err != nil {
		return nil
	}

	dr.Start = t

	t, err = time.Parse("2006-01-02", te)
	if err != nil {
		return nil
	}

	dr.End = t

	return &dr
}

func ParsingFilterFromURL(u url.Values) Filtering {
	f := Filtering{
		OrderBy:   "created_at",
		OrderType: SortDescending,
	}

	ob := u.Get("order_by")
	if len(ob) != 0 {
		f.OrderBy = strings.ToLower(strings.ToLower(ob))
	}

	ot := u.Get("order_type")
	if len(ot) != 0 {
		ot = strings.TrimSpace(strings.ToLower(ot))
		if ot == SortAscending {
			f.OrderType = SortAscending
		}
	}

	search := strings.TrimSpace(u.Get("search"))
	if len(search) == 0 {
		search = strings.TrimSpace(u.Get("keyword"))
	}

	if len(search) != 0 {
		f.Search = search
	}

	searchBy := strings.TrimSpace(u.Get("search_by"))
	if len(search) == 0 {
		search = strings.TrimSpace(u.Get("search"))
	}

	if len(search) != 0 {
		f.SearchBy = searchBy
	}

	filter := strings.TrimSpace(u.Get("filter"))
	if len(filter) == 0 {
		filter = strings.TrimSpace(u.Get("keyword"))
	}

	if len(filter) != 0 {
		f.Filter = filter
	}

	filterBy := strings.TrimSpace(u.Get("filter_by"))
	if len(filter) == 0 {
		filter = strings.TrimSpace(u.Get("filter"))
	}

	if len(filter) != 0 {
		f.FilterBy = filterBy
	}

	return f
}

func hasDuplicatesArrayOfString(slice []string) bool {
	encountered := make(map[string]bool)

	for _, value := range slice {
		if encountered[value] {
			return true
		}
		encountered[value] = true
	}
	return false
}

type DateTimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func EmptyMetadata() Metadata {
	return Metadata{}
}

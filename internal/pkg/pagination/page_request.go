package pagination

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PageRequest struct {
	MaxResultCount int    `query:"maxResultCount" json:"maxResultCount,omitempty"`
	SkipCount      int    `query:"skipCount" json:"skipCount,omitempty"`
	Sorting        string `query:"sorting" json:"sorting,omitempty"`
	Filters        string `query:"filters" json:"filters,omitempty"`
}

func GetPageRequestFromCtx(c echo.Context) (*PageRequest, error) {
	res := &PageRequest{}

	//https://echo.labstack.com/guide/binding/#fast-binding-with-dedicated-helpers
	err := echo.QueryParamsBinder(c).
		Int("maxResultCount", &res.MaxResultCount).
		Int("skipCount", &res.SkipCount).
		String("filters", &res.Filters).
		String("sorting", &res.Sorting).
		BindError() // returns first binding error

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PageRequest) SanitizeSorting(validSortFields ...string) error {
	if p.Sorting == "" {
		return nil
	}

	sortFieldsMap := make(map[string]struct{}) // Initialize the map

	// Populate the map with the strings
	for _, str := range validSortFields {
		sortFieldsMap[strings.ToLower(str)] = struct{}{} // Use struct{} to save memory
	}

	// Remove empty space, asc, desc keywords
	sorting := strings.Replace(strings.ToLower(p.Sorting), " ", "", -1)
	sorting = strings.Replace(sorting, "asc", "", -1)
	sorting = strings.Replace(sorting, "desc", "", -1)

	// Check if the requested field is valid
	if _, ok := sortFieldsMap[sorting]; !ok {
		return errors.New("invalid sorting")
	}

	return nil
}

func (p *PageRequest) BuildFiltersExpr(likeFields ...string) clause.Expr {
	searchPattern := "%" + p.Filters + "%"
	likeConditions := make([]string, len(likeFields))

	for i, field := range likeFields {
		likeConditions[i] = field + " LIKE ?"
	}

	// Create a slice for the parameters
	params := make([]interface{}, len(likeConditions))
	for i := range likeConditions {
		params[i] = searchPattern
	}

	// Join conditions with OR
	return gorm.Expr(strings.Join(likeConditions, " OR "), params...)
}

func (p *PageRequest) Paginate(query *gorm.DB) *gorm.DB {
	return query.Offset(p.SkipCount).Limit(p.MaxResultCount)
}

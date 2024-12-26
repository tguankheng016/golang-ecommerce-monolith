package pagination

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/commerce-mono/pkg/core/helpers"
)

type PageRequest struct {
	MaxResultCount int    `query:"maxResultCount" json:"maxResultCount,omitempty"`
	SkipCount      int    `query:"skipCount" json:"skipCount,omitempty"`
	Sorting        string `query:"sorting" json:"sorting,omitempty"`
	Filters        string `query:"filters" json:"filters,omitempty"`
}

func GetPageRequestFromCtx(c echo.Context) (PageRequest, error) {
	res := PageRequest{}

	err := echo.QueryParamsBinder(c).
		Int("maxResultCount", &res.MaxResultCount).
		Int("skipCount", &res.SkipCount).
		String("filters", &res.Filters).
		String("sorting", &res.Sorting).
		BindError() // returns first binding error

	if err != nil {
		return res, err
	}

	return res, nil
}

func (p PageRequest) SanitizeSorting(validSortFields ...string) error {
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
	sorting = helpers.ReplaceLast(sorting, "asc", "")
	sorting = helpers.ReplaceLast(sorting, "desc", "")

	// Check if the requested field is valid
	if _, ok := sortFieldsMap[sorting]; !ok {
		return errors.New("invalid sorting field")
	}

	return nil
}

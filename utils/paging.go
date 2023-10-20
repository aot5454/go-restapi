package utils

import (
	"go-restapi/app"
	"strconv"
)

func (u *utils) GetPage(ctx app.Context) (int, error) {
	p := ctx.GetQuery("page")
	if p == "" {
		return 1, nil
	}

	page, err := strconv.Atoi(p)
	if err != nil {
		return 1, err
	}

	if page < 1 {
		return 1, nil
	}
	return page, nil
}

func (u *utils) GetPageSize(ctx app.Context) (int, error) {
	ps := ctx.GetQuery("pageSize")
	if ps == "" {
		return 20, nil
	}
	return strconv.Atoi(ps)
}

func (u *utils) GetTotalPage(total, pageSize int) int {
	if total == 0 {
		return 0
	}

	if total%pageSize == 0 {
		return total / pageSize
	}

	return total/pageSize + 1
}
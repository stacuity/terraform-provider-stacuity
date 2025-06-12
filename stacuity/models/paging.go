// Copyright (c) HashiCorp, Inc.

package models

type PagingState struct {
	Offset int32  `url:"offset"`
	Limit  int32  `url:"limit"`
	SortBy string `url:"sortBy"`
	Filter string `url:"filter"`
}

package handlers

type QueryParams struct {
	Page  uint   `form:"page,default=1" binding:"required"`
	Limit uint   `form:"limit,default=1" binding:"required"`
	Query string `form:"query"`
}

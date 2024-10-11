package database

import (
	"fmt"

	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	"gorm.io/gorm"
)

func registerCallBacks(db *gorm.DB) {
	db.Callback().Create().After("gorm:create").Register("assigned_created_by", assignedCreatedBy)
	db.Callback().Update().After("gorm:update").Register("assigned_updated_by", assignedUpdatedBy)
}

func assignedCreatedBy(db *gorm.DB) {
	ctx := db.Statement.Context
	fmt.Println(ctx)
}

func assignedUpdatedBy(db *gorm.DB) {
	ctx := db.Statement.Context
	userId, ok := ctx.Value(constants.CtxKey(constants.CurrentUserContextKey)).(int64)
	//userId, ok := ctx.Value("currentUser:userId").(int64)
	fmt.Println(userId)
	fmt.Println(ok)
	fmt.Println(ctx)
}

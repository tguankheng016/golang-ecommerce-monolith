package database

import (
	"database/sql"

	"github.com/tguankheng016/golang-ecommerce-monolith/internal/pkg/constants"
	"gorm.io/gorm"
)

func registerCallBacks(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("assigned_created_by", assignedCreatedBy)
	db.Callback().Update().Before("gorm:update").Register("assigned_updated_by", assignedUpdatedBy)
	db.Callback().Delete().Before("gorm:before_delete").Register("assigned_deleted_by", assignedDeletedBy)
}

func assignedCreatedBy(db *gorm.DB) {
	field := db.Statement.Schema.LookUpField("created_by")
	if field != nil {
		ctx := db.Statement.Context
		userId, ok := ctx.Value(constants.CtxKey(constants.CurrentUserContextKey)).(int64)
		if !ok {
			return
		}

		err := field.Set(ctx, db.Statement.ReflectValue, sql.NullInt64{Int64: userId, Valid: true})
		if err != nil {
			db.AddError(err)
		}
	}
}

func assignedUpdatedBy(db *gorm.DB) {
	field := db.Statement.Schema.LookUpField("updated_by")
	if field != nil {
		ctx := db.Statement.Context
		userId, ok := ctx.Value(constants.CtxKey(constants.CurrentUserContextKey)).(int64)
		if !ok {
			return
		}

		err := field.Set(ctx, db.Statement.ReflectValue, sql.NullInt64{Int64: userId, Valid: true})
		if err != nil {
			db.AddError(err)
		}
	}
}

func assignedDeletedBy(db *gorm.DB) {
	field := db.Statement.Schema.LookUpField("deleted_by")
	if field != nil {
		ctx := db.Statement.Context
		userId, ok := ctx.Value(constants.CtxKey(constants.CurrentUserContextKey)).(int64)
		if !ok {
			return
		}

		err := field.Set(ctx, db.Statement.ReflectValue, sql.NullInt64{Int64: userId, Valid: true})
		if err != nil {
			db.AddError(err)
		}
	}
}

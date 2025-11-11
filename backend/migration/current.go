package migration

import (
	"github.com/Luke256/ducks/model"

	"github.com/go-gormigrate/gormigrate/v2"
)

// Migrations 全てのデータベースマイグレーション
//
// 新たなマイグレーションを行う場合は、この配列の末尾に必ず追加すること
func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{}
}

// AllTables 最新のスキーマの全テーブルモデル
//
// 最新のスキーマの全テーブルのモデル構造体を記述すること
func AllTables() []interface{} {
	return []interface{}{
		&model.Poster{},
		&model.Festival{},
	}
}
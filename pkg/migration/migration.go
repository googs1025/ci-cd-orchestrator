package migration

import (
	"database/sql"
	"log"
	"os"
	"strings"
)

// Run 执行数据库迁移
func Run(db *sql.DB) error {
	// 读取 schema.sql 文件
	// 尝试从多个位置读取 schema.sql
	var schema []byte
	var err error

	// 先尝试当前目录
	schema, err = os.ReadFile("schema.sql")
	if err != nil {
		// 再尝试上级目录
		schema, err = os.ReadFile("../schema.sql")
		if err != nil {
			// 最后尝试上上级目录
			schema, err = os.ReadFile("../../schema.sql")
			if err != nil {
				return err
			}
		}
	}

	// 分割 SQL 语句
	statements := strings.Split(string(schema), ";")

	// 执行每个 SQL 语句
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}

	log.Println("数据库迁移成功")
	return nil
}

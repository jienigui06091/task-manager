// package database 集中处理 PostgreSQL 数据库连接。
package database

// import (...) 导入创建连接池需要的标准库和第三方包。
import (
	"context" // context 用于传递取消信号、超时等请求范围的信息。
	"fmt"     // fmt 用于格式化错误消息。
	"time"    // time 用于设置连接存活时间。

	"github.com/jackc/pgx/v5/pgxpool" // pgxpool 提供 PostgreSQL 数据库连接池。
)

// NewPostgresPool 创建、配置并验证 PostgreSQL 连接池。
// 返回 *pgxpool.Pool 指向连接池，error 用于向调用者报告失败原因。
func NewPostgresPool() (*pgxpool.Pool, error) {
	// := 声明 dsn；DSN（数据源名称）描述连接数据库所需的信息。
	dsn := "postgres://postgres:postgres@localhost:5432/taskDB"

	// ParseConfig 根据 DSN 解析出可继续修改的连接池配置，并可能返回错误。
	config, err := pgxpool.ParseConfig(dsn)
	// err != nil 表示 DSN 解析失败，例如格式不正确。
	if err != nil {
		// return 返回 nil（没有连接池）以及原始错误。
		return nil, err
	}

	// 最大同时连接数设为 10。
	config.MaxConns = 10
	// 最小空闲连接数设为 2。
	config.MinConns = 2
	// 每条连接最长使用一小时，之后由连接池替换。
	config.MaxConnLifetime = time.Hour

	// context.Background 创建根上下文；NewWithConfig 依据配置实际创建连接池。
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	// 创建失败时进入该分支，例如数据库不可访问。
	if err != nil {
		// fmt.Errorf 格式化错误；%w 包装原始 err，仍可用 errors.Is/As 检查它。
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	// Ping 向数据库发送测试请求，确认连接池确实可用。
	if err := pool.Ping(context.Background()); err != nil {
		// 测试失败时关闭已创建的连接池，避免资源泄漏。
		pool.Close()
		// 返回带上下文信息的错误。
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	// 全部成功后返回可用连接池和 nil（没有错误）。
	return pool, nil
}

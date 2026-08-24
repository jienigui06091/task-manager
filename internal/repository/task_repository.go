// package repository 存放直接执行数据库查询的代码。
package repository

// import (...) 导入访问数据库和任务模型所需的包。
import (
	"context"                         // context 把请求的取消信号等信息传递给数据库操作。
	"github.com/jackc/pgx/v5/pgxpool" // pgxpool.Pool 是 PostgreSQL 连接池类型。

	"task-manager/internal/model" // model.Task 是项目定义的任务数据结构。
)

// TaskRepository 是任务数据访问对象；它持有数据库连接池。
type TaskRepository struct {
	// db 是指向连接池的指针；* 表示保存地址而不是复制整个连接池。
	db *pgxpool.Pool
}

// NewTaskRepository 用给定的连接池创建一个任务仓储对象。
func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	// &TaskRepository{...} 创建结构体并取其地址，返回 *TaskRepository 指针。
	return &TaskRepository{
		// 左侧 db 是字段名，右侧 db 是函数参数。
		db: db,
	}
}

// GetAll 查询所有任务；ctx 是请求上下文，返回任务切片或错误。
// (r *TaskRepository) 是方法接收者，表示通过 r 操作某个仓储实例。
func (r *TaskRepository) GetAll(ctx context.Context, userId int64, page int, pageSize int) ([]model.Task, int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM tasks where user_id = $1`, userId).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	// Query 执行可返回多行的 SQL；反引号定义原始字符串，适合书写多行 SQL。
	rows, err := r.db.Query(ctx, `
		-- SELECT 指定要读取的列。
		SELECT id, user_id,title, completed, created_at, updated_at
		-- FROM 指定读取 tasks 表。
		FROM tasks
		-- ORDER BY id DESC 按 id 降序排列，通常让较新的任务在前。
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`, userId, pageSize, offset)
	// 查询失败时没有结果可处理，立刻返回错误。
	if err != nil {
		// nil 表示没有任务切片，err 交给调用者处理。
		return nil, 0, err
	}
	// rows 用完后必须关闭；defer 保证函数结束时一定执行关闭。
	defer rows.Close()

	// var 只声明变量并赋予零值；nil 切片仍可用 append 追加元素。
	var tasks []model.Task

	// Next 每次移动到下一条查询结果；有下一行时返回 true。
	for rows.Next() {
		// 为当前数据库行声明一个空的 Task 结构体。
		var task model.Task

		// Scan 把当前行的列按顺序写入字段地址，顺序必须与 SELECT 一致。
		err := rows.Scan(
			// &task.ID 取字段地址，使 Scan 可以修改 ID。
			&task.ID,
			&task.UserId,
			// 将第二列 title 写入任务标题。
			&task.Title,
			// 将第三列 completed 写入完成状态。
			&task.Completed,
			// 将第四列 created_at 写入创建时间。
			&task.CreatedAt,
			// 将第五列 updated_at 写入更新时间。
			&task.UpdatedAt,
		)
		// 当前行的数据不能转换为 Go 类型时，返回错误。
		if err != nil {
			return nil, 0, err
		}

		// append 把当前任务追加到切片，并将可能扩容后的新切片赋回 tasks。
		tasks = append(tasks, task)
	}

	// rows.Err 检查遍历过程中出现的延迟错误，例如数据库连接中断。
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// 成功时返回任务列表和 nil 错误。
	return tasks, total, nil
}
func (r *TaskRepository) GetByID(ctx context.Context, userId int64, taskId int64) (model.Task, error) {
	var task model.Task

	err := r.db.QueryRow(ctx, `
		SELECT id, user_id,title, completed, created_at, updated_at
		FROM tasks
		WHERE id = $1 and user_id = $2
	`, taskId, userId).Scan(
		&task.ID,
		&task.UserId,
		&task.Title,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

// Create 将一个任务写入数据库，并返回数据库生成 ID 和时间后的完整任务。
func (r *TaskRepository) Create(
	// ctx 将取消信号和截止时间传给数据库操作。
	ctx context.Context,
	db DBTX,
	// task 是调用者提供的待创建任务；查询结果会写回这个局部变量。
	task model.Task,
	// 该方法返回一个 Task 和一个 error。
) (model.Task, error) {

	// QueryRow 执行期望只返回一行的 SQL；:= 在函数内创建 err 变量。
	err := db.QueryRow(ctx, `
		-- INSERT INTO 指定要向 tasks 表新增记录。
		INSERT INTO tasks (user_id,title, completed)
		-- $1、$2 是参数占位符，避免直接拼接用户输入造成 SQL 注入。
		VALUES ($1, $2, $3)
		-- RETURNING 要求 PostgreSQL 返回刚插入记录的完整字段。
		RETURNING id, user_id,title, completed, created_at, updated_at
	`,
		// 第一个 SQL 参数替换 $1，即任务标题。
		task.UserId,
		task.Title,
		// 第二个 SQL 参数替换 $2，即任务完成状态。
		task.Completed,
	// Scan 读取 INSERT ... RETURNING 返回的唯一一行，填充 task 的字段。
	).Scan(
		// 将数据库返回的 id 写入 task.ID。
		&task.ID,
		&task.UserId,
		// 将数据库返回的 title 写入 task.Title。
		&task.Title,
		// 将数据库返回的 completed 写入 task.Completed。
		&task.Completed,
		// 将数据库返回的 created_at 写入 task.CreatedAt。
		&task.CreatedAt,
		// 将数据库返回的 updated_at 写入 task.UpdatedAt。
		&task.UpdatedAt,
	)

	// 插入或读取返回行失败时进入该分支。
	if err != nil {
		// model.Task{} 是各字段均为零值的空任务；同时返回实际错误。
		return model.Task{}, err
	}

	// 成功时返回包含数据库生成字段的任务和 nil 错误。

	return task, nil
}
func (r *TaskRepository) Update(ctx context.Context, task model.Task) (model.Task, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE tasks
		SET title = $2, completed = $3, updated_at = NOW()
		where id = $1 and user_id = $4
		RETURNING id, user_id,title, completed, created_at, updated_at`, task.ID, task.Title, task.Completed, task.UserId).Scan(
		&task.ID,
		&task.UserId,
		&task.Title,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) DeleteByList(ctx context.Context, ids []int64, userId int64) error {
	_, err := r.db.Exec(ctx, `
	DELETE FROM tasks where id = ANY ($1) and user_id = $2
	`, ids, userId)
	return err
}

func (r *TaskRepository) Duplicate(ctx context.Context, userId int64, taskId int64) error {
	task, err := r.GetByID(ctx, userId, taskId)
	if err != nil {
		return err
	}
	taskEn := model.Task{
		UserId:    userId,
		Completed: task.Completed,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
	errT := r.db.QueryRow(ctx, `
		-- INSERT INTO 指定要向 tasks 表新增记录。
		INSERT INTO tasks (user_id,title, completed,created_at,updated_at)
		-- $1、$2 是参数占位符，避免直接拼接用户输入造成 SQL 注入。
		VALUES ($1, $2, $3,$4,$5)
		-- RETURNING 要求 PostgreSQL 返回刚插入记录的完整字段。
		RETURNING id, user_id,title, completed, created_at, updated_at
	`,
		// 第一个 SQL 参数替换 $1，即任务标题。
		taskEn.UserId,
		taskEn.Title,
		// 第二个 SQL 参数替换 $2，即任务完成状态。
		taskEn.Completed,
		taskEn.CreatedAt,
		taskEn.UpdatedAt,
	// Scan 读取 INSERT ... RETURNING 返回的唯一一行，填充 task 的字段。
	).Scan(
		// 将数据库返回的 id 写入 task.ID。
		&task.ID,
		&task.UserId,
		// 将数据库返回的 title 写入 task.Title。
		&task.Title,
		// 将数据库返回的 completed 写入 task.Completed。
		&task.Completed,
		// 将数据库返回的 created_at 写入 task.CreatedAt。
		&task.CreatedAt,
		// 将数据库返回的 updated_at 写入 task.UpdatedAt。
		&task.UpdatedAt,
	)
	if errT != nil {
		return errT
	}
	return nil
}

// package service 存放业务规则，并协调仓储层与 HTTP 处理器层。
package service

// import (...) 导入业务层需要的标准库和项目内部包。
import (
	"context" // context 用于把请求生命周期传递给仓储层。
	"errors"  // errors 用于创建普通错误值。

	"github.com/jackc/pgx/v5"

	"task-manager/internal/apperror"
	"task-manager/internal/model"      // model 包定义 Task 数据结构。
	"task-manager/internal/repository" // repository 包负责数据库操作。
)

var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrUpdateTaskFailed = errors.New("failed to update task")
)

// TaskService 是任务业务服务；它通过 repository 字段访问数据。
type TaskService struct {
	// repository 保存任务仓储的指针，使服务可调用它的方法。
	repository *repository.TaskRepository
}

// NewTaskService 创建业务服务，并把仓储依赖保存到结构体中。
func NewTaskService(repository *repository.TaskRepository) *TaskService {
	// 返回新结构体的地址；左侧 repository 是字段，右侧是函数参数。
	return &TaskService{
		repository: repository,
	}
}

// GetAllTasks 是获取全部任务的业务方法；目前没有额外规则，直接委托给仓储层。
func (s *TaskService) GetAllTasks(
	// ctx 会继续传给数据库操作，使请求取消能及时生效。
	ctx context.Context,
	// 返回任务切片和错误。
) ([]model.Task, error) {
	// s.repository 访问接收者字段，再调用其 GetAll 方法。
	return s.repository.GetAll(ctx)
}

// CreateTask 根据标题创建任务，并在写入数据库前验证输入。
func (s *TaskService) CreateTask(
	// ctx 会传给后续数据库操作。
	ctx context.Context,
	// title 是客户端传来的任务标题。
	title string,
	// 返回创建成功的任务或错误。
) (model.Task, error) {

	// == 比较两个值是否相等；空字符串不允许作为任务标题。
	if title == "" {
		// errors.New 创建错误；空结构体和错误一起返回以表示创建失败。
		return model.Task{}, errors.New("title is required")
	}

	// 使用结构体字面量创建任务；未列出的字段自动使用对应类型的零值。
	task := model.Task{
		// 设置任务标题为函数参数 title。
		Title: title,
		// 新任务默认尚未完成。
		Completed: false,
	}

	// 调用仓储层写入数据库，并直接返回其两个结果。
	return s.repository.Create(ctx, task)
}
func (s *TaskService) GetTaskByID(ctx context.Context, taskId int64) (model.Task, error) {
	task, err := s.repository.GetByID(ctx, taskId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, apperror.ErrNotFound
		}
		return model.Task{}, err
	}
	return task, nil
}

func (s *TaskService) UpdateTask(
	ctx context.Context,
	id int64,
	title string,
	completed bool,
) (model.Task, error) {
	task := model.Task{
		ID:        id,
		Title:     title,
		Completed: completed,
	}

	updatedTask, err := s.repository.Update(ctx, task)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, apperror.ErrNotFound
		}

		return model.Task{}, apperror.ErrInvalid
	}

	return updatedTask, nil
}

func (s *TaskService) DeleteByList(ctx context.Context, ids []int64) error {

	err := s.repository.DeleteByList(ctx, ids)
	if err != nil {
		return err
	}

	return nil

}

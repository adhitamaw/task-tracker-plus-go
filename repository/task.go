package repository

import (
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/model"
)

type TaskRepository interface {
	Store(task *model.Task) error
	Update(taskID int, task *model.Task) error
	Delete(id int) error
	GetByID(id int) (*model.Task, error)
	GetList(userID int) ([]model.Task, error)
	GetTaskCategory(id int) ([]model.TaskCategory, error)
	GetTaskCategoryByUser(categoryID, userID int) ([]model.TaskCategory, error)
}

type taskRepository struct {
	filebased *filebased.Data
}

func NewTaskRepo(filebasedDb *filebased.Data) *taskRepository {
	return &taskRepository{
		filebased: filebasedDb,
	}
}

func (t *taskRepository) Store(task *model.Task) error {
	err := t.filebased.StoreTask(*task)
	return err
}

func (t *taskRepository) Update(taskID int, task *model.Task) error {
	err := t.filebased.UpdateTask(task.ID, *task)
	return err
}

func (t *taskRepository) Delete(id int) error {
	err := t.filebased.DeleteTask(id)
	return err
}

func (t *taskRepository) GetByID(id int) (*model.Task, error) {
	task, err := t.filebased.GetTaskByID(id)

	return task, err // TODO: replace this
}

func (t *taskRepository) GetList(userID int) ([]model.Task, error) {
	tasks, err := t.filebased.GetTasksByUserID(userID)
	return tasks, err
}

func (t *taskRepository) GetTaskCategory(id int) ([]model.TaskCategory, error) {
	taskCategories, err := t.filebased.GetTaskListByCategory(id)

	return taskCategories, err // TODO: replace this
}

func (t *taskRepository) GetTaskCategoryByUser(categoryID, userID int) ([]model.TaskCategory, error) {
	taskCategories, err := t.filebased.GetTaskListByCategoryAndUser(categoryID, userID)
	return taskCategories, err
}

package todo

import "sync"

type List struct {
	tasksMap map[string]Task
	mtx      sync.RWMutex
}

func NewList() *List {
	return &List{
		tasksMap: make(map[string]Task),
	}
}

func (l *List) AddTask(task Task) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.tasksMap[task.Title]; ok {
		return ErrTaskAlreadyExists
	}

	l.tasksMap[task.Title] = task

	return nil
}

func (l *List) GetTask(title string) (Task, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	task, ok := l.tasksMap[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (l *List) ListTasks() map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	tmp := make(map[string]Task, len(l.tasksMap))

	for k, v := range l.tasksMap {
		tmp[k] = v
	}

	return tmp
}

func (l *List) ListNotCompletedTask() map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	notCompletedTask := make(map[string]Task)

	for title, task := range l.tasksMap {
		if !task.Completed {
			notCompletedTask[title] = task
		}
	}

	return notCompletedTask
}

func (l *List) CompleteTask(title string) (Task, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	task, ok := l.tasksMap[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	task.Complete()

	l.tasksMap[title] = task

	return l.tasksMap[title], nil
}

func (l *List) UncompleteTask(title string) (Task, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	task, ok := l.tasksMap[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	task.Uncomplete()

	l.tasksMap[title] = task

	return l.tasksMap[title], nil
}

func (l *List) DeleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	_, ok := l.tasksMap[title]
	if !ok {
		return ErrTaskNotFound
	}

	delete(l.tasksMap, title)

	return nil
}

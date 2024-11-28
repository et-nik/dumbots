package main

import "github.com/et-nik/metamod-go/vector"

type BotTask interface {
	Finished() bool
	SetFinished()

	StartedTime() float32
	SetStartedTime(time float32)
}

//
//// SequenceTask is a task that runs a sequence of tasks
//type SequenceTask struct {
//	Tasks []BotTask
//}
//
//func NewSequenceTask(tasks ...BotTask) *SequenceTask {
//	return &SequenceTask{
//		Tasks: tasks,
//	}
//}
//
//func (t *SequenceTask) Run(b *Bot) error {
//	for _, task := range t.Tasks {
//		if task.Finished() {
//			continue
//		}
//
//		err := task.Run(b)
//		if err != nil {
//			return errors.WithMessage(err, "failed to run task")
//		}
//	}
//
//	return nil
//}
//
//func (t *SequenceTask) Finished() bool {
//	for _, task := range t.Tasks {
//		if !task.Finished() {
//			return false
//		}
//	}
//
//	return true
//}
//
//type FindEnemyTask struct {
//	Globals *metamod.GlobalVars
//
//	Funcs  *PluginFuncs
//	Logger *Logger
//
//	startedTime float32
//	finished    bool
//}
//
//func NewFindEnemyTask(
//	globals *metamod.GlobalVars,
//	funcs *PluginFuncs,
//	logger *Logger,
//) *FindEnemyTask {
//	return &FindEnemyTask{
//		Globals: globals,
//		Funcs:   funcs,
//		Logger:  logger,
//
//		startedTime: globals.Time(),
//	}
//}
//
//func (t *FindEnemyTask) Run(b *Bot) error {
//	return nil
//}
//
//func (t *FindEnemyTask) Finished() bool {
//	return t.finished
//}

type baseTask struct {
	startedTime float32
	completed   bool
}

func (t *baseTask) Finished() bool {
	return t.completed
}

func (t *baseTask) SetFinished() {
	t.completed = true
}

func (t *baseTask) StartedTime() float32 {
	return t.startedTime
}

func (t *baseTask) SetStartedTime(time float32) {
	t.startedTime = time
}

type BotTaskMoveToPoint struct {
	baseTask

	TimeToFail float32
	Target     vector.Vector
}

type BotTaskMoveToGraphVertex struct {
	baseTask

	TimeToFail float32

	Target Vertex
}

type BotTaskMoveToSecret struct {
	baseTask

	Target vector.Vector
}

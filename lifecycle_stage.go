package pipeline

// BeforeStage is an alias for before hooks of a stage about to be executed with a given context.
// If the hook fails, the stage won't be executed
type BeforeStage func(stage Stage, ctx Context) error

// AfterStage is an alias for after hooks of a stage. If the stage fails, one can recover from here or fallback to a new error.
// Also, this stage can fail, thus failing the execution (note that this is a blob of a stage, so if a hook fails, the stage fails too).
// The provided context was the resulting one after the stage was executed
type AfterStage func(stage Stage, ctx Context, err error) error

// Blob structure that allows us to decorate a stage with pre/post hooks
// Note that we can compose many lifecycle stages if we want to have multiple hooks. Such as:
// lifecycleStage := CreateLifecycleStage(realStage, aBeforeHook, anAfterHook)
// lifecycleStage = CreateBeforeStageLifecycle(lifecycleStage, anotherBeforeHook)
// lifecycleStage = CreateAfterStageLifecycle(lifecycleStage, anotherAfterHook)
type lifecycleStage struct {
	Before BeforeStage
	After  AfterStage
	Stage  Stage
}

func (l *lifecycleStage) Draw(graph GraphDiagram) {
	l.Stage.Draw(graph)
}

// Run the hooks and the stage, validating errors along the way and mutating the stage error in case it failed.
func (l *lifecycleStage) Run(executor Executor, ctx Context) error {
	if l.Before != nil {
		err := l.Before(l.Stage, ctx)

		if err != nil {
			return err
		}
	}

	err := l.Stage.Run(executor, ctx)

	if l.After != nil {
		err = l.After(l.Stage, ctx, err)
	}
	return err
}

// CreateBeforeStageLifecycle creates a lifecycle stage with a before hook
func CreateBeforeStageLifecycle(stage Stage, before BeforeStage) Stage {
	return &lifecycleStage{
		Before: before,
		Stage:  stage,
	}
}

// CreateAfterStageLifecycle creates a lifecycle stage with an after hook
func CreateAfterStageLifecycle(stage Stage, after AfterStage) Stage {
	return &lifecycleStage{
		After: after,
		Stage: stage,
	}
}

// CreateStageLifecycle creates a lifecycle stage with a before and an after hook
func CreateStageLifecycle(stage Stage, before BeforeStage, after AfterStage) Stage {
	return &lifecycleStage{
		Before: before,
		After:  after,
		Stage:  stage,
	}
}

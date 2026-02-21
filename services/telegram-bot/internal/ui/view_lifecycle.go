package ui

import "context"

type ViewLifecycleStatus string

const (
	VLS_Empty   ViewLifecycleStatus = "empty"
	VLS_Created ViewLifecycleStatus = "created"
	VLS_Started ViewLifecycleStatus = "started"
)

type ViewLifecycle interface {
	OnCreate(context.Context) error
	OnStart(context.Context) error
	OnResume(context.Context) error
}

type ViewLifecycleController interface {
	NewLifeCycle(context.Context, View) error
	ContinueLifecycle(context.Context, View) error
}

type ViewLifecycleStore interface {
	GetViewLifecycleStatus(ctx context.Context, id ViewIdentifier) (ViewLifecycleStatus, error)
	StoreViewLifecycleStatus(ctx context.Context, id ViewIdentifier, vlc ViewLifecycleStatus) error
}

type viewLifecycleController struct {
	store ViewLifecycleStore
}

func NewViewLifecycleController(
	store ViewLifecycleStore,
) *viewLifecycleController {
	return &viewLifecycleController{
		store: store,
	}
}

func (c *viewLifecycleController) NewLifeCycle(ctx context.Context, v View) error {
	vid := v.ViewIdentifier()
	status, err := c.create(ctx, v)
	if err != nil {
		return err
	}
	err = c.store.StoreViewLifecycleStatus(ctx, vid, status)
	if err != nil {
		return err
	}
	return nil
}

func (c *viewLifecycleController) ContinueLifecycle(ctx context.Context, v View) error {
	vid := v.ViewIdentifier()
	status, err := c.store.GetViewLifecycleStatus(ctx, vid)
	if err != nil {
		return err
	}
	switch status {
	case VLS_Empty:
		status, err = c.create(ctx, v)
	case VLS_Started:
		status, err = c.resume(ctx, v)
	}
	if err != nil {
		return err
	}
	err = c.store.StoreViewLifecycleStatus(ctx, vid, status)
	if err != nil {
		return err
	}
	return nil
}

func (c *viewLifecycleController) create(ctx context.Context, v View) (ViewLifecycleStatus, error) {
	err := v.OnCreate(ctx)
	if err != nil {
		return VLS_Empty, err
	}
	status, err := c.start(ctx, v)
	if err != nil {
		return VLS_Created, err
	}
	return status, nil
}

func (c *viewLifecycleController) start(ctx context.Context, v View) (ViewLifecycleStatus, error) {
	err := v.OnStart(ctx)
	if err != nil {
		return VLS_Empty, err
	}
	return VLS_Started, nil
}

func (c *viewLifecycleController) resume(ctx context.Context, v View) (ViewLifecycleStatus, error) {
	err := v.OnResume(ctx)
	if err != nil {
		return VLS_Empty, err
	}
	return VLS_Started, nil
}

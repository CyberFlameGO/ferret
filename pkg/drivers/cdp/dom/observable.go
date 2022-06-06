package dom

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
)

type Observable struct {
	id   runtime.RemoteObjectID
	eval *eval.Runtime
}

func NewObservable(id runtime.RemoteObjectID, eval *eval.Runtime) *Observable {
	return &Observable{id, eval}
}

func (observer *Observable) WaitForElement(ctx context.Context, selector drivers.QuerySelector, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForElement(observer.id, selector, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (observer *Observable) WaitForElementAll(ctx context.Context, selector drivers.QuerySelector, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForElementAll(observer.id, selector, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (observer *Observable) WaitForClass(ctx context.Context, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForClass(observer.id, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (observer *Observable) WaitForClassBySelector(ctx context.Context, selector drivers.QuerySelector, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForClassBySelector(observer.id, selector, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (observer *Observable) WaitForClassBySelectorAll(ctx context.Context, selector drivers.QuerySelector, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForClassBySelectorAll(observer.id, selector, class, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (observer *Observable) WaitForAttribute(
	ctx context.Context,
	name values.String,
	value core.Value,
	when drivers.WaitEvent,
) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForAttribute(observer.id, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (observer *Observable) WaitForAttributeBySelector(ctx context.Context, selector drivers.QuerySelector, name values.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForAttributeBySelector(observer.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (observer *Observable) WaitForAttributeBySelectorAll(ctx context.Context, selector drivers.QuerySelector, name values.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForAttributeBySelectorAll(observer.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (observer *Observable) WaitForStyle(ctx context.Context, name values.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForStyle(observer.id, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (observer *Observable) WaitForStyleBySelector(ctx context.Context, selector drivers.QuerySelector, name values.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForStyleBySelector(observer.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (observer *Observable) WaitForStyleBySelectorAll(ctx context.Context, selector drivers.QuerySelector, name values.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		observer.eval,
		templates.WaitForStyleBySelectorAll(observer.id, selector, name, value, when),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

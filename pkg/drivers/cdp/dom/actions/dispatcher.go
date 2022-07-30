package actions

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/input"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/protocol/runtime"
	"time"
)

type Dispatcher struct {
	id    runtime.RemoteObjectID
	input *input.Manager
}

func NewDispatcher(id runtime.RemoteObjectID, input *input.Manager) *Dispatcher {
	return &Dispatcher{id, input}
}

func (dispatcher *Dispatcher) Dispatch(ctx context.Context, action events.Action) (core.Value, error) {
	opts := action.Options
	selector := getSelector(opts)

	switch action.Name {
	case "click":
		count := getCount(opts, 1)

		if drivers.IsValidSelector(selector) {
			return values.None, dispatcher.ClickBySelector(ctx, selector, count)
		}

		return values.None, dispatcher.Click(ctx, count)
	case "input":
		text := values.ToString(action.Args)
		delay := getDelay(opts, drivers.DefaultKeyboardDelay)

		if drivers.IsValidSelector(selector) {
			return values.None, dispatcher.InputBySelector(ctx, selector, text, delay)
		}

		return values.None, dispatcher.Input(ctx, text, delay)
	case "press":
		count := getCount(opts, 1)
		keys := values.ToStrings(values.ToArray(ctx, action.Args))

		if drivers.IsValidSelector(selector) {
			return values.None, dispatcher.PressBySelector(ctx, selector, keys, count)
		}

		return values.None, dispatcher.Press(ctx, keys, count)
	case "clear":
		if drivers.IsValidSelector(selector) {
			return values.None, dispatcher.ClearBySelector(ctx, selector)
		}

		return values.None, dispatcher.Clear(ctx)
	case "select":
		selectOptions := values.ToArray(ctx, action.Args)

		if drivers.IsValidSelector(selector) {
			return dispatcher.SelectBySelector(ctx, selector, selectOptions)
		}

		return dispatcher.Select(ctx, selectOptions)
	case "scroll":
		scrollOptions := getScrollOptions(opts)

		if drivers.IsValidSelector(selector) {
			return nil, dispatcher.ScrollBySelector(ctx, selector, scrollOptions)
		}

		return nil, dispatcher.Scroll(ctx, scrollOptions)
	case "focus":
		if drivers.IsValidSelector(selector) {
			return nil, dispatcher.FocusBySelector(ctx, selector)
		}

		return nil, dispatcher.Focus(ctx)
	case "blur":
		if drivers.IsValidSelector(selector) {
			return nil, dispatcher.BlurBySelector(ctx, selector)
		}

		return nil, dispatcher.Blur(ctx)
	case "hover":
		if drivers.IsValidSelector(selector) {
			return nil, dispatcher.HoverBySelector(ctx, selector)
		}

		return nil, dispatcher.Hover(ctx)
	default:
		return values.None, core.Errorf(core.ErrInvalidOperation, "unknown event: %s", action.Name)
	}
}

func (dispatcher *Dispatcher) Click(ctx context.Context, count values.Int) error {
	return dispatcher.input.Click(ctx, dispatcher.id, int(count))
}

func (dispatcher *Dispatcher) ClickBySelector(ctx context.Context, selector drivers.QuerySelector, count values.Int) error {
	return dispatcher.input.ClickBySelector(ctx, dispatcher.id, selector, count)
}

func (dispatcher *Dispatcher) Input(ctx context.Context, value core.Value, delay values.Int) error {
	return dispatcher.input.Type(ctx, dispatcher.id, input.TypeParams{
		Text:  value.String(),
		Clear: false,
		Delay: time.Duration(delay) * time.Millisecond,
	})
}

func (dispatcher *Dispatcher) InputBySelector(ctx context.Context, selector drivers.QuerySelector, value core.Value, delay values.Int) error {
	return dispatcher.input.TypeBySelector(ctx, dispatcher.id, selector, input.TypeParams{
		Text:  value.String(),
		Clear: false,
		Delay: time.Duration(delay) * time.Millisecond,
	})
}

func (dispatcher *Dispatcher) Press(ctx context.Context, keys []values.String, count values.Int) error {
	return dispatcher.input.Press(ctx, values.UnwrapStrings(keys), int(count))
}

func (dispatcher *Dispatcher) PressBySelector(ctx context.Context, selector drivers.QuerySelector, keys []values.String, count values.Int) error {
	return dispatcher.input.PressBySelector(ctx, dispatcher.id, selector, values.UnwrapStrings(keys), int(count))
}

func (dispatcher *Dispatcher) Clear(ctx context.Context) error {
	return dispatcher.input.Clear(ctx, dispatcher.id)
}

func (dispatcher *Dispatcher) ClearBySelector(ctx context.Context, selector drivers.QuerySelector) error {
	return dispatcher.input.ClearBySelector(ctx, dispatcher.id, selector)
}

func (dispatcher *Dispatcher) Select(ctx context.Context, value *values.Array) (*values.Array, error) {
	return dispatcher.input.Select(ctx, dispatcher.id, value)
}

func (dispatcher *Dispatcher) SelectBySelector(ctx context.Context, selector drivers.QuerySelector, value *values.Array) (*values.Array, error) {
	return dispatcher.input.SelectBySelector(ctx, dispatcher.id, selector, value)
}

func (dispatcher *Dispatcher) Scroll(ctx context.Context, options drivers.ScrollOptions) error {
	return dispatcher.input.ScrollByXY(ctx, options)
}

func (dispatcher *Dispatcher) ScrollIntoView(ctx context.Context, options drivers.ScrollOptions) error {
	return dispatcher.input.ScrollIntoView(ctx, dispatcher.id, options)
}

func (dispatcher *Dispatcher) ScrollTop(ctx context.Context, options drivers.ScrollOptions) error {
	return dispatcher.input.ScrollTop(ctx, options)
}

func (dispatcher *Dispatcher) ScrollBottom(ctx context.Context, options drivers.ScrollOptions) error {
	return dispatcher.input.ScrollBottom(ctx, options)
}

func (dispatcher *Dispatcher) ScrollBySelector(ctx context.Context, selector drivers.QuerySelector, options drivers.ScrollOptions) error {
	return dispatcher.input.ScrollIntoViewBySelector(ctx, dispatcher.id, selector, options)
}

func (dispatcher *Dispatcher) Focus(ctx context.Context) error {
	return dispatcher.input.Focus(ctx, dispatcher.id)
}

func (dispatcher *Dispatcher) FocusBySelector(ctx context.Context, selector drivers.QuerySelector) error {
	return dispatcher.input.FocusBySelector(ctx, dispatcher.id, selector)
}

func (dispatcher *Dispatcher) Blur(ctx context.Context) error {
	return dispatcher.input.Blur(ctx, dispatcher.id)
}

func (dispatcher *Dispatcher) BlurBySelector(ctx context.Context, selector drivers.QuerySelector) error {
	return dispatcher.input.BlurBySelector(ctx, dispatcher.id, selector)
}

func (dispatcher *Dispatcher) Hover(ctx context.Context) error {
	return dispatcher.input.MoveMouse(ctx, dispatcher.id)
}

func (dispatcher *Dispatcher) HoverBySelector(ctx context.Context, selector drivers.QuerySelector) error {
	return dispatcher.input.MoveMouseBySelector(ctx, dispatcher.id, selector)
}

func (dispatcher *Dispatcher) MoveMouseByXY(ctx context.Context, x, y values.Float) error {
	return dispatcher.input.MoveMouseByXY(ctx, x, y)
}

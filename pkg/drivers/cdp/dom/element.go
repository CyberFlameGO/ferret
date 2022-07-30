package dom

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/dom/actions"
	"hash/fnv"
	"strings"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/input"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTMLElement struct {
	*Observable
	*actions.Dispatcher
	logger   zerolog.Logger
	client   *cdp.Client
	dom      *Manager
	id       runtime.RemoteObjectID
	nodeType *common.LazyValue
	nodeName *common.LazyValue
}

func NewHTMLElement(
	logger zerolog.Logger,
	client *cdp.Client,
	domManager *Manager,
	input *input.Manager,
	exec *eval.Runtime,
	id runtime.RemoteObjectID,
) *HTMLElement {
	el := new(HTMLElement)
	el.Observable = NewObservable(id, exec)
	el.Dispatcher = actions.NewDispatcher(id, input)
	el.logger = logging.
		WithName(logger.With(), "dom_element").
		Str("object_id", string(id)).
		Logger()
	el.client = client
	el.dom = domManager
	el.eval = exec
	el.id = id
	el.nodeType = common.NewLazyValue(func(ctx context.Context) (core.Value, error) {
		return el.eval.EvalValue(ctx, templates.GetNodeType(el.id))
	})
	el.nodeName = common.NewLazyValue(func(ctx context.Context) (core.Value, error) {
		return el.eval.EvalValue(ctx, templates.GetNodeName(el.id))
	})

	return el
}

func (el *HTMLElement) RemoteID() runtime.RemoteObjectID {
	return el.id
}

func (el *HTMLElement) Close() error {
	return nil
}

func (el *HTMLElement) Type() core.Type {
	return drivers.HTMLElementType
}

func (el *HTMLElement) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(el.String(), jettison.NoHTMLEscaping())
}

func (el *HTMLElement) String() string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(drivers.DefaultWaitTimeout)*time.Millisecond)
	defer cancel()

	res, err := el.GetInnerHTML(ctx)

	if err != nil {
		el.logError(errors.Wrap(err, "HTMLElement.String"))

		return ""
	}

	return res.String()
}

func (el *HTMLElement) Compare(other core.Value) int64 {
	switch other.Type() {
	case drivers.HTMLElementType:
		other := other.(drivers.HTMLElement)

		return int64(strings.Compare(el.String(), other.String()))
	default:
		return drivers.Compare(el.Type(), other.Type())
	}
}

func (el *HTMLElement) Unwrap() interface{} {
	return el
}

func (el *HTMLElement) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(el.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(el.id))

	return h.Sum64()
}

func (el *HTMLElement) Copy() core.Value {
	return values.None
}

func (el *HTMLElement) Iterate(_ context.Context) (core.Iterator, error) {
	return common.NewIterator(el)
}

func (el *HTMLElement) GetIn(ctx context.Context, path []core.Value) (core.Value, core.PathError) {
	return common.GetInElement(ctx, path, el)
}

func (el *HTMLElement) SetIn(ctx context.Context, path []core.Value, value core.Value) core.PathError {
	return common.SetInElement(ctx, path, el, value)
}

func (el *HTMLElement) GetValue(ctx context.Context) (core.Value, error) {
	return el.eval.EvalValue(ctx, templates.GetValue(el.id))
}

func (el *HTMLElement) SetValue(ctx context.Context, value core.Value) error {
	return el.eval.Eval(ctx, templates.SetValue(el.id, value))
}

func (el *HTMLElement) GetNodeType(ctx context.Context) (values.Int, error) {
	out, err := el.nodeType.Read(ctx)

	if err != nil {
		return values.ZeroInt, err
	}

	return values.ToInt(out), nil
}

func (el *HTMLElement) GetNodeName(ctx context.Context) (values.String, error) {
	out, err := el.nodeName.Read(ctx)

	if err != nil {
		return values.EmptyString, err
	}

	return values.ToString(out), nil
}

func (el *HTMLElement) Length() values.Int {
	value, err := el.eval.EvalValue(context.Background(), templates.GetChildrenCount(el.id))

	if err != nil {
		el.logError(err)

		return 0
	}

	return values.ToInt(value)
}

func (el *HTMLElement) GetStyles(ctx context.Context) (*values.Object, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetStyles(el.id))

	if err != nil {
		return values.NewObject(), err
	}

	return values.ToObject(ctx, out), nil
}

func (el *HTMLElement) GetStyle(ctx context.Context, name values.String) (core.Value, error) {
	return el.eval.EvalValue(ctx, templates.GetStyle(el.id, name))
}

func (el *HTMLElement) SetStyles(ctx context.Context, styles *values.Object) error {
	return el.eval.Eval(ctx, templates.SetStyles(el.id, styles))
}

func (el *HTMLElement) SetStyle(ctx context.Context, name, value values.String) error {
	return el.eval.Eval(ctx, templates.SetStyle(el.id, name, value))
}

func (el *HTMLElement) RemoveStyle(ctx context.Context, names ...values.String) error {
	return el.eval.Eval(ctx, templates.RemoveStyles(el.id, names))
}

func (el *HTMLElement) GetAttributes(ctx context.Context) (*values.Object, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetAttributes(el.id))

	if err != nil {
		return values.NewObject(), err
	}

	return values.ToObject(ctx, out), nil
}

func (el *HTMLElement) GetAttribute(ctx context.Context, name values.String) (core.Value, error) {
	return el.eval.EvalValue(ctx, templates.GetAttribute(el.id, name))
}

func (el *HTMLElement) SetAttributes(ctx context.Context, attrs *values.Object) error {
	return el.eval.Eval(ctx, templates.SetAttributes(el.id, attrs))
}

func (el *HTMLElement) SetAttribute(ctx context.Context, name, value values.String) error {
	return el.eval.Eval(ctx, templates.SetAttribute(el.id, name, value))
}

func (el *HTMLElement) RemoveAttribute(ctx context.Context, names ...values.String) error {
	return el.eval.Eval(ctx, templates.RemoveAttributes(el.id, names))
}

func (el *HTMLElement) GetChildNodes(ctx context.Context) (*values.Array, error) {
	return el.eval.EvalElements(ctx, templates.GetChildren(el.id))
}

func (el *HTMLElement) GetChildNode(ctx context.Context, idx values.Int) (core.Value, error) {
	return el.eval.EvalElement(ctx, templates.GetChildByIndex(el.id, idx))
}

func (el *HTMLElement) GetParentElement(ctx context.Context) (core.Value, error) {
	return el.eval.EvalElement(ctx, templates.GetParent(el.id))
}

func (el *HTMLElement) GetPreviousElementSibling(ctx context.Context) (core.Value, error) {
	return el.eval.EvalElement(ctx, templates.GetPreviousElementSibling(el.id))
}

func (el *HTMLElement) GetNextElementSibling(ctx context.Context) (core.Value, error) {
	return el.eval.EvalElement(ctx, templates.GetNextElementSibling(el.id))
}

func (el *HTMLElement) QuerySelector(ctx context.Context, selector drivers.QuerySelector) (core.Value, error) {
	return el.eval.EvalElement(ctx, templates.QuerySelector(el.id, selector))
}

func (el *HTMLElement) QuerySelectorAll(ctx context.Context, selector drivers.QuerySelector) (*values.Array, error) {
	return el.eval.EvalElements(ctx, templates.QuerySelectorAll(el.id, selector))
}

func (el *HTMLElement) XPath(ctx context.Context, expression values.String) (result core.Value, err error) {
	return el.eval.EvalValue(ctx, templates.XPath(el.id, expression))
}

func (el *HTMLElement) GetInnerText(ctx context.Context) (values.String, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerText(el.id))

	if err != nil {
		return values.EmptyString, err
	}

	return values.ToString(out), nil
}

func (el *HTMLElement) SetInnerText(ctx context.Context, innerText values.String) error {
	return el.eval.Eval(
		ctx,
		templates.SetInnerText(el.id, innerText),
	)
}

func (el *HTMLElement) GetInnerTextBySelector(ctx context.Context, selector drivers.QuerySelector) (values.String, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerTextBySelector(el.id, selector))

	if err != nil {
		return values.EmptyString, err
	}

	return values.ToString(out), nil
}

func (el *HTMLElement) SetInnerTextBySelector(ctx context.Context, selector drivers.QuerySelector, innerText values.String) error {
	return el.eval.Eval(
		ctx,
		templates.SetInnerTextBySelector(el.id, selector, innerText),
	)
}

func (el *HTMLElement) GetInnerTextBySelectorAll(ctx context.Context, selector drivers.QuerySelector) (*values.Array, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerTextBySelectorAll(el.id, selector))

	if err != nil {
		return values.EmptyArray(), err
	}

	return values.ToArray(ctx, out), nil
}

func (el *HTMLElement) GetInnerHTML(ctx context.Context) (values.String, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerHTML(el.id))

	if err != nil {
		return values.EmptyString, err
	}

	return values.ToString(out), nil
}

func (el *HTMLElement) SetInnerHTML(ctx context.Context, innerHTML values.String) error {
	return el.eval.Eval(ctx, templates.SetInnerHTML(el.id, innerHTML))
}

func (el *HTMLElement) GetInnerHTMLBySelector(ctx context.Context, selector drivers.QuerySelector) (values.String, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerHTMLBySelector(el.id, selector))

	if err != nil {
		return values.EmptyString, err
	}

	return values.ToString(out), nil
}

func (el *HTMLElement) SetInnerHTMLBySelector(ctx context.Context, selector drivers.QuerySelector, innerHTML values.String) error {
	return el.eval.Eval(ctx, templates.SetInnerHTMLBySelector(el.id, selector, innerHTML))
}

func (el *HTMLElement) GetInnerHTMLBySelectorAll(ctx context.Context, selector drivers.QuerySelector) (*values.Array, error) {
	out, err := el.eval.EvalValue(ctx, templates.GetInnerHTMLBySelectorAll(el.id, selector))

	if err != nil {
		return values.EmptyArray(), err
	}

	return values.ToArray(ctx, out), nil
}

func (el *HTMLElement) CountBySelector(ctx context.Context, selector drivers.QuerySelector) (values.Int, error) {
	out, err := el.eval.EvalValue(ctx, templates.CountBySelector(el.id, selector))

	if err != nil {
		return values.ZeroInt, err
	}

	return values.ToInt(out), nil
}

func (el *HTMLElement) ExistsBySelector(ctx context.Context, selector drivers.QuerySelector) (values.Boolean, error) {
	out, err := el.eval.EvalValue(ctx, templates.ExistsBySelector(el.id, selector))

	if err != nil {
		return values.False, err
	}

	return values.ToBoolean(out), nil
}

func (el *HTMLElement) logError(err error) *zerolog.Event {
	return el.logger.
		Error().
		Timestamp().
		Str("objectID", string(el.id)).
		Err(err)
}

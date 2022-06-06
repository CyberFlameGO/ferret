package actions

import (
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const selectorOptKey = "selector"

func getSelector(opts *values.Object) drivers.QuerySelector {
	if opts != nil && opts.Has(selectorOptKey) {
		return drivers.NewSelector(values.ToString(opts.MustGet(selectorOptKey)))
	}

	return drivers.QuerySelector{}
}

const countOptKey = "count"

func getCount(opts *values.Object, defaultValue values.Int) values.Int {
	if opts != nil && opts.Has(countOptKey) {
		return values.ToInt(opts.MustGet(countOptKey))
	}

	return defaultValue
}

const delayOptKey = "delay"

func getDelay(opts *values.Object, defaultValue values.Int) values.Int {
	if opts != nil && opts.Has(delayOptKey) {
		return values.ToInt(opts.MustGet(delayOptKey))
	}

	return defaultValue
}

const scrollTopKey = "top"
const scrollLeftKey = "left"
const scrollBehaviorKey = "behavior"
const scrollBlockKey = "block"

func getScrollOptions(opts *values.Object) drivers.ScrollOptions {
	scrollOpts := drivers.ScrollOptions{}

	if opts == nil {
		return scrollOpts
	}

	if opts.Has(scrollTopKey) {
		scrollOpts.Top = values.ToFloat(opts.MustGet(scrollTopKey))
	}

	if opts.Has(scrollLeftKey) {
		scrollOpts.Left = values.ToFloat(opts.MustGet(scrollLeftKey))
	}

	if opts.Has(scrollBehaviorKey) {
		scrollOpts.Behavior = drivers.NewScrollBehavior(opts.MustGet(scrollBehaviorKey).String())
	}

	if opts.Has(scrollBlockKey) {
		scrollOpts.Behavior = drivers.NewScrollBehavior(opts.MustGet(scrollBehaviorKey).String())
	}

	return scrollOpts
}

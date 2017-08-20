package gopi_test

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/default/layout"
	_ "github.com/djthorpe/gopi/sys/default/logger"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE LAYOUT MODULE

func TestLayout_000(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true

	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	if app == nil {
		t.Fatal("Expecting app object")
	}
	if app.Logger == nil {
		t.Fatal("Expecting app.Logger object")
	}
	if app.Layout == nil {
		t.Fatal("Expecting app.Layout object")
	}
	app.Logger.Info("layout=%v", app.Layout)
}

func TestLayout_001(t *testing.T) {
	// Check direction default
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true

	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()

	// Check layout direction defaults to LEFTRIGHT
	layout := app.Layout
	default_direction := gopi.LAYOUT_DIRECTION_LEFTRIGHT
	if layout.Direction() != default_direction {
		t.Errorf("Layout direction is %v, expected %v", layout.Direction(), default_direction)
	}
}

func TestLayout_002(t *testing.T) {
	// Create a root node with tag 1
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true

	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()

	// Create a view
	layout := app.Layout
	view1 := layout.NewRootView(1, "root")
	if view1 == nil {
		t.Error("NewRootView failed")
	}
	if view1.Tag() != 1 {
		t.Errorf("view1.Tag() expected 1, received %v", view1.Tag())
	}
	if view1.Class() != "root" {
		t.Errorf("view1.Tag() expected root, received %v", view1.Class())
	}

	// Attempt to create a view with tag TagNone
	view2 := layout.NewRootView(gopi.TagNone, "")
	if view2 != nil {
		t.Error("NewRootView succeeded but should have failed")
	}
}

func TestLayout_003(t *testing.T) {
	// Check class names
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true

	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()

	class_name_tests := map[string]bool{
		"":       false,
		"a":      true,
		"-":      false,
		"0":      false,
		"test":   true,
		"t0":     true,
		"t-":     true,
		"t-test": true,
		"t!test": false,
	}

	// Create root view with particular class names
	tag := uint(1)
	for k, v := range class_name_tests {
		view := app.Layout.NewRootView(tag, k)
		failed := (view == nil)
		if failed == v {
			t.Errorf("class %v => %v, expected %v", k, !failed, v)
		}
		if view != nil {
			if view.Tag() != tag {
				t.Errorf("view.Tag() expected %v, received %v", tag, view.Tag())
			}
			if view.Class() != k {
				t.Errorf("view.Class() expected %v, received %v", k, view.Class())
			}
		}
		tag = tag + 1
	}
}

func TestLayout_004(t *testing.T) {
	// Check layout starts as absolute with auto edges
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true

	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()

	// Create a view
	layout := app.Layout
	view := layout.NewRootView(1, "root")
	if view == nil {
		t.Error("NewRootView failed")
	}
	if view.Positioning() != gopi.VIEW_POSITIONING_ABSOLUTE {
		t.Error("Expected positioning on root element to be absolute")
	}
	app.Logger.Info("view=%v", view)
}

func TestLayout_005(t *testing.T) {
	m := map[gopi.ViewDirection]string{
		gopi.VIEW_DIRECTION_COLUMN:         "VIEW_DIRECTION_COLUMN",
		gopi.VIEW_DIRECTION_COLUMN_REVERSE: "VIEW_DIRECTION_COLUMN_REVERSE",
		gopi.VIEW_DIRECTION_ROW:            "VIEW_DIRECTION_ROW",
		gopi.VIEW_DIRECTION_ROW_REVERSE:    "VIEW_DIRECTION_ROW_REVERSE",
	}
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true
	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	// Create a view and test direction property
	if view := app.Layout.NewRootView(1, "root"); view == nil {
		t.Error("View could not be created")
	} else {
		// Check direction is 'row' by default
		if view.Direction() != gopi.VIEW_DIRECTION_ROW {
			t.Errorf("Expected default direction to be VIEW_DIRECTION_ROW but it returned %v", view.Direction())
		}
		// Now check setting each direction
		for k, v := range m {
			view.SetDirection(k)
			if k.String() != v {
				t.Errorf("Expected string to return %v but it returned %v", v, k.String())
			}
			if view.Direction() != k {
				t.Errorf("Expected Direction() to return %v but it returned %v", k, view.Direction())
			}
		}
	}
}

func TestLayout_006(t *testing.T) {
	m := map[gopi.ViewWrap]string{
		gopi.VIEW_WRAP_ON:      "VIEW_WRAP_ON",
		gopi.VIEW_WRAP_OFF:     "VIEW_WRAP_OFF",
		gopi.VIEW_WRAP_REVERSE: "VIEW_WRAP_REVERSE",
	}
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true
	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	// Create a view and test direction property
	if view := app.Layout.NewRootView(1, "root"); view == nil {
		t.Error("View could not be created")
	} else {
		// Check wrap is 'nowrap' by default
		if view.Wrap() != gopi.VIEW_WRAP_OFF {
			t.Errorf("Expected default wrap to be VIEW_WRAP_OFF but it returned %v", view.Wrap())
		}
		// Now check setting each direction
		for k, v := range m {
			view.SetWrap(k)
			if k.String() != v {
				t.Errorf("Expected string to return %v but it returned %v", v, k.String())
			}
			if view.Wrap() != k {
				t.Errorf("Expected Wrap() to return %v but it returned %v", k, view.Wrap())
			}
		}
	}
}

func TestLayout_007(t *testing.T) {
	m := map[gopi.ViewJustify]string{
		gopi.VIEW_JUSTIFY_FLEX_START:    "VIEW_JUSTIFY_FLEX_START",
		gopi.VIEW_JUSTIFY_FLEX_END:      "VIEW_JUSTIFY_FLEX_END",
		gopi.VIEW_JUSTIFY_CENTER:        "VIEW_JUSTIFY_CENTER",
		gopi.VIEW_JUSTIFY_SPACE_BETWEEN: "VIEW_JUSTIFY_SPACE_BETWEEN",
		gopi.VIEW_JUSTIFY_SPACE_AROUND:  "VIEW_JUSTIFY_SPACE_AROUND",
	}
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true
	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	// Create a view and test direction property
	if view := app.Layout.NewRootView(1, "root"); view == nil {
		t.Error("View could not be created")
	} else {
		// Check justify is 'flex-start' by default
		if view.JustifyContent() != gopi.VIEW_JUSTIFY_FLEX_START {
			t.Errorf("Expected default justify-content to be VIEW_JUSTIFY_FLEX_START but it returned %v", view.JustifyContent())
		}
		// Now check setting each direction
		for k, v := range m {
			view.SetJustifyContent(k)
			if k.String() != v {
				t.Errorf("Expected string to return %v but it returned %v", v, k.String())
			}
			if view.JustifyContent() != k {
				t.Errorf("Expected Justify() to return %v but it returned %v", k, view.JustifyContent())
			}
		}
	}
}

func TestLayout_008(t *testing.T) {
	m := map[gopi.ViewAlign]string{
		gopi.VIEW_ALIGN_FLEX_START:    "VIEW_ALIGN_FLEX_START",
		gopi.VIEW_ALIGN_CENTER:        "VIEW_ALIGN_CENTER",
		gopi.VIEW_ALIGN_FLEX_END:      "VIEW_ALIGN_FLEX_END",
		gopi.VIEW_ALIGN_STRETCH:       "VIEW_ALIGN_STRETCH",
		gopi.VIEW_ALIGN_SPACE_BETWEEN: "VIEW_ALIGN_SPACE_BETWEEN",
		gopi.VIEW_ALIGN_SPACE_AROUND:  "VIEW_ALIGN_SPACE_AROUND",
	}
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true
	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	// Create a view and test direction property
	if view := app.Layout.NewRootView(1, "root"); view == nil {
		t.Error("View could not be created")
	} else {
		// Check align-content is 'stretch' by default
		if view.AlignContent() != gopi.VIEW_ALIGN_STRETCH {
			t.Errorf("Expected default align-content to be VIEW_ALIGN_STRETCH but it returned %v", view.AlignContent())
		}
		// Now check setting
		for k, v := range m {
			view.SetAlignContent(k)
			if k.String() != v {
				t.Errorf("Expected string to return %v but it returned %v", v, k.String())
			}
			if view.AlignContent() != k {
				t.Errorf("Expected AlignContent() to return %v but it returned %v", k, view.AlignContent())
			}
		}
	}
}

func TestLayout_009(t *testing.T) {
	m := map[gopi.ViewAlign]string{
		gopi.VIEW_ALIGN_FLEX_START: "VIEW_ALIGN_FLEX_START",
		gopi.VIEW_ALIGN_CENTER:     "VIEW_ALIGN_CENTER",
		gopi.VIEW_ALIGN_FLEX_END:   "VIEW_ALIGN_FLEX_END",
		gopi.VIEW_ALIGN_STRETCH:    "VIEW_ALIGN_STRETCH",
		gopi.VIEW_ALIGN_BASELINE:   "VIEW_ALIGN_BASELINE",
	}
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true
	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	// Create a view and test direction property
	if view := app.Layout.NewRootView(1, "root"); view == nil {
		t.Error("View could not be created")
	} else {
		// Check align-items is 'stretch' by default
		if view.AlignItems() != gopi.VIEW_ALIGN_STRETCH {
			t.Errorf("Expected default align-items to be VIEW_ALIGN_STRETCH but it returned %v", view.AlignItems())
		}
		// Now check setting
		for k, v := range m {
			view.SetAlignItems(k)
			if k.String() != v {
				t.Errorf("Expected string to return %v but it returned %v", v, k.String())
			}
			if view.AlignItems() != k {
				t.Errorf("Expected AlignItems() to return %v but it returned %v", k, view.AlignItems())
			}
		}
	}
}

func TestLayout_010(t *testing.T) {
	m := map[gopi.ViewAlign]string{
		gopi.VIEW_ALIGN_AUTO:       "VIEW_ALIGN_AUTO",
		gopi.VIEW_ALIGN_FLEX_START: "VIEW_ALIGN_FLEX_START",
		gopi.VIEW_ALIGN_CENTER:     "VIEW_ALIGN_CENTER",
		gopi.VIEW_ALIGN_FLEX_END:   "VIEW_ALIGN_FLEX_END",
		gopi.VIEW_ALIGN_STRETCH:    "VIEW_ALIGN_STRETCH",
		gopi.VIEW_ALIGN_BASELINE:   "VIEW_ALIGN_BASELINE",
	}
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true
	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	// Create a view and test direction property
	if view := app.Layout.NewRootView(1, "root"); view == nil {
		t.Error("View could not be created")
	} else {
		// Check align-self is 'auto' by default
		if view.AlignSelf() != gopi.VIEW_ALIGN_AUTO {
			t.Errorf("Expected default align-self to be VIEW_ALIGN_AUTO but it returned %v", view.AlignSelf())
		}
		// Now check setting
		for k, v := range m {
			view.SetAlignSelf(k)
			if k.String() != v {
				t.Errorf("Expected string to return %v but it returned %v", v, k.String())
			}
			if view.AlignSelf() != k {
				t.Errorf("Expected AlignSelf() to return %v but it returned %v", k, view.AlignSelf())
			}
		}
	}
}

func TestLayout_011(t *testing.T) {
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true
	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	// Check default grow and shrink values
	view := app.Layout.NewRootView(1, "root")
	if view == nil {
		t.Fatal("View could not be created")
	}
	if view.Grow() != 0.0 {
		t.Errorf("Expected default grow value of 0, returned %v", view.Grow())
	}
	if view.Shrink() != 1.0 {
		t.Errorf("Expected default shrink value of 1, returned %v", view.Shrink())
	}

	v := float32(12345.0)
	view.SetGrow(v)
	view.SetShrink(v)
	if view.Grow() != v {
		t.Errorf("Expected default grow value of %v, returned %v", v, view.Grow())
	}
	if view.Shrink() != v {
		t.Errorf("Expected default shrink value of %v, returned %v", v, view.Shrink())
	}

	// Output view
	app.Logger.Info("view=%v", view)
}

func TestLayout_012(t *testing.T) {
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true
	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	// Check basis value
	view := app.Layout.NewRootView(1, "root")
	if view == nil {
		t.Fatal("View could not be created")
	}
	if view.BasisString() != "auto" {
		t.Errorf("Expected default basis value of auto, returned %v", view.BasisString())
	}
	view.SetBasisValue(0)
	if view.BasisString() != "0" {
		t.Errorf("Expected default basis value of 0, returned %v", view.BasisString())
	}
	view.SetBasisValue(1)
	if view.BasisString() != "1" {
		t.Errorf("Expected default basis value of 1, returned %v", view.BasisString())
	}
	view.SetBasisPercent(0)
	if view.BasisString() != "0%" {
		t.Errorf("Expected default basis value of 0%%, returned %v", view.BasisString())
	}
	view.SetBasisPercent(100)
	if view.BasisString() != "100%" {
		t.Errorf("Expected default basis value of 100%%, returned %v", view.BasisString())
	}
	view.SetBasisAuto()
	if view.BasisString() != "auto" {
		t.Errorf("Expected default basis value of auto, returned %v", view.BasisString())
	}

	// Output view
	app.Logger.Info("view=%v", view)
}

func TestLayout_020(t *testing.T) {
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_LAYOUT)
	config.Debug = true
	config.Verbose = true
	// Create an application with a layout module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			t.Error(err)
		}
	}()
	// Check basis value
	view := app.Layout.NewRootView(1, "root")
	if view == nil {
		t.Fatal("View could not be created")
	}

	encoder := xml.NewEncoder(os.Stdout)
	encoder.Indent("  ", "    ")
	if err := encoder.Encode(view); err != nil {
		t.Errorf("error: %v\n", err)
	}
	if err := encoder.Flush(); err != nil {
		t.Errorf("error: %v\n", err)
	}
}

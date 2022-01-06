package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func main() {
	lbl := labels.Set{"foo": "bar", "baz": "qux"}

	sel := labels.NewSelector()
	req, err := labels.NewRequirement("foo", selection.Equals, []string{"bar"})
	if err != nil {
		panic(err.Error())
	}
	sel.Add(*req)
	if sel.Matches(lbl) {
		fmt.Printf("Selector %v matched label set %v\n", sel, lbl)
	} else {
		panic("Selector should have been matched labels")
	}

	// Selector from string expression.
	sel, err = labels.Parse("foo==bar")
	if err != nil {
		panic(err.Error())
	}
	if sel.Matches(lbl) {
		fmt.Printf("Selector %v matched label set %v\n", sel, lbl)
	} else {
		panic("Selector should have been matched labels")
	}
}

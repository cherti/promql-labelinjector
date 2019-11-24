package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
)

var (
	// logger
	logInfo  = log.New(os.Stdout, "", 0)
	logDebug = log.New(os.Stdout, "DEBUG: ", 0)
	logError = log.New(os.Stdout, "ERROR: ", 0)

	// operation
	injectTarget = flag.String("l", "job", "label to inject or overwrite")
	injectValue  = flag.String("v", "prometheus", "value to write to target-label")
	expression   = flag.String("e", "", "expression to inject into")
	neqExpr      = flag.Bool("n", false, "inject with != instead of =")
)

// modifyQuery modifies a given Prometheus-query-expression to contain the required
// labelmatchers.
func modifyQuery(e string) string {
	expr, err := promql.ParseExpr(e)
	if err != nil {
		log.Fatal("ERROR, invalid query:", err)
	}
	// closure is actually unnecessary, but logic consistent with PromAuthProxy
	// as the code originated there
	promql.Inspect(expr, rewriteLabelsets)
	return expr.String()
}

// rewriteLabelsets returns the function that will be used to walk the
// Prometheus-query-expression-tree and rewrites the necessary selectors with
// to the specified username before the query is handed over to Prometheus.
func rewriteLabelsets(n promql.Node, path []promql.Node) error {

	// decide on matcher-type
	matcherType := labels.MatchEqual
	if *neqExpr {
		matcherType = labels.MatchNotEqual
	}

	switch n := n.(type) {
	case *promql.VectorSelector:
		// check if label is already present, replace in this case
		found := false
		for i, l := range n.LabelMatchers {
			if l.Name == *injectTarget {
				if l.Type == matcherType {
					l.Value = *injectValue
					found = true
				} else { // drop matcher if not of matcherType
					if len(n.LabelMatchers) == i {
						n.LabelMatchers = n.LabelMatchers[:i]
					} else {
						n.LabelMatchers = append(n.LabelMatchers[:i], n.LabelMatchers[i+1:]...)
					}
				}
			}
		}

		// if label is not present, inject it
		if !found {
			joblabel, err := labels.NewMatcher(matcherType, *injectTarget, *injectValue)
			if err != nil {
				//handle
			}
			n.LabelMatchers = append(n.LabelMatchers, joblabel)

		}
	case *promql.MatrixSelector:
		// check if label is already present, replace in this case
		found := false
		for i, l := range n.LabelMatchers {
			if l.Name == *injectTarget {
				if l.Type == matcherType {
					l.Value = *injectValue
					found = true
				} else { // drop matcher if not of matcherType
					if len(n.LabelMatchers) == i {
						n.LabelMatchers = n.LabelMatchers[:i]
					} else {
						n.LabelMatchers = append(n.LabelMatchers[:i], n.LabelMatchers[i+1:]...)
					}
				}
			}
		}
		// if label is not present, inject it
		if !found {
			joblabel, err := labels.NewMatcher(matcherType, *injectTarget, *injectValue)
			if err != nil {
				//handle
			}
			n.LabelMatchers = append(n.LabelMatchers, joblabel) // this doesn't compile with compiler error
		}
	}
	return nil
}

func main() {
	flag.Parse()
	fmt.Println(modifyQuery(*expression))
}

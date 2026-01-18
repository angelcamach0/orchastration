package agent

import (
	"fmt"
	"os"
	"strings"
)

// ReviewerAgent validates outputs for quality and correctness.
type ReviewerAgent struct{}

func (a *ReviewerAgent) Name() string {
	return "ReviewerAgent"
}

func (a *ReviewerAgent) Capabilities() []string {
	return []string{"Review outputs for quality and correctness"}
}

func (a *ReviewerAgent) Execute(ctx *OrchContext) error {
	if ctx == nil {
		return nil
	}

	value, ok := ctx.Get(ctxKeyBuildOutputs)
	if !ok {
		ctx.Set(ctxKeyReviewStatus, "skipped")
		ctx.Set(ctxKeyReviewReport, "no outputs to review")
		return nil
	}

	outputs, ok := value.([]string)
	if !ok || len(outputs) == 0 {
		ctx.Set(ctxKeyReviewStatus, "skipped")
		ctx.Set(ctxKeyReviewReport, "no outputs to review")
		return nil
	}

	missing := []string{}
	for _, output := range outputs {
		if _, err := os.Stat(output); err != nil {
			missing = append(missing, output)
		}
	}

	if len(missing) > 0 {
		report := fmt.Sprintf("missing outputs: %s", strings.Join(missing, ", "))
		ctx.Set(ctxKeyReviewStatus, "failed")
		ctx.Set(ctxKeyReviewReport, report)
		return fmt.Errorf(report)
	}

	ctx.Set(ctxKeyReviewStatus, "passed")
	ctx.Set(ctxKeyReviewReport, "all outputs present")
	return nil
}

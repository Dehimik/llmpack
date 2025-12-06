package pricing

import (
	"fmt"
	"strings"
)

// Prices on 2025
var models = map[string]float64{
	// OpenAI
	"gpt-4o":      2.50, // $2.50 / 1M input
	"gpt-4o-mini": 0.15, // $0.15 / 1M input
	"o1-preview":  15.00,
	"o1-mini":     3.00,

	// Anthropic
	"claude-3-5-sonnet": 3.00,
	"claude-3-opus":     15.00,
	"claude-3-haiku":    0.25,

	// Google
	"gemini-1.5-pro":   3.50, // before 128K context length
	"gemini-1.5-flash": 0.35,
}

// Estimate returns formatted price: "$0.0045"
func Estimate(tokens int, modelName string) string {
	pricePerMillion, exists := models[strings.ToLower(modelName)]
	if !exists {
		pricePerMillion = models["gpt-4o"]
		return fmt.Sprintf("unknown model, assuming gpt-4o (~$%.4f)", calculate(tokens, pricePerMillion))
	}

	cost := calculate(tokens, pricePerMillion)
	return fmt.Sprintf("$%.5f", cost)
}

func calculate(tokens int, pricePerMillion float64) float64 {
	return (float64(tokens) / 1_000_000.0) * pricePerMillion
}

// ListModels return list of models
func ListModels() string {
	keys := make([]string, 0, len(models))
	for k := range models {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

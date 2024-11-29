package fair

import "strings"

type ArkPluginResultType uint

const (
	ARKPluginCannotHandle ArkPluginResultType = iota
	ARKPluginData
	ARKPluginRedirect
)

type PluginResult struct {
	Type ArkPluginResultType
	Data []byte
	Mime string
}

type Plugin interface {
	Handle(fair *Fair, pid string, item *ItemData) (*PluginResult, error)
}

func insertHyphens(qualifier string) string {
	var result strings.Builder
	for i, char := range qualifier {
		if i > 0 && i%4 == 0 {
			result.WriteRune('-')
		}
		result.WriteRune(char)
	}
	return result.String()
}

func ARKBeautifier(ark string) string {
	naan, qualifier, components, variants, inflection, err := ArkParts(ark)
	if err != nil {
		return ark
	}
	result := "ark:/" + naan + "/" + insertHyphens(qualifier)
	if components != "" {
		result += "/" + components
	}
	if variants != "" {
		result += "." + variants
	}
	if inflection != "" {
		result += "?" + inflection
	}
	return result
}

package fair

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

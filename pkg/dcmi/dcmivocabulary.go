package dcmi

type DCMIType string

const (
	DCMITypeCollection          DCMIType = "Collection"
	DCMITypeDataset             DCMIType = "Dataset"
	DCMITypeEvent               DCMIType = "Event"
	DCMITypeImage               DCMIType = "Image"
	DCMITypeMovingImage         DCMIType = "MovingImage"
	DCMITypeStillImage          DCMIType = "StillImage"
	DCMITypeInteractiveResource DCMIType = "InteractiveResource"
	DCMITypeService             DCMIType = "Service"
	DCMITypeSoftware            DCMIType = "Software"
	DCMITypeSound               DCMIType = "Sound"
	DCMITypeText                DCMIType = "Text"
	DCMITypePhysicalObject      DCMIType = "PhysicalObject"
)

var DCMITypeReverse = map[string]DCMIType{
	"Collection":          DCMITypeCollection,
	"Dataset":             DCMITypeDataset,
	"Event":               DCMITypeEvent,
	"Image":               DCMITypeImage,
	"MovingImage":         DCMITypeMovingImage,
	"StillImage":          DCMITypeStillImage,
	"InteractiveResource": DCMITypeInteractiveResource,
	"Service":             DCMITypeService,
	"Software":            DCMITypeSoftware,
	"Sound":               DCMITypeSound,
	"Text":                DCMITypeText,
	"PhysicalObject":      DCMITypePhysicalObject,
}

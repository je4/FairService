package dcmi

type Type string

const (
	TypeCollection          Type = "Collection"
	TypeDataset             Type = "Dataset"
	TypeEvent               Type = "Event"
	TypeImage               Type = "Image"
	TypeMovingImage         Type = "MovingImage"
	TypeStillImage          Type = "StillImage"
	TypeInteractiveResource Type = "InteractiveResource"
	TypeService             Type = "Service"
	TypeSoftware            Type = "Software"
	TypeSound               Type = "Sound"
	TypeText                Type = "Text"
	TypePhysicalObject      Type = "PhysicalObject"
)

var DCMITypeReverse = map[string]Type{
	"Collection":          TypeCollection,
	"Dataset":             TypeDataset,
	"Event":               TypeEvent,
	"Image":               TypeImage,
	"MovingImage":         TypeMovingImage,
	"StillImage":          TypeStillImage,
	"InteractiveResource": TypeInteractiveResource,
	"Service":             TypeService,
	"Software":            TypeSoftware,
	"Sound":               TypeSound,
	"Text":                TypeText,
	"PhysicalObject":      TypePhysicalObject,
}

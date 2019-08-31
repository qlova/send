package send

//Type is a send type.
type Type interface {
	Name() string
	Detect(to string) bool
}

//Types are implicitly available types.
var Types []Type

//RegisterType registers a new type.
func RegisterType(t Type) {
	Types = append(Types, t)
}

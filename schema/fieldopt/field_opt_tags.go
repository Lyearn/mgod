package fieldopt

// FieldOptionTag is the BSON tag by which the properties of a field option is configured.
type FieldOptionTag string

const (
	FieldOptionTagRequired FieldOptionTag = "bson"
	FieldOptionTagXID      FieldOptionTag = "mgoID"
	FieldOptionTagDefault  FieldOptionTag = "mgoDefault"
)

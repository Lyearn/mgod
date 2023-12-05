package schema

type EntityModelSchemaOptions struct {
	xidRequired      bool
	bsonInlineParent bool
	parentBSONFields []string
}

func NewEntityModelSchemaOptions() *EntityModelSchemaOptions {
	return &EntityModelSchemaOptions{}
}

func (o *EntityModelSchemaOptions) SetXIDRequired(xidRequired bool) *EntityModelSchemaOptions {
	o.xidRequired = xidRequired
	return o
}

func (o *EntityModelSchemaOptions) SetBSONInlineParent(bsonInlineParent bool) *EntityModelSchemaOptions {
	o.bsonInlineParent = bsonInlineParent
	return o
}

func (o *EntityModelSchemaOptions) SetParentBSONFields(parentBSONFields []string) *EntityModelSchemaOptions {
	o.parentBSONFields = parentBSONFields
	return o
}

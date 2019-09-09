package pbast

type Enum struct {
	Name    string       `json:"name"`
	Options []*Option    `json:"options"`
	Comment Comment      `json:"comment"`
	Fields  []*EnumField `json:"fields"`
}

func NewEnum(name string) *Enum {
	return &Enum{
		Name: name,
	}
}

func (e *Enum) AddField(f *EnumField) *Enum {
	if f == nil {
		return e
	}
	e.Fields = append(e.Fields, f)
	return e
}

func (e *Enum) AddOption(o *Option) *Enum {
	if o == nil {
		return e
	}
	e.Options = append(e.Options, o)
	return e
}

func (e *Enum) identifiers() stringSet {
	if len(e.Fields) == 0 {
		return newStringSet()
	}

	set := newStringSet()
	for _, f := range e.Fields {
		set.add(f.Name)
	}
	return set
}

type EnumField struct {
	Name    string `json:"name"`
	Index   int `json:"index"`
	Comment []string `json:"comment"`
	Options []*EnumValueOption `json:"option"`
}

type EnumValueOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewEnumField(name string, index int) *EnumField {
	return &EnumField{
		Name:  name,
		Index: index,
	}
}

func NewEnumValueOption(name, value string) *EnumValueOption {
	return &EnumValueOption{
		Name:  name,
		Value: value,
	}
}

func (f *EnumField) AddOption(o *EnumValueOption) *EnumField {
	if o == nil {
		return f
	}
	f.Options = append(f.Options, o)
	return f
}

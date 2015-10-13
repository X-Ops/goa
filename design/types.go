// Package design defines types which describe the data types used by action controllers.
// These are the data structures of the request payloads and parameters as well as the response
// payloads.
// There are primitive types corresponding to the JSON primitive types (bool, string, integer and
// number), array types which represent a collection of another type and object types corresponding
// to JSON objects (i.e. a map indexed by strings where each value may be any of the data types).
// On top of these the package also defines "user types" and "media types". Both these types are
// named objects with additional properties (a description and for media types the media type
// identifier, links and views).
package design

import (
	"fmt"
	"reflect"
)

type (
	// A Kind defines the JSON type that a DataType represents.
	Kind uint

	// DataType is the common interface to all types.
	DataType interface {
		// Kind of data type, one of the Kind enum.
		Kind() Kind
		// Name returns the type name.
		Name() string
		// IsObject returns true if the underlying type is an object, a user type which
		// is an object or a media type whose type is an object.
		IsObject() bool
		// IsArray returns true if the underlying type is an array, a user type which
		// is an array or a media type whose type is an array.
		IsArray() bool
		// ToObject returns the underlying object if any (i.e. if IsObject returns true),
		// nil otherwise.
		ToObject() Object
		// ToArray returns the underlying array if any (i.e. if IsArray returns true),
		// nil otherwise.
		ToArray() *Array
		// IsCompatible checks whether val has a Go type that is
		// compatible with the data type.
		IsCompatible(val interface{}) bool
	}

	// DataStructure is the interface implemented by all data structure types.
	// That is attribute definitions, user types and media types.
	DataStructure interface {
		// Definition returns the data structure definition.
		Definition() *AttributeDefinition
	}

	// Primitive is the type for null, boolean, integer, number and string.
	Primitive Kind

	// Array is the type for a JSON array.
	Array struct {
		ElemType *AttributeDefinition
	}

	// Object is the type for a JSON object.
	Object map[string]*AttributeDefinition

	// UserTypeDefinition is the type for user defined types that are not media types
	// (e.g. payload types).
	UserTypeDefinition struct {
		// A user type is an attribute definition.
		*AttributeDefinition
		// Name of type
		TypeName string
		// DSL contains the DSL used to create this definition if any.
		DSL func()
	}

	// MediaTypeDefinition describes the rendering of a resource using property and link
	// definitions. A property corresponds to a single member of the media type,
	// it has a name and a type as well as optional validation rules. A link has a
	// name and a URL that points to a related resource.
	// Media types also define views which describe which members and links to render when
	// building the response body for the corresponding view.
	MediaTypeDefinition struct {
		// A media type is a type
		*UserTypeDefinition
		// Identifier is the RFC 6838 media type identifier.
		Identifier string
		// Links list the rendered links indexed by name.
		Links map[string]*LinkDefinition
		// Views list the supported views indexed by name.
		Views map[string]*ViewDefinition
		// Resource this media type is the canonical representation for if any
		Resource *ResourceDefinition
	}
)

const (
	// BooleanKind represents a JSON bool.
	BooleanKind = iota + 1
	// IntegerKind represents a JSON integer.
	IntegerKind
	// NumberKind represents a JSON number including integers.
	NumberKind
	// StringKind represents a JSON string.
	StringKind
	// ArrayKind represents a JSON array.
	ArrayKind
	// ObjectKind represents a JSON object.
	ObjectKind
	// UserTypeKind represents a user type.
	UserTypeKind
	// MediaTypeKind represents a media type.
	MediaTypeKind
)

const (
	// Boolean is the type for a JSON boolean.
	Boolean = Primitive(BooleanKind)

	// Integer is the type for a JSON number without a fraction or exponent part.
	Integer = Primitive(IntegerKind)

	// Number is the type for any JSON number, including integers.
	Number = Primitive(NumberKind)

	// String is the type for a JSON string.
	String = Primitive(StringKind)
)

// DataType implementation

// Kind implements DataKind.
func (p Primitive) Kind() Kind { return Kind(p) }

// Name returns the type name.
func (p Primitive) Name() string {
	switch p {
	case Boolean:
		return "boolean"
	case Integer:
		return "integer"
	case Number:
		return "number"
	case String:
		return "string"
	default:
		panic("unknown primitive type") // bug
	}
}

// IsObject returns false.
func (p Primitive) IsObject() bool { return false }

// IsArray returns false.
func (p Primitive) IsArray() bool { return false }

// ToObject returns nil.
func (p Primitive) ToObject() Object { return nil }

// ToArray returns nil.
func (p Primitive) ToArray() *Array { return nil }

// IsCompatible returns true if val is compatible with p.
func (p Primitive) IsCompatible(val interface{}) (ok bool) {
	switch p {
	case Boolean:
		_, ok = val.(bool)
	case Integer:
		_, ok = val.(int)
		if !ok {
			_, ok = val.(int8)
		}
		if !ok {
			_, ok = val.(int16)
		}
		if !ok {
			_, ok = val.(int32)
		}
		if !ok {
			_, ok = val.(int64)
		}
		if !ok {
			_, ok = val.(uint)
		}
		if !ok {
			_, ok = val.(uint8)
		}
		if !ok {
			_, ok = val.(uint16)
		}
		if !ok {
			_, ok = val.(uint32)
		}
		if !ok {
			_, ok = val.(uint64)
		}
	case Number:
		ok = Integer.IsCompatible(val)
		if !ok {
			_, ok = val.(float32)
		}
		if !ok {
			_, ok = val.(float64)
		}
	case String:
		_, ok = val.(string)
	default:
		panic("unknown primitive type") // bug
	}
	return
}

// Kind implements DataKind.
func (a *Array) Kind() Kind { return ArrayKind }

// Name returns the type name.
func (a *Array) Name() string {
	if a.ElemType != nil && a.ElemType.Type != nil {
		return fmt.Sprintf("array of %s", a.ElemType.Type.Name())
	}
	return fmt.Sprintf("untyped array")
}

// IsObject returns false.
func (a *Array) IsObject() bool { return false }

// IsArray returns true.
func (a *Array) IsArray() bool { return true }

// ToObject returns nil.
func (a *Array) ToObject() Object { return nil }

// ToArray returns a.
func (a *Array) ToArray() *Array { return a }

// IsCompatible returns true if val is compatible with p.
func (a *Array) IsCompatible(val interface{}) bool {
	k := reflect.TypeOf(val).Kind()
	return k == reflect.Array || k == reflect.Slice
}

// Dup creates a shallow copy of a.
func (a *Array) Dup() *Array {
	return &Array{ElemType: a.ElemType.Dup()}
}

// Kind implements DataKind.
func (o Object) Kind() Kind { return ObjectKind }

// Name returns the type name.
func (o Object) Name() string { return "object" }

// IsObject returns true.
func (o Object) IsObject() bool { return true }

// IsArray returns false.
func (o Object) IsArray() bool { return false }

// ToObject returns the underlying object.
func (o Object) ToObject() Object { return o }

// ToArray returns nil.
func (o Object) ToArray() *Array { return nil }

// Dup creates a shallow copy of o.
func (o Object) Dup() Object {
	res := make(Object, len(o))
	for n, att := range o {
		res[n] = att.Dup()
	}
	return res
}

// Merge copies other's attributes into o overridding any pre-existing attribute with the same name.
func (o Object) Merge(other Object) {
	for n, att := range other {
		o[n] = att.Dup()
	}
}

// IsCompatible returns true if val is compatible with p.
func (o Object) IsCompatible(val interface{}) bool {
	k := reflect.TypeOf(val).Kind()
	return k == reflect.Map || k == reflect.Struct
}

// NewUserTypeDefinition creates a user type definition but does not
// execute the DSL.
func NewUserTypeDefinition(name string, dsl func()) *UserTypeDefinition {
	return &UserTypeDefinition{
		TypeName:            name,
		AttributeDefinition: &AttributeDefinition{},
		DSL:                 dsl,
	}
}

// Kind implements DataKind.
func (u *UserTypeDefinition) Kind() Kind { return UserTypeKind }

// Name returns the type name.
func (u *UserTypeDefinition) Name() string { return u.TypeName }

// IsObject calls IsObject on the user type underlying data type.
func (u *UserTypeDefinition) IsObject() bool { return u.Type.IsObject() }

// IsArray calls IsArray on the user type underlying data type.
func (u *UserTypeDefinition) IsArray() bool { return u.Type.IsArray() }

// ToObject calls ToObject on the user type underlying data type.
func (u *UserTypeDefinition) ToObject() Object { return u.Type.ToObject() }

// ToArray calls ToArray on the user type underlying data type.
func (u *UserTypeDefinition) ToArray() *Array { return u.Type.ToArray() }

// IsCompatible returns true if val is compatible with p.
func (u *UserTypeDefinition) IsCompatible(val interface{}) bool {
	return u.Type.IsCompatible(val)
}

// NewMediaTypeDefinition creates a media type definition but does not
// execute the DSL.
func NewMediaTypeDefinition(name, identifier string, dsl func()) *MediaTypeDefinition {
	return &MediaTypeDefinition{
		UserTypeDefinition: &UserTypeDefinition{
			AttributeDefinition: &AttributeDefinition{},
			TypeName:            name,
			DSL:                 dsl,
		},
		Identifier: identifier,
	}
}

// Kind implements DataKind.
func (m *MediaTypeDefinition) Kind() Kind { return MediaTypeKind }

// DataStructure implementation

// Definition returns the underlying attribute definition.
// Note that this function is "inherited" by both UserTypeDefinition and
// MediaTypeDefinition.
func (a *AttributeDefinition) Definition() *AttributeDefinition {
	return a
}

package internal

type Datatype string

const (
	TypeString Datatype = "string"
	TypeInt    Datatype = "int"
	TypeFloat  Datatype = "float"
	TypeBool   Datatype = "bool"
	TypeFile   Datatype = "file"
)

type Store struct {
	Namespaces map[string]*Namespace `json:"namespaces"` // map for fast lookup
}

type Namespace struct {
	Name string              `json:"name"`
	Data map[string]Metadata `json:"data"` // key -> Metadata
}

type Metadata struct {
	Type  Datatype    `json:"type"`
	Value interface{} `json:"value"`
}

func InferType(v interface{}) string {
	switch v.(type) {
	case int:
		return "int"
	case float64:
		return "float"
	case float32:
		return "float"
	case string:
		return "string"
	case bool:
		return "bool"
	default:
		return "file"
	}
}

// db
// └── Namespaces (map[string]*Namespace)
//     ├── "users" → *Namespace
//     │       Name: "users"
//     │       Data:
//     │           "username" → Metadata{Type: "string", Value: "ammar"}
//     │           "age"      → Metadata{Type: "int", Value: 23}
//     └── "products" → *Namespace
//             Name: "products"
//             Data:
//                 "laptop" → Metadata{Type: "string", Value: "MacBook"}
//                 "price"  → Metadata{Type: "float", Value: 2999.99}

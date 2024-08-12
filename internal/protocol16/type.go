package protocol16

// alias for byte
type Protocol16Type byte

// data types in protocol16
const (
	Unknown           Protocol16Type = 0   // \0
	Null              Protocol16Type = 42  // *
	Dictionary        Protocol16Type = 68  // D
	StringArray       Protocol16Type = 97  // a
	Byte              Protocol16Type = 98  // b
	Double            Protocol16Type = 100 // d
	EventData         Protocol16Type = 101 // e
	Float             Protocol16Type = 102 // f
	Integer           Protocol16Type = 105 // i
	Hashtable         Protocol16Type = 104 // j
	Short             Protocol16Type = 107 // k
	Long              Protocol16Type = 108 // l
	IntegerArray      Protocol16Type = 110 // n
	Boolean           Protocol16Type = 111 // o
	OperationResponse Protocol16Type = 112 // p
	OperationRequest  Protocol16Type = 113 // q
	String            Protocol16Type = 115 // s
	ByteArray         Protocol16Type = 120 // x
	Array             Protocol16Type = 121 // y
	ObjectArray       Protocol16Type = 122 // z
)

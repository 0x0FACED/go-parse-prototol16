package photon

type OperationResponse struct {
	Code     byte
	RespCode Short
	DebugMsg string
	Params   map[byte]any
}

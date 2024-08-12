package photon

type OperationRequest struct {
	Code   byte
	Params map[byte]any
}

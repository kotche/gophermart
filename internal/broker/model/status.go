package model

type Status int

const (
	StatusUNKNOWN Status = iota
	StatusNEW
	StatusPROCESSING
	StatusINVALID
	StatusPROCESSED
)

func ConvertStatus(status string) Status {
	switch status {
	case "NEW", "REGISTERED":
		return StatusNEW
	case "PROCESSING":
		return StatusPROCESSING
	case "INVALID":
		return StatusINVALID
	case "PROCESSED":
		return StatusPROCESSED
	default:
		return StatusUNKNOWN
	}
}

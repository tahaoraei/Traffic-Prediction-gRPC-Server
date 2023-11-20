package param

type ETARequest struct {
	CurrentETA    int32
	Distance      int32
	Sx            int32
	Sy            int32
	Dx            int32
	Dy            int32
	Time          int32
	TrafficLength int32
}

type ETAResponse struct {
	ETA int32
}

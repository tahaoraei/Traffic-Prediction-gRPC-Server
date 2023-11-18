package param

type ETARequest struct {
	CurrentETA  int
	StraightETA int
	Distance    int
	Sx          int
	Sy          int
	Dx          int
	Dy          int
	Time        int
}

type ETAResponse struct {
	ETA int
}

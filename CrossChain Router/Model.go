package main
type Register struct {
	Id int
	Name string
	Identity string
	Address string
	Port string
	PublicKey string
	CalcType string
}

type Message struct {
	UUID string
	SCID string
	TCID string
	CalcType string
	TimeStamp int64
	Sign string
	Proof LpProof
	Type string
}

type LpProof struct {
	C []float64 `json:"C"`
	X []float64 `json:"X"`
	B []float64 `json:"B"`
	Y []float64 `json:"Y"`
	A []float64 `json:"A"`
}

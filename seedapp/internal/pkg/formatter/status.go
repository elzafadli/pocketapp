package formatter

type Status string

var (
	Success             Status = "00"
	CacheError          Status = "KAI01"
	DatabaseError       Status = "KAI02"
	InvalidRequest      Status = "KAI03"
	DataNotFound        Status = "KAI04"
	InternalServerError Status = "KAI05"
	DataConflict        Status = "KAI06"
	Unauthorized        Status = "KAI07"
)

func (s Status) String() string {
	return string(s)
}

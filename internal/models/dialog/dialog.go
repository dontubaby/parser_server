package models

type Task struct {
	Profession    string
	TaskID        uint64
	TaskGroupID   uint64
	GroupWeight   uint32
	Position      uint32
	ProfessionXP  uint32
	TmpXP         uint32
	PlayerXp      uint32
	ActivationXP  uint32
	DialogID      string
	ActivationAct string
	NextActXP     uint32
	NextChapterXP uint32
}

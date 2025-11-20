package contracts

type SubjectInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TeacherIDs  []uint `json:"teacher_ids"`
}

type SubjectDTO struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Teachers    []string `json:"teachers"`
}

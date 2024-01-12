package table

type Environment struct {
	GormBase

	Name       string    `json:"name"`
	ProjectIDs UUIDSlice `gorm:"type:jsonb" json:"project_ds"`
}

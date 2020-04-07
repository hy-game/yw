package app

type RoleInDb struct {
	Guid  uint32 `gorm:"primary_key"`
	Acc   string `gorm:"index:idx_acc"`
	Name  string
	Level uint32
	Data  []byte `gorm:"type:blob(104857600)"`
}
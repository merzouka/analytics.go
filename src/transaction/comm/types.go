package comm

type Product struct {
	ID    uint   `gorm:"primaryKey,autoIncrement" json:"id"`
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

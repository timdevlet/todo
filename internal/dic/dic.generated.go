package dic

func (repo *DicRepository) FetchCiti() ([]Citi, error) {
	items := []Citi{}
	result := repo.gorm.DB.Table("Cities").Scan(&items)
	if result.Error != nil {
		return []Citi{}, result.Error
	}
	return items, nil
}

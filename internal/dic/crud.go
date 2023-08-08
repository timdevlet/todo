package dic

func FetchAll[T any](repo *DicRepository, uuid string, s T) ([]T, error) {

	items := []T{}

	result := repo.gorm.DB.Table("users").Scan(&items)

	if result.Error != nil {
		return []T{}, result.Error
	}

	return items, nil
}

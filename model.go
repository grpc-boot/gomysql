package gomysql

type Model interface {
	NewModel() Model
	Assemble(br BytesRecord)
}

func BytesRecords2Models[T Model](brs []BytesRecord, model T) []T {
	models := make([]T, 0, len(brs))
	if len(brs) == 0 {
		return models
	}

	for _, br := range brs {
		m := model.NewModel().(T)
		m.Assemble(br)
		models = append(models, m)
	}

	return models
}

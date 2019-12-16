package models

type FactoryInfo struct {
	Message string
}

func (fi *FactoryInfo) Unmarshal(data []byte) error {
	fi.Message = string(data)
	return nil
}

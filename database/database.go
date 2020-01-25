package database

type MockDB struct {
	Values map[string]interface{}
}

func NewDB() *MockDB {
	return &MockDB{
		Values: map[string]interface{}{},
	}
}

func (db *MockDB) Get(val string) interface{} {
	if data, ok := db.Values[val]; ok {
		return data
	}
	return nil
}

func (db *MockDB) Insert(key string, val interface{}) {
	db.Values[key] = val
}

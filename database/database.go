package database

type MockDB struct {
	Values map[string]interface{}
}

var db = &MockDB{}

func NewDB() *MockDB {
	return &MockDB{
		Values: map[string]interface{}{},
	}
}

func SetDB(_db *MockDB) {
	db = _db
}

func Get(val string) interface{} {
	if data, ok := db.Values[val]; ok {
		return data
	}
	return nil
}

func Insert(key string, val interface{}) {
	db.Values[key] = val
}

package dbhelper

import "database/sql"

func Int32OrNil(val sql.NullInt32, nullValue int32) interface{} {
	if val.Int32 == nullValue {
		return nil
	}

	return val.Int32
}

func NullableToInt(val sql.NullInt32, nullValue int32) int {
	var value = int(nullValue)

	if val.Valid {
		value = int(val.Int32)
	}

	return value
}

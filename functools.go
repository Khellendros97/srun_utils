package utils

/*
* MapGet create a function that could get value from a map, if key not exists, return the default value
* Usage:
* m := map[string]interface{} {
* 	"key1": "value1",
* 	"key2": 2,
* }
* getOrElse := MapGet(m) // make a "getOrElse" function
* val1 := getOrElse("key1", "default value").(string) // val1 == "value1"
* val2 := getOrElse("key2", 0).(int) // val2 == 2
* val3 := getOrElse("key3(not exists)", "unknown") // no such key in the map, so val3 == "unknown"
 */
func MapGet(m map[string]interface{}) func(string, interface{}) interface{} {
	return func(key string, dft interface{}) interface{} {
		v, ok := m[key]
		if !ok {
			v = dft
		}
		return v
	}
}
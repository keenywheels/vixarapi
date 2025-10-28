package stopwords

// DictRegistry holds all available stopword dictionaries
var DictRegistry = []map[string]string{
	English,
	Russian,
}

// All is a merged dictionary of all stopwords
var All = map[string]string{}

// init merges all dictionaries into the All map
func init() {
	for _, dict := range DictRegistry {
		for k, v := range dict {
			All[k] = v
		}
	}
}

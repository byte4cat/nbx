package pbconv

func OverrideDefaultFields(fields ...string) {
	defaultFields = fields
	defaultFieldsMap = make(map[string]struct{})
	for _, field := range fields {
		defaultFieldsMap[field] = struct{}{}
	}
}

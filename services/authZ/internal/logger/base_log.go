package logger

type LogKey string

func BaseLogContext(pairs ...string) map[string]interface{} {
	logContext := make(map[string]interface{})

	for i := 0; i < len(pairs)-1; i += 2 {
		key := pairs[i]
		value := pairs[i+1]
		logContext[key] = value
	}
	return logContext
}

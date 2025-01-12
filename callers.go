package faket

import "runtime"

const (
	withSelf = 0
	skipSelf = 1
)

func getCallers(skip int) []uintptr {
	skip += 2 // skip runtime.Callers and self.
	depth := 32
	for {
		pc := make([]uintptr, depth)
		n := runtime.Callers(skip, pc)
		if n < len(pc) {
			return pc[:n]
		}
		depth *= 2
	}
}

func getCaller(skip int) uintptr {
	skip += 2 // skip runtime.Callers and this function
	var pc [1]uintptr
	n := runtime.Callers(skip, pc[:])
	if n == 0 {
		return 0
	}

	return pc[0]
}

func pcToFunction(pc uintptr) string {
	frames := runtime.CallersFrames([]uintptr{pc})
	f, _ := frames.Next()
	return f.Function
}

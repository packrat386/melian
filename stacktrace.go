package melian

import (
	"go/build"
	"path/filepath"
	"runtime"
	"strings"
)

type stacktrace struct {
	Frames        []frame `json:"frames,omitempty"`
	FramesOmitted []uint  `json:"frames_omitted,omitempty"`
}

func newStacktrace() *stacktrace {
	pcs := make([]uintptr, 100)
	n := runtime.Callers(1, pcs)

	if n == 0 {
		return nil
	}

	return &stacktrace{
		Frames: extractFrames(pcs[:n]),
	}
}

func extractFrames(pcs []uintptr) []frame {
	var frames []frame
	callersFrames := runtime.CallersFrames(pcs)

	for {
		callerFrame, more := callersFrames.Next()

		frames = append([]frame{
			newFrame(callerFrame),
		}, frames...)

		if !more {
			break
		}
	}

	return frames
}

// https://docs.sentry.io/development/sdk-dev/event-payloads/stacktrace/
type frame struct {
	Function    string                 `json:"function,omitempty"`
	Symbol      string                 `json:"symbol,omitempty"`
	Module      string                 `json:"module,omitempty"`
	Package     string                 `json:"package,omitempty"`
	Filename    string                 `json:"filename,omitempty"`
	AbsPath     string                 `json:"abs_path,omitempty"`
	Lineno      int                    `json:"lineno,omitempty"`
	Colno       int                    `json:"colno,omitempty"`
	PreContext  []string               `json:"pre_context,omitempty"`
	ContextLine string                 `json:"context_line,omitempty"`
	PostContext []string               `json:"post_context,omitempty"`
	InApp       bool                   `json:"in_app,omitempty"`
	Vars        map[string]interface{} `json:"vars,omitempty"`
}

func newFrame(f runtime.Frame) frame {
	abspath := f.File
	filename := f.File
	function := f.Function
	var pkg string

	if filename != "" {
		filename = filepath.Base(filename)
	} else {
		filename = "unknown"
	}

	if abspath == "" {
		abspath = "unknown"
	}

	if function != "" {
		pkg, function = splitQualifiedFunctionName(function)
	}

	frame := frame{
		AbsPath:  abspath,
		Filename: filename,
		Lineno:   f.Line,
		Module:   pkg,
		Function: function,
	}

	frame.InApp = isInAppFrame(frame)

	return frame
}

func isInAppFrame(frame frame) bool {
	if strings.HasPrefix(frame.AbsPath, build.Default.GOROOT) ||
		strings.Contains(frame.Module, "vendor") ||
		strings.Contains(frame.Module, "third_party") {
		return false
	}

	return true
}

func splitQualifiedFunctionName(name string) (pkg string, fun string) {
	pkg = packageName(name)
	fun = strings.TrimPrefix(name, pkg+".")
	return
}

func packageName(name string) string {
	// A prefix of "type." and "go." is a compiler-generated symbol that doesn't belong to any package.
	// See variable reservedimports in cmd/compile/internal/gc/subr.go
	if strings.HasPrefix(name, "go.") || strings.HasPrefix(name, "type.") {
		return ""
	}

	pathend := strings.LastIndex(name, "/")
	if pathend < 0 {
		pathend = 0
	}

	if i := strings.Index(name[pathend:], "."); i != -1 {
		return name[:pathend+i]
	}
	return ""
}

func baseName(name string) string {
	if i := strings.LastIndex(name, "."); i != -1 {
		return name[i+1:]
	}
	return name
}

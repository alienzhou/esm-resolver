package esm_resolver

import (
	"strings"

	"github.com/silenceper/log"
)

func PatternKeyCompare(a string, b string) int {
	var baseLenA int
	var baseLenB int
	aPatternIndex := strings.Index(a, "*")
	bPatternIndex := strings.Index(b, "*")
	if aPatternIndex == -1 {
		baseLenA = len(a)
	} else {
		baseLenA = aPatternIndex + 1
	}
	if bPatternIndex == -1 {
		baseLenB = len(b)
	} else {
		baseLenB = bPatternIndex + 1
	}
	if baseLenA > baseLenB {
		return -1
	}
	if baseLenB > baseLenA {
		return 1
	}
	if aPatternIndex == -1 {
		return 1
	}
	if bPatternIndex == -1 {
		return -1
	}
	if len(a) > len(b) {
		return -1
	}
	if len(b) > len(a) {
		return 1
	}
	return 0
}

func IsKeyExist(m map[string]interface{}, k string) bool {
	_, ok := m[k]
	return ok
}

func PackageTargetResolve(packageURL string, target string, subpath string, pattern string, internal bool) {

}

func PackageImportsExportsResolve(subpath string, exports map[string]interface{}, packageUrl string, isImport bool) {
	if IsKeyExist(exports, subpath) && strings.Index(subpath, "*") == -1 {
		target := exports[subpath]
		PackageTargetResolve(target)
	}
}

func PackageExportsResolve(exports map[string]interface{}, subpath string) bool {
	// 1. If exports is an Object with both a key starting with "." and a key not starting with ".", throw an Invalid Package Configuration error.

	if subpath == "." {
		mainExport := ""
		exports.(string)
	}

	var bestMatch string
	var bestMatchSubpath string
	for name := range exports {
		/**
		match ->
		exports: {
			"./lib/languages/*": {
				"require": "./lib/languages/*.js",
				"import": "./es/languages/*.js"
			},
		}
		or ->
		exports: {
			"./lib/languages/*.js": {
				"require": "./lib/languages/*.js",
				"import": "./es/languages/*.js"
			},
		}
		*/
		patternIndex := strings.Index(name, "*")
		if patternIndex != -1 && strings.HasPrefix(subpath, name[:patternIndex]) {
			if strings.HasSuffix(subpath, "/") {
				// deprecated
				log.Infof("exports pattern %s is deprecated (no '/' suffix)", subpath)
			}

			patternTrailer := name[patternIndex+1:]
			if len(subpath) >= len(name) &&
				strings.HasSuffix(subpath, patternTrailer) &&
				PatternKeyCompare(bestMatch, name) == 1 &&
				strings.LastIndex(name, "*") == patternIndex {
				bestMatch = name
				bestMatchSubpath = subpath[patternIndex : len(subpath)-len(patternTrailer)]
			}
		}
	}

	if bestMatchSubpath != "" {
		return true
	}

	return false
}

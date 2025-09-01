package stdlib

import "strings"

// ImportAliases maps convenient import names to actual Go package paths
// This provides Python/Node.js-like convenience for common packages
var ImportAliases = map[string]string{
	// Core packages
	"json":   "encoding/json",
	"base64": "encoding/base64",
	"hex":    "encoding/hex",
	"xml":    "encoding/xml",
	"csv":    "encoding/csv",

	// File system and I/O
	"fs":       "io/fs",
	"io":       "io",
	"ioutil":   "io/ioutil", // Deprecated but still used
	"path":     "path",
	"filepath": "path/filepath",
	"os":       "os",

	// Network and HTTP
	"http": "net/http",
	"url":  "net/url",
	"net":  "net",
	"mail": "net/mail",

	// Crypto and security
	"crypto": "crypto",
	"md5":    "crypto/md5",
	"sha1":   "crypto/sha1",
	"sha256": "crypto/sha256",
	"rand":   "crypto/rand",
	"tls":    "crypto/tls",

	// Time and date
	"time": "time",

	// String processing
	"strings": "strings",
	"strconv": "strconv",
	"regexp":  "regexp",

	// Math and numbers
	"math":      "math",
	"big":       "math/big",
	"rand_math": "math/rand", // Avoid conflict with crypto/rand

	// Compression
	"gzip": "compress/gzip",
	"zip":  "archive/zip",
	"tar":  "archive/tar",

	// Database
	"sql": "database/sql",

	// Logging and debugging
	"log": "log",
	"fmt": "fmt",

	// Concurrency
	"sync":    "sync",
	"context": "context",

	// Reflection and runtime
	"reflect": "reflect",
	"runtime": "runtime",

	// Testing
	"testing": "testing",

	// Sorting and containers
	"sort": "sort",
	"heap": "container/heap",
	"list": "container/list",
	"ring": "container/ring",

	// Image processing
	"image": "image",
	"color": "image/color",
	"png":   "image/png",
	"jpeg":  "image/jpeg",
	"gif":   "image/gif",

	// Template engines
	"template":      "text/template",
	"html_template": "html/template",

	// HTML and text processing
	"html":      "html",
	"scanner":   "text/scanner",
	"tabwriter": "text/tabwriter",

	// System and OS specific
	"exec":   "os/exec",
	"signal": "os/signal",
	"user":   "os/user",

	// Buffers and bytes
	"bytes": "bytes",
	"bufio": "bufio",

	// Error handling
	"errors": "errors",

	// Unicode and character handling
	"unicode": "unicode",
	"utf8":    "unicode/utf8",
	"utf16":   "unicode/utf16",
}

// GetRealPackagePath returns the actual Go package path for a given alias
// If no alias exists, returns the original name
func GetRealPackagePath(alias string) string {
	if realPath, exists := ImportAliases[alias]; exists {
		return realPath
	}
	return alias
}

// IsStandardLibrary checks if a package is part of Go's standard library
func IsStandardLibrary(packagePath string) bool {
	// Common standard library prefixes
	stdPrefixes := []string{
		"archive/", "bufio", "builtin", "bytes", "compress/", "container/",
		"context", "crypto/", "database/", "debug/", "embed", "encoding/",
		"errors", "expvar", "flag", "fmt", "go/", "hash/", "html/", "image/",
		"index/", "io/", "log/", "math/", "mime/", "net/", "os/", "path/",
		"plugin", "reflect", "regexp", "runtime/", "sort", "strconv", "strings",
		"sync/", "syscall", "testing/", "text/", "time", "unicode/", "unsafe",
	}

	// Direct matches
	directMatches := []string{
		"bufio", "builtin", "bytes", "context", "embed", "errors", "expvar",
		"flag", "fmt", "plugin", "reflect", "regexp", "sort", "strconv",
		"strings", "syscall", "time", "unsafe",
	}

	for _, prefix := range stdPrefixes {
		if strings.HasPrefix(packagePath, prefix) {
			return true
		}
	}

	for _, match := range directMatches {
		if packagePath == match {
			return true
		}
	}

	return false
}

// GetCommonAliases returns a list of commonly used aliases for documentation
func GetCommonAliases() map[string][]string {
	categories := map[string][]string{
		"Data Processing": {
			"json - JSON encoding and decoding",
			"xml - XML encoding and decoding",
			"csv - CSV file processing",
			"base64 - Base64 encoding",
			"hex - Hexadecimal encoding",
		},
		"File System": {
			"fs - File system interface",
			"io - Basic I/O primitives",
			"path - File path utilities",
			"filepath - File path manipulation",
			"os - Operating system interface",
		},
		"Network & HTTP": {
			"http - HTTP client and server",
			"url - URL parsing",
			"net - Network I/O",
			"mail - Mail parsing",
		},
		"Crypto & Security": {
			"crypto - Cryptographic functions",
			"md5 - MD5 hash algorithm",
			"sha1 - SHA1 hash algorithm",
			"sha256 - SHA256 hash algorithm",
			"rand - Cryptographically secure random numbers",
			"tls - TLS/SSL support",
		},
		"String & Text": {
			"strings - String manipulation",
			"strconv - String conversions",
			"regexp - Regular expressions",
		},
		"Math & Numbers": {
			"math - Mathematical functions",
			"big - Arbitrary precision arithmetic",
			"rand_math - Pseudo-random numbers",
		},
		"Time & Date": {
			"time - Time and date functions",
		},
		"Compression": {
			"gzip - Gzip compression",
			"zip - ZIP archive format",
			"tar - TAR archive format",
		},
		"Logging & Debug": {
			"log - Logging utilities",
			"fmt - Formatted I/O",
		},
		"Concurrency": {
			"sync - Synchronization primitives",
			"context - Request contexts",
		},
		"System": {
			"exec - External command execution",
			"signal - Signal handling",
			"user - User account information",
		},
	}

	return categories
}

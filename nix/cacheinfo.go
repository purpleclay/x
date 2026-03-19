package nix

import (
	"fmt"
	"strconv"
	"strings"
)

// CacheInfoName is the resource path Nix requests from a binary cache substituter.
const CacheInfoName = "nix-cache-info"

// CacheInfoMIMEType is the MIME type for the nix-cache-info resource.
const CacheInfoMIMEType = "text/x-nix-cache-info"

// CacheInfo holds the metadata served at the [CacheInfoName] endpoint of a Nix
// binary cache. When Nix connects to a substituter, this is the first resource
// it fetches to determine the store directory, whether the cache supports bulk
// queries, and what priority to assign when multiple binary caches are configured.
//
//	c := nix.CacheInfo{StoreDir: "/nix/store", WantMassQuery: true, Priority: 39}
//	text, _ := c.MarshalText()
//	// "StoreDir: /nix/store\nWantMassQuery: 1\nPriority: 39\n"
type CacheInfo struct {
	// StoreDir is the absolute path of the Nix store served by this cache.
	// If empty, it defaults to [DefaultStoreDirectory] during marshaling.
	StoreDir StoreDirectory

	// WantMassQuery reports whether the cache supports efficient bulk queries.
	// Serialised as 1 (true) or 0 (false).
	WantMassQuery bool

	// Priority determines substituter preference when multiple caches are
	// configured. Lower values take precedence. cache.nixos.org uses a
	// priority of 40; a custom cache with a lower value will be queried first.
	Priority int
}

// MarshalText implements [encoding.TextMarshaler]. It serialises CacheInfo in
// the Nix Key: Value line format. An empty StoreDir is marshaled as
// [DefaultStoreDirectory].
//
//	c := nix.CacheInfo{StoreDir: "/nix/store", WantMassQuery: true, Priority: 39}
//	text, _ := c.MarshalText()
//	// "StoreDir: /nix/store\nWantMassQuery: 1\nPriority: 39\n"
func (c CacheInfo) MarshalText() ([]byte, error) {
	storeDir := c.StoreDir
	if storeDir == "" {
		storeDir = DefaultStoreDirectory
	}
	wantMassQuery := "0"
	if c.WantMassQuery {
		wantMassQuery = "1"
	}
	return fmt.Appendf(nil, "StoreDir: %s\nWantMassQuery: %s\nPriority: %d\n",
		storeDir, wantMassQuery, c.Priority), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]. It parses CacheInfo from
// the Nix Key: Value line format. Unknown keys are silently ignored for forward
// compatibility.
//
//	var c nix.CacheInfo
//	_ = c.UnmarshalText([]byte("StoreDir: /nix/store\nWantMassQuery: 1\nPriority: 39\n"))
//	// c == nix.CacheInfo{StoreDir: "/nix/store", WantMassQuery: true, Priority: 39}
func (c *CacheInfo) UnmarshalText(text []byte) error {
	for line := range strings.SplitSeq(string(text), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ": ")
		if !ok {
			continue
		}
		switch key {
		case "StoreDir":
			c.StoreDir = StoreDirectory(value)
		case "WantMassQuery":
			switch value {
			case "1":
				c.WantMassQuery = true
			case "0":
				c.WantMassQuery = false
			default:
				return fmt.Errorf("nix: CacheInfo: invalid WantMassQuery value %q", value)
			}
		case "Priority":
			n, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("nix: CacheInfo: invalid Priority value %q", value)
			}
			c.Priority = n
		}
	}
	return nil
}

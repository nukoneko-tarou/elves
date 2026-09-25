package tree

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func BenchmarkParseJSON(b *testing.B) {
	// Read sample.json
	data, err := os.ReadFile("../../sample.json")
	if err != nil {
		b.Fatalf("failed to read sample.json: %v", err)
	}

	b.Run("json.NewDecoder", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r := bytes.NewReader(data)
			var tree Tree
			if err := json.NewDecoder(r).Decode(&tree); err != nil {
				b.Fatalf("decode failed: %v", err)
			}
		}
	})

	b.Run("json.Unmarshal", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var tree Tree
			if err := json.Unmarshal(data, &tree); err != nil {
				b.Fatalf("unmarshal failed: %v", err)
			}
		}
	})
}

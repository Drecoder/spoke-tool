package benchmark

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Basic Benchmarks
// ============================================================================

// BenchmarkStringConcatenation benchmarks different string concatenation methods
func BenchmarkStringConcatenation(b *testing.B) {
	parts := []string{"hello", "world", "this", "is", "a", "test"}

	b.Run("plus operator", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := ""
			for _, p := range parts {
				result += p
			}
			_ = result
		}
	})

	b.Run("strings.Join", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := strings.Join(parts, "")
			_ = result
		}
	})

	b.Run("bytes.Buffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			for _, p := range parts {
				buf.WriteString(p)
			}
			_ = buf.String()
		}
	})

	b.Run("strings.Builder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			for _, p := range parts {
				builder.WriteString(p)
			}
			_ = builder.String()
		}
	})
}

// ============================================================================
// Integer Operations Benchmarks
// ============================================================================

// BenchmarkIntegerOperations benchmarks various integer operations
func BenchmarkIntegerOperations(b *testing.B) {
	values := make([]int, 1000)
	for i := range values {
		values[i] = i
	}

	b.Run("addition", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sum := 0
			for _, v := range values {
				sum += v
			}
			_ = sum
		}
	})

	b.Run("multiplication", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			product := 1
			for _, v := range values {
				product *= v
			}
			_ = product
		}
	})

	b.Run("division", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			quotient := 1
			for _, v := range values {
				if v != 0 {
					quotient /= v
				}
			}
			_ = quotient
		}
	})

	b.Run("modulo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			remainder := 0
			for _, v := range values {
				if v != 0 {
					remainder = i % v
				}
			}
			_ = remainder
		}
	})
}

// ============================================================================
// Slice Operations Benchmarks
// ============================================================================

// BenchmarkSliceOperations benchmarks various slice operations
func BenchmarkSliceOperations(b *testing.B) {
	size := 10000
	data := make([]int, size)
	for i := range data {
		data[i] = i
	}

	b.Run("slice creation-make", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := make([]int, size)
			_ = slice
		}
	})

	b.Run("slice creation-literal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := []int{}
			_ = slice
		}
	})

	b.Run("slice copy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dst := make([]int, size)
			copy(dst, data)
		}
	})

	b.Run("slice append-single", func(b *testing.B) {
		slice := make([]int, 0, size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			slice = append(slice, i)
		}
	})

	b.Run("slice append-multiple", func(b *testing.B) {
		slice := make([]int, 0, size)
		values := []int{1, 2, 3, 4, 5}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			slice = append(slice, values...)
		}
	})

	b.Run("slice iteration-range", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sum := 0
			for _, v := range data {
				sum += v
			}
			_ = sum
		}
	})

	b.Run("slice iteration-index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sum := 0
			for j := 0; j < len(data); j++ {
				sum += data[j]
			}
			_ = sum
		}
	})
}

// ============================================================================
// Map Operations Benchmarks
// ============================================================================

// BenchmarkMapOperations benchmarks various map operations
func BenchmarkMapOperations(b *testing.B) {
	size := 10000
	data := make(map[int]int, size)
	for i := 0; i < size; i++ {
		data[i] = i * 2
	}

	b.Run("map creation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m := make(map[int]int)
			_ = m
		}
	})

	b.Run("map insertion", func(b *testing.B) {
		m := make(map[int]int)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			m[i] = i
		}
	})

	b.Run("map lookup", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = data[i%size]
		}
	})

	b.Run("map deletion", func(b *testing.B) {
		m := make(map[int]int)
		for i := 0; i < 1000; i++ {
			m[i] = i
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			delete(m, i%1000)
		}
	})

	b.Run("map iteration", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sum := 0
			for k, v := range data {
				sum += k + v
			}
			_ = sum
		}
	})
}

// ============================================================================
// Sorting Benchmarks
// ============================================================================

// BenchmarkSorting benchmarks different sorting algorithms
func BenchmarkSorting(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			data := make([]int, size)
			for i := range data {
				data[i] = size - i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tmp := make([]int, size)
				copy(tmp, data)
				sort.Ints(tmp)
			}
		})
	}

	b.Run("sort-ints", func(b *testing.B) {
		data := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
		for i := 0; i < b.N; i++ {
			tmp := make([]int, len(data))
			copy(tmp, data)
			sort.Ints(tmp)
		}
	})

	b.Run("sort-strings", func(b *testing.B) {
		data := []string{"z", "y", "x", "w", "v", "u", "t", "s", "r", "q"}
		for i := 0; i < b.N; i++ {
			tmp := make([]string, len(data))
			copy(tmp, data)
			sort.Strings(tmp)
		}
	})

	b.Run("sort-stable", func(b *testing.B) {
		type item struct {
			key   int
			value string
		}
		data := make([]item, 100)
		for i := range data {
			data[i] = item{key: i % 10, value: fmt.Sprintf("value-%d", i)}
		}
		for i := 0; i < b.N; i++ {
			tmp := make([]item, len(data))
			copy(tmp, data)
			sort.SliceStable(tmp, func(i, j int) bool {
				return tmp[i].key < tmp[j].key
			})
		}
	})
}

// ============================================================================
// String Conversion Benchmarks
// ============================================================================

// BenchmarkStringConversion benchmarks string conversion methods
func BenchmarkStringConversion(b *testing.B) {
	num := 123456789

	b.Run("strconv.Itoa", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = strconv.Itoa(num)
		}
	})

	b.Run("fmt.Sprint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprint(num)
		}
	})

	b.Run("strconv.FormatInt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = strconv.FormatInt(int64(num), 10)
		}
	})

	b.Run("string conversion with base", func(b *testing.B) {
		bases := []int{2, 8, 10, 16}
		for _, base := range bases {
			b.Run(fmt.Sprintf("base-%d", base), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_ = strconv.FormatInt(int64(num), base)
				}
			})
		}
	})
}

// ============================================================================
// JSON Benchmarks
// ============================================================================

type testStruct struct {
	ID       int                    `json:"id"`
	Name     string                 `json:"name"`
	Tags     []string               `json:"tags"`
	Active   bool                   `json:"active"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}

func generateTestData() testStruct {
	return testStruct{
		ID:     42,
		Name:   "test object",
		Tags:   []string{"go", "benchmark", "json", "performance"},
		Active: true,
		Score:  3.14159,
		Metadata: map[string]interface{}{
			"created": "2024-01-01",
			"version": 1,
			"source":  "benchmark",
		},
	}
}

// BenchmarkJSONMarshal benchmarks JSON marshaling
func BenchmarkJSONMarshal(b *testing.B) {
	data := generateTestData()

	b.Run("marshal-struct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("marshal-indent", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("marshal-map", func(b *testing.B) {
		m := map[string]interface{}{
			"id":     42,
			"name":   "test",
			"values": []int{1, 2, 3, 4, 5},
		}
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(m)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkJSONUnmarshal benchmarks JSON unmarshaling
func BenchmarkJSONUnmarshal(b *testing.B) {
	data := generateTestData()
	jsonData, err := json.Marshal(data)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("unmarshal-struct", func(b *testing.B) {
		var result testStruct
		for i := 0; i < b.N; i++ {
			err := json.Unmarshal(jsonData, &result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("unmarshal-map", func(b *testing.B) {
		var result map[string]interface{}
		for i := 0; i < b.N; i++ {
			err := json.Unmarshal(jsonData, &result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("unmarshal-interface", func(b *testing.B) {
		var result interface{}
		for i := 0; i < b.N; i++ {
			err := json.Unmarshal(jsonData, &result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// ============================================================================
// Cryptography Benchmarks
// ============================================================================

// BenchmarkHashFunctions benchmarks various hash functions
func BenchmarkHashFunctions(b *testing.B) {
	data := []byte("The quick brown fox jumps over the lazy dog")

	b.Run("md5", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = md5.Sum(data)
		}
	})

	b.Run("sha256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sha256.Sum256(data)
		}
	})

	b.Run("md5-variable-size", func(b *testing.B) {
		sizes := []int{10, 100, 1000, 10000, 100000}
		for _, size := range sizes {
			b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
				input := make([]byte, size)
				for i := range input {
					input[i] = byte(i)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_ = md5.Sum(input)
				}
			})
		}
	})
}

// ============================================================================
// Parallel Benchmarks
// ============================================================================

// BenchmarkParallel benchmarks parallel operations
func BenchmarkParallel(b *testing.B) {
	b.Run("sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := expensiveOperation(i)
			_ = result
		}
	})

	b.Run("parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				result := expensiveOperation(42)
				_ = result
			}
		})
	})

	b.Run("parallel-setup", func(b *testing.B) {
		data := make([]int, 1000)
		for i := range data {
			data[i] = i
		}
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sum := 0
				for _, v := range data {
					sum += v
				}
				_ = sum
			}
		})
	})
}

// ============================================================================
// Memory Allocation Benchmarks
// ============================================================================

// BenchmarkAllocation benchmarks memory allocation patterns
func BenchmarkAllocation(b *testing.B) {
	b.Run("allocate-once", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := make([]int, 1000)
			_ = slice
		}
	})

	b.Run("allocate-reuse", func(b *testing.B) {
		slice := make([]int, 1000)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := range slice {
				slice[j] = j
			}
			_ = slice[0]
		}
	})

	b.Run("allocate-growing", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := make([]int, 0, 10)
			for j := 0; j < 1000; j++ {
				slice = append(slice, j)
			}
		}
	})

	b.Run("allocate-prealloc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := make([]int, 0, 1000)
			for j := 0; j < 1000; j++ {
				slice = append(slice, j)
			}
		}
	})
}

// ============================================================================
// Function Call Benchmarks
// ============================================================================

func emptyFunction()                   {}
func functionWithArgs(a, b, c int) int { return a + b + c }

// BenchmarkFunctionCall benchmarks function call overhead
func BenchmarkFunctionCall(b *testing.B) {
	b.Run("empty-function", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			emptyFunction()
		}
	})

	b.Run("function-with-args", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = functionWithArgs(1, 2, 3)
		}
	})

	b.Run("closure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			func() {
				_ = i
			}()
		}
	})

	b.Run("defer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			func() {
				defer emptyFunction()
			}()
		}
	})

	b.Run("method-call", func(b *testing.B) {
		t := testStruct{}
		for i := 0; i < b.N; i++ {
			t.String()
		}
	})
}

func (t testStruct) String() string {
	return fmt.Sprintf("testStruct{ID:%d, Name:%s}", t.ID, t.Name)
}

// ============================================================================
// Helper Functions
// ============================================================================

func expensiveOperation(n int) int {
	result := 0
	for i := 0; i < 1000; i++ {
		result += n * i
	}
	return result
}

// ============================================================================
// Benchmark with Different Input Sizes
// ============================================================================

// BenchmarkVariableInput benchmarks functions with different input sizes
func BenchmarkVariableInput(b *testing.B) {
	sizes := []int{1, 10, 100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			input := make([]int, size)
			for i := range input {
				input[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				processVariableInput(input)
			}
		})
	}
}

func processVariableInput(input []int) int {
	sum := 0
	for _, v := range input {
		sum += v
	}
	return sum
}

// ============================================================================
// Benchmark Setups and Cleanups
// ============================================================================

// BenchmarkWithSetup demonstrates benchmark setup and cleanup
func BenchmarkWithSetup(b *testing.B) {
	// Setup
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Prepare test data
		tmp := make([]int, len(data))
		copy(tmp, data)
		b.StartTimer()

		// Actual benchmark
		sort.Ints(tmp)
	}
}

// ============================================================================
// Benchmark Reporting
// ============================================================================

// BenchmarkReporting demonstrates custom benchmark reporting
func BenchmarkReporting(b *testing.B) {
	b.Run("custom-reporting", func(b *testing.B) {
		var total time.Duration
		for i := 0; i < b.N; i++ {
			start := time.Now()
			expensiveOperation(42)
			total += time.Since(start)
		}
		b.ReportMetric(float64(total)/float64(b.N), "ns/op")
	})

	b.Run("memory-reporting", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = make([]int, 1000)
		}
	})

	b.Run("custom-metric", func(b *testing.B) {
		count := 0
		for i := 0; i < b.N; i++ {
			if expensiveOperation(i) > 1000 {
				count++
			}
		}
		b.ReportMetric(float64(count)/float64(b.N), "ratio>1000")
	})
}

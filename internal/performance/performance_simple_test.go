package performance

import (
	"context"
	"testing"
	"time"
)

func TestCacheManager_NewCacheManager(t *testing.T) {
	ctx := context.Background()
	cacheManager := NewCacheManager(ctx)

	if cacheManager == nil {
		t.Fatal("NewCacheManager returned nil")
	}
}

func TestMetricsCollectorV2_NewMetricsCollectorV2(t *testing.T) {
	ctx := context.Background()
	config := &MetricsConfig{
		CollectionInterval:  time.Second,
		RetentionDuration:   24 * time.Hour,
		MaxMetrics:          1000,
		BufferSize:          100,
		FlushInterval:       10 * time.Second, // Important: Add FlushInterval
		BatchSize:           50,
		SamplingRate:        1.0,
		AnomalyDetection:    false, // Disable for simple test
		TrendAnalysis:       false,
		CorrelationAnalysis: false,
		ExportEnabled:       false,
	}

	collector := NewMetricsCollectorV2(ctx, config)

	if collector == nil {
		t.Fatal("NewMetricsCollectorV2 returned nil")
	}
}

func TestQueryOptimizer_NewQueryOptimizer(t *testing.T) {
	config := CacheConfig{
		MaxSize:        1000,
		TTL:            time.Hour,
		EvictionPolicy: "lru",
	}

	optimizer := NewQueryOptimizer(config)

	if optimizer == nil {
		t.Fatal("NewQueryOptimizer returned nil")
	}
}

func TestResourceManager_NewResourceManager(t *testing.T) {
	ctx := context.Background()
	config := &ResourceManagerConfig{
		GlobalMaxResources:  10,
		GlobalIdleTimeout:   5 * time.Minute,
		HealthCheckInterval: time.Minute,
		MetricsInterval:     30 * time.Second,
		CleanupInterval:     time.Minute, // Add required field
	}

	resourceManager := NewResourceManager(ctx, config)

	if resourceManager == nil {
		t.Fatal("NewResourceManager returned nil")
	}
}

func TestPerformanceOptimizer_NewPerformanceOptimizer(t *testing.T) {
	optimizer := NewPerformanceOptimizer()

	if optimizer == nil {
		t.Fatal("NewPerformanceOptimizer returned nil")
	}
}

func TestPerformanceOptimizer_WithContext(t *testing.T) {
	ctx := context.Background()
	optimizer := NewPerformanceOptimizerWithContext(ctx)

	if optimizer == nil {
		t.Fatal("NewPerformanceOptimizerWithContext returned nil")
	}
}

// Basic functionality test
func TestCacheBasicFunctionality(t *testing.T) {
	ctx := context.Background()

	cache := NewCacheManager(ctx)
	if cache == nil {
		t.Fatal("NewCacheManager returned nil")
	}

	// Test basic operations
	key := "test-key"
	value := "test-value"

	_ = cache.Set(key, value)

	if result, found := cache.Get(key); !found {
		t.Error("Expected to find cached value")
	} else if result != value {
		t.Errorf("Expected value %s, got %s", value, result)
	}
}

// Benchmark tests for performance validation
func BenchmarkCacheSet(b *testing.B) {
	ctx := context.Background()
	cache := NewCacheManager(ctx)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "benchmark-key"
		value := "benchmark-value"
		_ = cache.Set(key, value)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	ctx := context.Background()
	cache := NewCacheManager(ctx)

	// Pre-populate cache
	key := "benchmark-key"
	value := "benchmark-value"
	_ = cache.Set(key, value)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cache.Get(key)
	}
}

func BenchmarkPerformanceOptimizer_Creation(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		optimizer := NewPerformanceOptimizer()
		_ = optimizer
	}
}

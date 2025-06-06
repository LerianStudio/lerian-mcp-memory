package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"lerian-mcp-memory/internal/storage"
	"lerian-mcp-memory/pkg/types"
)

// MockStore implements a simple in-memory storage for testing
type MockStore struct {
	chunks map[string]*types.ConversationChunk
}

func NewMockStore() *MockStore {
	return &MockStore{
		chunks: make(map[string]*types.ConversationChunk),
	}
}

func (m *MockStore) Store(_ context.Context, chunk *types.ConversationChunk) error {
	return m.storeChunk(chunk)
}

func (m *MockStore) storeChunk(chunk *types.ConversationChunk) error {
	m.chunks[chunk.ID] = chunk
	return nil
}

func (m *MockStore) GetByID(_ context.Context, id string) (*types.ConversationChunk, error) {
	chunk, exists := m.chunks[id]
	if !exists {
		return nil, errors.New("chunk not found")
	}
	return chunk, nil
}

func (m *MockStore) Update(_ context.Context, chunk *types.ConversationChunk) error {
	return m.updateChunk(chunk)
}

func (m *MockStore) updateChunk(chunk *types.ConversationChunk) error {
	m.chunks[chunk.ID] = chunk
	return nil
}

func (m *MockStore) Search(_ context.Context, query *types.MemoryQuery, _ []float64) (*types.SearchResults, error) {
	return m.searchWithQuery(query)
}

func (m *MockStore) searchWithQuery(query *types.MemoryQuery) (*types.SearchResults, error) {
	var chunks []types.ConversationChunk
	for _, chunk := range m.chunks {
		if query.Repository != nil && chunk.Metadata.Repository == *query.Repository {
			chunks = append(chunks, *chunk)
			if len(chunks) >= query.Limit {
				break
			}
		}
	}

	results := make([]types.SearchResult, 0, len(chunks))
	for i := range chunks {
		results = append(results, types.SearchResult{
			Chunk: chunks[i],
			Score: 0.8, // Fixed score for testing
		})
	}

	return &types.SearchResults{
		Results:   results,
		Total:     len(results),
		QueryTime: time.Millisecond * 10,
	}, nil
}

func (m *MockStore) SearchByTags(_ context.Context, _ []string, _ int) ([]types.ConversationChunk, error) {
	return nil, nil
}

func (m *MockStore) SearchByTimeRange(_ context.Context, _, _ time.Time, _ int) ([]types.ConversationChunk, error) {
	return nil, nil
}

func (m *MockStore) SearchSimilar(_ context.Context, _ []float32, _ int) ([]types.ConversationChunk, error) {
	return nil, nil
}

func (m *MockStore) ListRepositories(_ context.Context) ([]string, error) {
	return nil, nil
}

func (m *MockStore) Initialize(_ context.Context) error {
	return nil
}

func (m *MockStore) ListByRepository(_ context.Context, repository string, limit, _ int) ([]types.ConversationChunk, error) {
	var results []types.ConversationChunk
	for _, chunk := range m.chunks {
		if chunk.Metadata.Repository == repository {
			results = append(results, *chunk)
			if len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

func (m *MockStore) ListBySession(_ context.Context, sessionID string) ([]types.ConversationChunk, error) {
	var results []types.ConversationChunk
	for _, chunk := range m.chunks {
		if chunk.SessionID == sessionID {
			results = append(results, *chunk)
		}
	}
	return results, nil
}

func (m *MockStore) Delete(_ context.Context, id string) error {
	delete(m.chunks, id)
	return nil
}

func (m *MockStore) Cleanup(_ context.Context, _ int) (int, error) {
	return 0, nil
}

func (m *MockStore) GetStats(_ context.Context) (*storage.StoreStats, error) {
	return &storage.StoreStats{}, nil
}

func (m *MockStore) Backup(_ context.Context, _ string) error {
	return nil
}

func (m *MockStore) Close() error {
	return nil
}

func (m *MockStore) HealthCheck(_ context.Context) error {
	return nil
}

// New interface methods for updated VectorStore interface
func (m *MockStore) GetAllChunks(_ context.Context) ([]types.ConversationChunk, error) {
	chunks := make([]types.ConversationChunk, 0, len(m.chunks))
	for _, chunk := range m.chunks {
		chunks = append(chunks, *chunk)
	}
	return chunks, nil
}

func (m *MockStore) DeleteCollection(_ context.Context, _ string) error {
	m.chunks = make(map[string]*types.ConversationChunk)
	return nil
}

func (m *MockStore) ListCollections(_ context.Context) ([]string, error) {
	return []string{"claude_memory"}, nil
}

func (m *MockStore) FindSimilar(_ context.Context, _ string, _ *types.ChunkType, _ int) ([]types.ConversationChunk, error) {
	return nil, nil // Simplified mock
}

func (m *MockStore) StoreChunk(ctx context.Context, chunk *types.ConversationChunk) error {
	return m.storeChunk(chunk)
}

func (m *MockStore) BatchStore(ctx context.Context, chunks []*types.ConversationChunk) (*storage.BatchResult, error) {
	result := &storage.BatchResult{
		Success:      0,
		Failed:       0,
		Errors:       []string{},
		ProcessedIDs: []string{},
	}

	for i := range chunks {
		if err := m.storeChunk(chunks[i]); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Success++
		}
		result.ProcessedIDs = append(result.ProcessedIDs, chunks[i].ID)
	}

	return result, nil
}

func (m *MockStore) BatchDelete(ctx context.Context, ids []string) (*storage.BatchResult, error) {
	result := &storage.BatchResult{
		Success:      0,
		Failed:       0,
		Errors:       []string{},
		ProcessedIDs: ids,
	}

	for _, id := range ids {
		if err := m.Delete(ctx, id); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Success++
		}
	}

	return result, nil
}

// Relationship management methods (required by VectorStore interface)
func (m *MockStore) StoreRelationship(_ context.Context, _, _ string, _ types.RelationType, _ float64, _ types.ConfidenceSource) (*types.MemoryRelationship, error) {
	return nil, errors.New("not implemented in mock")
}

func (m *MockStore) GetRelationships(_ context.Context, query *types.RelationshipQuery) ([]types.RelationshipResult, error) {
	return m.getRelationshipsWithQuery(query)
}

func (m *MockStore) getRelationshipsWithQuery(_ *types.RelationshipQuery) ([]types.RelationshipResult, error) {
	return nil, errors.New("not implemented in mock")
}

func (m *MockStore) TraverseGraph(_ context.Context, _ string, _ int, _ []types.RelationType) (*types.GraphTraversalResult, error) {
	return nil, errors.New("not implemented in mock")
}

func (m *MockStore) UpdateRelationship(_ context.Context, _ string, _ float64, _ types.ConfidenceFactors) error {
	return errors.New("not implemented in mock")
}

func (m *MockStore) DeleteRelationship(_ context.Context, _ string) error {
	return errors.New("not implemented in mock")
}

func (m *MockStore) GetRelationshipByID(_ context.Context, _ string) (*types.MemoryRelationship, error) {
	return nil, errors.New("not implemented in mock")
}

func TestMemoryAnalytics_RecordAccess(t *testing.T) {
	store := NewMockStore()
	analytics := NewMemoryAnalytics(store)
	defer analytics.Stop()

	ctx := context.Background()
	chunkID := "test-chunk-1"

	// Create a test chunk
	chunk := &types.ConversationChunk{
		ID:        chunkID,
		Content:   "Test content",
		Timestamp: time.Now(),
		Metadata: types.ChunkMetadata{
			Repository: "test-repo",
		},
	}
	if err := store.Store(ctx, chunk); err != nil {
		t.Fatalf("Failed to store chunk: %v", err)
	}

	// Record multiple accesses
	for i := 0; i < 5; i++ {
		err := analytics.RecordAccess(ctx, chunkID)
		if err != nil {
			t.Fatalf("Failed to record access: %v", err)
		}
	}

	// Check cache
	analytics.mu.RLock()
	metrics, exists := analytics.accessCache[chunkID]
	analytics.mu.RUnlock()

	if !exists {
		t.Fatal("Expected metrics to be cached")
	}

	if metrics.AccessCount != 5 {
		t.Errorf("Expected access count 5, got %d", metrics.AccessCount)
	}
}

func TestMemoryAnalytics_RecordUsage(t *testing.T) {
	store := NewMockStore()
	analytics := NewMemoryAnalytics(store)
	defer analytics.Stop()

	ctx := context.Background()
	chunkID := "test-chunk-2"

	// Record mixed usage
	_ = analytics.RecordUsage(ctx, chunkID, true)
	_ = analytics.RecordUsage(ctx, chunkID, true)
	_ = analytics.RecordUsage(ctx, chunkID, false)
	_ = analytics.RecordUsage(ctx, chunkID, true)

	// Check cache
	analytics.mu.RLock()
	metrics, exists := analytics.accessCache[chunkID]
	analytics.mu.RUnlock()

	if !exists {
		t.Fatal("Expected metrics to be cached")
	}

	if metrics.TotalUses != 4 {
		t.Errorf("Expected total uses 4, got %d", metrics.TotalUses)
	}

	if metrics.SuccessfulUses != 3 {
		t.Errorf("Expected successful uses 3, got %d", metrics.SuccessfulUses)
	}
}

func TestMemoryAnalytics_CalculateEffectivenessScore(t *testing.T) {
	analytics := NewMemoryAnalytics(nil)
	defer analytics.Stop()

	tests := []struct {
		name     string
		chunk    types.ConversationChunk
		minScore float64
		maxScore float64
	}{
		{
			name: "Solution with high success rate",
			chunk: types.ConversationChunk{
				Type:      types.ChunkTypeSolution,
				Timestamp: time.Now().Add(-24 * time.Hour),
				Metadata: types.ChunkMetadata{
					ExtendedMetadata: map[string]interface{}{
						types.EMKeySuccessRate:  0.9,
						types.EMKeyAccessCount:  10,
						types.EMKeyLastAccessed: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
					},
				},
			},
			minScore: 0.7,
			maxScore: 1.0,
		},
		{
			name: "Old problem with no usage",
			chunk: types.ConversationChunk{
				Type:      types.ChunkTypeProblem,
				Timestamp: time.Now().Add(-60 * 24 * time.Hour),
				Metadata: types.ChunkMetadata{
					ExtendedMetadata: map[string]interface{}{},
				},
			},
			minScore: 0.0,
			maxScore: 0.4,
		},
		{
			name: "Recent decision",
			chunk: types.ConversationChunk{
				Type:      types.ChunkTypeArchitectureDecision,
				Timestamp: time.Now().Add(-1 * time.Hour),
				Metadata: types.ChunkMetadata{
					ExtendedMetadata: map[string]interface{}{
						types.EMKeyAccessCount: 5,
					},
				},
			},
			minScore: 0.4,
			maxScore: 0.7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := analytics.CalculateEffectivenessScore(&tt.chunk)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Score %f is outside expected range [%f, %f]", score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestMemoryAnalytics_UpdateChunkAnalytics(t *testing.T) {
	store := NewMockStore()
	analytics := NewMemoryAnalytics(store)
	defer analytics.Stop()

	ctx := context.Background()
	chunkID := "test-chunk-3"

	// Create a test chunk
	chunk := &types.ConversationChunk{
		ID:        chunkID,
		Content:   "Test content",
		Type:      types.ChunkTypeSolution,
		Timestamp: time.Now(),
		Metadata: types.ChunkMetadata{
			Repository: "test-repo",
		},
	}
	if err := store.Store(ctx, chunk); err != nil {
		t.Fatalf("Failed to store chunk: %v", err)
	}

	// Record some usage
	_ = analytics.RecordAccess(ctx, chunkID)
	_ = analytics.RecordAccess(ctx, chunkID)
	_ = analytics.RecordUsage(ctx, chunkID, true)
	_ = analytics.RecordUsage(ctx, chunkID, true)
	_ = analytics.RecordUsage(ctx, chunkID, false)

	// Update analytics
	err := analytics.UpdateChunkAnalytics(ctx, chunkID)
	if err != nil {
		t.Fatalf("Failed to update analytics: %v", err)
	}

	// Verify updates
	updatedChunk, _ := store.GetByID(ctx, chunkID)

	accessCount, ok := updatedChunk.Metadata.ExtendedMetadata[types.EMKeyAccessCount].(int)
	if !ok || accessCount != 2 {
		t.Errorf("Expected access count 2, got %v", accessCount)
	}

	successRate, ok := updatedChunk.Metadata.ExtendedMetadata[types.EMKeySuccessRate].(float64)
	if !ok || successRate < 0.6 || successRate > 0.7 {
		t.Errorf("Expected success rate ~0.67, got %v", successRate)
	}

	effectivenessScore, ok := updatedChunk.Metadata.ExtendedMetadata[types.EMKeyEffectivenessScore].(float64)
	if !ok || effectivenessScore < 0.5 {
		t.Errorf("Expected effectiveness score > 0.5, got %v", effectivenessScore)
	}
}

func TestMemoryAnalytics_MarkObsolete(t *testing.T) {
	store := NewMockStore()
	analytics := NewMemoryAnalytics(store)
	defer analytics.Stop()

	ctx := context.Background()
	chunkID := "test-chunk-4"

	// Create a test chunk
	chunk := &types.ConversationChunk{
		ID:        chunkID,
		Content:   "Outdated content",
		Timestamp: time.Now(),
		Metadata: types.ChunkMetadata{
			Repository: "test-repo",
		},
	}
	if err := store.Store(ctx, chunk); err != nil {
		t.Fatalf("Failed to store chunk: %v", err)
	}

	// Mark as obsolete
	err := analytics.MarkObsolete(ctx, chunkID, "Superseded by newer solution")
	if err != nil {
		t.Fatalf("Failed to mark obsolete: %v", err)
	}

	// Verify
	updatedChunk, _ := store.GetByID(ctx, chunkID)

	isObsolete, ok := updatedChunk.Metadata.ExtendedMetadata[types.EMKeyIsObsolete].(bool)
	if !ok || !isObsolete {
		t.Error("Expected chunk to be marked obsolete")
	}

	archivedAt, ok := updatedChunk.Metadata.ExtendedMetadata[types.EMKeyArchivedAt].(string)
	if !ok || archivedAt == "" {
		t.Error("Expected archived_at timestamp")
	}

	reason, ok := updatedChunk.Metadata.ExtendedMetadata["obsolete_reason"].(string)
	if !ok || reason != "Superseded by newer solution" {
		t.Errorf("Expected obsolete reason, got %v", reason)
	}
}

func TestMemoryAnalytics_GetTopMemories(t *testing.T) {
	store := NewMockStore()
	analytics := NewMemoryAnalytics(store)
	defer analytics.Stop()

	ctx := context.Background()
	repo := "test-repo"

	// Create chunks with different effectiveness
	chunks := []struct {
		id          string
		chunkType   types.ChunkType
		successRate float64
		accessCount int
		isObsolete  bool
	}{
		{"chunk1", types.ChunkTypeSolution, 0.9, 20, false},
		{"chunk2", types.ChunkTypeArchitectureDecision, 0.8, 15, false},
		{"chunk3", types.ChunkTypeProblem, 0.5, 5, false},
		{"chunk4", types.ChunkTypeSolution, 0.95, 25, true}, // Obsolete, should be filtered
		{"chunk5", types.ChunkTypeDiscussion, 0.7, 10, false},
	}

	for _, c := range chunks {
		chunk := &types.ConversationChunk{
			ID:        c.id,
			Type:      c.chunkType,
			Timestamp: time.Now().Add(-24 * time.Hour),
			Metadata: types.ChunkMetadata{
				Repository: repo,
				ExtendedMetadata: map[string]interface{}{
					types.EMKeySuccessRate:  c.successRate,
					types.EMKeyAccessCount:  c.accessCount,
					types.EMKeyIsObsolete:   c.isObsolete,
					types.EMKeyLastAccessed: time.Now().Format(time.RFC3339),
				},
			},
		}
		if err := store.Store(ctx, chunk); err != nil {
			t.Fatalf("Failed to store chunk: %v", err)
		}
	}

	// Get top 3 memories
	topMemories, err := analytics.GetTopMemories(ctx, repo, 3)
	if err != nil {
		t.Fatalf("Failed to get top memories: %v", err)
	}

	if len(topMemories) != 3 {
		t.Errorf("Expected 3 memories, got %d", len(topMemories))
	}

	// Verify chunk1 is in top (high success rate and access count)
	foundChunk1 := false
	for i := range topMemories {
		if topMemories[i].ID == "chunk1" {
			foundChunk1 = true
			break
		}
	}
	if !foundChunk1 {
		t.Error("Expected chunk1 to be in top memories")
	}

	// Verify chunk4 (obsolete) is not included
	for i := range topMemories {
		if topMemories[i].ID == "chunk4" {
			t.Error("Obsolete chunk should not be in top memories")
		}
	}
}

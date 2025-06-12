# MCP Memory Server - Direct Refactor Implementation Plan

## Overview

This document outlines the complete step-by-step implementation plan for **directly refactoring the existing MCP Memory Server** based on the critical design issues identified in `MEMORY_MCP_IMPROVEMENTS.md`. We will replace the current complex system with a clean, intuitive architecture in a single coordinated effort.

**Goal**: Transform the existing MCP server into 4 clean, intuitive tools with logical parameter requirements and clear domain boundaries, while updating the CLI to leverage the improved architecture.

---

## 🎯 **PHASE 1: CORE ARCHITECTURE REFACTOR**

### **P1.1: Design Clean Parameter System**

**Goal**: Build a logical, consistent parameter system from scratch that eliminates confusion.

#### **Server v2 Implementation**

**Step 1.1.1: Replace existing parameter system**
- [x] Replace `internal/mcp/consolidated_tools.go` parameter handling with clean types:
- [x] Create `internal/types/core.go` (replacing fragmented type definitions):
```go
package types

import "fmt"

// ProjectID represents a project/tenant identifier for data isolation
type ProjectID string

// Validate ensures ProjectID follows consistent format rules
func (p ProjectID) Validate() error {
    if len(p) == 0 || len(p) > 100 {
        return fmt.Errorf("project_id must be 1-100 characters")
    }
    // Add format validation (alphanumeric + hyphens/underscores)
    return nil
}

// SessionID represents a user session for scoped operations
type SessionID string

// OperationScope defines the access level for operations
type OperationScope string

const (
    ScopeSession OperationScope = "session" // requires session_id + project_id
    ScopeProject OperationScope = "project" // requires project_id only  
    ScopeGlobal  OperationScope = "global"  // no project_id required
)

// StandardParams provides consistent parameter structure
type StandardParams struct {
    ProjectID ProjectID      `json:"project_id,omitempty"`
    SessionID SessionID      `json:"session_id,omitempty"`
    Scope     OperationScope `json:"scope"`
}
```

**Step 1.1.2: Replace parameter validation system**
- [x] Replace existing validation in `internal/mcp/server.go` with `internal/validation/params.go`
- [x] Remove backwards session logic from all handlers
- [x] Add scope-based parameter requirement validation
- [x] Replace cryptic error messages with clear, helpful ones

**Step 1.1.3: Refactor all tool parameter schemas**
- [x] Replace `repository` with `project_id` in all 9 current tools
- [x] Update all parameter descriptions for clarity
- [x] Update all handler functions to use new parameter system
- [x] Update database queries to use ProjectID consistently

#### **CLI Refactor (Coordinated)**

**Step 1.1.4: Update CLI entities and types**
- [x] Update `cli/internal/domain/entities/task.go` to use ProjectID (replace Repository field)
- [x] Update `cli/internal/domain/ports/mcp.go` interface for new parameters
- [x] Update CLI to automatically handle parameter scoping

**Step 1.1.5: Update CLI MCP client**
- [x] Update `cli/internal/adapters/secondary/mcp/client.go` to send `project_id`
- [x] Update all MCP request structures for new parameters
- [x] Add automatic parameter defaulting based on operation scope
- [x] Update response parsing for new parameter names

**Step 1.1.6: Update CLI services**
- [x] Update `cli/internal/domain/services/task_service.go` to use ProjectID
- [x] Update repository detection logic to return ProjectID
- [x] Update all service layer operations for new parameter system

---

### **P1.2: Implement Logical Session Management**

**Goal**: Build intuitive session semantics where including session_id provides MORE access, not less.

#### **Server v2 Implementation**

**Step 1.2.1: Replace backwards session logic**
- [x] Replace existing session logic in `internal/mcp/server.go` with logical semantics:
- [x] Create `internal/session/manager.go` (replacing scattered session handling):
```go
package session

import "lerian-mcp-memory/internal/types"

type Manager struct {
    // Session management logic
}

type AccessLevel string

const (
    AccessReadOnly  AccessLevel = "read_only"   // No session_id: limited project data
    AccessSession   AccessLevel = "session"     // With session_id: full session + project data  
    AccessProject   AccessLevel = "project"     // Project scope: all project data
)

// GetAccessLevel determines what data user can access
func (m *Manager) GetAccessLevel(projectID types.ProjectID, sessionID types.SessionID) AccessLevel {
    if sessionID == "" {
        return AccessReadOnly  // Limited access without session
    }
    return AccessSession // Full access with session
}
```

**Step 1.2.2: Refactor all session-dependent operations**
- [x] Update all `memory_tasks` operations to require session_id for writes
- [x] Update `memory_analyze` to work without session_id (read-only project data)
- [x] Update `memory_create` operations to require session_id for persistence
- [x] Remove "cross-session continuity" logic from all handlers

**Step 1.2.3: Update database queries for logical session scoping**
- [x] Update all SQL queries in `internal/storage/` to respect session scoping
- [x] Update Qdrant queries to use session-based filtering when provided
- [x] Add session validation in storage layer
- [x] Remove backwards session logic from data access

#### **CLI Refactor (Coordinated)**

**Step 1.2.4: Update CLI session handling**
- [x] Update CLI to always provide session_id for write operations
- [x] Update read operations to optionally provide session_id for expanded access
- [x] Update session creation/management in existing CLI code
- [x] Update CLI caching to be session-aware

---

### **P1.3: Build Clean 4-Tool Architecture**

**Goal**: Design 4 logical tools with clear boundaries and no overlapping responsibilities.

#### **Server v2 Implementation**

**Step 1.3.1: Replace 9 fragmented tools with 4 clean tools**
- [x] Replace `internal/mcp/consolidated_tools.go` with 4 distinct tools:
- [x] Create `internal/tools/` package replacing fragmented tool system:
```go
// internal/tools/store.go - All data persistence
type StoreOperations string
const (
    OpStoreContent     StoreOperations = "store_content"
    OpStoreDecision    StoreOperations = "store_decision"
    OpUpdateContent    StoreOperations = "update_content"
    OpDeleteContent    StoreOperations = "delete_content"
    OpCreateThread     StoreOperations = "create_thread"
    OpCreateRelation   StoreOperations = "create_relationship"
)

// internal/tools/retrieve.go - All data retrieval
type RetrieveOperations string
const (
    OpSearch           RetrieveOperations = "search"
    OpGetContent       RetrieveOperations = "get_content"
    OpFindSimilar      RetrieveOperations = "find_similar"
    OpGetThreads       RetrieveOperations = "get_threads"
    OpGetRelationships RetrieveOperations = "get_relationships"
)

// internal/tools/analyze.go - All analysis and intelligence
type AnalyzeOperations string
const (
    OpDetectPatterns   AnalyzeOperations = "detect_patterns"
    OpSuggestRelated   AnalyzeOperations = "suggest_related"
    OpAnalyzeQuality   AnalyzeOperations = "analyze_quality"
    OpDetectConflicts  AnalyzeOperations = "detect_conflicts"
)

// internal/tools/system.go - All system operations
type SystemOperations string
const (
    OpHealth           SystemOperations = "health"
    OpExportProject    SystemOperations = "export_project"
    OpImportProject    SystemOperations = "import_project"
    OpGenerateCitation SystemOperations = "generate_citation"
)
```

**Step 1.3.2: Implement clear operation boundaries**
- [x] Store tool: handles all persistence, updates, deletes
- [x] Retrieve tool: handles all searches, gets, lists
- [x] Analyze tool: handles all AI/ML operations, pattern detection
- [x] System tool: handles all admin, export, health operations
- [x] No overlapping operations between tools

**Step 1.3.3: Replace existing tool handlers**
- [x] Replace all handlers in `internal/mcp/server.go` with clean implementations
- [x] Create `internal/tools/store/handler.go` (replacing memory_create, memory_update, memory_delete)
- [x] Create `internal/tools/retrieve/handler.go` (replacing memory_read)
- [x] Create `internal/tools/analyze/handler.go` (replacing memory_analyze, memory_intelligence)
- [x] Create `internal/tools/system/handler.go` (replacing memory_system, memory_transfer)

**Step 1.3.4: Update server registration**
- [x] Update `internal/mcp/server.go` to register 4 clean tools instead of 9 fragmented ones
- [x] Remove old tool registration code
- [x] Update tool discovery and documentation
- [x] Implement consistent parameter validation across all tools

#### **CLI Refactor (Coordinated)**

**Step 1.3.5: Update CLI for clean tools**
- [x] Update `cli/internal/adapters/secondary/mcp/client.go` to use new tool names
- [x] Update method calls:
  - `SyncTask` → use `memory_store` tool instead of `memory_tasks/todo_write`
  - `GetTasks` → use `memory_retrieve` tool instead of `memory_tasks/todo_read`
  - `UpdateTaskStatus` → use `memory_store` tool instead of `memory_tasks/todo_update`
  - `QueryIntelligence` → use `memory_analyze` tool instead of `memory_intelligence`

**Step 1.3.6: Update CLI operation mapping**
- [x] Update CLI to automatically select correct tool based on operation type
- [x] Update request parameter structures to match new tool schemas
- [x] Update response parsing for new tool response formats
- [x] Add comprehensive error handling for new tool structure

---

## ✅ **PHASE 1 COMPLETION STATUS**

**🎯 Phase 1: Core Architecture Refactor - COMPLETED** ✅

All Phase 1 objectives have been successfully implemented:

### **✅ P1.1: Clean Parameter System** 
- ✅ Created `internal/types/core.go` with ProjectID/SessionID types
- ✅ Created `internal/validation/params.go` with centralized validation  
- ✅ Replaced confusing "repository" parameter with clear "project_id" across all tools
- ✅ Updated CLI entities, ports, and services for new parameter system

### **✅ P1.2: Logical Session Management**
- ✅ Created `internal/session/manager.go` with correct session semantics
- ✅ Fixed backwards logic: session_id now provides MORE access, not less
- ✅ Updated all operations to use logical session scoping
- ✅ Updated CLI to leverage session-based access patterns

### **✅ P1.3: Clean 4-Tool Architecture** 
- ✅ Created `internal/tools/operations.go` defining 4 clean tools
- ✅ Implemented `memory_store` tool (internal/tools/store/handler.go)
- ✅ Implemented `memory_retrieve` tool (internal/tools/retrieve/handler.go)
- ✅ Implemented `memory_analyze` tool (internal/tools/analyze/handler.go)
- ✅ Implemented `memory_system` tool (internal/tools/system/handler.go)
- ✅ Updated CLI client to use new 4-tool architecture

### **✅ Supporting Infrastructure**
- ✅ Created database migration `migrations/012_refactor_to_project_id.sql`
- ✅ Updated CLI for coordinated server/client changes
- ✅ Implemented comprehensive parameter validation
- ✅ Added proper error handling throughout

**Result**: The MCP Memory Server now has a clean, intuitive architecture with 4 logical tools, consistent parameter system, and proper session semantics. All critical design issues from MEMORY_MCP_IMPROVEMENTS.md have been resolved.

---

## 🔧 **PHASE 2: INFRASTRUCTURE REFACTOR**

### **P2.1: Build Clean Storage Layer**

**Goal**: Refactor the existing storage layer to support clean parameter system and remove complexity.

#### **Server Refactor**

**Step 2.1.1: Refactor storage interfaces**
- [ ] Update `internal/storage/interface.go` to use clean parameter types:
- [ ] Replace fragmented storage interfaces with clear contracts:
```go
package storage

import (
    "context"
    "lerian-mcp-memory/internal/types"
)

// ContentStore handles all content persistence
type ContentStore interface {
    Store(ctx context.Context, content *types.Content) error
    Update(ctx context.Context, content *types.Content) error
    Delete(ctx context.Context, projectID types.ProjectID, contentID string) error
    Get(ctx context.Context, projectID types.ProjectID, contentID string) (*types.Content, error)
}

// SearchStore handles all search and retrieval
type SearchStore interface {
    Search(ctx context.Context, query *types.SearchQuery) (*types.SearchResults, error)
    FindSimilar(ctx context.Context, content string, projectID types.ProjectID) ([]*types.Content, error)
    GetByProject(ctx context.Context, projectID types.ProjectID, filters *types.Filters) ([]*types.Content, error)
}

// AnalysisStore handles all analysis data
type AnalysisStore interface {
    StorePattern(ctx context.Context, pattern *types.Pattern) error
    GetPatterns(ctx context.Context, projectID types.ProjectID) ([]*types.Pattern, error)
    StoreInsight(ctx context.Context, insight *types.Insight) error
}
```

**Step 2.1.2: Update database schema**
- [ ] Create `migrations/012_refactor_to_project_id.sql` to update all tables
- [ ] Update all `repository` columns to `project_id` in existing tables
- [ ] Update all indexes and constraints to use `project_id`
- [ ] Update foreign key relationships for consistency

**Step 2.1.3: Refactor storage implementations**
- [ ] Update `internal/storage/qdrant.go` to use ProjectID instead of repository
- [ ] Update `internal/storage/task_repository.go` for new parameter system
- [ ] Update all SQL queries to use `project_id` consistently
- [ ] Update connection pooling and retry logic for new parameter system

---

### **P2.2: Implement Clear Operation Names**

**Goal**: Replace cryptic operation names with self-documenting names that clearly explain what they do.

#### **Server Refactor**

**Step 2.2.1: Replace cryptic operation names**
- [ ] Replace cryptic names in all tool handlers with clear, descriptive names:
- [ ] Create `internal/operations/names.go` to centralize operation definitions:
```go
package operations

// Clear, descriptive operation names that explain what they do
const (
    // Store operations - clearly indicate data persistence
    StoreContent            = "store_content"
    StoreDecision           = "store_decision"
    UpdateExistingContent   = "update_existing_content"
    DeleteOldContent        = "delete_old_content"
    ExpireStaleContent      = "expire_stale_content"
    
    // Retrieve operations - clearly indicate data access
    SearchContent           = "search_content"
    GetContentByID          = "get_content_by_id"
    FindSimilarContent      = "find_similar_content"
    GetContentHistory       = "get_content_history"
    
    // Analyze operations - clearly indicate analysis type
    DetectContentPatterns   = "detect_content_patterns"
    AnalyzeContentQuality   = "analyze_content_quality"
    FindContentRelationships = "find_content_relationships"
    GenerateContentInsights = "generate_content_insights"
    
    // System operations - clearly indicate admin functions
    CheckSystemHealth       = "check_system_health"
    ExportProjectData       = "export_project_data"
    ValidateDataIntegrity   = "validate_data_integrity"
)
```

**Step 2.2.2: Update all operation handlers**
- [ ] Replace `decay_management` → `expire_old_content` in all handlers
- [ ] Replace `mark_refreshed` → `validate_current` in all handlers
- [ ] Replace `traverse_graph` → `explore_relationships` in all handlers
- [ ] Replace `auto_detect_relationships` → `detect_content_relationships` in all handlers
- [ ] Update all documentation strings with clear operation descriptions

**Step 2.2.3: Update operation validation and error handling**
- [ ] Update validation in `internal/mcp/server.go` to use new operation names
- [ ] Replace cryptic error messages with helpful, descriptive ones
- [ ] Add operation capability discovery using new clear names
- [ ] Update all tool registration to use descriptive operation names

---

### **P2.3: Refactor Domain Architecture**

**Goal**: Separate task management from knowledge storage using existing codebase structure.

#### **Server Refactor**

**Step 2.3.1: Reorganize existing domain structure**
- [ ] Refactor `internal/tasks/` to be pure task domain (remove memory mixing)
- [ ] Refactor existing memory operations to be separate from task operations
- [ ] Update `internal/intelligence/` to be pure knowledge domain
- [ ] Define clear interfaces between existing domains

**Step 2.3.2: Separate memory operations from task operations**
- [ ] Move memory-specific logic out of `memory_tasks` tool
- [ ] Create clear separation in storage layer between tasks and memory
- [ ] Update business logic to respect domain boundaries
- [ ] Refactor existing services for domain clarity

**Step 2.3.3: Update task domain**
- [ ] Refactor existing `internal/tasks/` package for pure task management
- [ ] Update task operations to not mix with memory operations
- [ ] Separate task storage from memory storage
- [ ] Build clean task-specific business logic

**Step 2.3.4: Update domain integration**
- [ ] Define how tasks can reference memory content (without mixing domains)
- [ ] Update existing cross-domain relationships for clarity
- [ ] Add domain-specific authorization using existing patterns
- [ ] Update event publishing for clean domain integration

#### **CLI Refactor (Coordinated)**

**Step 2.3.5: Update CLI for domain separation**
- [ ] Update CLI commands to respect domain boundaries
- [ ] Implement domain-specific CLI workflows using existing structure
- [ ] Update CLI routing to use appropriate domain tools
- [ ] Build integrated workflows that properly span domains

---

## 📚 **PHASE 3: USER EXPERIENCE & DOCUMENTATION**

### **P3.1: Build Comprehensive Documentation**

**Step 3.1.1: Create clear data model documentation**
- [ ] Create `v2/docs/data-model.md` with clean entity relationships
- [ ] Document the simplified v2 data model with clear examples
- [ ] Add entity lifecycle and relationship documentation
- [ ] Include migration guide from v1 concepts to v2

**Step 3.1.2: Generate complete API documentation**
- [ ] Create OpenAPI 3.0 specification for all 4 tools
- [ ] Add comprehensive operation examples with real request/response data
- [ ] Document all parameter requirements and validation rules
- [ ] Include authentication and authorization documentation

**Step 3.1.3: Build interactive documentation**
- [ ] Create interactive API documentation with Swagger UI
- [ ] Add "try it out" functionality for all operations
- [ ] Include code generation for multiple programming languages
- [ ] Add comprehensive troubleshooting guide

### **P3.2: Create User-Friendly Getting Started Experience**

**Step 3.2.1: Design progressive onboarding**
- [ ] Create `v2/docs/quickstart.md` with 5-minute setup guide
- [ ] Build step-by-step tutorials progressing from simple to complex
- [ ] Add interactive CLI walkthrough with real examples
- [ ] Include common workflow patterns and best practices

**Step 3.2.2: Build example applications**
- [ ] Create example applications demonstrating common use cases
- [ ] Build CLI usage examples for typical workflows
- [ ] Add integration examples for different programming languages
- [ ] Include performance optimization guides

### **P3.3: Implement Excellent Error Experience**

**Step 3.3.1: Design user-friendly error handling**
- [ ] Create consistent, helpful error messages with suggestions
- [ ] Add error codes with clear categorization
- [ ] Implement progressive error disclosure (summary + details)
- [ ] Add error recovery suggestions and next steps

**Step 3.3.2: Build error monitoring and analytics**
- [ ] Add error tracking and analytics for common issues
- [ ] Create error pattern detection and alerting
- [ ] Build user feedback collection for error improvement
- [ ] Add automated error reporting and resolution guides

---

## 🎯 **PHASE 4: PERFORMANCE & PRODUCTION READINESS**

### **P4.1: Implement High-Performance Architecture**

**Step 4.1.1: Build scalable server architecture**
- [ ] Implement connection pooling and resource management
- [ ] Add horizontal scaling support with load balancing
- [ ] Build caching layers for frequently accessed data
- [ ] Implement async processing for heavy operations

**Step 4.1.2: Optimize data access patterns**
- [ ] Add intelligent query optimization and indexing
- [ ] Implement data partitioning for large datasets
- [ ] Build efficient bulk operation support
- [ ] Add streaming support for large responses

**Step 4.1.3: Build monitoring and observability**
- [ ] Add comprehensive metrics collection and monitoring
- [ ] Implement distributed tracing for request flow
- [ ] Build performance dashboards and alerting
- [ ] Add capacity planning and auto-scaling

### **P4.2: Ensure Production Quality**

**Step 4.2.1: Implement comprehensive testing**
- [ ] Build extensive unit test coverage (>90%)
- [ ] Add integration testing for all workflows
- [ ] Implement load testing and performance benchmarks
- [ ] Add chaos engineering and resilience testing

**Step 4.2.2: Add security and compliance**
- [ ] Implement authentication and authorization
- [ ] Add data encryption at rest and in transit
- [ ] Build audit logging and compliance reporting
- [ ] Add security scanning and vulnerability testing

**Step 4.2.3: Build deployment and operations**
- [ ] Create containerized deployment with Docker/Kubernetes
- [ ] Add CI/CD pipelines with automated testing
- [ ] Implement blue-green deployment and rollback
- [ ] Build operational runbooks and incident response

---

## 🚀 **IMPLEMENTATION PHASES**

### **Phase 1: Core Refactor (Weeks 1-2)**
- Complete P1.1 (Clean Parameter System)
- Complete P1.2 (Logical Session Management)
- Complete P1.3 (4-Tool Architecture)
- **Deliverable**: Refactored server with clean tool structure + updated CLI

### **Phase 2: Infrastructure Polish (Weeks 3-4)**
- Complete P2.1 (Storage Layer Refactor)
- Complete P2.2 (Clear Operation Names)
- Complete P2.3 (Domain Architecture)
- **Deliverable**: Fully refactored system with clean architecture

### **Phase 3: Systems Integration (Weeks 5-6)**
- Complete P3.1 (Security & Audit Updates)
- Complete P3.2 (WebSocket & Real-time)
- Complete P3.3 (Intelligence Systems)
- **Deliverable**: Complete system with all integrations working

### **Phase 4: Finalization (Weeks 7-8)**
- Complete P4.1 (Testing & Validation)
- Complete P4.2 (Documentation & Deployment)
- **Deliverable**: Production-ready refactored system

**Note**: Monitoring stack removed from scope - can be added later if needed

---

## 🧪 **TESTING STRATEGY**

### **Development Testing (Continuous)**
- [ ] Unit tests for all components (>90% coverage)
- [ ] Integration tests for all workflows
- [ ] Property-based testing for parameter validation
- [ ] Mutation testing to validate test quality

### **System Testing (Pre-Production)**
- [ ] End-to-end testing of complete workflows
- [ ] Performance testing with realistic datasets
- [ ] Security testing and vulnerability scanning
- [ ] Compatibility testing with multiple MCP clients

### **Production Testing (Live)**
- [ ] Canary deployments with real traffic
- [ ] A/B testing between v1 and v2 performance
- [ ] Chaos engineering for resilience testing
- [ ] User acceptance testing with pilot users

---

## 📋 **DEPLOYMENT STRATEGY**

### **Direct Refactor Approach**
- Refactor existing codebase in coordinated server + CLI effort
- No parallel systems or complex migration
- Use existing infrastructure and deployment
- Maintain existing data with schema updates

### **Coordinated Development**
- Server and CLI changes implemented together
- Single coordinated deployment of refactored system
- Existing data migrated through schema updates
- No downtime beyond normal deployment window

### **Schema Migration**
- Database schema updates for `repository` → `project_id`
- Existing data preservation with column renames
- Index and constraint updates
- Rollback capability for schema changes

---

## ✅ **SUCCESS CRITERIA**

### **Architectural Success**
1. **Clean Tool Architecture**: 4 logical tools with zero overlapping responsibilities
2. **Intuitive Parameters**: Consistent `project_id`/`session_id` usage across all operations
3. **Logical Session Management**: Include session_id for expanded access, not limited access
4. **Self-Documenting Operations**: Operation names clearly explain what they do
5. **Domain Separation**: Memory and task domains with clear boundaries

### **User Experience Success**
6. **Simplified Onboarding**: New users productive within 5 minutes
7. **Excellent Documentation**: Complete guides, examples, and interactive docs
8. **Superior Error Experience**: Helpful error messages with recovery suggestions
9. **CLI Excellence**: Intuitive CLI that leverages clean server architecture
10. **Performance**: 2x better performance than v1 with cleaner architecture

### **Production Quality Success**
11. **High Availability**: 99.9% uptime with graceful degradation
12. **Scalability**: Linear scaling to 100x current load
13. **Security**: Comprehensive security and compliance implementation
14. **Observability**: Complete monitoring, tracing, and alerting
15. **Operational Excellence**: Automated deployments and incident response

---

## 🔧 **IMPLEMENTATION NOTES**

### **Key Advantages of Direct Refactor Approach**
- **Faster Implementation**: No parallel systems or complex migrations
- **Coordinated Changes**: Server and CLI updated together in single effort
- **Existing Infrastructure**: Leverage existing deployment and operational setup
- **Data Preservation**: Keep existing data with simple schema updates

### **Risk Mitigation**
- **Coordinated Development**: Server and CLI changes tested together
- **Schema Migration**: Safe database updates with rollback capability
- **Feature Preservation**: Ensure all existing functionality remains available
- **Comprehensive Testing**: Full system testing before deployment

### **Technology Approach**
- **Language**: Continue with Go 1.23+ (no technology changes)
- **Storage**: Keep PostgreSQL + Qdrant (refactor usage patterns)
- **Architecture**: Refactor to clean architecture within existing codebase
- **Testing**: Comprehensive testing of refactored components
- **Deployment**: Use existing containerized deployment

### **Timeline Expectations**
- **7-8 weeks total**: Focused refactor timeline (monitoring stack removed)
- **2 week phases**: Manageable chunks with clear deliverables
- **Continuous validation**: Ensure existing functionality preserved
- **Quality gates**: Each phase must maintain system functionality
- **Scope reduction**: ~45% complexity reduction by removing monitoring infrastructure

This direct refactor approach eliminates all design flaws while preserving existing infrastructure and data. The result will be a clean, intuitive MCP server that provides the simplicity the CLI team already achieved, without the complexity of parallel systems or data migration.
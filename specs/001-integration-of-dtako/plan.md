# Implementation Plan: dtako_mod Integration

**Branch**: `001-integration-of-dtako` | **Date**: 2025-09-12 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-integration-of-dtako/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → Successfully loaded spec for dtako_mod integration
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detected: Production DB auth, sync frequency, performance requirements
   → Set Project Type: single (Go module/library)
   → Set Structure Decision: Option 1 (single project)
3. Evaluate Constitution Check section below
   → No violations detected
   → Update Progress Tracking: Initial Constitution Check PASS
4. Execute Phase 0 → research.md
   → Researching clarification items
5. Execute Phase 1 → contracts, data-model.md, quickstart.md
   → Will generate after research complete
6. Re-evaluate Constitution Check section
   → To be checked after Phase 1
7. Plan Phase 2 → Describe task generation approach
8. STOP - Ready for /tasks command
```

## Summary
Integration of dtako_mod as a submodule for ryohi_sub_cal2 router to import production data from dtako_rows, dtako_events, and dtako_ferry tables. The module provides REST API endpoints for importing and querying vehicle operational data, event logs, and ferry operations from production databases to local storage.

## Technical Context
**Language/Version**: Go 1.21  
**Primary Dependencies**: go-chi/chi v5 (router), go-sql-driver/mysql (database), godotenv (config)  
**Storage**: MySQL (production source + local destination)  
**Testing**: Go testing package with integration tests  
**Target Platform**: Linux/Windows servers  
**Project Type**: single - Go module/library with API endpoints  
**Performance Goals**: Import operations complete within 5 minutes for 1 month of data  
**Constraints**: Minimize production DB load, handle network interruptions gracefully  
**Scale/Scope**: ~10k records per day across 3 tables, support concurrent imports

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Simplicity**:
- Projects: 1 (dtako_mod Go module)
- Using framework directly? Yes (chi router, no wrappers)
- Single data model? Yes (direct table mappings)
- Avoiding patterns? Yes (simple repository pattern for DB access)

**Architecture**:
- EVERY feature as library? Yes (dtako_mod is a library)
- Libraries listed: 
  - dtako_mod: Import/query production data for dtako tables
- CLI per library: API endpoints serve as interface
- Library docs: README.md with API documentation

**Testing (NON-NEGOTIABLE)**:
- RED-GREEN-Refactor cycle enforced? Yes
- Git commits show tests before implementation? Will enforce
- Order: Contract→Integration→E2E→Unit strictly followed? Yes
- Real dependencies used? Yes (actual MySQL databases)
- Integration tests for: new libraries, contract changes, shared schemas? Yes
- FORBIDDEN: Implementation before test, skipping RED phase

**Observability**:
- Structured logging included? Yes (import results, errors)
- Frontend logs → backend? N/A (API only)
- Error context sufficient? Yes (detailed error responses)

**Versioning**:
- Version number assigned? 1.0.0
- BUILD increments on every change? Yes
- Breaking changes handled? Will use semantic versioning

## Project Structure

### Documentation (this feature)
```
specs/001-integration-of-dtako/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
# Option 1: Single project (DEFAULT)
handlers/
├── dtako_rows.go
├── dtako_events.go
└── dtako_ferry.go

services/
├── dtako_rows_service.go
├── dtako_events_service.go
└── dtako_ferry_service.go

repositories/
├── database.go
├── dtako_rows_repository.go
├── dtako_events_repository.go
└── dtako_ferry_repository.go

models/
└── models.go

config/
└── config.go

tests/
├── contract/
├── integration/
└── unit/
```

**Structure Decision**: Option 1 - Single Go module project

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context**:
   - Production database authentication method
   - Data synchronization frequency requirements
   - Deduplication strategy for imports
   - Maximum date range for imports
   - Production database connection details

2. **Generate and dispatch research agents**:
   ```
   Task: "Research MySQL authentication best practices for production access"
   Task: "Find optimal batch import strategies for MySQL"
   Task: "Research deduplication patterns for data imports"
   Task: "Determine reasonable date range limits for bulk imports"
   ```

3. **Consolidate findings** in `research.md`:
   - Decision: Use environment variables for DB credentials
   - Rationale: Standard practice, secure, configurable
   - Alternatives: Certificate auth (more complex), OAuth (overkill)

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - DtakoRow: vehicle operation records
   - DtakoEvent: system/operational events
   - DtakoFerry: ferry operation records
   - ImportResult: import operation results

2. **Generate API contracts** from functional requirements:
   - GET /dtako/rows - List rows
   - POST /dtako/rows/import - Import rows
   - GET /dtako/rows/{id} - Get row by ID
   - GET /dtako/events - List events
   - POST /dtako/events/import - Import events
   - GET /dtako/events/{id} - Get event by ID
   - GET /dtako/ferry - List ferry records
   - POST /dtako/ferry/import - Import ferry records
   - GET /dtako/ferry/{id} - Get ferry record by ID

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Schema validation tests
   - Tests must fail initially

4. **Extract test scenarios** from user stories:
   - Import with date range
   - Filter by event type
   - Filter by ferry route
   - Handle production DB failure
   - Prevent duplicate imports

5. **Update agent file incrementally**:
   - Add Go, MySQL, chi router context
   - Document API endpoints
   - Include test requirements

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, CLAUDE.md

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Generate from contracts and data model
- Each endpoint → contract test task [P]
- Each entity → model + repository task [P]
- Each service → business logic task
- Integration tests for import flows
- Configuration and database setup tasks

**Ordering Strategy**:
- Database schema first
- Models and repositories [P]
- Services layer
- Handlers/controllers
- Contract tests [P]
- Integration tests
- Documentation

**Estimated Output**: 25-30 numbered, ordered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*No violations detected - section empty*

## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (none)

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
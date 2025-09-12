# Tasks: dtako_mod Integration

**Input**: Design documents from `/specs/001-integration-of-dtako/`
**Prerequisites**: plan.md (✓), research.md (✓), data-model.md (✓), contracts/ (✓)

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Tech stack: Go 1.21, chi router, MySQL
   → Structure: Single Go module project
2. Load optional design documents:
   → data-model.md: 3 entities (DtakoRow, DtakoEvent, DtakoFerry)
   → contracts/openapi.yaml: 9 endpoints (3 per entity)
   → research.md: Environment variables for auth
3. Generate tasks by category:
   → Setup: go mod init, dependencies, database
   → Tests: 9 contract tests, 5 integration tests
   → Core: 3 models, 3 services, 3 handlers
   → Integration: DB connections, middleware
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001-T040)
6. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Phase 3.1: Setup
- [ ] T001 Verify Go 1.21+ installed and initialize go.mod if not exists
- [ ] T002 Install dependencies: go get github.com/go-chi/chi/v5 github.com/go-sql-driver/mysql github.com/joho/godotenv
- [ ] T003 [P] Create project directory structure: handlers/, services/, repositories/, models/, config/, tests/
- [ ] T004 [P] Create .env from .env.example with test database configuration
- [ ] T005 Execute schema.sql to create local test database tables (dtako_rows, dtako_events, dtako_ferry)

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests (API endpoints)
- [ ] T006 [P] Contract test GET /dtako/rows in tests/contract/rows_list_test.go
- [ ] T007 [P] Contract test POST /dtako/rows/import in tests/contract/rows_import_test.go  
- [ ] T008 [P] Contract test GET /dtako/rows/{id} in tests/contract/rows_get_test.go
- [ ] T009 [P] Contract test GET /dtako/events in tests/contract/events_list_test.go
- [ ] T010 [P] Contract test POST /dtako/events/import in tests/contract/events_import_test.go
- [ ] T011 [P] Contract test GET /dtako/events/{id} in tests/contract/events_get_test.go
- [ ] T012 [P] Contract test GET /dtako/ferry in tests/contract/ferry_list_test.go
- [ ] T013 [P] Contract test POST /dtako/ferry/import in tests/contract/ferry_import_test.go
- [ ] T014 [P] Contract test GET /dtako/ferry/{id} in tests/contract/ferry_get_test.go

### Integration Tests (User scenarios)
- [ ] T015 [P] Integration test: Import dtako_rows with date range in tests/integration/rows_import_scenario_test.go
- [ ] T016 [P] Integration test: Filter dtako_events by type in tests/integration/events_filter_scenario_test.go
- [ ] T017 [P] Integration test: Filter dtako_ferry by route in tests/integration/ferry_route_scenario_test.go
- [ ] T018 [P] Integration test: Handle production DB failure gracefully in tests/integration/db_failure_test.go
- [ ] T019 [P] Integration test: Prevent duplicate imports (UPSERT) in tests/integration/duplicate_handling_test.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Models
- [ ] T020 [P] Create DtakoRow model in models/models.go (already exists, verify structure)
- [ ] T021 [P] Create DtakoEvent model in models/models.go (already exists, verify structure)
- [ ] T022 [P] Create DtakoFerry model in models/models.go (already exists, verify structure)
- [ ] T023 [P] Create ImportResult model in models/models.go (already exists, verify structure)

### Database Layer
- [ ] T024 Create database connection manager in repositories/database.go (already exists, verify)
- [ ] T025 [P] Implement DtakoRowsRepository in repositories/dtako_rows_repository.go (already exists, verify)
- [ ] T026 [P] Implement DtakoEventsRepository in repositories/dtako_events_repository.go (already exists, verify)
- [ ] T027 [P] Implement DtakoFerryRepository in repositories/dtako_ferry_repository.go (already exists, verify)

### Service Layer
- [ ] T028 [P] Implement DtakoRowsService business logic in services/dtako_rows_service.go (already exists, verify)
- [ ] T029 [P] Implement DtakoEventsService business logic in services/dtako_events_service.go (already exists, verify)
- [ ] T030 [P] Implement DtakoFerryService business logic in services/dtako_ferry_service.go (already exists, verify)

### API Handlers
- [ ] T031 [P] Implement DtakoRowsHandler endpoints in handlers/dtako_rows.go (already exists, verify)
- [ ] T032 [P] Implement DtakoEventsHandler endpoints in handlers/dtako_events.go (already exists, verify)
- [ ] T033 [P] Implement DtakoFerryHandler endpoints in handlers/dtako_ferry.go (already exists, verify)

### Router Registration
- [ ] T034 Implement RegisterRoutes function in main.go to wire all endpoints (already exists, verify)

## Phase 3.4: Integration
- [ ] T035 Add configuration loading from environment in config/config.go (already exists, verify)
- [ ] T036 Add structured logging middleware for request/response tracking
- [ ] T037 Add error recovery middleware to handle panics gracefully
- [ ] T038 Add timeout middleware for long-running import operations (5 minutes)

## Phase 3.5: Polish
- [ ] T039 [P] Add unit tests for validation logic in tests/unit/validation_test.go
- [ ] T040 [P] Add unit tests for date range handling in tests/unit/date_range_test.go
- [ ] T041 [P] Performance test: Import 10k records within 30 seconds in tests/performance/import_benchmark_test.go
- [ ] T042 [P] Update README.md with actual API examples and deployment instructions
- [ ] T043 Create Makefile with common commands (test, build, run, clean)
- [ ] T044 Run quickstart.md scenarios for manual validation

## Dependencies
- Setup (T001-T005) must complete first
- All tests (T006-T019) before ANY implementation (T020-T034)
- Database layer (T024-T027) before service layer (T028-T030)
- Service layer before handlers (T031-T033)
- All implementation before integration (T035-T038)
- Everything before polish (T039-T044)

## Parallel Execution Examples

### Batch 1: Contract Tests (after setup)
```bash
# Launch T006-T014 together (all different files):
Task: "Contract test GET /dtako/rows in tests/contract/rows_list_test.go"
Task: "Contract test POST /dtako/rows/import in tests/contract/rows_import_test.go"
Task: "Contract test GET /dtako/rows/{id} in tests/contract/rows_get_test.go"
Task: "Contract test GET /dtako/events in tests/contract/events_list_test.go"
Task: "Contract test POST /dtako/events/import in tests/contract/events_import_test.go"
Task: "Contract test GET /dtako/events/{id} in tests/contract/events_get_test.go"
Task: "Contract test GET /dtako/ferry in tests/contract/ferry_list_test.go"
Task: "Contract test POST /dtako/ferry/import in tests/contract/ferry_import_test.go"
Task: "Contract test GET /dtako/ferry/{id} in tests/contract/ferry_get_test.go"
```

### Batch 2: Integration Tests
```bash
# Launch T015-T019 together (all different files):
Task: "Integration test: Import dtako_rows with date range in tests/integration/rows_import_scenario_test.go"
Task: "Integration test: Filter dtako_events by type in tests/integration/events_filter_scenario_test.go"
Task: "Integration test: Filter dtako_ferry by route in tests/integration/ferry_route_scenario_test.go"
Task: "Integration test: Handle production DB failure gracefully in tests/integration/db_failure_test.go"
Task: "Integration test: Prevent duplicate imports (UPSERT) in tests/integration/duplicate_handling_test.go"
```

### Batch 3: Repository Layer (after tests fail)
```bash
# Launch T025-T027 together (all different files):
Task: "Implement DtakoRowsRepository in repositories/dtako_rows_repository.go"
Task: "Implement DtakoEventsRepository in repositories/dtako_events_repository.go"
Task: "Implement DtakoFerryRepository in repositories/dtako_ferry_repository.go"
```

### Batch 4: Service Layer
```bash
# Launch T028-T030 together (all different files):
Task: "Implement DtakoRowsService business logic in services/dtako_rows_service.go"
Task: "Implement DtakoEventsService business logic in services/dtako_events_service.go"
Task: "Implement DtakoFerryService business logic in services/dtako_ferry_service.go"
```

### Batch 5: Handlers
```bash
# Launch T031-T033 together (all different files):
Task: "Implement DtakoRowsHandler endpoints in handlers/dtako_rows.go"
Task: "Implement DtakoEventsHandler endpoints in handlers/dtako_events.go"
Task: "Implement DtakoFerryHandler endpoints in handlers/dtako_ferry.go"
```

## Notes
- Many files already exist from initial scaffold - tasks should verify and fix to make tests pass
- [P] tasks use different files and have no dependencies
- Commit after each task with descriptive message
- Run `go test ./...` after each implementation task to verify progress
- Use `go test -v` for detailed test output during debugging

## Validation Checklist
- [x] All 9 API endpoints have corresponding contract tests
- [x] All 3 entities have model tasks
- [x] All 4 entities (including ImportResult) have model definitions
- [x] All tests (T006-T019) come before implementation (T020+)
- [x] Parallel tasks are truly independent (different files)
- [x] Each task specifies exact file path
- [x] No [P] task modifies same file as another [P] task

## Estimated Completion Time
- Setup: 15 minutes
- Tests: 2 hours (writing failing tests)
- Implementation: 3 hours (making tests pass)
- Integration: 1 hour
- Polish: 1 hour
- **Total**: ~7-8 hours for experienced Go developer
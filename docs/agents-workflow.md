# Agent Workflows

Multi-agent workflows enable complex tasks to be completed by coordinating multiple specialized agents. The Mimir orchestrator defines and executes workflows that span multiple agents, handling handoffs, dependencies, and error recovery.

## Overview

Agent workflows:
- Define sequences of agent tasks
- Coordinate handoffs between agents
- Handle parallel and sequential execution
- Manage dependencies and conditions
- Track progress via GitHub issues
- Support looping and branching patterns

## Workflow Patterns

### Sequential Workflow

Agents execute one after another, with each passing results to the next.

```mermaid
graph LR
    A[Agent 1] --> B[Agent 2]
    B --> C[Agent 3]
    C --> D[Agent 4]
```

**Use Cases**:
- Code review → Fix → Re-review
- Analysis → Implementation → Testing
- Documentation → Review → Publish

**Example**: Sentry Issue Resolution
1. Sentry Triage Agent analyzes event
2. Backend Golang Agent fixes bug
3. Code Review Agent reviews changes
4. Orchestrator merges PR

### Parallel Workflow

Multiple agents execute simultaneously on independent tasks.

```mermaid
graph TD
    Start[Orchestrator] --> A[Agent 1]
    Start --> B[Agent 2]
    Start --> C[Agent 3]

    A --> End[Orchestrator]
    B --> End
    C --> End
```

**Use Cases**:
- Multiple independent bug fixes
- Simultaneous code reviews
- Parallel documentation updates

**Example**: Multi-Project Fix
1. Orchestrator spawns 3 agents for different sub-projects
2. All agents work in parallel
3. Orchestrator waits for all to complete

### Conditional Workflow

Agent selection based on analysis results or conditions.

```mermaid
graph TD
    Start[Start] --> A[Agent 1: Analyze]
    A --> Decision{Condition?}

    Decision -->|Path A| B[Agent 2a]
    Decision -->|Path B| C[Agent 2b]

    B --> D[Agent 3]
    C --> D
    D --> End[End]
```

**Use Cases**:
- Sentry triage determines if action needed
- Code review determines fix complexity
- Task analysis determines required agents

**Example**: Sentry Triage
1. Sentry Triage Agent analyzes event
2. If critical → Spawn Backend Golang Agent immediately
3. If low severity → Wait for human approval first

### Looping Workflow

Agent iterates until a condition is met (e.g., test-fix loop).

```mermaid
graph TD
    Start[Start] --> A[Agent: Execute Task]
    A --> B[Check Condition]
    B -->|Not Complete| A
    B -->|Complete| End[End]
```

**Use Cases**:
- Test-fix iteration until tests pass
- Code review iterations until approved
- Documentation refinement until complete

**Example**: Fix and Test Loop
1. Backend Golang Agent fixes bug
2. Run tests
3. If tests fail → Fix again
4. If tests pass → Create PR

### Hybrid Workflow

Combination of sequential, parallel, and conditional patterns.

```mermaid
graph TD
    Start[Orchestrator] --> A[Agent 1: Triage]
    A --> Decision{Action Needed?}

    Decision -->|Yes| B[Agent 2: Fix]
    Decision -->|No| End[End]

    B --> C[Agent 3: Review]
    C --> Approved{Approved?}

    Approved -->|No| B
    Approved -->|Yes| D[Agent 4: Merge]
    D --> E[Agent 5: Cleanup]
    E --> End
```

**Use Cases**:
- Complex Sentry workflows
- Multi-step feature development
- Complete issue lifecycle

## Sentry Workflow Example

Complete end-to-end Sentry issue resolution workflow:

```mermaid
sequenceDiagram
    participant Sentry as Sentry
    participant O as Orchestrator
    participant ST as Sentry Triage Agent
    participant GH as GitHub
    participant BG as Backend Golang Agent
    participant CR as Code Review Agent

    Sentry->>O: Webhook: Error event
    O->>GH: Create issue #123
    O->>O: Analyze event
    O->>ST: Spawn session

    ST->>ST: Analyze error
    ST->>O: Analysis complete
    ST->>GH: Update issue #123

    O->>O: Determine action needed
    O->>BG: Spawn session

    BG->>BG: Create worktree
    BG->>BG: Fix bug
    BG->>BG: Run tests
    BG->>GH: Create PR #456
    BG->>GH: Update issue #123
    BG->>O: Fix complete

    O->>CR: Spawn session
    CR->>CR: Review PR #456
    CR->>GH: Update issue #123

    CR->>O: Review complete
    O->>O: Merge PR #456
    O->>GH: Merge PR #456
    O->>GH: Close issue #123
    O->>O: Cleanup sessions/worktrees
```

### Workflow Steps

| Step | Agent | Action | Output |
|------|-------|--------|--------|
| 1 | Orchestrator | Receive Sentry webhook | Event data |
| 2 | Orchestrator | Create GitHub issue | Issue #123 |
| 3 | Orchestrator | Spawn Sentry Triage Agent | Session created |
| 4 | Sentry Triage | Analyze error event | Analysis report |
| 5 | Sentry Triage | Update GitHub issue | Comment added |
| 6 | Orchestrator | Determine action needed | Decision: yes/no |
| 7 | Orchestrator | Spawn Backend Golang Agent | Session created |
| 8 | Backend Golang | Create git worktree | Worktree path |
| 9 | Backend Golang | Fix bug | Code changes |
| 10 | Backend Golang | Run tests | Test results |
| 11 | Backend Golang | Create PR | PR #456 |
| 12 | Backend Golang | Update GitHub issue | PR link added |
| 13 | Backend Golang | Complete task | Session done |
| 14 | Orchestrator | Spawn Code Review Agent | Session created |
| 15 | Code Review | Review PR | Review comments |
| 16 | Code Review | Update GitHub issue | Review added |
| 17 | Code Review | Complete task | Session done |
| 18 | Orchestrator | Merge PR | PR merged |
| 19 | Orchestrator | Close GitHub issue | Issue #123 closed |
| 20 | Orchestrator | Cleanup sessions/worktrees | Cleanup complete |

## Scheduled Code Review Workflow

```mermaid
graph TD
    Start[Cron Job] --> Orch[Orchestrator]
    Orch --> CreateIssue[Create GitHub Issue]
    CreateIssue --> SelectRepos[Select Repositories]
    SelectRepos --> Parallel[Spawn Agents in Parallel]

    Parallel --> Repo1[Agent for Repo 1]
    Parallel --> Repo2[Agent for Repo 2]
    Parallel --> Repo3[Agent for Repo 3]

    Repo1 --> Review1[Review Code]
    Repo2 --> Review2[Review Code]
    Repo3 --> Review3[Review Code]

    Review1 --> Update1[Update Issue]
    Review2 --> Update2[Update Issue]
    Review3 --> Update3[Update Issue]

    Update1 --> Wait1[Wait for All]
    Update2 --> Wait1
    Update3 --> Wait1

    Wait1 --> Aggregate[Aggregate Findings]
    Aggregate --> Final[Final Report]
    Final --> Close[Close Issue]
```

## Multi-Agent Task Workflow

Task requiring coordination between multiple agents:

```mermaid
sequenceDiagram
    participant User as User
    participant O as Orchestrator
    participant DA as Documentation Agent
    participant BG as Backend Golang Agent
    participant FT as Frontend TypeScript Agent

    User->>O: Task: Add user documentation
    O->>GH: Create issue
    O->>O: Analyze task

    O->>DA: Spawn session
    DA->>DA: Analyze API
    DA->>GH: Update issue with API specs

    DA->>O: API analysis complete
    O->>BG: Spawn session for backend docs
    O->>FT: Spawn session for frontend docs

    par Parallel execution
        BG->>BG: Write backend docs
    and
        FT->>FT: Write frontend docs
    end

    BG->>GH: Create PR for backend docs
    FT->>GH: Create PR for frontend docs

    BG->>O: Backend docs complete
    FT->>O: Frontend docs complete

    O->>O: Wait for both agents
    O->>GH: Merge both PRs
    O->>GH: Close issue
```

## Workflow Definition

### Workflow DSL (Conceptual)

Workflows are defined by the orchestrator using a simple DSL:

```yaml
name: sentry-resolution
description: Resolve Sentry issues automatically
steps:
  - name: triage
    agent: sentry-triage
    input:
      event: ${sentryEvent}
    output:
      analysis: ${triageAnalysis}
    condition:
      field: analysis.actionNeeded
      operator: equals
      value: true

  - name: fix
    agent: ${analysis.agentType}
    input:
      worktree: true
      task: ${analysis.task}
    output:
      pr: ${fixPR}

  - name: review
    agent: code-review
    input:
      pr: ${fixPR}
    condition:
      field: ${fixPR.status}
      operator: equals
      value: created

  - name: merge
    agent: orchestrator
    input:
      pr: ${fixPR}
      condition:
        field: review.approved
        operator: equals
        value: true
```

### Dynamic Workflow Generation

Workflows can be dynamically generated based on:
- Task analysis results
- Agent capabilities
- Project structure
- User configuration

## Handoff Mechanism

### Handoff via GitHub Issues

Agents communicate via GitHub issue comments:

```mermaid
sequenceDiagram
    participant A as Agent A
    participant GH as GitHub Issue
    participant O as Orchestrator
    participant B as Agent B

    A->>GH: Comment: Task complete, needs Agent B
    A->>GH: Comment: Analysis results...
    A->>O: Signal: Handoff requested

    O->>GH: Read handoff details
    O->>B: Spawn Agent B
    B->>GH: Read handoff from issue
    B->>B: Execute next task
```

### Handoff Data Structure

Handoff comments use a structured format:

```markdown
---
**Handoff from {Agent Name}**

**Task Completed**: {Description of completed task}
**Analysis**: {Key findings}
**Recommendation**: {What should happen next}
**Agent Requested**: {agent-name}
**Context**: {Relevant context, file paths, etc.}

**Artifacts**:
- PR: #{pr-number}
- Branch: mimir-{session-id}
- Files changed: {list of files}
---
```

## Workflow State Management

### State Tracking

The orchestrator tracks workflow state:

```mermaid
stateDiagram-v2
    [*] --> Created: Workflow started
    Created --> Running: First agent spawned
    Running --> Running: Agent handoff
    Running --> Paused: Awaiting input
    Paused --> Running: Input received
    Running --> Completed: All steps done
    Running --> Failed: Agent failure
    Failed --> Retrying: Retry attempt
    Retrying --> Running: Retry started
    Retrying --> Failed: Retry exhausted
```

### State Storage

Workflow state is stored in:
- **GitHub Issue**: Progress comments
- **SQLite Database**: Workflow metadata
- **Session Storage**: Agent session data

## Error Handling in Workflows

### Error Recovery Strategies

| Error Type | Recovery Strategy |
|------------|-------------------|
| Agent spawn failure | Retry up to 3 times, then fail |
| Agent execution timeout | Terminate, log error, create issue |
| Handoff failure | Use fallback handoff method |
| PR rejected | Restart agent with feedback |
| Merge conflict | Agent resolves, retries merge |
| Worktree creation failure | Use fallback directory |

### Workflow Rollback

When workflow fails:
1. Terminate all active agent sessions
2. Create GitHub issue with error details
3. Clean up worktrees and branches
4. Preserve logs for debugging
5. Notify human for intervention

## Workflow Monitoring

### Progress Tracking

Monitor workflow progress via:
- **GitHub Issue**: Real-time progress updates
- **Event Stream**: Agent events via opencode SDK
- **Dashboard**: Workflow status dashboard (future)

### Key Metrics

- **Workflow Duration**: Total time to complete
- **Agent Utilization**: Agent session time vs. idle time
- **Success Rate**: Percentage of successful workflows
- **Average Handoffs**: Number of handoffs per workflow
- **Error Rate**: Percentage of failed workflows

## Best Practices

### Workflow Design

1. **Keep workflows simple**: Break complex tasks into smaller sub-workflows
2. **Define clear handoff points**: Explicit data passing between agents
3. **Handle all error cases**: Define recovery strategies
4. **Track everything**: Log all state changes and decisions
5. **Use conditional branching**: Avoid unnecessary agent spawns

### Handoff Design

1. **Use structured comments**: Follow handoff data structure format
2. **Include all context**: Pass relevant files, analysis, recommendations
3. **Be explicit**: Clearly state what agent is needed next
4. **Link artifacts**: Include PR links, branch names, file paths

### Error Handling

1. **Fail fast**: Detect errors early
2. **Retry intelligently**: Exponential backoff, max retries
3. **Preserve state**: Don't lose progress on failure
4. **Notify humans**: Create issues when automation fails

## Next Steps

- [Agent Isolation](./agents-isolation.md) - Learn how agents work in isolated worktrees
- [GitHub Integration](./agents-github-integration.md) - Understand issue-based tracking
- [Examples](./agents-examples/) - See complete workflow implementations

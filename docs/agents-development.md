# Agent Development

Create custom agents for Mimir by following opencode's agent definition format. Agents are defined in markdown files with YAML frontmatter, making them 100% compatible with opencode.

## Overview

Agent development covers:
- **Agent definition format**: Following opencode's `.opencode/agent/*.md` structure
- **Agent capabilities**: Permissions, tools, and instructions
- **Plugin development**: Adding tools and agent definitions
- **Agent registration**: Registering agents in Mimir's registry
- **Testing agents**: Local testing and validation

## Agent Definition Format

### File Structure

Agent definitions follow opencode's format:

```markdown
---
description: Agent description
color: "#RRGGBB"
model: anthropic/claude-sonnet-4
mode: primary
temperature: 0.7
steps: 100
permission:
  edit:
    "**/*.go": allow
    "**/*.md": allow
    "*": ask
  bash:
    "go build *": allow
    "go test *": allow
    "*": ask
  external_directory: deny
---

You are a [agent type] with expertise in [area].

## Expertise Areas
- [Expertise area 1]
- [Expertise area 2]
- [Expertise area 3]

## Guidelines
- [Guideline 1]
- [Guideline 2]
- [Guideline 3]

## Available Tools
- [Tool 1]: Description
- [Tool 2]: Description
```

### Frontmatter Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `description` | string | Yes | Agent description |
| `color` | string | No | Color for UI display (hex) |
| `model` | string | No | Model identifier |
| `mode` | string | No | Mode: primary, secondary |
| `temperature` | number | No | Model temperature (0-2) |
| `steps` | number | No | Max steps/iterations |
| `permission` | object | No | Permission rules |

### Permission Structure

```yaml
permission:
  edit:
    "pattern1": allow/deny/ask
    "pattern2": allow/deny/ask
  bash:
    "command1": allow/deny/ask
    "command2": allow/deny/ask
  external_directory: deny/ask
  doom_loop: deny/ask
```

## Agent Types

### Orchestrator Agents

Example: Mimir Orchestrator

```markdown
---
description: Main orchestrator agent for Mimir - manages workflows, event routing, and agent coordination
color: "#9B59B6"
model: anthropic/claude-sonnet-4
mode: primary
temperature: 0.6
steps: 50
permission:
  external_directory: deny
---

You are the Mimir Orchestrator, responsible for coordinating multiple opencode agents.

## Responsibilities
- Manage multi-agent workflows
- Route events to appropriate agents
- Create and manage opencode sessions
- Track progress via GitHub issues
- Handle agent lifecycles and handoffs

## Available Tools
- `opencode_*`: opencode session management
- `github_*`: GitHub issue and PR management
- `sentry_*`: Sentry event analysis
- `git_*`: Git and worktree management
```

### Language-Specific Agents

Example: Backend Golang Agent

```markdown
---
description: Go backend development specialist with expertise in Gin, Echo, and microservices
color: "#00ADD8"
model: anthropic/claude-sonnet-4
mode: primary
temperature: 0.7
steps: 100
permission:
  edit:
    "**/*.go": allow
    "**/*.md": allow
    "go.mod": allow
    "go.sum": allow
    "*": ask
  bash:
    "go build *": allow
    "go test *": allow
    "go run *": allow
    "go mod *": allow
    "*": ask
  external_directory: deny
---

You are a Go backend development specialist.

## Expertise Areas
- Go language best practices
- Gin, Echo, Fiber frameworks
- gRPC and REST APIs
- Microservices architecture
- Testing with testify
- Error handling patterns

## Guidelines
- Follow Go best practices and effective Go
- Use proper error handling (wrap errors)
- Write tests for all new code
- Use interfaces for abstraction
- Prefer composition over inheritance
- Follow project's gofmt standards

## Testing Requirements
- Write unit tests for all functions
- Use table-driven tests where appropriate
- Aim for >80% code coverage
- Mock external dependencies

## Code Style
- Use `gofmt` and `goimports`
- Exported functions and types must have documentation
- Use descriptive variable names
- Keep functions focused and small
```

### Domain-Specific Agents

Example: Sentry Triage Agent

```markdown
---
description: Sentry event analysis and triage specialist
color: "#F6821F"
model: anthropic/claude-sonnet-4
mode: primary
temperature: 0.5
steps: 30
permission:
  edit:
    "**/*.md": allow
    "*": deny
  bash:
    "*": deny
  external_directory: deny
---

You are a Sentry event triage specialist.

## Responsibilities
- Analyze Sentry error events
- Determine severity and impact
- Identify root cause when possible
- Recommend appropriate action
- Create GitHub issues for actionable events

## Expertise Areas
- Sentry event analysis
- Error classification
- Root cause analysis
- Impact assessment
- Action prioritization

## Workflow
1. Receive Sentry event data
2. Analyze error stack trace
3. Identify affected components
4. Assess severity and impact
5. Determine if action is needed
6. Create GitHub issue with analysis
7. Request appropriate agent if action needed

## Available Tools
- `sentry_get_event`: Fetch Sentry event details
- `sentry_get_issue`: Fetch Sentry issue history
- `github_create_issue`: Create GitHub issue
- `github_add_comment`: Add comment to issue
```

## Agent Registration

### Centralized Registry

Agents are registered in `/docs/agents-examples/registry/`:

```
/docs/agents-examples/registry/
├── mimir-orchestrator.md
├── sentry-triage-agent.md
├── backend-golang-agent.md
├── backend-python-agent.md
├── backend-typescript-agent.md
├── frontend-typescript-agent.md
├── backend-rust-agent.md
├── documentation-agent.md
└── research-agent.md
```

### Registration Process

1. **Create agent definition**: Create `.md` file in registry
2. **Validate configuration**: Ensure frontmatter is valid
3. **Test agent**: Test agent in a sandbox environment
4. **Register**: Agent is automatically discovered by Mimir

### Project-Specific Overrides

Agents can be overridden per project in `.opencode/agent/`:

```
project/
└── .opencode/
    └── agent/
        └── backend-golang-agent.md  # Override for this project
```

Override merges with registry defaults:
- Registry defaults loaded first
- Project overrides applied on top
- Merge rules: Override takes precedence for most fields
- Permissions: Merged (union of registry + override)

## Plugin Development

### Plugin Structure

Plugins add tools and agent definitions:

```typescript
// .opencode/plugin/my-plugin.ts
import { Plugin } from "@opencode-ai/plugin"
import { tool } from "@opencode-ai/plugin"

export const MyPlugin: Plugin = async (input) => {
  const { client, directory, $ } = input

  return {
    // Define custom tools
    tool: {
      "my-tool": tool({
        description: "Tool description",
        args: {
          arg1: tool.schema.string().describe("Argument description")
        },
        async execute(args, context) {
          // Tool implementation
          return {
            title: "Result title",
            output: "Result output"
          }
        }
      })
    }
  }
}
```

### Agent Definition via Plugin

Plugins can also define agents:

```typescript
export const MyPlugin: Plugin = async (input) => {
  return {
    agent: {
      "my-custom-agent": {
        description: "Custom agent description",
        model: {
          modelID: "claude-sonnet-4",
          providerID: "anthropic"
        },
        temperature: 0.7,
        steps: 100,
        permission: {
          edit: {
            "**/*.ts": "allow",
            "*": "ask"
          }
        }
      }
    }
  }
}
```

### Tool Development

Tools provide specific capabilities to agents:

```typescript
tool({
  description: "Create GitHub issue with labels and assignees",
  args: {
    title: tool.schema.string().describe("Issue title"),
    body: tool.schema.string().describe("Issue description"),
    labels: tool.schema.array(tool.schema.string()).optional(),
    assignees: tool.schema.array(tool.schema.string()).optional()
  },
  async execute(args, context) {
    const { sessionID, messageID, abort } = context

    // Request permission
    const allowed = await context.ask({
      permission: "github.issue.create",
      pattern: args.title,
      message: `Create GitHub issue: ${args.title}?`
    })

    if (!allowed) {
      return { output: "Issue creation cancelled" }
    }

    // Create issue via GitHub API
    const result = await fetch("https://api.github.com/repos/OWNER/REPO/issues", {
      method: "POST",
      headers: {
        "Authorization": `token ${process.env.GITHUB_TOKEN}`,
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        title: args.title,
        body: args.body,
        labels: args.labels || [],
        assignees: args.assignees || []
      }),
      signal: abort
    })

    const issue = await result.json()

    return {
      title: `Created issue #${issue.number}`,
      output: `Issue created successfully\nURL: ${issue.html_url}`,
      metadata: {
        number: issue.number,
        url: issue.html_url
      }
    }
  }
})
```

## Agent Capabilities

### Reading and Analysis

Agents can read and analyze code:
- **Read files**: Use `read` tool
- **Search code**: Use `grep` or search tools
- **Analyze structure**: Use MCP filesystem tools
- **Parse code**: Use language-specific tools

### Execution

Agents can execute commands (with permissions):
- **Build**: Run build commands
- **Test**: Run test suites
- **Deploy**: Run deployment scripts
- **Lint**: Run linters and formatters

### Git Operations

Agents can work with git:
- **Create branches**: Create feature branches
- **Commit changes**: Commit to worktree branch
- **Push to remote**: Push to origin
- **Create PRs**: Create pull requests for review

### GitHub Integration

Agents can interact with GitHub:
- **Create issues**: Create tracking issues
- **Update issues**: Add comments and updates
- **Create PRs**: Create pull requests
- **Manage labels**: Add/remove labels

## Best Practices

### Agent Design

1. **Be specific**: Focus on narrow expertise
2. **Clear instructions**: Provide clear guidelines and workflows
3. **Define permissions**: Explicitly set permission rules
4. **Test thoroughly**: Test agent in various scenarios
5. **Document well**: Document capabilities and limitations

### Permission Management

1. **Principle of least privilege**: Only grant necessary permissions
2. **Explicit rules**: Use specific patterns, not wildcard `*`
3. **Ask by default**: Use `ask` for uncertain operations
4. **Deny dangerous operations**: Block `rm -rf`, `git push --force`, etc.
5. **External directory**: Always deny by default

### Tool Development

1. **Validate inputs**: Validate all arguments
2. **Request permissions**: Ask for dangerous operations
3. **Handle errors**: Graceful error handling
4. **Provide context**: Return useful metadata
5. **Be idempotent**: Tools should be safe to retry

### Plugin Development

1. **Modular design**: Split into focused plugins
2. **Clear naming**: Use descriptive tool names
3. **Document tools**: Provide clear descriptions
4. **Type safety**: Use TypeScript types
5. **Error handling**: Handle all error cases

## Testing Agents

### Local Testing

Test agents locally before deployment:

```bash
# Start opencode with agent
opencode --agent my-agent

# Or specify agent config
opencode --config .opencode/config.json

# Test agent capabilities
opencode --agent my-agent --interactive
```

### Testing Checklist

- [ ] Agent loads successfully
- [ ] Permissions work as expected
- [ ] Tools execute correctly
- [ ] Error handling works
- [ ] Integration with GitHub works
- [ ] Integration with opencode SDK works
- [ ] Agent follows guidelines

## Agent Deployment

### Registry Deployment

Deploy agent to registry:

1. **Create agent definition**: Add to `/docs/agents-examples/registry/`
2. **Validate**: Test in local environment
3. **Commit**: Commit to repository
4. **Deploy**: Agent is automatically discovered

### Project-Specific Deployment

Deploy agent to specific project:

1. **Create `.opencode/agent/` directory** in project
2. **Copy agent definition**: Add agent `.md` file
3. **Configure**: Adjust permissions and instructions as needed
4. **Test**: Test agent in project context

## Troubleshooting

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| Agent not discovered | Wrong file location | Ensure agent is in `registry/` or `.opencode/agent/` |
| Invalid frontmatter | YAML syntax error | Validate YAML syntax |
| Permissions too restrictive | Agent can't work | Loosen permission rules |
| Agent not following guidelines | Unclear instructions | Improve agent documentation |
| Tools not available | Plugin not loaded | Ensure plugin is installed and enabled |

## Next Steps

- [Examples](./agents-examples/) - See complete agent definitions
- [Registry](./agents-examples/registry/) - Production-ready agent templates

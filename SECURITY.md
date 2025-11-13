# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| main    | :white_check_mark: |
| < 1.0   | :x:                |

**Note**: CAMI is currently in pre-1.0 development. Security updates will be applied to the `main` branch.

## Reporting a Vulnerability

We take the security of CAMI seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Please Do Not

- Open a public GitHub issue for security vulnerabilities
- Disclose the vulnerability publicly before it has been addressed

### Please Do

**Report security vulnerabilities to**: [security@lando.com](mailto:security@lando.com)

Please include the following information in your report:

- Type of vulnerability (e.g., code injection, path traversal, etc.)
- Full paths of source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit or direct URL)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if available)
- Impact of the vulnerability and potential attack scenarios

### What to Expect

1. **Acknowledgment**: You will receive an acknowledgment of your report within 48 hours
2. **Assessment**: We will investigate and assess the vulnerability
3. **Updates**: We will keep you informed about the progress of fixing the issue
4. **Resolution**: Once fixed, we will notify you and publicly disclose the vulnerability

### Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Varies based on severity and complexity
  - Critical: Within 7 days
  - High: Within 14 days
  - Medium: Within 30 days
  - Low: Within 60 days

## Security Considerations

### Agent Sources

CAMI manages agent sources from Git repositories. Be aware that:

- **Trust your sources**: Only add agent sources from trusted repositories
- **Review agents**: Review agent content before deploying to projects
- **Git security**: Agent sources use Git, which could execute hooks
- **File permissions**: CAMI creates `.claude/agents/` directories with standard permissions

### Configuration

- **Config location**: CAMI stores configuration in `~/.cami/config.yaml`
- **No secrets**: Do not store sensitive information in agent files or config
- **File access**: CAMI requires read/write access to project directories

### MCP Server (Future)

When the MCP server is implemented:

- The server will have access to your filesystem where allowed
- The server will execute deployment operations on your behalf
- Only configure the MCP server to access directories you trust

## Best Practices

### For Users

1. **Review agents before deployment**
   ```bash
   cami list --detail
   ```

2. **Verify agent sources**
   ```bash
   cami source list
   ```

3. **Keep CAMI updated**
   ```bash
   git pull
   go build -o cami ./cmd/cami
   ```

4. **Check deployment targets**
   ```bash
   cami locations list
   ```

### For Contributors

1. **Validate user input**: Always validate and sanitize user input
2. **Path traversal**: Use `filepath.Clean()` and validate paths
3. **Command injection**: Never execute user input directly in shell commands
4. **File permissions**: Set appropriate file permissions (0644 for files, 0755 for directories)
5. **Dependencies**: Keep dependencies updated and review security advisories
6. **Error messages**: Don't expose sensitive information in error messages

### Security Checklist for Code Review

- [ ] User input is validated and sanitized
- [ ] File paths are cleaned and validated
- [ ] No shell command injection vulnerabilities
- [ ] File permissions are set appropriately
- [ ] Errors don't expose sensitive information
- [ ] Dependencies are from trusted sources
- [ ] No hardcoded credentials or secrets
- [ ] Proper error handling for edge cases

## Known Security Considerations

### Git Operations

CAMI executes `git` commands for:
- Cloning agent sources (`git clone`)
- Updating sources (`git pull`)
- Checking status (`git status`)

**Mitigation**:
- Only add trusted Git repositories as sources
- Review the source repository before adding
- CAMI does not execute arbitrary Git commands from user input

### File System Access

CAMI requires access to:
- `~/.cami/` for configuration
- `vc-agents/` for agent sources
- Project directories for deployment

**Mitigation**:
- CAMI only writes to `.claude/agents/` within target directories
- Users explicitly configure deployment locations
- All file operations use validated paths

### YAML Parsing

Agent files and configuration use YAML frontmatter.

**Mitigation**:
- YAML parser is from a trusted library (`gopkg.in/yaml.v3`)
- CAMI does not execute code from YAML files
- Agent content is treated as data, not code

## Security Updates

Security updates will be announced via:
- GitHub Security Advisories
- Release notes
- Repository README

## Acknowledgments

We appreciate the security research community's efforts in responsibly disclosing vulnerabilities. Security researchers who report valid vulnerabilities will be acknowledged in our release notes (unless they prefer to remain anonymous).

---

Thank you for helping keep CAMI and its users safe!

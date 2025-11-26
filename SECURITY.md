# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.4.x   | :white_check_mark: |
| < 0.4   | :x:                |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability in CAMI, please report it responsibly.

### How to Report

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please report security vulnerabilities by emailing:

**security@lando-labs.com**

Include the following information:
- Type of vulnerability
- Full path to the affected source file(s)
- Step-by-step instructions to reproduce
- Proof-of-concept or exploit code (if possible)
- Impact assessment

### What to Expect

1. **Acknowledgment**: We will acknowledge receipt of your report within 48 hours
2. **Assessment**: We will investigate and assess the vulnerability within 7 days
3. **Resolution**: We aim to release a fix within 30 days for critical vulnerabilities
4. **Disclosure**: We will coordinate with you on public disclosure timing

### Safe Harbor

We consider security research conducted in accordance with this policy to be:
- Authorized concerning any applicable anti-hacking laws
- Authorized concerning any relevant anti-circumvention laws
- Exempt from restrictions in our Terms of Service that would interfere with conducting security research

We will not pursue civil action or file a complaint with law enforcement for accidental, good-faith violations of this policy.

## Security Best Practices for Users

### Agent Sources

- Only add agent sources from trusted repositories
- Review agent code before deployment, especially from third-party sources
- Use `.camiignore` to exclude sensitive files from agent loading

### Configuration

- Keep your `config.yaml` private if it contains sensitive paths
- Use environment variables for sensitive configuration when possible
- Regularly update CAMI to get security patches

### MCP Server

- The MCP server runs locally and communicates via stdio
- No network ports are opened by default
- Be cautious when configuring global MCP access in `~/.claude/settings.json`

## Known Security Considerations

### File System Access

CAMI reads and writes files to:
- `~/cami-workspace/` (configuration and sources)
- Project `.claude/agents/` directories (deployment targets)

Ensure these directories have appropriate permissions.

### Git Operations

CAMI can clone repositories when adding sources. Only add sources from URLs you trust.

### No Remote Code Execution

CAMI does not execute code from agents - it only manages markdown files that Claude Code interprets. However, malicious agent instructions could potentially guide Claude to perform unintended actions.

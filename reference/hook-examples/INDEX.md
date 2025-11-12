# CAMI Hooks Documentation Index

Complete documentation for Claude Code hooks integration with CAMI.

## Documents Overview

### 1. **[QUICK_START.md](./QUICK_START.md)** ⚡ START HERE
**5-minute setup guide**
- Get agent scanning working immediately
- Minimal configuration
- Troubleshooting basics
- **Best for:** Getting started quickly

### 2. **[README.md](./README.md)** 📖
**Feature overview and usage**
- What each hook does
- Installation options
- Hook behavior examples
- Output samples
- Basic customization
- **Best for:** Understanding features

### 3. **[IMPLEMENTATION_GUIDE.md](./IMPLEMENTATION_GUIDE.md)** 🔧
**Complete implementation guide**
- Step-by-step instructions
- Plugin integration
- Customization recipes
- Advanced features
- Troubleshooting in-depth
- Best practices
- **Best for:** Building plugins or deep customization

### 4. **[HOOKS_FEATURE_RESEARCH.md](../HOOKS_FEATURE_RESEARCH.md)** 📚
**Complete technical reference**
- All 8 lifecycle events
- Hook configuration format
- JSON input/output specs
- Exit codes and control flow
- Security considerations
- Official documentation links
- **Best for:** Technical reference and understanding hooks in depth

## Configuration Files

### Basic Configuration
**[hooks.json](./hooks.json)**
- Simple 4-hook setup
- SessionStart, PostToolUse, SessionEnd
- Good for most use cases

### Complete Configuration
**[hooks-complete.json](./hooks-complete.json)**
- All lifecycle events
- PreToolUse validation included
- Production-ready

## Script Files

All scripts are executable and ready to use:

| Script | Purpose | Hook Event |
|--------|---------|------------|
| **cami-session-start.sh** | Scan & display agent status | SessionStart |
| **cami-pre-tool-validate.sh** | Validate before deployment | PreToolUse |
| **cami-post-tool.sh** | Auto-update CLAUDE.md | PostToolUse |
| **cami-session-end.sh** | Log session activity | SessionEnd |

## Quick Reference

### Installation Paths

**Global:** `~/.claude/settings.json`
```json
{
  "hooks": {
    "SessionStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "/Users/lando/Development/cami/reference/hook-examples/cami-session-start.sh"
          }
        ]
      }
    ]
  }
}
```

**Project:** `.claude/settings.json` (same format)

**Plugin:** Reference via `${CLAUDE_PLUGIN_ROOT}`

### Testing Commands

```bash
# Test SessionStart
echo '{"session_id":"test","cwd":"'$(pwd)'"}' | ./cami-session-start.sh

# Test PostToolUse
echo '{"tool_name":"mcp__cami__deploy_agents","cwd":"'$(pwd)'"}' | ./cami-post-tool.sh

# Test PreToolUse
echo '{"tool_name":"mcp__cami__deploy_agents","tool_input":{"agent_names":["architect"]}}' | ./cami-pre-tool-validate.sh

# Test SessionEnd
echo '{"session_id":"test","cwd":"'$(pwd)'"}' | ./cami-session-end.sh
```

## Reading Path

### For Users (Want to use CAMI hooks)
1. ⚡ [QUICK_START.md](./QUICK_START.md) - Get started (5 min)
2. 📖 [README.md](./README.md) - Understand features (15 min)
3. 🔧 [IMPLEMENTATION_GUIDE.md](./IMPLEMENTATION_GUIDE.md) - Customize (optional)

### For Developers (Building CAMI plugin)
1. 📚 [HOOKS_FEATURE_RESEARCH.md](../HOOKS_FEATURE_RESEARCH.md) - Learn hooks (30 min)
2. 🔧 [IMPLEMENTATION_GUIDE.md](./IMPLEMENTATION_GUIDE.md) - Implementation details (45 min)
3. 📖 [README.md](./README.md) - Feature specs (15 min)

### For Integrators (Adding hooks to existing tools)
1. 📚 [HOOKS_FEATURE_RESEARCH.md](../HOOKS_FEATURE_RESEARCH.md) - Technical specs
2. Study script files - See practical examples
3. 🔧 [IMPLEMENTATION_GUIDE.md](./IMPLEMENTATION_GUIDE.md) - Integration patterns

## File Sizes

| File | Lines | Size | Purpose |
|------|-------|------|---------|
| HOOKS_FEATURE_RESEARCH.md | 1192 | 58KB | Complete research |
| IMPLEMENTATION_GUIDE.md | 568 | 26KB | Implementation guide |
| README.md | 205 | 9KB | Feature overview |
| QUICK_START.md | 179 | 5KB | Quick setup |
| hooks-complete.json | 47 | 1.3KB | Full config |
| hooks.json | 31 | 864B | Basic config |
| cami-session-start.sh | 97 | 2.8KB | Session scanner |
| cami-pre-tool-validate.sh | 75 | 2.1KB | Validator |
| cami-post-tool.sh | 44 | 1.0KB | Doc updater |
| cami-session-end.sh | 29 | 858B | Logger |

## Features by Hook

### SessionStart Hook
✅ Scan deployed agents
✅ Show version status
✅ Update notifications
✅ Context injection
- 97 lines, robust error handling

### PreToolUse Hook
✅ Validate agent names
✅ Check VC repository
✅ Block invalid deployments
✅ Interactive approval
- 75 lines, JSON control flow

### PostToolUse Hook
✅ Auto-update CLAUDE.md
✅ Filter MCP calls
✅ Call CAMI CLI
✅ Success notifications
- 44 lines, lightweight

### SessionEnd Hook
✅ Log session info
✅ Track activity
✅ Cleanup old logs
✅ Audit trail
- 29 lines, minimal

## Dependencies

All scripts require:
- **bash** (built-in on macOS/Linux)
- **jq** (install: `brew install jq`)

Optional:
- **CAMI CLI** (for auto-update feature)
- **git** (for git integration)

## Compatibility

- ✅ macOS (tested)
- ✅ Linux (should work)
- ❓ Windows (needs testing with WSL/Git Bash)

## Security Notes

- All scripts validate input from stdin
- No external network calls
- File operations limited to project directory
- Logs stored in user home directory
- No credentials stored

## Version History

- **v1.0.0** (2025-10-09) - Initial release
  - SessionStart agent scanner
  - PostToolUse doc updater
  - PreToolUse validator
  - SessionEnd logger
  - Complete documentation

## Next Steps

1. **Try it:** Follow [QUICK_START.md](./QUICK_START.md)
2. **Customize:** See [IMPLEMENTATION_GUIDE.md](./IMPLEMENTATION_GUIDE.md)
3. **Build plugin:** Package as Claude Code plugin
4. **Share:** Contribute improvements back

## Support Resources

### Official Documentation
- [Claude Code Hooks Guide](https://docs.claude.com/en/docs/claude-code/hooks-guide)
- [Claude Code Hooks Reference](https://docs.claude.com/en/docs/claude-code/hooks)
- [Claude Code Plugins](https://docs.claude.com/en/docs/claude-code/plugins-reference)

### Community Examples
- [claude-code-hooks-mastery](https://github.com/disler/claude-code-hooks-mastery)
- [claude-git](https://github.com/listfold/claude-git)
- [GitButler Blog](https://blog.gitbutler.com/automate-your-ai-workflows-with-claude-code-hooks)

### CAMI Resources
- **Main repo:** github.com/yourusername/cami
- **MCP server:** Integrated with hooks
- **CLI tool:** Used by post-tool hook

## Contributing

Found a bug? Have an improvement?

1. Test your change
2. Update relevant docs
3. Submit PR with examples
4. Update this index if adding files

## License

MIT License - Same as CAMI project

---

**Last Updated:** October 9, 2025
**Research by:** Claude Code (Sonnet 4.5)
**Status:** Production Ready ✅

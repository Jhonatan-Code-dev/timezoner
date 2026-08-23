---
name: clean-doc-architect
description: >-
  Expert Technical Documentation & README Architect. Crafts clean, executive,
  and authoritative software documentation without generic AI emojis, buzzwords, or superficial decoration.
---

# Clean Technical Documentation & README Architect

This skill guides the creation of industry-standard, professional READMEs and technical documentation. It discards AI clichés (excessive emojis, hype words, generic decorative bullet points) in favor of high-signal, clean, and authoritative technical writing used by companies like Stripe, HashiCorp, Cloudflare, and the Go team.

---

## 🖋️ Core Tenets of Clean Technical Writing

### 1. No Superficial Emoji Pollution
- Eliminate decorative emojis (🚀, ✨, 🌍, ⚡, 📦, 🔥, 💡) in titles, headers, and bullet points.
- Use standard typography, clear Markdown hierarchies (`#`, `##`, `###`), and code formatting to establish visual structure.

### 2. High Signal-to-Noise Ratio
- Focus on practical, copy-pasteable code examples, architecture tables, benchmark figures, and exact error models.
- Avoid marketing filler ("blazing fast", "state of the art", "revolutionary"). Let benchmarks and type signatures speak for themselves.

### 3. Standard Technical README Structure
A professional library README contains:
1. **Project Title & Direct Description**: What the library does in one or two technical sentences.
2. **Key Capabilities**: Bullet points detailing concrete functionality (e.g., IANA timezone conversion, memory-cached lookups, DST detection).
3. **Installation**: Clean package manager command (`go get ...`).
4. **Quick Start / Usage**: Concise, idiomatic code snippets with comments.
5. **API Reference / Package Layout**: Table or overview of exported functions and structs.
6. **Performance & Benchmarks**: Structured table with operations/sec, ns/op, and bytes/op.
7. **License & Copyright**: Formal notice.

### 4. Code Sample Excellence
- Include complete, runnable examples with realistic variables.
- Show both the input and the expected output in clean code comments.
- Demonstrate idiomatic error handling (`if err != nil`) rather than hiding it.

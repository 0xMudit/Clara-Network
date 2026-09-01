# README Additions — Draft for Owner Review

These are proposed additions to `README.md`. **Review, customize, and merge
into the README at your discretion.**

---

## Badges (add after the `# Clara Network` heading)

```markdown
[![CI](https://github.com/0xMudit/Clara-Network/actions/workflows/ci.yml/badge.svg)](https://github.com/0xMudit/Clara-Network/actions/workflows/ci.yml)
[![CodeQL](https://github.com/0xMudit/Clara-Network/actions/workflows/codeql.yml/badge.svg)](https://github.com/0xMudit/Clara-Network/actions/workflows/codeql.yml)
[![Go Version](https://img.shields.io/badge/go-1.26+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](./LICENSE)
[![Open Issues](https://img.shields.io/github/issues/0xMudit/Clara-Network)](https://github.com/0xMudit/Clara-Network/issues)
```

---

## Contributing section (add before the `## License` section)

```markdown
## Contributing

Contributions are welcome! Whether you're fixing a bug, adding a feature,
improving documentation, or helping with translations — we'd love your help.

- **New to the project?** Start with the [Contributing Guide](CONTRIBUTING.md)
  and the [Contributor Architecture Guide](docs/28-contributor-architecture-guide.md).
- **Looking for something to work on?** Check
  [good first issues](https://github.com/0xMudit/Clara-Network/labels/good%20first%20issue)
  and the [Roadmap](ROADMAP.md).
- **Have a question?** Open a
  [Discussion](https://github.com/0xMudit/Clara-Network/discussions).
- **Found a security issue?** See [SECURITY.md](SECURITY.md) for responsible
  disclosure instructions.

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before participating.
```

---

## "Fork and build your own" callout (add after the `## Building & running` section)

```markdown
### Fork and build your own payment system

Clara Network is designed so that **county and bank engineers can fork it
and build their own payment system** with a few tweaks. The main things to
customize:

1. **BIN ranges** — configure `CLARA_BIN_TABLE` with your country's BIN
   prefixes and member institution IDs.
2. **Currency** — set your local currency in the schema and simulator configs.
3. **Member institutions** — map your banks' sort codes to `CLARA_ISSUER_ROUTES`.
4. **Settlement rules** — configure prefund accounts and the default fund.

See the [Contributing Guide](CONTRIBUTING.md#how-to-customize-for-your-own-payment-system)
for a step-by-step walkthrough.
```

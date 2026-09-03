# Owner Action Items — Manual Steps Required

These items require **repo owner permissions** (branch protection, GitHub
settings, account setup). They cannot be done via file changes alone.

---

## Branch Protection Rules (apply in GitHub Settings → Branches → main)

- [ ] **Require pull request reviews** before merging (1 review minimum;
      2 for `main` once the team grows)
- [ ] **Require status checks** before merging:
  - [ ] `lint` (golangci-lint)
  - [ ] `build-test` (go vet + go build + go test -race)
  - [ ] `docker-smoke` (Docker Compose stack health check)
- [ ] **Require branches to be up to date** before merging
- [ ] **Restrict force pushes** to `main`
- [ ] **Restrict deletions** of `main`
- [ ] **Require conversation resolution** before merging

---

## Codecov Setup

- [ ] Create a [Codecov](https://codecov.io) account and link it to the repo
- [ ] Add the `CODECOV_TOKEN` as a repository secret in Settings → Secrets →
      Actions
- [ ] Verify the coverage upload works by checking the Codecov dashboard after
      the next CI run

---

## GitHub Discussions

- [ ] Enable [Discussions](https://docs.github.com/en/discussions) in
      Settings → General → Discussions
- [ ] Create initial categories:
  - **General** — open-ended questions and ideas
  - **Architecture** — design decisions and trade-offs
  - **Show and Tell** — forks, deployments, demos
  - **Q&A** — get help from the community

---

## GitHub Security Advisories

- [ ] Enable [Security Advisories](https://docs.github.com/en/code-security/security-advisories)
      in Settings → Code security and analysis
- [ ] Verify the SECURITY.md reporting instructions work (test by visiting
      the "Security" tab on the repo)

---

## Issue Labels

Create these labels if they don't already exist:

- [ ] `bug` — Something isn't working
- [ ] `enhancement` — New feature or request
- [ ] `good first issue` — Good for newcomers
- [ ] `help wanted` — Extra attention is needed
- [ ] `documentation` — Improvements or additions to documentation
- [ ] `question` — Further information is requested
- [ ] `wontfix` — This will not be worked on
- [ ] `dependencies` — Pull requests that update a dependency (used by Dependabot)
- [ ] `ci` — CI/infrastructure changes (used by Dependabot)
- [ ] `phase/11-fraud` — Fraud detection & ML
- [ ] `phase/12-3ds` — 3-D Secure
- [ ] `phase/13-fx` — Cross-border & FX
- [ ] `phase/14-cards` — Card production
- [ ] `phase/15-regulatory` — Regulatory jurisdictions
- [ ] `phase/16-infra` — Kubernetes & deployment
- [ ] `phase/17-dashboard` — Web dashboard

---

## Starter Issues

- [ ] Create the 15 starter issues from `.github/STARTER_ISSUES.md`
- [ ] Tag them `good first issue`
- [ ] Delete `.github/STARTER_ISSUES.md` after creating the issues

---

## Community Health Files

Verify these files exist and are visible on the repo's community profile
(Settings → Community):

- [ ] `CODE_OF_CONDUCT.md`
- [ ] `CONTRIBUTING.md`
- [ ] `SECURITY.md`
- [ ] `LICENSE` (MIT)
- [ ] `.github/ISSUE_TEMPLATE/bug_report.md`
- [ ] `.github/ISSUE_TEMPLATE/feature_request.md`
- [ ] `.github/PULL_REQUEST_TEMPLATE.md`

---

## Optional (when ready)

- [ ] **Discord/Slack** — Create a community chat channel and link it in the
      README
- [ ] **CLA/DCO** — Decide whether to require a Contributor License Agreement
      or Developer Certificate of Origin for contributions
- [ ] **"Show HN" post** — Draft and publish to Hacker News (see Phase 4
      visibility assets)
- [ ] **Technical blog post** — Write about a distinctive component (HSM key
      ceremonies, stand-in chaos drill, etc.)
- [ ] **Code coverage badge** — Add to README after Codecov is set up
- [ ] **Release workflow** — Add a GitHub Actions workflow for tagging and
      publishing releases (`goreleaser` or manual)

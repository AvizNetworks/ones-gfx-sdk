# Contributing to ONES GFX SDK

Thanks for your interest in contributing! This repository hosts two SDKs — Python and Go — for the AVIZ ONES Spectrum-X tenant management API. This guide covers the parts common to both. Language-specific conventions live alongside each SDK:

- [ones-gfx-sdk-python/README.md](./ones-gfx-sdk-python/README.md)
- [ones-gfx-sdk-go/README.md](./ones-gfx-sdk-go/README.md)

## Reporting Bugs

Please use the [Bug Report template](./.github/ISSUE_TEMPLATE/bug_report.md) and include:
- Steps to reproduce
- Expected vs. actual behavior
- Environment details (OS, version, relevant config)

For **security vulnerabilities**, do not open a public issue — email security@aviznetworks.com instead.

## Suggesting Features

Please use the [Feature Request template](./.github/ISSUE_TEMPLATE/feature_request.md)
and describe the problem you're trying to solve, not just the solution — it helps
us evaluate alternatives.

## Development Workflow 

### Python (`ones-gfx-sdk-python/`)
```bash
cd ones-gfx-sdk-python
make install-dev   # editable install + pytest/ruff
make check         # sanity import check
```

### Go (`ones-gfx-sdk-go/`)
```bash
cd ones-gfx-sdk-go
make build-lib      # build the core library
make check          # sanity import check
```

## Changes Addition and Testing

1. Fork the repository and create a branch off `master`:
   ```bash
   git checkout -b issue-no-short-description
   ```
2. Keep changes scoped to one SDK where possible. If a change must touch both (e.g. an API contract update), say so in the PR description and keep the two implementations at parity.
3. Match the existing style of the file you're editing rather than introducing a new convention.
4. Update the relevant README/API reference docs when behavior, endpoints, or public interfaces change.
5. Add or update tests for the bug/feature you develop and Verify the changes.
6. Refer ./ones-gfx-sdk-python/README.md | ./ones-gfx-sdk-go/README.md for more information on development/testing

Both SDKs wrap the same ONES Spectrum-X API — if you change request/response handling in one language, check whether the other needs the equivalent fix (see the Feature Comparison table in the [root README](./README.md)).

## What the transport owns (read before adding a resource method)

Several request details are set centrally by the transport, not by the method
you are writing. Setting them again in a resource is at best redundant and at
worst a conflict, because the transport writes last.

| Set by the transport | Where | Driven by |
|---|---|---|
| `Prefer: respond-sync` / `respond-async` | `transport.go` — `send` and `sendRaw` | the `OperationMode` argument |
| `Authorization` | `auth.Apply` on every request | the client's auth provider |
| `Content-Type`, `Accept` | `sendRaw` | the call, with `ReqOpts.Accept` overriding |
| `Idempotency-Key` | `sendRaw` | `ReqOpts.IdempotencyKey` |
| `x-request-origin` | `sendRaw` | `ReqOpts.RequestOrigin` |
| `/api/fm/` prefix, or direct-FM routing | `buildURL` | the path, via `fmdirect.go` |
| Query string | `doRaw` | `ReqOpts.Query` |
| `{"data": ...}` unwrapping | `Call` | always |
| Auth refresh and one retry on 401 | `doRaw` | always |

### Sync vs async is an argument, not a header

A resource method chooses behaviour by passing an `OperationMode`:

```go
ones_gfx.Call[ApiResponseMessage](r.transport, "PATCH", path, body,
    ones_gfx.OperationModeSynchronous,   // transport turns this into Prefer: respond-sync
    rc.ReqOpts())
```

Do **not** write `headers["Prefer"] = ...` in a resource. If a method needs a
mode the enum does not cover, add it to `OperationMode` and teach the transport,
so every call site stays consistent.

The paired `...Async` methods differ only in this argument — compare
`TenantsResource.Update` with `UpdateAsync`.

### Checklist for a new resource method

1. Give it a `// Maps to VERB /path` comment. `fmdirect.go` and the coverage
   checks both read these.
2. Pass an `OperationMode`; never set `Prefer` yourself.
3. Pass `rc.ReqOpts()` so idempotency, request origin, timeout and query
   options keep working.
4. Return a typed model rather than `map[string]interface{}`. Keep unmodelled
   keys in an `Extra` field if the payload varies by fabric type.
5. If the route is one the ONES gateway does not proxy, add it to
   `DirectFMPaths` in `fmdirect.go`.

## Commit Messages

Write commit messages that explain *why* a change was made, not just what changed. Reference the related issue number when one exists (e.g. `Fixes #42`).

## Submitting a Pull Request

1. Ensure `make check` (or the language equivalent) passes.
2. Push your branch and open a PR against `master`. GitHub will pre-fill the [PR template](./.github/PULL_REQUEST_TEMPLATE.md) — fill it in rather than deleting it.
3. Describe what changed, why, and how you tested it. Note any impact on the other language SDK.
4. One of the maintainers will review and may request changes.

## Pull Request Process

1. Fill out the PR template completely — incomplete PRs may be asked for more info before review.
2. Link the issue your PR resolves (e.g. `Closes #123`).
3. Keep PRs small and focused; large PRs take longer to review and are more likely
  to be asked to split.
4. A maintainer will review, request changes if needed, and merge once approved.

## Versioning

This repository follows [Semantic Versioning (SemVer)](https://semver.org/) and maintains its version history independently of other platform releases. 
See the [Versioning Policy](./README.md#versioning-policy) in the README for details.

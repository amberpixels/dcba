# DCBA

A small, enforceable layering convention for Go applications. This repo is the
**canonical spec** — reference it from a project's `internal/doc.go` instead of
copying the whole explanation into every codebase.

DCBA splits an application's `internal/` tree into four sibling layers and fixes
the **direction** dependencies may point. It's the same idea as Hexagonal /
Clean / Onion architecture, boiled down to four directories and one rule you can
enforce with a linter.

## The four layers

| | Layer | What lives here |
|---|---|---|
| **D** | `dtos/` | wire shapes — request/response DTOs; maps domain ⇄ JSON (or protobuf, etc.) |
| **C** | `components/` | infrastructure — db, auth, external API clients, caches, storage |
| **B** | `business/` | the domain — entities and domain rules. The dependency **sink**. |
| **A** | `application/` | composition root + transport (handlers/routing) + use cases (services) |

### Why "DCBA" (reverse-alphabetical)?

The letters are deliberately **not** in top-to-bottom order. Alphabetical "ABCD"
reads like a stack where A is the foundation — but here **A is the outermost
layer, not the base**. Naming them D-C-B-A stops that wrong mental picture. The
label is just a mnemonic for the four directories; the real contract is the
dependency rule.

## The dependency rule

Dependencies point **inward, toward the domain**. The only allowed directions:

```
A ──▶ C, D, B        application may import anything
C ──▶ B              components may import the domain only
D ──▶ B              dtos may import the domain only
B ──▶ (nothing)      business imports no sibling layer*
```

`*` business/* packages may import **other** business/* packages (e.g. `raid`
references `activist`). That intra-domain graph must stay **acyclic**.

### Consequences worth remembering

- **Business never imports dtos, components, or application.** If a domain type
  needs a wire shape, the DTO depends on the domain — not the reverse.
- **Components is infrastructure only:** it knows domain types, but not wire
  shapes (no `dtos`) and not the application layer.
- **Application is the composition root.** It's the only layer allowed to depend
  on everything, so it's where wiring and transport live. A common pattern:
  service **interfaces** live in `application/appdeps` (a neutral seam) because
  both the transport handlers and other services consume them.

## Adopting it in a project

1. **Structure `internal/`** into `dtos/`, `components/`, `business/`,
   `application/`.

2. **Drop in a short package doc** so the layout is documented at the source.
   Copy [`templates/doc.go.tmpl`](templates/doc.go.tmpl) to `internal/doc.go` and
   fill in the project's specifics (it links back here for the full spec).

3. **Enforce it mechanically.** Merge the depguard rules from
   [`templates/golangci-dcba.yml`](templates/golangci-dcba.yml) into your
   `.golangci.yml` and replace `MODULE` with your module path. An inverted import
   then fails `golangci-lint` instead of silently rotting the architecture.

The rule is enforced by **depguard** today (config, copied per project). See the
[roadmap](#roadmap) for the plan to replace the copy-paste with a dedicated
linter.

## Roadmap

- **`dcbacheck` — a custom `go/analysis` linter.** A single, configurable
  analyzer (usable as a golangci-lint plugin/module) that enforces the layering
  from a tiny per-project config (which package = which layer), so projects stop
  copying the depguard block. Tracked as an issue in this repo.

## License

MIT — see [LICENSE](LICENSE).

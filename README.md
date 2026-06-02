# Family Tree
A small web app to view and edit a family tree. It stores people and parent
relations in SQLite and renders the tree in HTML with HTMX.

# Tech stack
- go (gin)
- sqlite
- htmx

# Run locally
Prereqs: Go 1.24+, sqlite3 CLI (optional).

If you want a fresh database:

```bash
rm -f database.db
sqlite3 database.db < schema.sql
```

Start server:

```bash
go run ./main.go
```

Open http://localhost:8181

# Release binary

Push tag like `v1.0.0` to trigger GitHub Actions release workflow.

```bash
git tag v1.0.0
git push origin v1.0.0
```

Workflow builds `family-tree-linux-amd64`, creates checksums, uploads both to GitHub Release.

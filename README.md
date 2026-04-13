# Family Tree
build your family tree

# Tech stack
- grpc
- sqllite
- htmx

# Release binary

Push tag like `v1.0.0` to trigger GitHub Actions release workflow.

```bash
git tag v1.0.0
git push origin v1.0.0
```

Workflow builds `family-tree-linux-amd64`, creates checksums, uploads both to GitHub Release.

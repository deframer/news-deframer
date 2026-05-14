# Plan

## Goal

Fail fast when admin feed categories are invalid, using the English category list as the source of truth.

## Steps

- [ ] Move `category.go` into a shared `pkg/category` package.
- [ ] Keep validation in the existing admin commands, not a new tool.
- [ ] Update the admin feed commands in `cmd/admin/feed.go` to validate categories.
- [ ] Update the admin import path in `cmd/admin/importexport.go` to validate categories.
- [ ] Reuse the existing think category validator with language `en`.
- [ ] Return clear errors for unknown categories and stop the command.
- [ ] Run `make lint` and relevant tests after the change.

## Expected Result

- Admin rejects unknown feed categories immediately.
- Feed imports do not accept invalid category names.
- English categories remain the canonical admin format.

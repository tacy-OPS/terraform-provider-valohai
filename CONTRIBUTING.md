# Contributing to terraform-provider-valohai

Thank you for your interest in contributing to the Valohai Terraform provider!

## How to Contribute

1. **Fork the repository** and clone your fork locally.
2. **Create a new branch** for your feature or bugfix:
   ```sh
   git checkout -b my-feature
   ```
3. **Write clear, concise code** and include tests when possible.
4. **Run the linter and tests** before submitting:
   ```sh
   make clean build check-binary
   make test
   ```
5. **Document your changes** in the relevant Markdown files (README, docs/).
6. **Open a Pull Request** with a clear description of your changes and the motivation.

## Development Environment

- Go 1.21+
- Terraform 1.0+
- Valohai account and API token

## Local Provider Build

- Use `make install-local` to build and install the provider locally.
- Update your `~/.terraformrc` or `%APPDATA%/terraform.rc` to use the local provider (see README).

## Code Style

- Follow Go best practices and idioms.
- Use `gofmt` and `goimports`.
- Keep code and documentation in English.

## Testing

- Unit tests: `go test -v ./tests/resource_projects_unit_test.go`
- Acceptance tests (require a real Valohai token):
  ```sh
  export TF_ACC=1
  export VALOHAI_API_TOKEN=your_token
  export VALOHAI_OWNER=your_org
  go test -v ./tests/...
  ```

## Issues & Feature Requests

- Please use GitHub Issues to report bugs or request features.
- Provide as much detail as possible (logs, Terraform config, provider version, etc.).

## Code of Conduct

Be respectful and constructive. We welcome all contributions and contributors.

---

Thank you for helping improve the Valohai Terraform provider!

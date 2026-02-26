# terraformify

`terraformify` is an experimental CLI that generates Terraform configuration files for managing **existing Fastly services** (VCL and Compute).

The tool works by:
- scaffolding minimal Terraform configuration,
- importing existing Fastly resources into state, and
- rewriting that state into usable `.tf` files, including variables for sensitive values.

---

## Terraform Compatibility

`terraformify` supports **modern Terraform 1.x versions**, including **Terraform 1.4.6 and later**.

Earlier versions of this project were incompatible with newer Terraform releases due to changes in how Terraform renders state and redacts sensitive values. These limitations have been addressed.

### Nested Block & Sensitive Attribute Handling

Newer versions of Terraform may fully redact **nested blocks** if any attribute inside the block is marked as sensitive. This often results in output like:

```hcl
backend {
  # At least one attribute in this block is (or was) sensitive,
  # so its contents will not be displayed.
}
```

This behavior commonly affects:
- `backend {}` blocks
- logging blocks (e.g. `logging_s3`, `logging_bigquery`)
- other nested resource blocks containing credentials or secrets

**Current behavior in terraformify:**
- Nested blocks are preserved structurally
- Sensitive fields are externalized into variables
- Generated `.tf` files no longer contain broken or empty nested blocks

You may still need to populate sensitive values in `terraform.tfvars` or your preferred secrets workflow before running `terraform plan`.

#### How terraformify reconstructs redacted nested blocks

Terraform may redact entire nested blocks in `terraform show` when any child attribute is (or was) marked sensitive, which can produce incomplete HCL output.

To avoid generating broken configuration from redacted output, `terraformify` reconstructs nested blocks by reading the underlying resource structure from Terraform state and re-emitting those attributes back into HCL. When sensitive attributes are encountered, `terraformify` preserves the block structure but replaces those values with variable references and emits the corresponding variables into `variables.tf` and `terraform.tfvars` (or placeholders), so the generated configuration remains valid and reviewable.

---

## Installation / Upgrade

### Go install

```bash
go install github.com/hrmsk66/terraformify@latest
```

### Prebuilt binaries

Prebuilt binaries are not currently published for this version. Please build from source using `go install`.

---

## Configuration

`terraformify` requires read access to your Fastly account.

Provide your Fastly API token using one of the following methods:

- Pass the token on each command with `--api-key` or `-k`
- Set the environment variable:

```bash
export FASTLY_API_KEY="your-token-here"
```

---

## Usage

Run `terraformify` in an **empty directory**, or in an existing Terraform directory **only if you are comfortable with state and file modifications**.

> **Important**
> Running `terraformify` in a directory with existing Terraform configuration will:
> - modify `terraform.tfstate`
> - potentially rewrite `variables.tf` and `terraform.tfvars`
>
> Back up your files before importing a service.

---

### Importing a VCL Service

```bash
terraformify service vcl <service-id>
```

---

### Importing a Compute Service

For Compute services, provide the service ID and the path to the WASM package:

```bash
terraformify service compute <service-id> -p <path-to-package>
```

---

### Documentation

For full command usage, flags, and advanced options, see:

```
docs/USAGE.md
```

---

## Supported Resources

`terraformify` supports importing both **VCL** and **Compute** services and their associated Fastly resources.

### VCL Resources

- fastly_service_vcl
- fastly_service_acl_entries
- fastly_service_dictionary_items
- fastly_service_dynamic_snippet_content
- fastly_service_waf_configuration

### Compute Resources

- fastly_service_compute
- fastly_configstore
- fastly_configstore_entries
- fastly_secretstore
- fastly_kvstore
- fastly_service_dictionary_items

---

## Troubleshooting

### Empty or Redacted Nested Blocks

If you see nested blocks commented as sensitive in generated output, this is Terraform redaction behavior—not a Fastly or terraformify error.

Ensure that:
- sensitive variables are populated in `terraform.tfvars`, or
- values are sourced from your secrets manager

After doing so, `terraform plan` should converge cleanly.

---

## License

MIT License

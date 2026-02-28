# Security Policy

## Reporting a Vulnerability
If you discover a security issue, please open a private report if possible or create an issue with minimal details and mark it as security-related.

## LLM External Transmission Policy

- Optional AI workflows are disabled by default and run only when `AI_ENABLED=true`.
- The PR summary workflow sends the PR title and up to the first 10,000 bytes of `git diff origin/<base>...HEAD` to an external LLM API.
- Lockfiles are excluded in that workflow (`*.lock`, `package-lock.json`), but all other changed text may be transmitted.
- Do not enable this workflow for repositories or periods that may include secrets, credentials, personal data, customer data, confidential contract information, or undisclosed vulnerability details.
- `LLM_API_KEY` must be stored only in GitHub Secrets and must never be committed.
- Repository maintainers are responsible for data classification, provider terms review, and disabling AI workflows when data handling requirements are not met.

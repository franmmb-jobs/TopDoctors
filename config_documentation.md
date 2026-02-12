# Configuration and Secrets Documentation

This document outlines the configuration variables and secrets required to run the TopDoctors Diagnostics API.

## Configuration Variables

These variables can be set via environment variables.

| Variable | Description | Default Value |
| :--- | :--- | :--- |
| `PORT` | The port on which the HTTP server will listen. | `8080` |

## Secrets

These are sensitive values that should be kept secure.

| Variable | Description | Default Value (Dev) |
| :--- | :--- | :--- |
| `JWT_SECRET` | Secret key used for signing and verifying JWT tokens. **Change this in production!** | `secret` |




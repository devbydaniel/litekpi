# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in LiteKPI, please report it responsibly.

**Do not open a public GitHub issue for security vulnerabilities.**

Instead, please send an email to the project maintainers with:

1. A description of the vulnerability
2. Steps to reproduce the issue
3. Potential impact of the vulnerability
4. Any suggested fixes (optional)

## Response Timeline

- We will acknowledge receipt within 48 hours
- We will provide an initial assessment within 7 days
- We aim to release a fix within 30 days for critical vulnerabilities

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |

## Security Best Practices for Self-Hosters

When deploying LiteKPI, please ensure:

1. **Strong secrets**: Use long, random values for `JWT_SECRET` and `POSTGRES_PASSWORD`
2. **HTTPS**: Always use TLS in production (via reverse proxy)
3. **Firewall**: Only expose necessary ports (80/443 for web, keep 5432/8080 internal)
4. **Updates**: Keep LiteKPI and its dependencies up to date
5. **Backups**: Regularly back up your PostgreSQL database
6. **Network isolation**: Run services on an internal Docker network

## Security Features

LiteKPI includes the following security measures:

- Passwords hashed with Argon2id
- API keys hashed with bcrypt (never stored in plain text)
- JWT tokens in httpOnly cookies
- CORS protection
- Input validation on all endpoints
- SQL injection protection via parameterized queries

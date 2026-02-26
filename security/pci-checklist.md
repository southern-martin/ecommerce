# PCI DSS Compliance Checklist for Payment Service

## Requirements Met:
- [ ] No credit card data stored (delegated to payment gateway)
- [ ] All API communication over TLS 1.2+
- [ ] JWT tokens expire (access: 15min, refresh: 7 days)
- [ ] Password hashing (bcrypt, cost 12)
- [ ] Rate limiting on auth endpoints
- [ ] SQL injection prevention (parameterized queries via GORM)
- [ ] Input validation on all endpoints
- [ ] Audit logging for payment operations
- [ ] Separate database for payment service
- [ ] Environment variables for secrets (no hardcoded credentials)

## Recommendations:
- Use AWS KMS / GCP KMS for key management
- Enable database encryption at rest
- Implement IP allowlisting for admin endpoints
- Set up alerting for failed payment attempts
- Regular dependency vulnerability scanning
- Penetration testing before launch

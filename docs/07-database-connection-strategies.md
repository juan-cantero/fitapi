# Database Connection Strategies

## Overview

Supabase offers two connection methods, each optimized for different use cases.

## Connection Types

### 1. Direct Connection (Port 5432)

**Connection String:**
```
postgresql://postgres.PROJECT:PASSWORD@HOST.supabase.com:5432/postgres
```

**Characteristics:**
- Direct connection to PostgreSQL
- Full PostgreSQL feature support
- No statement caching issues
- Limited to ~60 concurrent connections (Supabase free tier)

**Best For:**
- ✅ Local development
- ✅ Single long-running server (VPS, EC2, DigitalOcean)
- ✅ 1-5 application instances
- ✅ Traditional server deployments

**Capacity:**
- Handles 500-2000 concurrent users (with proper connection pooling)
- Uses 10-20 connections per server instance

### 2. Connection Pooler (Port 6543)

**Connection String:**
```
postgresql://postgres.PROJECT:PASSWORD@HOST.pooler.supabase.com:6543/postgres
```

**Characteristics:**
- Connection pooler (PgBouncer-like)
- Multiplexes connections
- Transaction mode by default
- Unlimited client connections

**Best For:**
- ✅ Serverless deployments (Vercel, AWS Lambda, Cloud Functions)
- ✅ Auto-scaling applications
- ✅ Many application instances (6+)
- ✅ High-traffic applications

**Capacity:**
- Handles 10,000+ concurrent users
- Efficiently reuses underlying PostgreSQL connections

## Comparison Table

| Deployment Type | Recommended | Port | Max Users | Notes |
|----------------|-------------|------|-----------|-------|
| **Local Development** | Direct | 5432 | N/A | Simpler, no caching issues |
| **Single VPS/Server** | Direct | 5432 | 500-2000 | Won't hit connection limit |
| **2-5 Servers** | Direct | 5432 | 1000-5000 | Still under 60 total connections |
| **6+ Servers** | Pooler | 6543 | 10,000+ | Would exceed 60 connections |
| **Serverless (Vercel, Lambda)** | Pooler | 6543 | 100,000+ | Essential - each function creates connections |
| **Auto-scaling (Kubernetes)** | Pooler | 6543 | 50,000+ | Dynamic instance count |

## Connection Pool Configuration

### Current Configuration (database.go)

```go
config.MaxConns = 10              // Max connections per instance
config.MinConns = 2               // Warm connections
config.MaxConnLifetime = 0        // No max (let pooler handle)
config.MaxConnIdleTime = 30min    // Close idle after 30min
config.HealthCheckPeriod = 1min   // Check health every minute
```

### Why These Numbers?

**MaxConns = 10:**
- Prevents flooding the pooler
- Single instance uses ~10 connections
- 6 instances × 10 = 60 total (at Supabase limit)

**MinConns = 2:**
- Keeps warm connections ready
- Reduces latency for first requests

**MaxConnIdleTime = 30min:**
- Releases unused connections
- Pooler has 15s idle timeout, we respect that

## Supabase Limits

### Free Tier
- **Direct connections:** ~60 concurrent
- **Pooler connections:** Unlimited (client-side)
- **Database storage:** 500MB

### Pro Tier
- **Direct connections:** ~200 concurrent
- **Pooler connections:** Unlimited
- **Database storage:** 8GB+

## Self-Hosted PostgreSQL

If you host your own PostgreSQL:

**Default Limits:**
- `max_connections = 100` (can increase to 200-500)
- Each connection uses ~10MB RAM

**Example with 8GB RAM:**
- Direct: ~200-300 concurrent connections
- With PgBouncer pooler: 10,000+ concurrent users

## Troubleshooting

### "Too many connections" Error

**Cause:** Exceeded PostgreSQL connection limit

**Solutions:**
1. Switch to Connection Pooler (port 6543)
2. Reduce MaxConns in your pool config
3. Upgrade Supabase tier
4. Ensure connections are properly closed

### "Prepared statement already exists" Error

**Cause:** Connection Pooler caching bug with pgx

**Solutions:**
1. Use Direct Connection (port 5432) - recommended for traditional servers
2. Wait for cache to expire (~1 hour)
3. Use transaction mode pooler (default in production)

## Best Practices

### 1. Use Correct Pooling Mode
- Transaction mode (default) for REST APIs
- Session mode only if you need persistent session state

### 2. Tune Pool Sizes
- Align client MaxConns with Supabase limits
- Don't over-provision (more isn't always better)

### 3. Release Connections Properly
```go
defer rows.Close()  // Always close rows
```

### 4. Avoid Long-Lived Connections
- Use context timeouts
- Set reasonable MaxConnIdleTime

### 5. Monitor Pool Metrics
- Track pool saturation
- Alert on connection acquisition timeouts
- Monitor query latency

### 6. Test Under Load
- Use staging environment
- Simulate production traffic
- Test failover scenarios

## Recommendations by Use Case

### Fitness API (This Project)

**Development:**
```bash
DATABASE_URL=postgresql://....:5432/postgres  # Direct connection
```

**Production (Single Server):**
```bash
DATABASE_URL=postgresql://....:5432/postgres  # Direct connection
```

**Production (Serverless):**
```bash
DATABASE_URL=postgresql://....:6543/postgres  # Use pooler
```

## Migration Path

### From Direct to Pooler

1. Update `DATABASE_URL` port from 5432 → 6543
2. Restart application
3. Monitor connection usage
4. Adjust MaxConns if needed

### From Pooler to Direct

1. Ensure total connections < 60 (free tier)
2. Update `DATABASE_URL` port from 6543 → 5432
3. Restart application
4. Verify no "too many connections" errors

## References

- [Supabase Connection Pooling](https://supabase.com/docs/guides/database/connecting-to-postgres#connection-pooler)
- [PostgreSQL Connection Limits](https://www.postgresql.org/docs/current/runtime-config-connection.html)
- [PgBouncer Documentation](https://www.pgbouncer.org/config.html)

package postgres

import "time"

// Option allows for managing postgres options.
type Option func(*Postgres)

// MaxPoolSize sets the max pool size for the connection pool.
func MaxPoolSize(size int) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

// ConnAttempts sets the number of connection attempts.
func ConnAttempts(attempts int) Option {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

// ConnTimeout sets the connection timeout.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}

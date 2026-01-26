package rate_limiter

type TokenBucketLimiter struct {
	tokenBucketCh chan struct{}
}

func NewTokenBucketLimiter() *TokenBucketLimiter {
	return nil
}
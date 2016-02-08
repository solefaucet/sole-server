package cache

// Cache defines interface that one should implement
type Cache interface {
	GetBitcoinPrice() int64
}

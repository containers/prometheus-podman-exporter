package pdcs

var (
	UpdateImages       = updateImages
	DedupImages        = dedupImages
	UpdateDiskUsage    = updateDiskUsage
	SummarizeDiskUsage = summarizeDiskUsage
	StoreContainerSize = storeContainerSize
)

// CachedContainerSize returns the cached rw and rootfs size of a container.
func CachedContainerSize(id string) (int64, int64, bool) {
	cntSizeCache.cacheLock.Lock()
	defer cntSizeCache.cacheLock.Unlock()

	size, ok := cntSizeCache.cache[getID(id)]

	return size.rwSize, size.rootFsSize, ok
}

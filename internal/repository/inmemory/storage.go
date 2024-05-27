package inmemory

import (
	"sync"
)

// postsStorage is a sync.Map that stores posts and comments.
var postsStorage = sync.Map{}

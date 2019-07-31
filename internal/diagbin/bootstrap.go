// +build bootstrap

package diagbin

// If you try to build an embedded binary for a new distribution, the build will fail, complaining
// that the following function is undefined. Change the build tag to reflect the distribution you're
// building for. When you have your embedded binary, change the build tag back.

func Asset(string) ([]byte, error) { return nil, nil }

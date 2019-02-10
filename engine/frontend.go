package engine

import (
	"fmt"
)

func (this *Engine) FrontEnd() bool {
	if this.Wrap.R.URL.Path == "/" {
		this.Wrap.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		this.Wrap.W.Header().Set("Content-Type", "text/html")

		counter := this.Wrap.S.GetInt("counter", 0)
		// this.Wrap.LogAccess(fmt.Sprintf("Counter value was: %d", counter))

		this.Wrap.W.Write([]byte(`Logic -> (` + fmt.Sprintf("%d", counter) + `)`))

		counter++
		this.Wrap.S.SetInt("counter", counter)
		// this.Wrap.LogAccess(fmt.Sprintf("Counter value now: %d", counter))

		return true
	}
	return false
}

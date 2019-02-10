package engine

func (this *Engine) BackEnd() bool {
	return this.Mods.BackEnd(this.Wrap)
}

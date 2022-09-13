package inlineif

func If[V any](cond bool, val V) *V {
	if cond {
		return &val
	}
	return nil
}

func IfElse[V any](cond bool, yes, nope V) V {
	if cond {
		return yes
	} else {
		return nope
	}
}

func IfElsePtr[V any](cond bool, yes, nope V) *V {
	if cond {
		return &yes
	} else {
		return &nope
	}
}

func IfElseFn[V any](cond bool, yes func() V, nope func() V) V {
	if cond {
		return yes()
	} else {
		return nope()
	}
}

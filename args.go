package ioc

type Args []any

type T any
type Y any
type U any
type I any
type O any
type P any

func (a Args) Find(t T) T {
	for _, v := range a {
		res, ok := v.(T)
		if ok {
			return res
		}
	}
	return t
}

func (a Args) Error() error {
	for _, v := range a {
		res, ok := v.(error)
		if ok {
			return res
		}
	}
	return nil
}

func (a Args) Zero() {}

func (a Args) One(t T) T {
	t, _ = a[0].(T)
	return t
}

func (a Args) Two(t T, y Y) (T, Y) {
	t, _ = a[0].(T)
	y, _ = a[1].(Y)
	return t, y
}

func (a Args) Three(t T, y Y, u U) (T, Y, U) {
	t, _ = a[0].(T)
	y, _ = a[1].(Y)
	u, _ = a[2].(U)
	return t, y, u
}

func (a Args) Four(t T, y Y, u U, i I) (T, Y, U, I) {
	t, _ = a[0].(T)
	y, _ = a[1].(Y)
	u, _ = a[2].(U)
	i, _ = a[3].(I)
	return t, y, u, i
}

func (a Args) Five(t T, y Y, u U, i I, o O) (T, Y, U, I, O) {
	t, _ = a[0].(T)
	y, _ = a[1].(Y)
	u, _ = a[2].(U)
	i, _ = a[3].(I)
	o, _ = a[4].(O)
	return t, y, u, i, o
}

func (a Args) Six(t T, y Y, u U, i I, o O, p P) (T, Y, U, I, O, P) {
	t, _ = a[0].(T)
	y, _ = a[1].(Y)
	u, _ = a[2].(U)
	i, _ = a[3].(I)
	o, _ = a[4].(O)
	p, _ = a[5].(P)
	return t, y, u, i, o, p
}

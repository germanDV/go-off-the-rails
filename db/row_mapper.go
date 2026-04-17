package db

func mapRows[S any, D any](src []S, fn func(S) (D, error)) ([]D, error) {
	dst := make([]D, 0, len(src))
	for _, s := range src {
		d, err := fn(s)
		if err != nil {
			return nil, err
		}
		dst = append(dst, d)
	}
	return dst, nil
}

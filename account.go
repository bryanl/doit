package doit

// AccountGet retrieves an account.
func AccountGet(c DConfig) (interface{}, error) {
	a, _, err := c.GodoClient().Account.Get()
	if err != nil {
		return nil, err
	}

	return a, nil
}

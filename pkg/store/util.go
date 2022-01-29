package store

// CheckKeyAndValue returns an error if k is empty or v is nil
func CheckKeyAndValue(k string, v []byte) error {
	if err := CheckKey(k); err != nil {
		return err
	}
	return CheckVal(v)
}

// CheckKey returns an error if k is empty
func CheckKey(k string) error {
	if k == "" {
		return ErrKeyInvalid
	}
	return nil
}

// CheckVal returns an error if v == nil
func CheckVal(v []byte) error {
	if v == nil {
		return ErrValInvalid
	}
	return nil
}

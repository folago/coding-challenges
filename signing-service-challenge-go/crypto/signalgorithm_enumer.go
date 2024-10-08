// Code generated by "enumer -type=SignAlgorithm -json"; DO NOT EDIT.

package crypto

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _SignAlgorithmName = "RSAECC"

var _SignAlgorithmIndex = [...]uint8{0, 3, 6}

const _SignAlgorithmLowerName = "rsaecc"

func (i SignAlgorithm) String() string {
	if i < 0 || i >= SignAlgorithm(len(_SignAlgorithmIndex)-1) {
		return fmt.Sprintf("SignAlgorithm(%d)", i)
	}
	return _SignAlgorithmName[_SignAlgorithmIndex[i]:_SignAlgorithmIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _SignAlgorithmNoOp() {
	var x [1]struct{}
	_ = x[RSA-(0)]
	_ = x[ECC-(1)]
}

var _SignAlgorithmValues = []SignAlgorithm{RSA, ECC}

var _SignAlgorithmNameToValueMap = map[string]SignAlgorithm{
	_SignAlgorithmName[0:3]:      RSA,
	_SignAlgorithmLowerName[0:3]: RSA,
	_SignAlgorithmName[3:6]:      ECC,
	_SignAlgorithmLowerName[3:6]: ECC,
}

var _SignAlgorithmNames = []string{
	_SignAlgorithmName[0:3],
	_SignAlgorithmName[3:6],
}

// SignAlgorithmString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func SignAlgorithmString(s string) (SignAlgorithm, error) {
	if val, ok := _SignAlgorithmNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _SignAlgorithmNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to SignAlgorithm values", s)
}

// SignAlgorithmValues returns all values of the enum
func SignAlgorithmValues() []SignAlgorithm {
	return _SignAlgorithmValues
}

// SignAlgorithmStrings returns a slice of all String values of the enum
func SignAlgorithmStrings() []string {
	strs := make([]string, len(_SignAlgorithmNames))
	copy(strs, _SignAlgorithmNames)
	return strs
}

// IsASignAlgorithm returns "true" if the value is listed in the enum definition. "false" otherwise
func (i SignAlgorithm) IsASignAlgorithm() bool {
	for _, v := range _SignAlgorithmValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for SignAlgorithm
func (i SignAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for SignAlgorithm
func (i *SignAlgorithm) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("SignAlgorithm should be a string, got %s", data)
	}

	var err error
	*i, err = SignAlgorithmString(s)
	return err
}

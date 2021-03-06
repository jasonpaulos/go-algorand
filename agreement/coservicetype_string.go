// Code generated by "stringer -type=coserviceType"; DO NOT EDIT.

package agreement

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[demuxCoserviceType-0]
	_ = x[tokenizerCoserviceType-1]
	_ = x[cryptoVerifierCoserviceType-2]
	_ = x[pseudonodeCoserviceType-3]
	_ = x[clockCoserviceType-4]
	_ = x[networkCoserviceType-5]
}

const _coserviceType_name = "demuxCoserviceTypetokenizerCoserviceTypecryptoVerifierCoserviceTypepseudonodeCoserviceTypeclockCoserviceTypenetworkCoserviceType"

var _coserviceType_index = [...]uint8{0, 18, 40, 67, 90, 108, 128}

func (i coserviceType) String() string {
	if i < 0 || i >= coserviceType(len(_coserviceType_index)-1) {
		return "coserviceType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _coserviceType_name[_coserviceType_index[i]:_coserviceType_index[i+1]]
}

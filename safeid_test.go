package safeid

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	zeroUUIDString               = "00000000-0000-0000-0000-000000000000"
	zeroTypeSafeString           = "0000000000000000000000"
	zeroTypeSafeStringWithPrefix = "test_0000000000000000000000"

	leadingZeroesUUIDString               = "00000000-ffff-ffff-ffff-ffffffffffff"
	leadingZeroesTypeSafeString           = "000001F2si9ujpxVB7VDj1"
	leadingZeroesTypeSafeStringWithPrefix = "test_000001F2si9ujpxVB7VDj1"

	maxUUIDString               = "ffffffff-ffff-ffff-ffff-ffffffffffff"
	maxTypeSafeString           = "7N42dgm5tFLK9N8MT7fHC7"
	maxTypeSafeStringWithPrefix = "test_7N42dgm5tFLK9N8MT7fHC7"

	validUUIDString               = "0193508b-e85e-7812-ba4b-91d85495d7bc"
	validTypeSafeString           = "02Yjy1AYf1ckS6jBZ5zw3G"
	validTypeSafeStringWithPrefix = "test_02Yjy1AYf1ckS6jBZ5zw3G"
)

var (
	zeroUUIDBytes          = [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	leadingZeroesUUIDBytes = [16]byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	maxUUIDBytes           = [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	validUUIDBytes         = [16]byte{0x01, 0x93, 0x50, 0x8b, 0xe8, 0x5e, 0x78, 0x12, 0xba, 0x4b, 0x91, 0xd8, 0x54, 0x95, 0xd7, 0xbc}
)

type empty struct{}

func (empty) Prefix() string { return "" }

type test struct{}

func (test) Prefix() string { return "test" }

func TestIsZeroNil(t *testing.T) {
	var nilID ID[Generic]
	assert.True(t, IsZero(nilID))

	var nilTypedID ID[test]
	assert.True(t, IsZero(nilTypedID))
}

func TestIsZero(t *testing.T) {
	tt := []struct {
		name       string
		bytes      [16]byte
		expOutcome bool
	}{
		{
			"zero",
			zeroUUIDBytes,
			true,
		},
		{
			"set",
			validUUIDBytes,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expOutcome, IsZero(ID[Generic]{uuid.UUID(tc.bytes)}))
		})
	}
}

func TestIsZeroPrefix(t *testing.T) {
	tt := []struct {
		name       string
		bytes      [16]byte
		expOutcome bool
	}{
		{
			"zero",
			zeroUUIDBytes,
			true,
		},
		{
			"set",
			maxUUIDBytes,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expOutcome, IsZero(ID[test]{uuid.UUID(tc.bytes)}))
		})
	}
}

func TestMust(t *testing.T) {
	genericID := ID[Generic]{uuid.UUID(validUUIDBytes)}
	assert.Equal(t, genericID, Must[Generic](genericID, nil))

	custID := ID[test]{uuid.UUID(validUUIDBytes)}
	assert.Equal(t, custID, Must[test](custID, nil))
}

func TestMustPanic(t *testing.T) {
	assert.Panics(t, func() {
		genericID := ID[Generic]{uuid.UUID(validUUIDBytes)}
		Must[Generic](genericID, errors.New("error"))
	})

	assert.Panics(t, func() {
		custID := ID[test]{uuid.UUID(validUUIDBytes)}
		Must[test](custID, errors.New("error"))
	})
}

func TestNew(t *testing.T) {
	genericID, err := New[Generic]()
	assert.NotZero(t, genericID.uuid)
	assert.NoError(t, err)

	custID, err := New[test]()
	assert.NotZero(t, custID.uuid)
	assert.NoError(t, err)

	assert.Panics(t, func() {
		_, _ = New[empty]()
	})
}

func TestFromString(t *testing.T) {
	tt := []struct {
		name     string
		value    string
		expBytes [16]byte
	}{
		{
			"valid",
			validTypeSafeString,
			validUUIDBytes,
		},
		{
			"leading zeroes",
			leadingZeroesTypeSafeString,
			leadingZeroesUUIDBytes,
		},
		{
			"nil",
			"",
			zeroUUIDBytes,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			custID, err := FromString[Generic](tc.value)
			assert.EqualValues(t, tc.expBytes, custID.uuid)
			assert.NoError(t, err)
		})
	}
}

func TestFromStringPrefix(t *testing.T) {
	tt := []struct {
		name     string
		value    string
		expBytes [16]byte
	}{
		{
			"valid",
			validTypeSafeStringWithPrefix,
			validUUIDBytes,
		},
		{
			"leading zeroes",
			leadingZeroesTypeSafeStringWithPrefix,
			leadingZeroesUUIDBytes,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			custID, err := FromString[test](tc.value)
			assert.EqualValues(t, tc.expBytes, custID.uuid)
			assert.NoError(t, err)
		})
	}
}

func TestFromStringError(t *testing.T) {
	tt := []struct {
		name   string
		f      func() (any, error)
		expErr error
	}{
		{
			"prefix mismatch",
			func() (any, error) {
				return FromString[test]("other_02Yjy1AYf1ckS6jBZ5zw3G")
			},
			ParseError("invalid prefix"),
		},
		{
			"generic to prefix",
			func() (any, error) {
				return FromString[test](validTypeSafeString)
			},
			ParseError("invalid prefix"),
		},
		{
			"prefix to generic",
			func() (any, error) {
				return FromString[Generic](validTypeSafeStringWithPrefix)
			},
			ParseError("invalid format"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			errID, err := tc.f()
			assert.Empty(t, errID)
			assert.ErrorIs(t, tc.expErr, err)
		})
	}
}

func TestFromStringEmpty(t *testing.T) {
	tt := []struct {
		name  string
		value string
	}{
		{
			"valid",
			validTypeSafeString,
		},
		{
			"nil",
			"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Panics(t, func() {
				_, _ = FromString[empty](tc.value)
			})
		})
	}
}

func TestFromUUID(t *testing.T) {
	nilID, err := FromUUID[Generic]("")
	assert.Empty(t, nilID)
	assert.NoError(t, err)

	genericID, err := FromUUID[Generic](validUUIDString)
	assert.EqualValues(t, validUUIDBytes, genericID.uuid)
	assert.NoError(t, err)

	nilCustID, err := FromUUID[test]("")
	assert.Empty(t, nilCustID)
	assert.NoError(t, err)

	custID, err := FromUUID[test](validUUIDString)
	assert.EqualValues(t, validUUIDBytes, custID.uuid)
	assert.NoError(t, err)
}

func TestFromUUIDError(t *testing.T) {
	tt := []struct {
		name string
		f    func() (any, error)
	}{
		{
			"generic invalid uuid",
			func() (any, error) { return FromUUID[Generic]("000") },
		},
		{
			"prefix invalid uuid",
			func() (any, error) { return FromUUID[test]("000") },
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			errID, err := tc.f()
			assert.Empty(t, errID)
			assert.ErrorIs(t, ParseError("invalid format"), err)
		})
	}
}

func TestFromUUIDEmpty(t *testing.T) {
	tt := []struct {
		name  string
		value string
	}{
		{
			"valid",
			validUUIDString,
		},
		{
			"nil",
			"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Panics(t, func() {
				_, _ = FromUUID[empty](tc.value)
			})
		})
	}
}

func TestID_StringFormat(t *testing.T) {
	tt := []struct {
		name     string
		bytes    [16]byte
		expValue string
	}{
		{
			"min",
			zeroUUIDBytes,
			zeroTypeSafeString,
		},
		{
			"leading zeros",
			leadingZeroesUUIDBytes,
			leadingZeroesTypeSafeString,
		},
		{
			"max",
			maxUUIDBytes,
			maxTypeSafeString,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			genericID := ID[Generic]{uuid.UUID(tc.bytes)}
			assert.Equal(t, tc.expValue, genericID.String())
		})
	}
}

func TestIDPrefix_StringFormat(t *testing.T) {
	tt := []struct {
		name     string
		bytes    [16]byte
		expValue string
	}{
		{
			"min",
			zeroUUIDBytes,
			zeroTypeSafeStringWithPrefix,
		},
		{
			"leading zeros",
			leadingZeroesUUIDBytes,
			leadingZeroesTypeSafeStringWithPrefix,
		},
		{
			"max",
			maxUUIDBytes,
			maxTypeSafeStringWithPrefix,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			genericID := ID[test]{uuid.UUID(tc.bytes)}
			assert.Equal(t, tc.expValue, genericID.String())
		})
	}
}

func TestID_UUIDFormat(t *testing.T) {
	tt := []struct {
		name    string
		bytes   [16]byte
		expUUID string
	}{
		{
			"min",
			zeroUUIDBytes,
			zeroUUIDString,
		},
		{
			"leading zeros",
			leadingZeroesUUIDBytes,
			leadingZeroesUUIDString,
		},
		{
			"max",
			maxUUIDBytes,
			maxUUIDString,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			genericID := ID[Generic]{uuid.UUID(tc.bytes)}
			assert.Equal(t, tc.expUUID, genericID.UUID())

			custID := ID[test]{uuid.UUID(tc.bytes)}
			assert.Equal(t, tc.expUUID, custID.UUID())
		})
	}
}

func TestMarshalText(t *testing.T) {
	genericID := ID[Generic]{uuid.UUID(validUUIDBytes)}
	res, err := genericID.MarshalText()
	assert.Equal(t, validTypeSafeString, string(res))
	assert.NoError(t, err)

	custID := ID[test]{uuid.UUID(validUUIDBytes)}
	res, err = custID.MarshalText()
	assert.Equal(t, validTypeSafeStringWithPrefix, string(res))
	assert.NoError(t, err)
}

func TestID_UnmarshalText(t *testing.T) {
	var genericID ID[Generic]
	err := genericID.UnmarshalText([]byte(validTypeSafeString))
	assert.NoError(t, err)
	assert.EqualValues(t, validUUIDBytes, genericID.uuid)

	err = genericID.UnmarshalText(nil)
	assert.NoError(t, err)
	assert.EqualValues(t, zeroUUIDBytes, genericID.uuid)

	var custID ID[test]
	err = custID.UnmarshalText([]byte(validTypeSafeStringWithPrefix))
	assert.NoError(t, err)
	assert.EqualValues(t, validUUIDBytes, custID.uuid)

	err = custID.UnmarshalText(nil)
	assert.NoError(t, err)
	assert.EqualValues(t, zeroUUIDBytes, custID.uuid)

	assert.Panics(t, func() {
		var custID ID[empty]
		_ = custID.UnmarshalText([]byte(validTypeSafeString))
	})

	assert.Panics(t, func() {
		var custID ID[empty]
		_ = custID.UnmarshalText(nil)
	})
}

func TestID_Value(t *testing.T) {
	genericID := ID[Generic]{uuid.UUID(validUUIDBytes)}
	res, err := genericID.Value()
	assert.Equal(t, validUUIDString, res)
	assert.NoError(t, err)

	custID := ID[test]{uuid.UUID(validUUIDBytes)}
	res, err = custID.MarshalText()
	assert.Equal(t, []byte(validTypeSafeStringWithPrefix), res)
	assert.NoError(t, err)
}

func TestID_Scan(t *testing.T) {
	var genericID ID[Generic]
	err := genericID.Scan(validUUIDString)
	assert.NoError(t, err)
	assert.EqualValues(t, validUUIDBytes, genericID.uuid)

	err = genericID.Scan(nil)
	assert.NoError(t, err)
	assert.EqualValues(t, zeroUUIDBytes, genericID.uuid)

	var custID ID[test]
	err = custID.Scan(validUUIDString)
	assert.NoError(t, err)
	assert.EqualValues(t, validUUIDBytes, custID.uuid)

	err = custID.Scan(nil)
	assert.NoError(t, err)
	assert.EqualValues(t, zeroUUIDBytes, custID.uuid)

	assert.Panics(t, func() {
		var custID ID[empty]
		_ = custID.Scan(validUUIDString)
	})

	assert.Panics(t, func() {
		var custID ID[empty]
		_ = custID.Scan(nil)
	})
}

package safeid

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
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
	assertEqual(t, true, IsZero(nilID))

	var nilTypedID ID[test]
	assertEqual(t, true, IsZero(nilTypedID))
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
			assertEqual(t, tc.expOutcome, IsZero(ID[Generic]{uuid.UUID(tc.bytes)}))
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
			assertEqual(t, tc.expOutcome, IsZero(ID[test]{uuid.UUID(tc.bytes)}))
		})
	}
}

func TestMust(t *testing.T) {
	genericID := ID[Generic]{uuid.UUID(validUUIDBytes)}
	assertEqual(t, genericID, Must[Generic](genericID, nil))

	custID := ID[test]{uuid.UUID(validUUIDBytes)}
	assertEqual(t, custID, Must[test](custID, nil))
}

func TestMustPanic(t *testing.T) {
	assertPanics(t, func() {
		genericID := ID[Generic]{uuid.UUID(validUUIDBytes)}
		Must[Generic](genericID, errors.New("error"))
	})

	assertPanics(t, func() {
		custID := ID[test]{uuid.UUID(validUUIDBytes)}
		Must[test](custID, errors.New("error"))
	})
}

func TestNew(t *testing.T) {
	genericID, err := New[Generic]()
	assertNotZero(t, genericID.uuid)
	assertErrorIs(t, err, nil)

	custID, err := New[test]()
	assertNotZero(t, custID.uuid)
	assertErrorIs(t, err, nil)

	assertPanics(t, func() {
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
			assertEqual(t, tc.expBytes[:], custID.uuid[:])
			assertErrorIs(t, err, nil)
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
			assertEqual(t, tc.expBytes[:], custID.uuid[:])
			assertErrorIs(t, err, nil)
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
			assertZero(t, errID)
			assertErrorIs(t, tc.expErr, err)
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
			assertPanics(t, func() {
				_, _ = FromString[empty](tc.value)
			})
		})
	}
}

func TestFromUUID(t *testing.T) {
	nilID, err := FromUUID[Generic]("")
	assertZero(t, nilID)
	assertErrorIs(t, err, nil)

	genericID, err := FromUUID[Generic](validUUIDString)
	assertEqual(t, validUUIDBytes[:], genericID.uuid[:])
	assertErrorIs(t, err, nil)

	nilCustID, err := FromUUID[test]("")
	assertZero(t, nilCustID)
	assertErrorIs(t, err, nil)

	custID, err := FromUUID[test](validUUIDString)
	assertEqual(t, validUUIDBytes[:], custID.uuid[:])
	assertErrorIs(t, err, nil)
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
			assertZero(t, errID)
			assertErrorIs(t, ParseError("invalid format"), err)
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
			assertPanics(t, func() {
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
			assertEqual(t, tc.expValue, genericID.String())
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
			assertEqual(t, tc.expValue, genericID.String())
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
			assertEqual(t, tc.expUUID, genericID.UUID())

			custID := ID[test]{uuid.UUID(tc.bytes)}
			assertEqual(t, tc.expUUID, custID.UUID())
		})
	}
}

func TestMarshalText(t *testing.T) {
	genericID := ID[Generic]{uuid.UUID(validUUIDBytes)}
	res, err := genericID.MarshalText()
	assertEqual(t, validTypeSafeString, string(res))
	assertErrorIs(t, err, nil)

	custID := ID[test]{uuid.UUID(validUUIDBytes)}
	res, err = custID.MarshalText()
	assertEqual(t, validTypeSafeStringWithPrefix, string(res))
	assertErrorIs(t, err, nil)
}

func TestID_UnmarshalText(t *testing.T) {
	var genericID ID[Generic]
	err := genericID.UnmarshalText([]byte(validTypeSafeString))
	assertErrorIs(t, err, nil)
	assertEqual(t, validUUIDBytes[:], genericID.uuid[:])

	err = genericID.UnmarshalText(nil)
	assertErrorIs(t, err, nil)
	assertEqual(t, zeroUUIDBytes[:], genericID.uuid[:])

	var custID ID[test]
	err = custID.UnmarshalText([]byte(validTypeSafeStringWithPrefix))
	assertErrorIs(t, err, nil)
	assertEqual(t, validUUIDBytes[:], custID.uuid[:])

	err = custID.UnmarshalText(nil)
	assertErrorIs(t, err, nil)
	assertEqual(t, zeroUUIDBytes[:], custID.uuid[:])

	assertPanics(t, func() {
		var custID ID[empty]
		_ = custID.UnmarshalText([]byte(validTypeSafeString))
	})

	assertPanics(t, func() {
		var custID ID[empty]
		_ = custID.UnmarshalText(nil)
	})
}

func TestID_Value(t *testing.T) {
	genericID := ID[Generic]{uuid.UUID(validUUIDBytes)}
	res, err := genericID.Value()
	assertEqual(t, validUUIDString, res)
	assertErrorIs(t, err, nil)

	custID := ID[test]{uuid.UUID(validUUIDBytes)}
	res, err = custID.MarshalText()
	assertEqual(t, []byte(validTypeSafeStringWithPrefix), res)
	assertErrorIs(t, err, nil)
}

func TestID_Scan(t *testing.T) {
	var genericID ID[Generic]
	err := genericID.Scan(validUUIDString)
	assertErrorIs(t, err, nil)
	assertEqual(t, validUUIDBytes[:], genericID.uuid[:])

	err = genericID.Scan(nil)
	assertErrorIs(t, err, nil)
	assertEqual(t, zeroUUIDBytes[:], genericID.uuid[:])

	var custID ID[test]
	err = custID.Scan(validUUIDString)
	assertErrorIs(t, err, nil)
	assertEqual(t, validUUIDBytes[:], custID.uuid[:])

	err = custID.Scan(nil)
	assertErrorIs(t, err, nil)
	assertEqual(t, zeroUUIDBytes[:], custID.uuid[:])

	assertPanics(t, func() {
		var custID ID[empty]
		_ = custID.Scan(validUUIDString)
	})

	assertPanics(t, func() {
		var custID ID[empty]
		_ = custID.Scan(nil)
	})
}

func assertErrorIs(t *testing.T, err error, expErr error) {
	t.Helper()
	if !errors.Is(err, expErr) {
		t.Errorf("expected error %#v, got %#v", expErr, err)
	}
}

func assertEqual(t *testing.T, exp, v any) {
	t.Helper()
	if !reflect.DeepEqual(exp, v) {
		t.Errorf("expected %#v, got %#v", exp, v)
	}
}

func assertPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic did not happened")
		}
	}()
	f()
}

func assertZero[T any](t *testing.T, v T) {
	zero := reflect.Zero(reflect.TypeOf(v)).Interface()
	if !reflect.DeepEqual(v, zero) {
		t.Errorf("expected zero-value, got %#v", v)
	}
}

func assertNotZero[T any](t *testing.T, v T) {
	t.Helper()
	zero := reflect.Zero(reflect.TypeOf(v)).Interface()
	if reflect.DeepEqual(v, zero) {
		t.Errorf("expected non zero-value, got %#v", v)
	}
}

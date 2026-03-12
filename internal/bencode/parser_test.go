package bencode

import "testing"

//
// INTEGER TESTS
//

func TestDecodeIntPositive(t *testing.T) {
	got, err := Decode([]byte("i42e"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if got != 42 {
		t.Fatalf("got %v, want 42", got)
	}
}

func TestDecodeIntZero(t *testing.T) {
	got, err := Decode([]byte("i0e"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if got != 0 {
		t.Fatalf("got %v, want 0", got)
	}
}

func TestDecodeIntNegative(t *testing.T) {
	got, err := Decode([]byte("i-7e"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if got != -7 {
		t.Fatalf("got %v, want -7", got)
	}
}

func TestDecodeIntLarge(t *testing.T) {
	got, err := Decode([]byte("i123456e"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if got != 123456 {
		t.Fatalf("got %v, want 123456", got)
	}
}

func TestDecodeIntEmpty(t *testing.T) {
	_, err := Decode([]byte("ie"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeIntNegativeWithoutDigits(t *testing.T) {
	_, err := Decode([]byte("i-e"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeIntInvalidCharacter(t *testing.T) {
	_, err := Decode([]byte("i4xe"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeIntUnterminated(t *testing.T) {
	_, err := Decode([]byte("i42"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeIntOnlyPrefix(t *testing.T) {
	_, err := Decode([]byte("i"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeIntTrailingData(t *testing.T) {
	_, err := Decode([]byte("i42ejunk"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

//
// STRING TESTS
//

func TestDecodeStringBasic(t *testing.T) {
	got, err := Decode([]byte("4:spam"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if got != "spam" {
		t.Fatalf("got %v, want spam", got)
	}
}

func TestDecodeStringHello(t *testing.T) {
	got, err := Decode([]byte("5:hello"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if got != "hello" {
		t.Fatalf("got %v, want hello", got)
	}
}

func TestDecodeStringSingleCharacter(t *testing.T) {
	got, err := Decode([]byte("1:a"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if got != "a" {
		t.Fatalf("got %v, want a", got)
	}
}

func TestDecodeStringEmpty(t *testing.T) {
	got, err := Decode([]byte("0:"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if got != "" {
		t.Fatalf("got %v, want empty string", got)
	}
}

func TestDecodeStringLong(t *testing.T) {
	got, err := Decode([]byte("11:hello world"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if got != "hello world" {
		t.Fatalf("got %v, want hello world", got)
	}
}

func TestDecodeStringMultiDigitLength(t *testing.T) {
	got, err := Decode([]byte("10:helloworld"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if got != "helloworld" {
		t.Fatalf("got %v, want helloworld", got)
	}
}

func TestDecodeStringMissingLength(t *testing.T) {
	_, err := Decode([]byte(":spam"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeStringMissingColon(t *testing.T) {
	_, err := Decode([]byte("4spam"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeStringTooShort(t *testing.T) {
	_, err := Decode([]byte("4:spa"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeStringDeclaredLengthTooLong(t *testing.T) {
	_, err := Decode([]byte("12:short"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeStringEmptyPayloadButNonZeroLength(t *testing.T) {
	_, err := Decode([]byte("4:"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeStringInvalidLength(t *testing.T) {
	_, err := Decode([]byte("x:spam"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeStringInvalidLengthCharacterAfterDigit(t *testing.T) {
	_, err := Decode([]byte("4a:spam"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeStringTrailingData(t *testing.T) {
	_, err := Decode([]byte("4:spamjunk"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeStringZeroLengthTrailingData(t *testing.T) {
	_, err := Decode([]byte("0:junk"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

//
// LIST TESTS
//

func TestDecodeListEmpty(t *testing.T) {
	got, err := Decode([]byte("le"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	list, ok := got.([]any)
	if !ok {
		t.Fatalf("got %T, want []any", got)
	}

	if len(list) != 0 {
		t.Fatalf("got length %d, want 0", len(list))
	}
}

func TestDecodeListStrings(t *testing.T) {
	got, err := Decode([]byte("l4:spam4:eggse"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	list := got.([]any)

	if len(list) != 2 {
		t.Fatalf("got length %d, want 2", len(list))
	}

	if list[0] != "spam" || list[1] != "eggs" {
		t.Fatalf("unexpected list contents: %#v", list)
	}
}

func TestDecodeListIntegers(t *testing.T) {
	got, err := Decode([]byte("li1ei2ei3ee"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	list := got.([]any)

	if len(list) != 3 {
		t.Fatalf("got length %d, want 3", len(list))
	}

	if list[0] != 1 || list[1] != 2 || list[2] != 3 {
		t.Fatalf("unexpected list contents: %#v", list)
	}
}

func TestDecodeListMixedTypes(t *testing.T) {
	got, err := Decode([]byte("l4:spami42e3:abce"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	list := got.([]any)

	if len(list) != 3 {
		t.Fatalf("got length %d, want 3", len(list))
	}

	if list[0] != "spam" || list[1] != 42 || list[2] != "abc" {
		t.Fatalf("unexpected list contents: %#v", list)
	}
}

func TestDecodeListUnterminated(t *testing.T) {
	_, err := Decode([]byte("l4:spam4:eggs"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

//
// DICTIONARY TESTS
//

func TestDecodeDictEmpty(t *testing.T) {
	got, err := Decode([]byte("de"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	dict, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("got %T, want map[string]any", got)
	}

	if len(dict) != 0 {
		t.Fatalf("got length %d, want 0", len(dict))
	}
}

func TestDecodeDictBasic(t *testing.T) {
	got, err := Decode([]byte("d3:cow3:moo4:spam4:eggse"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	dict, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("got %T, want map[string]any", got)
	}

	if len(dict) != 2 {
		t.Fatalf("got length %d, want 2", len(dict))
	}

	if dict["cow"] != "moo" {
		t.Fatalf("got %v, want moo", dict["cow"])
	}

	if dict["spam"] != "eggs" {
		t.Fatalf("got %v, want eggs", dict["spam"])
	}
}

func TestDecodeDictMixedValues(t *testing.T) {
	got, err := Decode([]byte("d3:fooi42e3:bar4:spame"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	dict, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("got %T, want map[string]any", got)
	}

	if len(dict) != 2 {
		t.Fatalf("got length %d, want 2", len(dict))
	}

	if dict["foo"] != 42 {
		t.Fatalf("got %v, want 42", dict["foo"])
	}

	if dict["bar"] != "spam" {
		t.Fatalf("got %v, want spam", dict["bar"])
	}
}

func TestDecodeDictWithList(t *testing.T) {
	got, err := Decode([]byte("d4:listli1ei2ei3eee"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	dict, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("got %T, want map[string]any", got)
	}

	rawList, ok := dict["list"]
	if !ok {
		t.Fatalf("missing key %q", "list")
	}

	list, ok := rawList.([]any)
	if !ok {
		t.Fatalf("got %T, want []any", rawList)
	}

	if len(list) != 3 {
		t.Fatalf("got length %d, want 3", len(list))
	}

	if list[0] != 1 || list[1] != 2 || list[2] != 3 {
		t.Fatalf("unexpected list contents: %#v", list)
	}
}

func TestDecodeDictWithNestedDict(t *testing.T) {
	got, err := Decode([]byte("d5:innerd3:fooi7eee"))
	if err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	dict, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("got %T, want map[string]any", got)
	}

	rawInner, ok := dict["inner"]
	if !ok {
		t.Fatalf("missing key %q", "inner")
	}

	inner, ok := rawInner.(map[string]any)
	if !ok {
		t.Fatalf("got %T, want map[string]any", rawInner)
	}

	if inner["foo"] != 7 {
		t.Fatalf("got %v, want 7", inner["foo"])
	}
}

func TestDecodeDictUnterminated(t *testing.T) {
	_, err := Decode([]byte("d3:cow3:moo"))
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

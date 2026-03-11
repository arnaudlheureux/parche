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

func TestDecodeStringInvalidLength(t *testing.T) {
	_, err := Decode([]byte("x:spam"))
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
package ones_gfx

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestParseGPUOperation(t *testing.T) {
	for in, want := range map[string]GPUOperation{"ADD": OperationAdd, "delete": OperationDelete, " Add ": OperationAdd} {
		got, err := ParseGPUOperation(in)
		if err != nil || got != want {
			t.Errorf("ParseGPUOperation(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"", "ADDD", "REMOVE"} {
		if _, err := ParseGPUOperation(in); err == nil {
			t.Errorf("ParseGPUOperation(%q): expected error", in)
		}
	}
	if GpuAction("ADDD").Valid() || !GpuActionDelete.Valid() {
		t.Error("GpuAction.Valid misreports")
	}
}

func TestParseOperationStatus(t *testing.T) {
	got, err := ParseOperationStatus("success")
	if err != nil || got != OperationStatusSuccess {
		t.Fatalf("got %q, %v", got, err)
	}
	if _, err := ParseOperationStatus("DONE"); err == nil {
		t.Error("expected error for DONE")
	}
	if OperationStatus("DONE").Valid() || !OperationStatusPending.Valid() {
		t.Error("OperationStatus.Valid misreports")
	}

	// FM sends FAILED for a failed operation; it must end a wait like FAILURE.
	for _, s := range []OperationStatus{OperationStatusFailure, OperationStatusFailed} {
		if !s.Valid() || !s.IsTerminal() || !s.IsFailure() {
			t.Errorf("%s: want valid, terminal and a failure", s)
		}
	}
	if OperationStatusSuccess.IsFailure() || OperationStatusRunning.IsTerminal() {
		t.Error("SUCCESS/RUNNING misreported")
	}
	if got, err := ParseOperationStatus("failed"); err != nil || got != OperationStatusFailed {
		t.Errorf("ParseOperationStatus(failed) = %q, %v", got, err)
	}
}

func TestValidationErrorIsONESError(t *testing.T) {
	var err error = &ValidationError{Field: "fabricName", Message: "must not be empty"}
	var onesErr ONESError
	if !errors.As(err, &onesErr) {
		t.Fatal("ValidationError does not implement ONESError")
	}
	if want := "invalid fabricName: must not be empty"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestTimestamp(t *testing.T) {
	var v struct {
		Local   Timestamp `json:"local"`
		Zoned   Timestamp `json:"zoned"`
		Epoch   Timestamp `json:"epoch"`
		Null    Timestamp `json:"null"`
		Weird   Timestamp `json:"weird"`
		Missing Timestamp `json:"missing"`
	}
	data := `{"local":"2026-04-18T14:03:32.632113","zoned":"2026-04-18T14:03:32Z",
		"epoch":1776521012631,"null":null,"weird":[2026,4,18]}`
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if tm, ok := v.Local.Time(); !ok || tm.Hour() != 14 || tm.Nanosecond() != 632113000 {
		t.Errorf("Local.Time() = %v, %v", tm, ok)
	}
	if tm, ok := v.Zoned.Time(); !ok || !tm.Equal(time.Date(2026, 4, 18, 14, 3, 32, 0, time.UTC)) {
		t.Errorf("Zoned.Time() = %v, %v", tm, ok)
	}
	if tm, ok := v.Epoch.Time(); !ok || tm.UnixMilli() != 1776521012631 {
		t.Errorf("Epoch.Time() = %v, %v", tm, ok)
	}
	if !v.Null.IsZero() || !v.Missing.IsZero() {
		t.Error("null/missing should be zero")
	}
	// An unexpected format must not fail decoding; it is kept verbatim.
	if v.Weird.String() != "[2026,4,18]" {
		t.Errorf("Weird = %q", v.Weird)
	}
	if _, ok := v.Weird.Time(); ok {
		t.Error("Weird.Time() should not parse")
	}
}

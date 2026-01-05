package snapshot

import (
	"errors"
	"testing"
	"time"
)

func TestStoreSaveAndLatest(t *testing.T) {
	dir := t.TempDir()
	store := Store{Directory: dir}

	early := Entry{
		PluginName:   "alpha",
		PreviousSpec: "alpha@1.0.0",
		Timestamp:    time.Now().Add(-time.Hour),
	}
	late := Entry{
		PluginName:   "alpha",
		PreviousSpec: "alpha@2.0.0",
		Timestamp:    time.Now(),
	}

	if err := store.Save(early); err != nil {
		t.Fatalf("save early: %v", err)
	}
	if err := store.Save(late); err != nil {
		t.Fatalf("save late: %v", err)
	}

	got, err := store.Latest("alpha")
	if err != nil {
		t.Fatalf("latest: %v", err)
	}
	if got.PreviousSpec != "alpha@2.0.0" {
		t.Fatalf("expected latest spec, got %s", got.PreviousSpec)
	}
}

func TestStoreLatestMissing(t *testing.T) {
	store := Store{Directory: t.TempDir()}
	_, err := store.Latest("missing")
	if !errors.Is(err, ErrSnapshotNotFound) {
		t.Fatalf("expected ErrSnapshotNotFound, got %v", err)
	}
}

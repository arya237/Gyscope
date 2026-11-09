package tui

import "Gyscope/internal/domain/system"

type SnapshotReader interface {
	Snapshot() system.Snapshot
}

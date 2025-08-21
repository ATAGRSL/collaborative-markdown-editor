package ot

import (
	"sync"
)

// Manager handles operational transformation for a document
type Manager struct {
	mu         sync.RWMutex
	currentDoc string
	version    int
}

// NewManager creates a new OT manager
func NewManager(initialDoc string) *Manager {
	return &Manager{
		currentDoc: initialDoc,
		version:    0,
	}
}

// ApplyOperation applies a new operation to the document
func (m *Manager) ApplyOperation(op *Operation) (*Operation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// For now, just apply the operation directly without transformation
	// since we're using simple content replacement
	if op == nil {
		return op, nil
	}

	// Update the document
	m.applyToDocument(op)

	// Update version
	m.version++
	op.Version = m.version

	return op, nil
}



// GetCurrentDocument returns the current document state
func (m *Manager) GetCurrentDocument() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentDoc
}

// SetCurrentDocument sets the current document state (for simple content replacement)
func (m *Manager) SetCurrentDocument(content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentDoc = content
	m.version = 0
}

// GetVersion returns the current version
func (m *Manager) GetVersion() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.version
}

// applyToDocument applies an operation to the current document
func (m *Manager) applyToDocument(op *Operation) {
	runes := []rune(m.currentDoc)

	switch op.Type {
	case Insert:
		if op.Position > len(runes) {
			op.Position = len(runes)
		}
		newRunes := make([]rune, 0, len(runes)+1)
		newRunes = append(newRunes, runes[:op.Position]...)
		newRunes = append(newRunes, []rune(op.Character)...)
		newRunes = append(newRunes, runes[op.Position:]...)
		m.currentDoc = string(newRunes)

	case Delete:
		if op.Position >= len(runes) {
			return
		}
		end := op.Position + op.Length
		if end > len(runes) {
			end = len(runes)
		}
		newRunes := make([]rune, 0, len(runes)-op.Length)
		newRunes = append(newRunes, runes[:op.Position]...)
		newRunes = append(newRunes, runes[end:]...)
		m.currentDoc = string(newRunes)
	}
}





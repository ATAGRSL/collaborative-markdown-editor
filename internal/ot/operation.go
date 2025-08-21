package ot

import (
	"encoding/json"
	"fmt"
)

// OperationType represents the type of operation
type OperationType string

const (
	Insert OperationType = "insert"
	Delete OperationType = "delete"
)

// Operation represents a single edit operation
type Operation struct {
	Type      OperationType `json:"type"`
	Position  int           `json:"position"`
	Character string        `json:"character,omitempty"` // Only for Insert operations
	Length    int           `json:"length,omitempty"`    // Only for Delete operations
	Version   int           `json:"version"`
	ClientID  string        `json:"clientId"`
}

// NewInsertOperation creates a new insert operation
func NewInsertOperation(position int, character string, version int, clientID string) *Operation {
	return &Operation{
		Type:      Insert,
		Position:  position,
		Character: character,
		Version:   version,
		ClientID:  clientID,
	}
}

// NewDeleteOperation creates a new delete operation
func NewDeleteOperation(position, length, version int, clientID string) *Operation {
	return &Operation{
		Type:     Delete,
		Position: position,
		Length:   length,
		Version:  version,
		ClientID: clientID,
	}
}

// Transform transforms this operation against another concurrent operation
// Returns the transformed operations for both this and the other operation
func (op *Operation) Transform(other *Operation) (*Operation, *Operation) {
	if op.Version < other.Version {
		// This operation is older, transform it against the newer one
		transformedOp, transformedOther := op.transformAgainst(other)
		return transformedOp, transformedOther
	} else if op.Version > other.Version {
		// Other operation is older, transform it against this one
		transformedOp, transformedOther := other.transformAgainst(op)
		return transformedOp, transformedOther
	}

	// Same version, no transformation needed
	return op, other
}

// transformAgainst transforms this operation against another operation
func (op *Operation) transformAgainst(other *Operation) (*Operation, *Operation) {
	var transformedOp *Operation

	switch op.Type {
	case Insert:
		transformedOp = op.transformInsert(other)
	case Delete:
		transformedOp = op.transformDelete(other)
	default:
		transformedOp = op
	}

	// For now, just return the transformed operation and the original other
	// In a more sophisticated implementation, both operations would be transformed
	return transformedOp, other
}

// transformInsert transforms an insert operation
func (op *Operation) transformInsert(other *Operation) *Operation {
	switch other.Type {
	case Insert:
		// If other insert is before this insert, shift position right
		if other.Position <= op.Position {
			return NewInsertOperation(op.Position+1, op.Character, op.Version, op.ClientID)
		}
		return op
	case Delete:
		// If delete is before insert, shift position left
		if other.Position < op.Position {
			newPos := op.Position - other.Length
			if newPos < other.Position {
				newPos = other.Position
			}
			return NewInsertOperation(newPos, op.Character, op.Version, op.ClientID)
		}
		return op
	}
	return op
}

// transformDelete transforms a delete operation
func (op *Operation) transformDelete(other *Operation) *Operation {
	switch other.Type {
	case Insert:
		// If insert is before delete, shift position right
		if other.Position <= op.Position {
			return NewDeleteOperation(op.Position+1, op.Length, op.Version, op.ClientID)
		}
		return op
	case Delete:
		// Handle overlapping deletes
		if other.Position < op.Position {
			// Other delete starts before this delete
			if other.Position+other.Length <= op.Position {
				// No overlap, shift left
				return NewDeleteOperation(op.Position-other.Length, op.Length, op.Version, op.ClientID)
			} else {
				// Overlap, adjust position and length
				newPos := other.Position
				newLen := op.Position + op.Length - (other.Position + other.Length)
				if newLen <= 0 {
					// Complete overlap, delete becomes no-op
					return nil
				}
				return NewDeleteOperation(newPos, newLen, op.Version, op.ClientID)
			}
		} else if other.Position < op.Position+op.Length {
			// Overlap, reduce length
			newLen := op.Length - (op.Position + op.Length - other.Position)
			if newLen <= 0 {
				return nil
			}
			return NewDeleteOperation(op.Position, newLen, op.Version, op.ClientID)
		}
		return op
	}
	return op
}

// ToJSON converts the operation to JSON
func (op *Operation) ToJSON() ([]byte, error) {
	return json.Marshal(op)
}

// FromJSON creates an operation from JSON
func OperationFromJSON(data []byte) (*Operation, error) {
	var op Operation
	err := json.Unmarshal(data, &op)
	return &op, err
}

// String returns a string representation of the operation
func (op *Operation) String() string {
	switch op.Type {
	case Insert:
		return fmt.Sprintf("Insert('%s', %d) v%d", op.Character, op.Position, op.Version)
	case Delete:
		return fmt.Sprintf("Delete(%d, %d) v%d", op.Position, op.Length, op.Version)
	}
	return "Unknown operation"
}

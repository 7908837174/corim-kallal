// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"errors"
	"fmt"
)

// DomainMembershipTriple represents a domain membership triple record as
// specified in the CoRIM specification. It links a domain identifier to its
// member environments, forming a subject-predicate-object triple where the
// domain-id is the subject and the members array contains the objects.
// This enables describing the topological structure of composite attesters.
type DomainMembershipTriple struct {
	_        struct{}      `cbor:",toarray"`
	DomainID Environment   `json:"domain-id"`
	Members  []Environment `json:"members"`
}

// Valid validates the domain membership triple according to CoRIM specification
func (o *DomainMembershipTriple) Valid() error {
	if err := o.DomainID.Valid(); err != nil {
		return fmt.Errorf("domain-id validation failed: %w", err)
	}

	if len(o.Members) == 0 {
		return errors.New("members validation failed: no member environments")
	}

	for i, member := range o.Members {
		if err := member.Valid(); err != nil {
			return fmt.Errorf("members validation failed: member at index %d: %w", i, err)
		}
	}

	return nil
}

// AddMember adds a member environment to the domain membership triple
func (o *DomainMembershipTriple) AddMember(env Environment) *DomainMembershipTriple {
	if o != nil {
		o.Members = append(o.Members, env)
	}
	return o
}

// DomainMembershipTriples is a container for DomainMembershipTriple instances.
// It represents the membership-triples array in the CoRIM specification.
type DomainMembershipTriples []DomainMembershipTriple

func NewDomainMembershipTriples() *DomainMembershipTriples {
	return &DomainMembershipTriples{}
}

func (o *DomainMembershipTriples) Valid() error {
	for i, triple := range *o {
		if err := triple.Valid(); err != nil {
			return fmt.Errorf("domain membership triple at index %d: %w", i, err)
		}
	}
	return nil
}

func (o *DomainMembershipTriples) IsEmpty() bool {
	return len(*o) == 0
}

func (o *DomainMembershipTriples) Add(val *DomainMembershipTriple) *DomainMembershipTriples {
	if o != nil && val != nil {
		*o = append(*o, *val)
	}
	return o
}

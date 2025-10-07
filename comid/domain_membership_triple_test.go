// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainMembershipTriple_Valid_Success(t *testing.T) {
	// Create a valid domain ID
	domainID := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Domain Vendor").
			SetModel("Domain Model"),
	}

	// Create valid member environments
	member1 := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Member 1 Vendor").
			SetModel("Member 1 Model"),
	}

	member2 := Environment{
		Instance: MustNewUEIDInstance(TestUEID),
	}

	triple := &DomainMembershipTriple{
		DomainID: domainID,
		Members:  []Environment{member1, member2},
	}

	err := triple.Valid()
	assert.NoError(t, err)
}

func TestDomainMembershipTriple_Valid_EmptyDomainID(t *testing.T) {
	member := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Member Vendor").
			SetModel("Member Model"),
	}

	triple := &DomainMembershipTriple{
		DomainID: Environment{}, // Empty domain ID
		Members:  []Environment{member},
	}

	err := triple.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "domain-id validation failed")
}

func TestDomainMembershipTriple_Valid_EmptyMembers(t *testing.T) {
	domainID := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Domain Vendor").
			SetModel("Domain Model"),
	}

	triple := &DomainMembershipTriple{
		DomainID: domainID,
		Members:  []Environment{}, // Empty members
	}

	err := triple.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no member environments")
}

func TestDomainMembershipTriple_Valid_InvalidMember(t *testing.T) {
	domainID := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Domain Vendor").
			SetModel("Domain Model"),
	}

	invalidMember := Environment{} // Invalid empty member

	triple := &DomainMembershipTriple{
		DomainID: domainID,
		Members:  []Environment{invalidMember},
	}

	err := triple.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "member at index 0")
}

func TestDomainMembershipTriple_AddMember(t *testing.T) {
	domainID := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Domain Vendor").
			SetModel("Domain Model"),
	}

	triple := &DomainMembershipTriple{
		DomainID: domainID,
		Members:  []Environment{},
	}

	member := Environment{
		Instance: MustNewUEIDInstance(TestUEID),
	}

	result := triple.AddMember(member)
	assert.Equal(t, triple, result)
	assert.Len(t, triple.Members, 1)
	assert.Equal(t, member, triple.Members[0])
}

func TestDomainMembershipTriples_NewDomainMembershipTriples(t *testing.T) {
	triples := NewDomainMembershipTriples()
	require.NotNil(t, triples)
	assert.True(t, triples.IsEmpty())
}

func TestDomainMembershipTriples_Add_Success(t *testing.T) {
	triples := NewDomainMembershipTriples()

	domainID := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Domain Vendor").
			SetModel("Domain Model"),
	}

	member := Environment{
		Instance: MustNewUEIDInstance(TestUEID),
	}

	triple := &DomainMembershipTriple{
		DomainID: domainID,
		Members:  []Environment{member},
	}

	result := triples.Add(triple)
	assert.Equal(t, triples, result)
	assert.False(t, triples.IsEmpty())
}

func TestDomainMembershipTriples_Valid_Success(t *testing.T) {
	triples := NewDomainMembershipTriples()

	domainID := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Domain Vendor").
			SetModel("Domain Model"),
	}

	member := Environment{
		Instance: MustNewUEIDInstance(TestUEID),
	}

	triple := &DomainMembershipTriple{
		DomainID: domainID,
		Members:  []Environment{member},
	}

	triples.Add(triple)

	err := triples.Valid()
	assert.NoError(t, err)
}

func TestDomainMembershipTriples_Valid_Empty(t *testing.T) {
	triples := NewDomainMembershipTriples()

	err := triples.Valid()
	assert.NoError(t, err) // Empty collection is valid
}

func TestDomainMembershipTriples_Valid_InvalidTriple(t *testing.T) {
	triples := NewDomainMembershipTriples()

	// Add an invalid triple (empty domain ID)
	triple := &DomainMembershipTriple{
		DomainID: Environment{}, // Invalid empty domain ID
		Members: []Environment{
			{Instance: MustNewUEIDInstance(TestUEID)},
		},
	}

	triples.Add(triple)

	err := triples.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "domain membership triple at index 0")
}

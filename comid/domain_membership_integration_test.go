// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTriples_AddDomainMembershipTriple_Success(t *testing.T) {
	triples := &Triples{}

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

	result := triples.AddDomainMembershipTriple(triple)
	assert.Equal(t, triples, result)
	assert.NotNil(t, triples.DomainMembershipTriples)
	assert.False(t, triples.DomainMembershipTriples.IsEmpty())
}

func TestTriples_AddDomainMembershipTriple_NilTriples(t *testing.T) {
	var triples *Triples = nil

	triple := &DomainMembershipTriple{
		DomainID: Environment{
			Class: NewClassUUID(TestUUID).
				SetVendor("Domain Vendor").
				SetModel("Domain Model"),
		},
		Members: []Environment{
			{Instance: MustNewUEIDInstance(TestUEID)},
		},
	}

	result := triples.AddDomainMembershipTriple(triple)
	assert.Nil(t, result)
}

func TestTriples_Valid_WithDomainMembershipTriples(t *testing.T) {
	triples := &Triples{}

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

	triples.AddDomainMembershipTriple(triple)

	err := triples.Valid()
	assert.NoError(t, err)
}

func TestTriples_Valid_WithInvalidDomainMembershipTriples(t *testing.T) {
	triples := &Triples{}

	// Add an invalid domain membership triple (empty domain ID)
	triple := &DomainMembershipTriple{
		DomainID: Environment{}, // Invalid empty domain ID
		Members: []Environment{
			{Instance: MustNewUEIDInstance(TestUEID)},
		},
	}

	triples.AddDomainMembershipTriple(triple)

	err := triples.Valid()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "domain membership triples")
}

func TestComid_AddDomainMembershipTriple_Success(t *testing.T) {
	comid := NewComid().
		SetLanguage("en-US").
		SetTagIdentity("test-domain-membership", 1).
		AddEntity("Test Corp", &TestRegID, RoleCreator, RoleTagCreator)

	domainID := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Test Vendor").
			SetModel("Test Model"),
	}

	member := Environment{
		Instance: MustNewUEIDInstance(TestUEID),
	}

	triple := &DomainMembershipTriple{
		DomainID: domainID,
		Members:  []Environment{member},
	}

	result := comid.AddDomainMembershipTriple(triple)
	assert.Equal(t, comid, result)
	assert.NotNil(t, comid.Triples.DomainMembershipTriples)
	assert.False(t, comid.Triples.DomainMembershipTriples.IsEmpty())
}

func TestComid_Full_Example_WithDomainMembershipTriple(t *testing.T) {
	comid := NewComid().
		SetLanguage("en-US").
		SetTagIdentity("domain-membership-test-comid", 1).
		AddEntity("Test Corp", &TestRegID, RoleCreator, RoleTagCreator)

	// Create a parent domain (composite device)
	parentDomain := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Test Vendor").
			SetModel("Composite Device").
			SetLayer(1),
		Instance: MustNewUEIDInstance(TestUEID),
	}

	// Create member environments (sub-components)
	tpmMember := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("TPM Vendor").
			SetModel("TPM 2.0"),
	}

	cpuMember := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("CPU Vendor").
			SetModel("Secure CPU"),
	}

	// Create domain membership triple
	triple := &DomainMembershipTriple{
		DomainID: parentDomain,
		Members:  []Environment{tpmMember, cpuMember},
	}

	comid.AddDomainMembershipTriple(triple)

	// Validation tests
	err := comid.Valid()
	require.NoError(t, err)

	// Serialization tests
	cborData, err := comid.ToCBOR()
	require.NoError(t, err)
	assert.NotEmpty(t, cborData)

	jsonData, err := comid.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	// Deserialization validation
	var roundtripComid Comid
	err = roundtripComid.FromCBOR(cborData)
	require.NoError(t, err)

	err = roundtripComid.Valid()
	require.NoError(t, err)

	// Verify the structure
	require.NotNil(t, roundtripComid.Triples.DomainMembershipTriples)
	assert.False(t, roundtripComid.Triples.DomainMembershipTriples.IsEmpty())
}

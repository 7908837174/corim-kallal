// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package comid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Example_domainMembershipTriple() {
	// Create a new Comid
	comid := NewComid().
		SetLanguage("en-US").
		SetTagIdentity("domain-membership-example", 1).
		AddEntity("ACME Corp", &TestRegID, RoleCreator, RoleTagCreator)

	// Create a parent domain (e.g., a composite device)
	parentDomain := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("ACME Corp").
			SetModel("Composite Secure Device v2.0"),
		Instance: MustNewUEIDInstance(TestUEID),
	}

	// Create member environments (sub-attesters)
	tpmMember := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("TPM Vendor").
			SetModel("TPM 2.0"),
	}

	secureElementMember := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Secure Element Vendor").
			SetModel("SE v1.0"),
	}

	// Create domain membership triple that describes the topology
	triple := &DomainMembershipTriple{
		DomainID: parentDomain,
		Members:  []Environment{tpmMember, secureElementMember},
	}

	// Add the domain membership triple to the Comid
	comid.AddDomainMembershipTriple(triple)

	// Validate the comid
	err := comid.Valid()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Convert to JSON for demonstration
	jsonData, err := comid.ToJSON()
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		return
	}

	fmt.Printf("Successfully created Comid with DomainMembershipTriple: %d bytes\n", len(jsonData))
	fmt.Println("DomainMembershipTriple includes:")
	fmt.Println("- Parent domain (composite device)")
	fmt.Println("- Two member environments (TPM and Secure Element)")

	// Output:
	// Successfully created Comid with DomainMembershipTriple: 678 bytes
	// DomainMembershipTriple includes:
	// - Parent domain (composite device)
	// - Two member environments (TPM and Secure Element)
}

func Example_domainMembershipTriple_hierarchical() {
	// Create a new Comid for hierarchical device topology
	comid := NewComid().
		SetLanguage("en-US").
		SetTagIdentity("hierarchical-domain-example", 1).
		AddEntity("Enterprise Solutions Ltd", &TestRegID, RoleCreator, RoleTagCreator)

	// Top-level system domain
	systemDomain := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Enterprise Solutions Ltd").
			SetModel("Enterprise Server Rack"),
		Instance: MustNewUEIDInstance(TestUEID),
	}

	// Server nodes in the rack
	server1 := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Server Vendor").
			SetModel("Blade Server 1"),
	}

	server2 := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Server Vendor").
			SetModel("Blade Server 2"),
	}

	networkSwitch := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Network Vendor").
			SetModel("Rack Switch"),
	}

	// Create domain membership triple for the rack topology
	triple := &DomainMembershipTriple{
		DomainID: systemDomain,
		Members:  []Environment{server1, server2, networkSwitch},
	}

	comid.AddDomainMembershipTriple(triple)

	err := comid.Valid()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Successfully created hierarchical domain membership")
	fmt.Println("System Domain: Enterprise Server Rack")
	fmt.Println("Members: 2 blade servers + 1 network switch")

	// Output:
	// Successfully created hierarchical domain membership
	// System Domain: Enterprise Server Rack
	// Members: 2 blade servers + 1 network switch
}

func TestExample_domainMembershipTriple(t *testing.T) {
	// This test ensures the example compiles and runs without error
	// The actual output verification is done by the example test framework
	comid := NewComid().
		SetLanguage("en-US").
		SetTagIdentity("test-domain-membership", 1).
		AddEntity("Test Corp", &TestRegID, RoleCreator, RoleTagCreator)

	parentDomain := Environment{
		Class: NewClassUUID(TestUUID).
			SetVendor("Test Corp").
			SetModel("Test Device"),
	}

	member := Environment{
		Instance: MustNewUEIDInstance(TestUEID),
	}

	triple := &DomainMembershipTriple{
		DomainID: parentDomain,
		Members:  []Environment{member},
	}

	comid.AddDomainMembershipTriple(triple)

	err := comid.Valid()
	require.NoError(t, err)

	jsonData, err := comid.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonData)
}

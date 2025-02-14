package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "80a95cf9-65b4-4bbe-9645-ce60c00e7572"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "daig0104",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// GetVirtualMachineNics gets a list of Network Interface names for a specifcied Azure Virtual Machine.
	// This function would fail the test if there is an error.
	nic := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)

	//get the virtual machine image
	vmImage := azure.GetVirtualMachineImage(t, vmName, resourceGroupName, subscriptionID)


	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
	//confirm NIC exists 
	assert.NotNil(t, nic)
	//confirm nic connection to vm
	assert.True(t, len(nic) == 1, "There should be one NIC attached to VM")



	
	//confirm that the image exists
	assert.NotNil(t, vmImage, "VM image should be found.")
	//confirm that the it is running an Ubuntu version
	assert.Equal(t, "Canonical", vmImage.Publisher, "VM image publisher should be 'Canonical'")
	assert.Equal(t, "0001-com-ubuntu-server-jammy", vmImage.Offer, "VM image offer should be '0001-com-ubuntu-server-jammy'")
	assert.Equal(t, "22_04-lts-gen2", vmImage.SKU, "VM image SKU should be '22_04-lts-gen2'")


}

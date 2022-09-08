using System.Collections.Generic;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CommonSchema;

namespace Pandora.Definitions.ResourceManager.Compute.Terraform;

public class VirtualMachineScaleSetResourceVirtualMachineScaleSetNetworkConfigurationPropertiesSchema
{

    [HclName("accelerated_networking_enabled")]
    [Optional]
    public bool AcceleratedNetworkingEnabled { get; set; }


    [HclName("delete_option")]
    [Optional]
    [PossibleValuesFromConstant(typeof(v2021_11_01.VirtualMachineScaleSets.DeleteOptionsConstant))]
    public string DeleteOption { get; set; }


    [HclName("dns_settings")]
    [Optional]
    public VirtualMachineScaleSetResourceVirtualMachineScaleSetNetworkConfigurationDnsSettingsSchema DnsSettings { get; set; }


    [HclName("fpga_enabled")]
    [Optional]
    public bool FpgaEnabled { get; set; }


    [HclName("ip_configuration")]
    [Required]
    public List<VirtualMachineScaleSetResourceVirtualMachineScaleSetIPConfigurationSchema> IPConfiguration { get; set; }


    [HclName("ip_forwarding_enabled")]
    [Optional]
    public bool IPForwardingEnabled { get; set; }


    [HclName("network_security_group")]
    [Optional]
    public VirtualMachineScaleSetResourceSubResourceSchema NetworkSecurityGroup { get; set; }


    [HclName("primary")]
    [Optional]
    public bool Primary { get; set; }

}

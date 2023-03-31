using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.Search.v2022_09_01.Services;

[ConstantType(ConstantTypeAttribute.ConstantType.String)]
internal enum ProvisioningStateConstant
{
    [Description("failed")]
    Failed,

    [Description("provisioning")]
    Provisioning,

    [Description("succeeded")]
    Succeeded,
}

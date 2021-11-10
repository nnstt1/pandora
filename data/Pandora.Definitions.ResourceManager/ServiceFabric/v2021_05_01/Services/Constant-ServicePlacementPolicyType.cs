using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.ServiceFabric.v2021_05_01.Services
{
    [ConstantType(ConstantTypeAttribute.ConstantType.String)]
    internal enum ServicePlacementPolicyTypeConstant
    {
        [Description("InvalidDomain")]
        InvalidDomain,

        [Description("NonPartiallyPlaceService")]
        NonPartiallyPlaceService,

        [Description("PreferredPrimaryDomain")]
        PreferredPrimaryDomain,

        [Description("RequiredDomain")]
        RequiredDomain,

        [Description("RequiredDomainDistribution")]
        RequiredDomainDistribution,
    }
}

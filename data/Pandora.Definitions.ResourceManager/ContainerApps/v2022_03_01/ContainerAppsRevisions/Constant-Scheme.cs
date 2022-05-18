using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.ContainerApps.v2022_03_01.ContainerAppsRevisions;

[ConstantType(ConstantTypeAttribute.ConstantType.String)]
internal enum SchemeConstant
{
    [Description("HTTP")]
    HTTP,

    [Description("HTTPS")]
    HTTPS,
}

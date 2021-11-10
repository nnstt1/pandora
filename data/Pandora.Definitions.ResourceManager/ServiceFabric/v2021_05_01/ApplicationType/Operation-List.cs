using Pandora.Definitions.Attributes;
using Pandora.Definitions.CustomTypes;
using Pandora.Definitions.Interfaces;
using Pandora.Definitions.Operations;
using System;
using System.Collections.Generic;
using System.Net;

namespace Pandora.Definitions.ResourceManager.ServiceFabric.v2021_05_01.ApplicationType
{
    internal class ListOperation : Operations.ListOperation
    {
        public override string? FieldContainingPaginationDetails() => "nextLink";

        public override ResourceID? ResourceId() => new ManagedClusterId();

        public override Type NestedItemType() => typeof(ApplicationTypeResourceModel);

        public override string? UriSuffix() => "/applicationTypes";


    }
}

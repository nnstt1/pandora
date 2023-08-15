using Pandora.Definitions.Attributes;
using Pandora.Definitions.CustomTypes;
using Pandora.Definitions.Interfaces;
using System;
using System.Collections.Generic;
using System.Net;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.CostManagement.v2023_08_01.CostDetails;

internal class GenerateCostDetailsReportCreateOperationOperation : Pandora.Definitions.Operations.PostOperation
{
    public override IEnumerable<HttpStatusCode> ExpectedStatusCodes() => new List<HttpStatusCode>
        {
                HttpStatusCode.Accepted,
                HttpStatusCode.NoContent,
                HttpStatusCode.OK,
        };

    public override bool LongRunning() => true;

    public override Type? RequestObject() => typeof(GenerateCostDetailsReportRequestDefinitionModel);

    public override ResourceID? ResourceId() => new ScopeId();

    public override Type? ResponseObject() => typeof(CostDetailsOperationResultsModel);

    public override string? UriSuffix() => "/providers/Microsoft.CostManagement/generateCostDetailsReport";


}

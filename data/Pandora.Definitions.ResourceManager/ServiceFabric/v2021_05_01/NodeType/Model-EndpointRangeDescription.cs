using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;

namespace Pandora.Definitions.ResourceManager.ServiceFabric.v2021_05_01.NodeType
{

    internal class EndpointRangeDescriptionModel
    {
        [JsonPropertyName("endPort")]
        [Required]
        public int EndPort { get; set; }

        [JsonPropertyName("startPort")]
        [Required]
        public int StartPort { get; set; }
    }
}

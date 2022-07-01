using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.Media.v2020_05_01.Media;


internal class CommonEncryptionCencModel
{
    [JsonPropertyName("clearTracks")]
    public List<TrackSelectionModel>? ClearTracks { get; set; }

    [JsonPropertyName("contentKeys")]
    public StreamingPolicyContentKeysModel? ContentKeys { get; set; }

    [JsonPropertyName("drm")]
    public CencDrmConfigurationModel? Drm { get; set; }

    [JsonPropertyName("enabledProtocols")]
    public EnabledProtocolsModel? EnabledProtocols { get; set; }
}

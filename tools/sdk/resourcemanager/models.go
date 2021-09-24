package resourcemanager

type ApiObjectDefinition struct {
	NestedItem    *ApiObjectDefinition    `json:"nestedItem,omitempty"`
	ReferenceName *string                 `json:"referenceName,omitempty"`
	Type          ApiObjectDefinitionType `json:"type"`
}

type ApiObjectDefinitionType string

const (
	BooleanApiObjectDefinitionType    ApiObjectDefinitionType = "Boolean"
	CsvApiObjectDefinitionType        ApiObjectDefinitionType = "Csv"
	DateTimeApiObjectDefinitionType   ApiObjectDefinitionType = "DateTime"
	DictionaryApiObjectDefinitionType ApiObjectDefinitionType = "Dictionary"
	IntegerApiObjectDefinitionType    ApiObjectDefinitionType = "Integer"
	FloatApiObjectDefinitionType      ApiObjectDefinitionType = "Float"
	ListApiObjectDefinitionType       ApiObjectDefinitionType = "List"
	RawFileApiObjectDefinitionType    ApiObjectDefinitionType = "RawFile"
	RawObjectApiObjectDefinitionType  ApiObjectDefinitionType = "RawObject"
	ReferenceApiObjectDefinitionType  ApiObjectDefinitionType = "Reference"
	StringApiObjectDefinitionType     ApiObjectDefinitionType = "String"

	// Custom Types
	LocationApiObjectDefinitionType                       ApiObjectDefinitionType = "Location"
	SystemAssignedIdentityApiObjectDefinitionType         ApiObjectDefinitionType = "SystemAssignedIdentity"
	SystemUserAssignedIdentityListApiObjectDefinitionType ApiObjectDefinitionType = "SystemUserAssignedIdentityList"
	SystemUserAssignedIdentityMapApiObjectDefinitionType  ApiObjectDefinitionType = "SystemUserAssignedIdentityMap"
	UserAssignedIdentityListApiObjectDefinitionType       ApiObjectDefinitionType = "UserAssignedIdentityList"
	UserAssignedIdentityMapApiObjectDefinitionType        ApiObjectDefinitionType = "UserAssignedIdentityMap"
	TagsApiObjectDefinitionType                           ApiObjectDefinitionType = "Tags"
)

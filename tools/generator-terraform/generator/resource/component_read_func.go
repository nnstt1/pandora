package resource

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"

	"github.com/hashicorp/pandora/tools/generator-terraform/generator/models"
)

type readFunctionComponents struct {
	resourceId     resourcemanager.ResourceIdDefinition
	terraformModel resourcemanager.TerraformSchemaModelDefinition
}

func readFunctionForResource(input models.ResourceInput) string {
	if !input.Details.ReadMethod.Generate {
		return ""
	}

	idParseLine, err := input.ParseResourceIdFuncName()
	if err != nil {
		// TODO: thread through errors
		panic(err)
	}

	readOperation, ok := input.Operations[input.Details.ReadMethod.MethodName]
	if !ok {
		// TODO: thread through errors
		panic(fmt.Sprintf("couldn't find read operation named %q", input.Details.ReadMethod.MethodName))
	}

	methodArguments := argumentsForApiOperationMethod(readOperation, input.SdkResourceName, input.Details.ReadMethod.MethodName, true)

	terraformModel, ok := input.SchemaModels[input.SchemaModelName]
	if !ok {
		// TODO: real errors
		panic(fmt.Errorf("the Schema Model named %q was not found", input.SchemaModelName))
	}

	resourceId, ok := input.ResourceIds[input.Details.ResourceIdName]
	if !ok {
		// TODO: real errors
		panic(fmt.Errorf("the Resource ID named %q was not found", input.Details.ResourceIdName))
	}

	components := readFunctionComponents{
		resourceId:     resourceId,
		terraformModel: terraformModel,
	}
	lines := make([]string, 0)

	// first map all of the Resource ID segments across
	resourceIdMappingLines, err := components.codeForResourceIdMappings()
	if err != nil {
		// TODO: real errors
		panic(fmt.Errorf("building code for Resource ID Mappings: %+v", err))
	}
	lines = append(lines, *resourceIdMappingLines)

	// then map each of the values across

	//readObjectName, err := readOperation.ResponseObject.GolangTypeName(&input.SdkResourceName)
	//if err != nil {
	//	// TODO: thread real errors through
	//	panic(fmt.Errorf("building golang type name for read operation: %+v", err))
	//}

	// TODO: marshal method

	return fmt.Sprintf(`
func (r %[1]sResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: %[2]d * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.%[3]s.%[1]sClient

			id, err := %[4]s(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.%[5]s(%[6]s)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %%s: %%+v", *id, err)
			}

			var schema %[7]s

			if model := resp.Model; model != nil {
				%[8]s
			}

			return nil
		},
	}
}
`, input.ResourceTypeName, input.Details.ReadMethod.TimeoutInMinutes, input.ServiceName, *idParseLine, input.Details.ReadMethod.MethodName, methodArguments, input.SchemaModelName, strings.Join(lines, "\n"))
}

func (c readFunctionComponents) codeForResourceIdMappings() (*string, error) {
	lines := make([]string, 0)

	// TODO: note that when there's a parent ID field present we'll need to call `parent.NewParentID(..).ID()`
	// to get the right URI
	for _, v := range c.resourceId.Segments {
		if v.Type == resourcemanager.ResourceProviderSegment || v.Type == resourcemanager.StaticSegment {
			continue
		}
		// Subscription ID is a special-case that isn't output
		if v.Type == resourcemanager.SubscriptionIdSegment {
			continue
		}

		topLevelFieldForResourceIdSegment, err := findTopLevelFieldForResourceIdSegment(v.Name, c.terraformModel)
		if err != nil {
			// TODO: error handling
			panic(fmt.Errorf("finding mapping for resource id segment %q: %+v", v.Name, err))
		}

		lines = append(lines, fmt.Sprintf("schema.%s = id.%s", *topLevelFieldForResourceIdSegment, strings.Title(v.Name)))
	}
	sort.Strings(lines)

	output := strings.Join(lines, "\n")
	return &output, nil
}

package pipeline

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type Resource struct {
	Name       string
	Category   string
	Version    string
	Service    string
	Id         *ResourceId
	Paths      []ResourceId
	Operations []Operation
}

type Operation struct {
	Name         string
	Type         OperationType
	Method       string
	UriSuffix    *string
	RequestModel *string
	Responses    []Response
	Tags         []string
}

type Response struct {
	Status      int
	ContentType *string
	Collection  bool
	ModelName   *string
}

type OperationType uint8

const (
	OperationTypeUnknown OperationType = iota
	OperationTypeList
	OperationTypeRead
	OperationTypeCreate
	OperationTypeCreateUpdate
	OperationTypeUpdate
	OperationTypeDelete
)

func NewOperationType(method string) OperationType {
	switch method {
	case http.MethodGet:
		return OperationTypeRead
	case http.MethodPost:
		return OperationTypeCreate
	case http.MethodPatch:
		return OperationTypeUpdate
	case http.MethodPut:
		return OperationTypeCreateUpdate
	case http.MethodDelete:
		return OperationTypeDelete
	}
	return OperationTypeUnknown
}

func (pipelineTask) parseResourcesForService(logger hclog.Logger, apiVersion, service string, serviceTags []string, paths openapi3.Paths, resourceIds ResourceIds, models Models) (resources map[string]*Resource) {
	resources = make(map[string]*Resource)
	for p, item := range paths {
		path := strings.Clone(p)
		operations := item.Operations()
		operationTags := make([]string, 0)

		// Check tags and skip
		skip := true
		for _, operation := range operations {
			if tagMatches(service, operation.Tags) {
				operationTags = append(operationTags, operation.Tags...)
				skip = false
			}
		}
		if skip {
			continue
		}

		parsedPath := NewResourceId(path, operationTags)
		lastSegment := parsedPath.Segments[len(parsedPath.Segments)-1]
		if lastSegment.Type == SegmentCast || lastSegment.Type == SegmentFunction {
			continue
		}

		resourceName := ""
		if r, ok := parsedPath.FindResourceName(); ok {
			resourceName = *r
		}

		// TODO: skip unknown operations for now
		if resourceName == "" {
			continue
		}

		// Resources by default go into their own category when their final URI segment is a label
		resourceCategory := ""
		if lastSegment.Type == SegmentLabel || lastSegment.Type == SegmentUserValue {
			resourceCategory = strings.Clone(resourceName)
		}

		// Determine resource ID and/or URI suffix
		var resourceId *ResourceId
		var uriSuffix *string
		match, ok := resourceIds.MatchIdOrParent(parsedPath)
		if ok {
			if match.Id != nil {
				resourceId = match.Id
			}
			if match.Remainder != nil && len(match.Remainder.Segments) > 0 {
				uriSuffix = pointerTo(match.Remainder.ID())

				// When last segment is not a label (e.g. an action, function or cast), adopt the parent resource category,
				// but only if the suffix has one segment, else this could indicate a different parent, in which case
				// we'll attempt a match after parsing all resources.
				if resourceCategory == "" && strings.Count(*uriSuffix, "/") == 1 {
					resourceCategory = resourceId.Name
				}
			}
		} else {
			uriSuffix = &path
		}

		// Create a new resource if not already encountered
		if _, ok := resources[resourceName]; !ok {
			logger.Info(fmt.Sprintf("found new resource %q for service %q in API version %q", resourceName, service, apiVersion))

			resources[resourceName] = &Resource{
				Name:       resourceName,
				Category:   resourceCategory,
				Version:    apiVersion,
				Service:    cleanName(service),
				Id:         resourceId,
				Paths:      []ResourceId{parsedPath},
				Operations: make([]Operation, 0, len(operations)),
			}
		} else {
			resources[resourceName].Paths = append(resources[resourceName].Paths, parsedPath)
		}

		for method, operation := range operations {
			if !tagMatches(service, operation.Tags) {
				continue
			}
			listOperation := false
			responses := make([]Response, 0)
			if operation.Responses != nil {
				for stat, resp := range operation.Responses {
					var status int
					var contentType, responseModel *string
					var collection bool
					// TODO: expand status codes so we handle more than 200, 300 etc
					if s, err := strconv.Atoi(strings.ReplaceAll(stat, "X", "0")); err == nil {
						status = s
					}
					if resp.Value != nil && resp.Value.Content != nil {
						for t, m := range resp.Value.Content {
							contentType = &t
							if s := parseSchemaRef(m.Schema); s != nil {
								f, _ := flattenSchema(s, nil)
								if f.Title != "" {
									if strings.HasPrefix(strings.ToLower(f.Title), "collection of ") {
										f.Title = f.Title[14:]
										collection = true
										listOperation = true
									}
									if f.Title != "string" {
										if modelName := cleanName(f.Title); models.Found(modelName) {
											responseModel = &modelName
										}
									}
								}
							}
							break
						}
					}
					responses = append(responses, Response{
						Status:      status,
						ContentType: contentType,
						Collection:  collection,
						ModelName:   responseModel,
					})
				}
			}

			operationType := NewOperationType(method)
			if listOperation {
				operationType = OperationTypeList
			}

			operationName := ""
			lastSegment := parsedPath.Segments[len(parsedPath.Segments)-1]
			shortResourceName := strings.TrimPrefix(resourceName, singularize(cleanName(service)))

			switch operationType {
			case OperationTypeList:
				if _, ok = verbs.match(shortResourceName); ok {
					operationName = resourceName
				} else {
					operationName = fmt.Sprintf("List%s", pluralize(resourceName))
				}
			case OperationTypeRead:
				operationName = fmt.Sprintf("Get%s", resourceName)
			case OperationTypeCreate:
				if _, ok = verbs.match(shortResourceName); ok {
					operationName = resourceName
				} else if lastSegment.Type == SegmentODataReference {
					operationName = fmt.Sprintf("Add%s", singularize(resourceName))
				} else {
					operationName = fmt.Sprintf("Create%s", singularize(resourceName))
				}
			case OperationTypeCreateUpdate:
				operationName = fmt.Sprintf("CreateUpdate%s", singularize(resourceName))
			case OperationTypeUpdate:
				operationName = fmt.Sprintf("Update%s", singularize(resourceName))
			case OperationTypeDelete:
				if lastSegment.Type == SegmentODataReference {
					operationName = fmt.Sprintf("Remove%s", resourceName)
				} else {
					operationName = fmt.Sprintf("Delete%s", resourceName)
				}
			}

			// Trim the "Ref" suffix from operation names
			operationName = strings.TrimSuffix(operationName, "Ref")

			var requestModel *string
			if operation.RequestBody != nil && operation.RequestBody.Value != nil {
				for contentType, content := range operation.RequestBody.Value.Content {
					if strings.HasPrefix(strings.ToLower(contentType), "application/json") {
						if content.Schema != nil {
							schema, _ := flattenSchema(content.Schema.Value, nil)
							if schema.Title != "" {
								if modelName := cleanName(schema.Title); models.Found(modelName) {
									requestModel = &modelName
								}
							}
						}
					}
				}
			}

			resources[resourceName].Operations = append(resources[resourceName].Operations, Operation{
				Name:         operationName,
				Type:         operationType,
				Method:       method,
				UriSuffix:    uriSuffix,
				RequestModel: requestModel,
				Responses:    responses,
				Tags:         operation.Tags,
			})
		}
	}

	// Look for resources without a category, then iterate the known paths of it and all potential parent resources
	// to find a match by trimming the path to the preceding label segment, then adopt the parent resource category to ensure they are grouped together
	for _, resource := range resources {
		if pathsLen := len(resource.Paths); resource.Category == "" && pathsLen > 0 {
			for _, path := range resource.Paths {
				if trimmedPath := path.TruncateToLastSegmentOfTypeBeforeSegment([]ResourceIdSegmentType{SegmentLabel}, -1); trimmedPath != nil {
					for _, parentResource := range resources {
						for _, parentPath := range parentResource.Paths {
							if parentPath.ID() == trimmedPath.ID() {
								resource.Category = parentResource.Category
								break
							}
						}
					}
				}
			}
		}
	}

	return
}

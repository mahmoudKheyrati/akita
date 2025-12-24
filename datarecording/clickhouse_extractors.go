package datarecording

import (
	"fmt"
	"reflect"
)

// These extractors use minimal reflection ONLY as a fallback
// The fast path uses type assertions in the convert functions

func extractTaskTableEntry(entry any) taskTableEntryDB {
	v := reflect.ValueOf(entry)
	t := v.Type()

	// Validate this is the right structure
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected struct for task entry, got %T", entry))
	}

	result := taskTableEntryDB{}

	// Extract fields by index (MUCH faster than FieldByName)
	// Expected fields: ID, ParentID, Kind, What, Location, StartTime, EndTime
	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		switch field.Name {
		case "ID":
			result.ID = fieldValue.String()
		case "ParentID":
			result.ParentID = fieldValue.String()
		case "Kind":
			result.Kind = fieldValue.String()
		case "What":
			result.What = fieldValue.String()
		case "Location":
			result.Location = fieldValue.String()
		case "StartTime":
			result.StartTime = fieldValue.Float()
		case "EndTime":
			result.EndTime = fieldValue.Float()
		}
	}

	return result
}

func extractMilestoneTableEntry(entry any) milestoneTableEntryDB {
	v := reflect.ValueOf(entry)
	t := v.Type()

	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected struct for milestone entry, got %T", entry))
	}

	result := milestoneTableEntryDB{}

	// Extract fields by index for better performance
	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		switch field.Name {
		case "ID":
			result.ID = fieldValue.String()
		case "TaskID":
			result.TaskID = fieldValue.String()
		case "Time":
			result.Time = fieldValue.Float()
		case "Kind":
			result.Kind = fieldValue.String()
		case "What":
			result.What = fieldValue.String()
		case "Location":
			result.Location = fieldValue.String()
		}
	}

	return result
}

func extractSegmentTableEntry(entry any) segmentTableEntryDB {
	v := reflect.ValueOf(entry)
	t := v.Type()

	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected struct for segment entry, got %T", entry))
	}

	result := segmentTableEntryDB{}

	// Extract fields by index for better performance
	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		switch field.Name {
		case "StartTime":
			result.StartTime = fieldValue.Float()
		case "EndTime":
			result.EndTime = fieldValue.Float()
		}
	}

	return result
}

func extractMemoryTransactionEntry(entry any) memoryTransactionEntryDB {
	v := reflect.ValueOf(entry)
	t := v.Type()

	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected struct for memory transaction entry, got %T", entry))
	}

	result := memoryTransactionEntryDB{}

	// Extract fields by index for better performance
	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		switch field.Name {
		case "ID":
			result.ID = fieldValue.String()
		case "Location":
			result.Location = fieldValue.String()
		case "What":
			result.What = fieldValue.String()
		case "StartTime":
			result.StartTime = fieldValue.Float()
		case "EndTime":
			result.EndTime = fieldValue.Float()
		case "Address":
			result.Address = fieldValue.Uint()
		case "ByteSize":
			result.ByteSize = fieldValue.Uint()
		}
	}

	return result
}

func extractMemoryStepEntry(entry any) memoryStepEntryDB {
	v := reflect.ValueOf(entry)
	t := v.Type()

	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected struct for memory step entry, got %T", entry))
	}

	result := memoryStepEntryDB{}

	// Extract fields by index for better performance
	numFields := v.NumField()
	for i := 0; i < numFields; i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		switch field.Name {
		case "ID":
			result.ID = fieldValue.String()
		case "TaskID":
			result.TaskID = fieldValue.String()
		case "Time":
			result.Time = fieldValue.Float()
		case "What":
			result.What = fieldValue.String()
		}
	}

	return result
}

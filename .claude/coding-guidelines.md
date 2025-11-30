# Coding Standards: Avoid Magic Values and Hardcoded Strings

## Proto Message Names
- Use `proto.MessageName((*Type)(nil))` instead of hardcoded proto full names
- Ensures compile-time safety when proto definitions change

## Enum Values
- Use proto-generated enum constants instead of string literals
- Example: `ffv1.ClusterState_CLUSTER_STATE_READY` not `"CLUSTER_STATE_READY"`

## Magic Numbers
- Define named constants for all numeric values
- Use descriptive names that explain the purpose

## Object Type Comparisons
- Use proto descriptors (`helper.Descriptor()`) when available
- Avoid string-based type switches

## Testing Framework
- Use Ginkgo for all unit tests, not standard Go testing
- Use table-driven tests with `DescribeTable` and `Entry`
- Use proto message names in test data, not hardcoded strings

## When Hardcoded Values Are Acceptable
- User-facing display messages
- Proto field names from generated code (when unavoidable)

## Pre-Commit Checks
- Run code formatting before committing changes
- Run linter checks to ensure code quality
- Ensure all tests pass before committing

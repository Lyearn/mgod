package mongomodel

import "reflect"

// Used only for debugging in tests.
// Reason: reflect.DeepEqual doesn't provide enough information about the difference.
// Replace reflect.DeepEqual with this function in test file and start a debugger to get more information.
func (s *EntityModelSchema) Compare(c *EntityModelSchema) bool {
	// schema tree
	sTree := &s.Root.Children

	// to compare schema tree
	cTree := &c.Root.Children

	return s.compare(sTree, cTree)
}

func (s *EntityModelSchema) compare(sTree, cTree *[]TreeNode) bool {
	if len(*sTree) != len(*cTree) {
		return false
	}

	for i := 0; i < len(*sTree); i++ {
		sNode := (*sTree)[i]
		cNode := (*cTree)[i]

		if sNode.Path != cNode.Path || sNode.BSONKey != cNode.BSONKey || sNode.Key != cNode.Key {
			return false
		}

		if !reflect.DeepEqual(sNode.Props, cNode.Props) {
			return false
		}

		if !s.compare(&sNode.Children, &cNode.Children) {
			return false
		}
	}

	return true
}

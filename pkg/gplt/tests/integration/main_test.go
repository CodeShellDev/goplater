package tests

import (
	"github.com/codeshelldev/goplater/pkg/gplt/collections"
	"github.com/codeshelldev/goplater/pkg/gplt/utils/testutils"
)

var tester = testutils.NewTester(collections.Unsafe, collections.Safe)
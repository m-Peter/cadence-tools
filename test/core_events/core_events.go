/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright 2019-2022 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core_events

import (
	_ "embed"
	"fmt"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/cadence/runtime/sema"
	"github.com/onflow/cadence/runtime/stdlib"
)

//go:embed core_events.cdc
var CoreEvents []byte

const CoreEventsLocation = common.IdentifierLocation("CoreEvents")

func CoreEventsChecker() *sema.Checker {
	program, err := parser.ParseProgram(
		nil,
		CoreEvents,
		parser.Config{},
	)
	if err != nil {
		panic(err)
	}

	importHandler := func(
		checker *sema.Checker,
		importedLocation common.Location,
		importRange ast.Range,
	) (sema.Import, error) {
		var elaboration *sema.Elaboration
		switch importedLocation {
		case stdlib.TestContractLocation:
			testChecker := stdlib.GetTestContractType().Checker
			elaboration = testChecker.Elaboration
		default:
			return nil, fmt.Errorf("import not supported")
		}

		return sema.ElaborationImport{
			Elaboration: elaboration,
		}, nil
	}

	activation := sema.NewVariableActivation(sema.BaseValueActivation)
	activation.DeclareValue(stdlib.AssertFunction)
	activation.DeclareValue(stdlib.PanicFunction)

	checker, err := sema.NewChecker(
		program,
		CoreEventsLocation,
		nil,
		&sema.Config{
			BaseValueActivationHandler: func(_ common.Location) *sema.VariableActivation {
				return activation
			},
			AccessCheckMode: sema.AccessCheckModeStrict,
			ImportHandler:   importHandler,
		},
	)
	if err != nil {
		panic(err)
	}

	err = checker.Check()
	if err != nil {
		panic(err)
	}

	return checker
}

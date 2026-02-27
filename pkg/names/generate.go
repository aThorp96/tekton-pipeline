/*
Copyright 2019 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package names

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"

	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

// NameGenerator generates names for objects. Some backends may have more information
// available to guide selection of new names and this interface hides those details.
type NameGenerator interface {
	// RestrictLengthWithRandomSuffix generates a valid name from the base name, adding a random suffix to
	// the base. If base is valid, the returned name must also be valid. The generator is
	// responsible for knowing the maximum valid name length.
	RestrictLengthWithRandomSuffix(base string) string

	// RestrictLength generates a valid name from the name of a step specified in a Task,
	// shortening it to the maximum valid name length if needed.
	RestrictLength(base string) string
}

const (
	// TODO: make this flexible for non-core resources with alternate naming rules.
	maxNameLength             = 63
	defaultRandomSuffixLength = 5
)

// simpleNameGenerator generates random names.
type simpleNameGenerator struct {
	randomSuffixLength int
}

// SimpleNameGenerator is a generator that returns the name plus a random suffix of five alphanumerics
// when a name is requested. The string is guaranteed to not exceed the length of a standard Kubernetes
// name (63 characters)
var SimpleNameGenerator NameGenerator = simpleNameGenerator{randomSuffixLength: defaultRandomSuffixLength}

func NewSimpleNameGenerator(randomSuffixLength int) NameGenerator {
	return simpleNameGenerator{randomSuffixLength}
}

func (sng simpleNameGenerator) maxGeneratedNameLength() int {
	return maxNameLength - sng.randomSuffixLength - 1
}

// RestrictLengthWithRandomSuffix takes a base name and returns a potentially shortened version of that name with
// a random suffix, with the whole string no longer than 63 characters.
func (sng simpleNameGenerator) RestrictLengthWithRandomSuffix(base string) string {
	if len(base) > sng.maxGeneratedNameLength() {
		base = base[:sng.maxGeneratedNameLength()]
	}
	return fmt.Sprintf("%s-%s", base, utilrand.String(sng.randomSuffixLength))
}

var alphaNumericRE = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// RestrictLength takes a base name and returns a potentially shortened version of that name, no longer than 63 characters.
func (simpleNameGenerator) RestrictLength(base string) string {
	if len(base) > maxNameLength {
		base = base[:maxNameLength]
	}

	for !alphaNumericRE.MatchString(base[len(base)-1:]) {
		base = base[:len(base)-1]
	}
	return base
}

// GenerateHashedName creates a unique name with a hashed suffix.
func GenerateHashedName(prefix, name string, hashedLength int) string {
	if hashedLength <= 0 {
		hashedLength = defaultRandomSuffixLength
	}
	h := fnv.New32a()
	h.Write([]byte(name))
	suffix := strconv.FormatUint(uint64(h.Sum32()), 16)
	if ln := len(suffix); ln > hashedLength {
		suffix = suffix[:hashedLength]
	} else if ln < hashedLength {
		suffix += strings.Repeat("0", hashedLength-ln)
	}
	return fmt.Sprintf("%s-%s", prefix, suffix)
}

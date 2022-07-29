package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
)

type HelperVersion struct {
	Library string
	Tag     semver.Version
}

func NewHelperVersionFrom(image string) (*HelperVersion, error) {
	splitted := strings.Split(image, ":")
	if len(splitted) != 2 {
		return nil, errors.New("image doesn't contain repository and tag")
	}
	library := splitted[0]
	if(len(splitted[1])==0){
		return nil, errors.New("image tag was not found")
	}

	tag, err := semver.NewVersion(splitted[1])
	if err != nil {
		return nil, err
	}
	return &HelperVersion{Library: library, Tag: *tag}, err
}

func (h HelperVersion) Compare(second HelperVersion) int {
	if h.Tag.Major > second.Tag.Major {
		return 1
	} else if h.Tag.Major == second.Tag.Major {
		if h.Tag.Minor > second.Tag.Minor {
			return 1
		} else if h.Tag.Minor == second.Tag.Minor {
			if h.Tag.Patch > second.Tag.Patch {
				return 1
			} else if h.Tag.Patch == second.Tag.Patch {
				return 0
			} else {
				return -1
			}
		} else {
			return -1
		}
	} else {
		return -1
	}
}

func (h HelperVersion) String() string {
	return fmt.Sprintf("%s:%s", h.Library, h.Tag.String())
}

func AreEqual(first, second HelperVersion) bool {
	return first.Library == second.Library && first.Tag.Equal(second.Tag)
}

func (h HelperVersion) EqualTo(second HelperVersion) bool {
	return AreEqual(h, second)
}

func (h HelperVersion) SmallerThan(second HelperVersion) bool {
	return h.Tag.LessThan(second.Tag)
}

func (h HelperVersion) BiggerThan(second HelperVersion) bool {
	return second.Tag.LessThan(h.Tag)
}

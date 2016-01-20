package core

import (
	"fmt"
	"path"
	"strings"
)

type ReleaseMeta struct {
	Path       string
	Mountpoint string
	Name       string
}

func GetAllReleases() []*ReleaseMeta {
	list := make([]*ReleaseMeta, 0)
	out := ZFSMust(
		fmt.Errorf("No releases found"),
		"list", "-H", "-d", "1",
		"-o", "name,mountpoint",
		GetReleasesPath())
	lines := SplitOutput(out)
	// discard first line, as that is the jail dir itself
	for _, line := range lines[1:] {
		list = append(list, &ReleaseMeta{
			Path:       line[0],
			Mountpoint: line[1],
			Name:       path.Base(line[1]),
		})
	}
	return list
}

func FindRelease(lookup string) (*ReleaseMeta, error) {
	out, err := ZFS(
		"list", "-H", "-d", "1",
		"-o", "name,mountpoint",
		GetReleasesPath())
	if err != nil {
		return nil, err
	}
	lines := SplitOutput(out)
	for _, line := range lines {
		name := path.Base(line[1])
		if name == lookup || strings.HasPrefix(name, lookup) {
			return &ReleaseMeta{
				Path:       line[0],
				Mountpoint: line[1],
				Name:       name,
			}, nil
		}
	}
	return nil, fmt.Errorf("No release found")
}

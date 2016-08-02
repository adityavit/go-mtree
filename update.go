package mtree

import (
	"os"
	"sort"
)

var DefaultUpdateKeywords = []string{
	"uid",
	"gid",
	"mode",
	"time",
}

// Update attempts to set the attributes of root directory path, given the values of `keywords` in dh DirectoryHierarchy.
func Update(root string, dh *DirectoryHierarchy, keywords []Keyword) ([]Result, error) {
	creator := dhCreator{DH: dh}
	curDir, err := os.Getwd()
	if err == nil {
		defer os.Chdir(curDir)
	}

	if err := os.Chdir(root); err != nil {
		return nil, err
	}
	sort.Sort(byPos(creator.DH.Entries))

	results := []Result{}
	for i, e := range creator.DH.Entries {
		switch e.Type {
		case SpecialType:
			if e.Name == "/set" {
				creator.curSet = &creator.DH.Entries[i]
			} else if e.Name == "/unset" {
				creator.curSet = nil
			}
			Debugf("%#v", e)
			continue
		case RelativeType, FullType:
			pathname, err := e.Path()
			if err != nil {
				return nil, err
			}
			var toCheck []KeyVal
			if creator.curSet != nil {
				toCheck = append(toCheck, creator.curSet.Keywords...)
			}
			toCheck = append(toCheck, e.Keywords...)

			for _, kv := range toCheck {
				if !InKeywordSlice(kv.Keyword(), keywords) {
					continue
				}
				ukFunc, ok := UpdateKeywordFuncs[kv.Keyword()]
				if !ok {
					Debugf("no UpdateKeywordFunc for %s; skipping", kv.Keyword())
					continue
				}
				if _, err := ukFunc(pathname, kv.Value()); err != nil {
					results = append(results, Result{Path: pathname, Keyword: kv.Keyword(), Got: err.Error()})
				}
			}
		}
	}

	return results, nil
}

type Result struct {
	Path    string
	Keyword Keyword
	Got     string
}

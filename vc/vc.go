package vc

import set "github.com/deckarep/golang-set"

type VersionControl interface{
	//获取差异文件集合
	GetDiffSet() (set.Set, error)
}


package verdeps

type sortableRevisions []*revision

func (sr sortableRevisions) Len() int {
	return len(sr)
}

func (sr sortableRevisions) Less(i, j int) bool {
	return sr[i].toIndex > sr[j].fromIndex
}

func (sr sortableRevisions) Swap(i, j int) {
	oldI := sr[i]
	sr[i] = sr[j]
	sr[j] = oldI
}

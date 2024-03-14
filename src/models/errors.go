package models

type EmptyVersionListError struct{}

func (e *EmptyVersionListError) Error() string {
	return "error: version list is empty"
}

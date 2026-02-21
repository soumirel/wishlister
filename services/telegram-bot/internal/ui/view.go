package ui

type ViewIdentifier string

type View interface {
	ViewIdentifier() ViewIdentifier
	ViewLifecycle
}

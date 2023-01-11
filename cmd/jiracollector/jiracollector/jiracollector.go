package jiracollector

type JiraCollector interface {
	Execute(project string, jql string) error
}

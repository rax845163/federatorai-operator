package server

import (
	"fmt"

	"sync"
	"time"

	"github.com/containers-ai/alameda/admission-controller/pkg/recommendator/resource"
)

type recommendationState int

const (
	recommendationIsSynchronized     recommendationState = 0
	recommendationNeedsSynchronizing recommendationState = 1
	recommendationWaitsSynchronizing recommendationState = 2
)

type controllerRecommendation struct {
	lock               *sync.Mutex
	syncChan           chan bool
	state              recommendationState
	podRecommendations []*resource.PodResourceRecommendation
}

func NewControllerPodResourceRecommendation() *controllerRecommendation {
	return &controllerRecommendation{
		lock:               &sync.Mutex{},
		syncChan:           make(chan bool),
		state:              recommendationIsSynchronized,
		podRecommendations: make([]*resource.PodResourceRecommendation, 0),
	}
}

func (c *controllerRecommendation) waitOrSync() recommendationState {
	c.lock.Lock()
	s := c.state
	if c.state != recommendationWaitsSynchronizing {
		c.syncChan = make(chan bool)
		c.state = recommendationWaitsSynchronizing
	}
	c.lock.Unlock()
	return s
}

func (c *controllerRecommendation) renewSyncChan() {
	c.lock.Lock()
	close(c.syncChan)
	c.syncChan = make(chan bool)
	c.lock.Unlock()
}

func (c *controllerRecommendation) finishSync() {
	c.lock.Lock()
	c.state = recommendationIsSynchronized
	close(c.syncChan)
	c.lock.Unlock()
}

func (c *controllerRecommendation) setState(state recommendationState) {
	c.lock.Lock()
	c.state = state
	c.lock.Unlock()
}

func (c *controllerRecommendation) appendPodRecommendations(recommendations []*resource.PodResourceRecommendation) {
	c.lock.Lock()
	c.podRecommendations = append(c.podRecommendations, recommendations...)
	c.lock.Unlock()
}

func (c *controllerRecommendation) setPodRecommendations(recommendations []*resource.PodResourceRecommendation) {
	c.lock.Lock()
	c.podRecommendations = recommendations
	c.lock.Unlock()
}

func (c *controllerRecommendation) getPodRecommendations() []*resource.PodResourceRecommendation {
	c.lock.Lock()
	recommendations := c.podRecommendations
	c.lock.Unlock()
	return recommendations
}

// DispatchOneValidRecommendation dispatch one recommendation that timestamp is in time interval between recommendation.ValidStartTime and recommendation.ValidEndTime,
// return nil if no valid recommendation can provide
func (c *controllerRecommendation) dispatchOneValidPodRecommendation(timestamp time.Time) *resource.PodResourceRecommendation {

	var recommendation *resource.PodResourceRecommendation

	if len(c.podRecommendations) == 0 {
		return nil
	}

	c.lock.Lock()
	var indexOfValidRecommendation *int
	for i, recommendation := range c.podRecommendations {
		if recommendation.ValidStartTime.UnixNano() <= timestamp.UnixNano() && timestamp.UnixNano() <= recommendation.ValidEndTime.UnixNano() {
			indexOfValidRecommendation = &i
			break
		}
	}
	if indexOfValidRecommendation != nil {
		recommendation = c.podRecommendations[*indexOfValidRecommendation]
		c.podRecommendations = c.podRecommendations[*indexOfValidRecommendation+1:]
	}
	c.lock.Unlock()

	return recommendation
}

type namespaceKindName struct {
	namespace, kind, name string
}

func newNamespaceKindName(namespace, kind, name string) namespaceKindName {
	return namespaceKindName{
		namespace: namespace,
		name:      name,
		kind:      kind,
	}
}

func (n namespaceKindName) String() string {
	return fmt.Sprintf("%s.%s.%s", n.namespace, n.kind, n.name)
}

func (n namespaceKindName) getNamespace() string {
	return n.namespace
}

func (n namespaceKindName) getKind() string {
	return n.kind
}

func (n namespaceKindName) getName() string {
	return n.name
}

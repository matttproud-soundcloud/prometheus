// Copyright 2013 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retrieval

import (
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/model"
	"github.com/prometheus/prometheus/retrieval/format"
	"log"
	"time"
)

type TargetManager interface {
	acquire()
	release()
	AddTarget(job config.JobConfig, t Target, defaultScrapeInterval time.Duration)
	ReplaceTargets(job config.JobConfig, newTargets []Target, defaultScrapeInterval time.Duration)
	Remove(t Target)
	AddTargetsFromConfig(config config.Config)
	Pools() map[string]*TargetPool
}

type targetManager struct {
	requestAllowance chan bool
	poolsByJob       map[string]*TargetPool
	results          chan format.Result
}

func NewTargetManager(results chan format.Result, requestAllowance int) TargetManager {
	return &targetManager{
		requestAllowance: make(chan bool, requestAllowance),
		results:          results,
		poolsByJob:       make(map[string]*TargetPool),
	}
}

func (m *targetManager) acquire() {
	m.requestAllowance <- true
}

func (m *targetManager) release() {
	<-m.requestAllowance
}

func (m *targetManager) TargetPoolForJob(job config.JobConfig, defaultScrapeInterval time.Duration) (targetPool *TargetPool) {
	targetPool, ok := m.poolsByJob[job.GetName()]

	if !ok {
		targetPool = NewTargetPool(m)
		log.Printf("Pool for job %s does not exist; creating and starting...", job.GetName())

		interval := job.ScrapeInterval()
		m.poolsByJob[job.GetName()] = targetPool
		go targetPool.Run(m.results, interval)
	}
	return
}

func (m *targetManager) AddTarget(job config.JobConfig, t Target, defaultScrapeInterval time.Duration) {
	targetPool := m.TargetPoolForJob(job, defaultScrapeInterval)
	targetPool.AddTarget(t)
	m.poolsByJob[job.GetName()] = targetPool
}

func (m *targetManager) ReplaceTargets(job config.JobConfig, newTargets []Target, defaultScrapeInterval time.Duration) {
	targetPool := m.TargetPoolForJob(job, defaultScrapeInterval)
	targetPool.replaceTargets(newTargets)
}

func (m targetManager) Remove(t Target) {
	panic("not implemented")
}

func (m *targetManager) AddTargetsFromConfig(config config.Config) {
	for _, job := range config.Jobs() {
		for _, targetGroup := range job.TargetGroup {
			baseLabels := model.LabelSet{
				model.JobLabel: model.LabelValue(job.GetName()),
			}
			if targetGroup.Labels != nil {
				for _, label := range targetGroup.Labels.Label {
					baseLabels[model.LabelName(label.GetName())] = model.LabelValue(label.GetValue())
				}
			}

			for _, endpoint := range targetGroup.Target {
				target := NewTarget(endpoint, time.Second*5, baseLabels)
				m.AddTarget(job, target, config.ScrapeInterval())
			}
		}
	}
}

// XXX: Not really thread-safe. Only used in /status page for now.
func (m *targetManager) Pools() map[string]*TargetPool {
	return m.poolsByJob
}

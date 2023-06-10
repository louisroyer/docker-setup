// Copyright 2023 Louis Royer and docker-setup contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package setup

import (
	"fmt"
	"log"
	"os"
)

// Configuration
type Conf struct {
	hooksList map[string]Hook
	oneshot   bool
}

// Create a new configuration from env variables
func NewConf() Conf {
	conf := Conf{
		hooksList: make(map[string]Hook, 0),
	}
	conf.AddHooks()
	conf.AddUserHooks("pre", "PRE")
	conf.AddUserHooks("post", "POST")
	conf.oneshot = false
	if oneshot, isset := os.LookupEnv("ONESHOT"); isset && oneshot == "true" {
		conf.oneshot = true
	}
	return conf
}

// Return true if Oneshot is set
func (conf Conf) Oneshot() bool {
	return conf.oneshot
}

// Run exit hooks
func (conf Conf) RunExitHooks() {
	conf.RunExitHook("pre")
	conf.RunExitHook("nat4")
	conf.RunExitHook("iproute")
	conf.RunExitHook("post")
}

// Run init hooks
func (conf Conf) RunInitHooks() {
	conf.RunInitHook("pre")
	conf.RunInitHook("iproute")
	conf.RunInitHook("nat4")
	conf.RunInitHook("post")
}

// Add a new hook to the configuration
func (conf Conf) AddUserHooks(name string, env string) {
	conf.hooksList[name] = NewUserHooks(
		fmt.Sprintf("%s init", name), fmt.Sprintf("%s_INIT_HOOK", env),
		fmt.Sprintf("%s exit", name), fmt.Sprintf("%s_EXIT_HOOK", env))
}

// Add default hooks
func (conf Conf) AddHooks() {
	conf.hooksList["iproute"] = NewIPRouteHooks(
		"iproute init", "ROUTES_INIT",
		"iproute exit", "ROUTES_EXIT")
	conf.hooksList["nat4"] = NewNat4Hooks()
}

// Run an init hook
func (conf Conf) RunInitHook(name string) {
	if conf.hooksList[name] != nil {
		if err := conf.hooksList[name].RunInit(); err != nil {
			log.Printf("Error while running %s init hook: %s", name, err)
		}
	}
}

// Run an exit hook
func (conf Conf) RunExitHook(name string) {
	if conf.hooksList[name] != nil {
		if err := conf.hooksList[name].RunExit(); err != nil {
			log.Printf("Error while running %s exit hook: %s", name, err)
		}
	}
}

// Log the configuration
func (conf Conf) Log() {
	log.Println("Current configuration:")
	if conf.oneshot {
		log.Printf("\t- mode oneshot is enabled")
	} else {
		log.Printf("\t- mode oneshot is disabled")
	}
	for _, h := range conf.hooksList {
		for _, s := range h.String() {
			log.Printf("\t- %s\n", s)
		}
	}

}

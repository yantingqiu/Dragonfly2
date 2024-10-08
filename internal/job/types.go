/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package job

import "d7y.io/dragonfly/v2/scheduler/resource"

// PreheatRequest defines the request parameters for preheating.
type PreheatRequest struct {
	URL                 string            `json:"url" validate:"required,url"`
	Tag                 string            `json:"tag" validate:"omitempty"`
	Digest              string            `json:"digest" validate:"omitempty"`
	FilteredQueryParams string            `json:"filtered_query_params" validate:"omitempty"`
	Headers             map[string]string `json:"headers" validate:"omitempty"`
	Application         string            `json:"application" validate:"omitempty"`
	Priority            int32             `json:"priority" validate:"omitempty"`
}

// PreheatResponse defines the response parameters for preheating.
type PreheatResponse struct {
	TaskID string `json:"task_id"`
}

// GetTaskRequest defines the request parameters for getting task.
type GetTaskRequest struct {
	TaskID string `json:"task_id" validate:"required"`
}

// GetTaskResponse defines the response parameters for getting task.
type GetTaskResponse struct {
	Peers []*resource.Peer `json:"peers"`
}

// DeleteTaskRequest defines the request parameters for deleting task.
type DeleteTaskRequest struct {
	TaskID string `json:"task_id" validate:"required"`
}

// DeleteTaskResponse defines the response parameters for deleting task.
type DeleteTaskResponse struct {
	SuccessPeers []*DeletePeerResponse `json:"success_peers"`
	FailurePeers []*DeletePeerResponse `json:"failure_peers"`
}

// DeletePeerResponse represents the response after attempting to delete a peer.
type DeletePeerResponse struct {
	Peer        *resource.Peer `json:"peer"`
	Description string         `json:"description"`
}

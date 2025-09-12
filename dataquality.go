// Copyright 2023 Specter Ops, Inc.
//
// Licensed under the Apache License, Version 2.0
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package bloodhound

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ADDataQualityStat struct {
	DomainSID              string    `json:"domain_sid"`
	Users                  int       `json:"users"`
	Groups                 int       `json:"groups"`
	Computers              int       `json:"computers"`
	OUs                    int       `json:"ous"`
	Containers             int       `json:"containers"`
	GPOs                   int       `json:"gpos"`
	AIACAs                 int       `json:"aiacas"`
	RootCAs                int       `json:"rootcas"`
	EnterpriseCAs          int       `json:"enterprisecas"`
	NTAuthStores           int       `json:"ntauthstores"`
	CertTemplates          int       `json:"certtemplates"`
	IssuancePolicies       int       `json:"issuancepolicies"`
	ACLs                   int       `json:"acls"`
	Sessions               int       `json:"sessions"`
	Relationships          int       `json:"relationships"`
	SessionCompleteness    float64   `json:"session_completeness"`
	LocalGroupCompleteness float64   `json:"local_group_completeness"`
	RunID                  string    `json:"run_id"`
	ID                     int       `json:"id"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type AzureDataQualityStat struct {
	TenantID             string    `json:"tenantid"`
	Users                int       `json:"users"`
	Groups               int       `json:"groups"`
	Apps                 int       `json:"apps"`
	ServicePrincipals    int       `json:"service_principals"`
	Devices              int       `json:"devices"`
	ManagementGroups     int       `json:"management_groups"`
	Subscriptions        int       `json:"subscriptions"`
	ResourceGroups       int       `json:"resource_groups"`
	VMs                  int       `json:"vms"`
	KeyVaults            int       `json:"key_vaults"`
	AutomationAccounts   int       `json:"automation_accounts"`
	ContainerRegistries  int       `json:"container_registries"`
	FunctionApps         int       `json:"function_apps"`
	LogicApps            int       `json:"logic_apps"`
	ManagedClusters      int       `json:"managed_clusters"`
	VMScaleSets          int       `json:"vm_scale_sets"`
	WebApps              int       `json:"web_apps"`
	Relationships        int       `json:"relationships"`
	RunID                string    `json:"run_id"`
	ID                   int       `json:"id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (s *Client) GetADDataQualityStats(domainID string) (ADDataQualityStat, error) {
	var stats ADDataQualityStat
	url := s.baseURL.JoinPath("api/v2/ad-domains", domainID, "data-quality-stats")

	if req, err := s.newAuthenticatedRequest(http.MethodGet, url.String(), nil); err != nil {
		return stats, err
	} else if resp, err := s.do(req, nil); err != nil {
		return stats, err
	} else {
		defer resp.Body.Close()
		var response struct {
			Data []ADDataQualityStat `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return stats, err
		}
		if len(response.Data) > 0 {
			return response.Data[0], nil
		}
		return stats, fmt.Errorf("no AD data quality stats found for domain %s", domainID)
	}
}

func (s *Client) GetAzureDataQualityStats(tenantID string) (AzureDataQualityStat, error) {
	var stats AzureDataQualityStat
	url := s.baseURL.JoinPath("api/v2/azure-tenants", tenantID, "data-quality-stats")

	if req, err := s.newAuthenticatedRequest(http.MethodGet, url.String(), nil); err != nil {
		return stats, err
	} else if resp, err := s.do(req, nil); err != nil {
		return stats, err
	} else {
		defer resp.Body.Close()
		var response struct {
			Data []AzureDataQualityStat `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return stats, err
		}
		if len(response.Data) > 0 {
			return response.Data[0], nil
		}
		return stats, fmt.Errorf("no Azure data quality stats found for tenant %s", tenantID)
	}
}
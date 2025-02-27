package gke

// ATTENTION!
// This rule was autogenerated!
// Before making changes, consider updating the generator.

import (
	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/provider"
	"github.com/aquasecurity/tfsec/pkg/result"
	"github.com/aquasecurity/tfsec/pkg/rule"
	"github.com/aquasecurity/tfsec/pkg/severity"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		Provider:  provider.GoogleProvider,
		Service:   "gke",
		ShortCode: "node-pool-uses-cos",
		Documentation: rule.RuleDocumentation{
			Summary:     "Ensure Container-Optimized OS (cos) is used for Kubernetes Engine Clusters Node image",
			Explanation: `GKE supports several OS image types but COS_CONTAINERD is the recommended OS image to use on cluster nodes for enhanced security`,
			Impact:      "COS_CONTAINERD is the recommended OS image to use on cluster nodes",
			Resolution:  "Use the COS_CONTAINERD image type",
			BadExample: []string{`
resource "google_service_account" "default" {
  account_id   = "service-account-id"
  display_name = "Service Account"
}

resource "google_container_cluster" "primary" {
  name     = "my-gke-cluster"
  location = "us-central1"

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "bad_example" {
  name       = "my-node-pool"
  cluster    = google_container_cluster.primary.id
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "e2-medium"

    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
    image_type = "something"
  }
}
`},
			GoodExample: []string{`
resource "google_service_account" "default" {
  account_id   = "service-account-id"
  display_name = "Service Account"
}

resource "google_container_cluster" "primary" {
  name     = "my-gke-cluster"
  location = "us-central1"

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "good_example" {
  name       = "my-node-pool"
  cluster    = google_container_cluster.primary.id
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "e2-medium"

    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
    image_type = "COS_CONTAINERD"
  }
}
`},
			Links: []string{
				"https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/container_cluster#image_type",
				"https://cloud.google.com/kubernetes-engine/docs/concepts/node-images",
			},
		},
		RequiredTypes: []string{
			"resource",
		},
		RequiredLabels: []string{
			"google_container_node_pool",
		},
		DefaultSeverity: severity.Low,
		CheckFunc: func(set result.Set, resourceBlock block.Block, _ block.Module) {
			if imageTypeAttr := resourceBlock.GetBlock("node_config").GetAttribute("image_type"); imageTypeAttr.IsNil() { // alert on use of default value
				set.AddResult().
					WithDescription("Resource '%s' uses default value for node_config.image_type", resourceBlock.FullName())
			} else if imageTypeAttr.NotEqual("COS_CONTAINERD", block.IgnoreCase) && imageTypeAttr.NotEqual("COS", block.IgnoreCase) {
				set.AddResult().
					WithDescription("Resource '%s' does not have node_config.image_type set to COS_CONTAINERD", resourceBlock.FullName()).
					WithAttribute(imageTypeAttr)
			}
		},
	})
}

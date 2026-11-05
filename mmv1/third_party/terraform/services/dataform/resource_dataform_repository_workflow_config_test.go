package dataform_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/dataform"
	_ "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/secretmanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/sourcerepo"
)

func TestAccDataformRepositoryWorkflowConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataformRepositoryWorkflowConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataformRepositoryWorkflowConfig_basic(context),
			},
			{
				ResourceName:      "google_dataform_repository_workflow_config.workflow",
				ImportState:       true,
				ImportStateVerify: true,
				// region and repository are url_param_only and not returned in the API response body
				ImportStateVerifyIgnore: []string{"region", "repository"},
			},
			{
				Config: testAccDataformRepositoryWorkflowConfig_full(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataform_repository_workflow_config.workflow", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_dataform_repository_workflow_config.workflow",
				ImportState:       true,
				ImportStateVerify: true,
				// region and repository are url_param_only and not returned in the API response body
				ImportStateVerifyIgnore: []string{"region", "repository"},
			},
			{
				Config: testAccDataformRepositoryWorkflowConfig_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataform_repository_workflow_config.workflow", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_dataform_repository_workflow_config.workflow",
				ImportState:       true,
				ImportStateVerify: true,
				// region and repository are url_param_only and not returned in the API response body
				ImportStateVerifyIgnore: []string{"region", "repository"},
			},
		},
	})
}

func testAccDataformRepositoryWorkflowConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sourcerepo_repository" "git_repository" {
  name = "my/repository%{random_suffix}"
}

resource "google_secret_manager_secret" "secret" {
  secret_id = "tf-test-secret%{random_suffix}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret_version" {
  secret      = google_secret_manager_secret.secret.id
  secret_data = "tf-test-secret-data%{random_suffix}"
}

resource "google_dataform_repository" "repository" {
  name   = "tf_test_dataform_repository%{random_suffix}"
  region = "us-central1"

  git_remote_settings {
    url                                 = google_sourcerepo_repository.git_repository.url
    default_branch                      = "main"
    authentication_token_secret_version = google_secret_manager_secret_version.secret_version.id
  }
}

resource "google_dataform_repository_release_config" "release_1" {
  project       = google_dataform_repository.repository.project
  region        = google_dataform_repository.repository.region
  repository    = google_dataform_repository.repository.name
  name          = "tf_test_release_1_%{random_suffix}"
  git_commitish = "main"
}

resource "google_service_account" "dataform_sa_1" {
  account_id   = "tf-df-sa1-%{random_suffix}"
  display_name = "Dataform SA 1"
}

resource "google_dataform_repository_workflow_config" "workflow" {
  project        = google_dataform_repository.repository.project
  region         = google_dataform_repository.repository.region
  repository     = google_dataform_repository.repository.name
  name           = "tf_test_workflow%{random_suffix}"
  release_config = google_dataform_repository_release_config.release_1.id

  invocation_config {
    service_account = google_service_account.dataform_sa_1.email
  }
}
`, context)
}

func testAccDataformRepositoryWorkflowConfig_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sourcerepo_repository" "git_repository" {
  name = "my/repository%{random_suffix}"
}

resource "google_secret_manager_secret" "secret" {
  secret_id = "tf-test-secret%{random_suffix}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret_version" {
  secret      = google_secret_manager_secret.secret.id
  secret_data = "tf-test-secret-data%{random_suffix}"
}

resource "google_dataform_repository" "repository" {
  name   = "tf_test_dataform_repository%{random_suffix}"
  region = "us-central1"

  git_remote_settings {
    url                                 = google_sourcerepo_repository.git_repository.url
    default_branch                      = "main"
    authentication_token_secret_version = google_secret_manager_secret_version.secret_version.id
  }
}

resource "google_dataform_repository_release_config" "release_1" {
  project       = google_dataform_repository.repository.project
  region        = google_dataform_repository.repository.region
  repository    = google_dataform_repository.repository.name
  name          = "tf_test_release_1_%{random_suffix}"
  git_commitish = "main"
}

resource "google_service_account" "dataform_sa_1" {
  account_id   = "tf-df-sa1-%{random_suffix}"
  display_name = "Dataform SA 1"
}

resource "google_dataform_repository_workflow_config" "workflow" {
  project        = google_dataform_repository.repository.project
  region         = google_dataform_repository.repository.region
  repository     = google_dataform_repository.repository.name
  name           = "tf_test_workflow%{random_suffix}"
  release_config = google_dataform_repository_release_config.release_1.id

  invocation_config {
    included_targets {
      database = "gcp-example-project"
      schema   = "example-dataset"
      name     = "target_1"
    }
    included_tags                            = ["tag_1"]
    transitive_dependencies_included         = true
    transitive_dependents_included           = true
    fully_refresh_incremental_tables_enabled = false
    service_account                          = google_service_account.dataform_sa_1.email
  }

  cron_schedule = "0 7 * * *"
  time_zone     = "America/New_York"
  disabled      = true
}
`, context)
}

func testAccDataformRepositoryWorkflowConfig_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sourcerepo_repository" "git_repository" {
  name = "my/repository%{random_suffix}"
}

resource "google_secret_manager_secret" "secret" {
  secret_id = "tf-test-secret%{random_suffix}"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret_version" {
  secret      = google_secret_manager_secret.secret.id
  secret_data = "tf-test-secret-data%{random_suffix}"
}

resource "google_dataform_repository" "repository" {
  name   = "tf_test_dataform_repository%{random_suffix}"
  region = "us-central1"

  git_remote_settings {
    url                                 = google_sourcerepo_repository.git_repository.url
    default_branch                      = "main"
    authentication_token_secret_version = google_secret_manager_secret_version.secret_version.id
  }
}

resource "google_dataform_repository_release_config" "release_1" {
  project       = google_dataform_repository.repository.project
  region        = google_dataform_repository.repository.region
  repository    = google_dataform_repository.repository.name
  name          = "tf_test_release_1_%{random_suffix}"
  git_commitish = "main"
}

resource "google_dataform_repository_release_config" "release_2" {
  project       = google_dataform_repository.repository.project
  region        = google_dataform_repository.repository.region
  repository    = google_dataform_repository.repository.name
  name          = "tf_test_release_2_%{random_suffix}"
  git_commitish = "dev"
}

resource "google_service_account" "dataform_sa_1" {
  account_id   = "tf-df-sa1-%{random_suffix}"
  display_name = "Dataform SA 1"
}

resource "google_service_account" "dataform_sa_2" {
  account_id   = "tf-df-sa2-%{random_suffix}"
  display_name = "Dataform SA 2"
}

resource "google_dataform_repository_workflow_config" "workflow" {
  project        = google_dataform_repository.repository.project
  region         = google_dataform_repository.repository.region
  repository     = google_dataform_repository.repository.name
  name           = "tf_test_workflow%{random_suffix}"
  release_config = google_dataform_repository_release_config.release_2.id

  invocation_config {
    included_targets {
      database = "gcp-updated-project"
      schema   = "updated-dataset"
      name     = "target_2"
    }
    included_tags                            = ["tag_1", "tag_2"]
    transitive_dependencies_included         = false
    transitive_dependents_included           = false
    fully_refresh_incremental_tables_enabled = true
    service_account                          = google_service_account.dataform_sa_2.email
  }

  cron_schedule = "0 9 * * *"
  time_zone     = "UTC"
  disabled      = false
}
`, context)
}

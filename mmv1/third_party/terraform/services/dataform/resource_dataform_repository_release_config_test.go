package dataform_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	_ "github.com/hashicorp/terraform-provider-google/google/services/dataform"
	_ "github.com/hashicorp/terraform-provider-google/google/services/secretmanager"
	_ "github.com/hashicorp/terraform-provider-google/google/services/sourcerepo"
)

func TestAccDataformRepositoryReleaseConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataformRepositoryReleaseConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataformRepositoryReleaseConfig_basic(context),
			},
			{
				ResourceName:      "google_dataform_repository_release_config.release",
				ImportState:       true,
				ImportStateVerify: true,
				// region and repository are url_param_only and not returned in the API response body
				ImportStateVerifyIgnore: []string{"region", "repository"},
			},
			{
				Config: testAccDataformRepositoryReleaseConfig_full(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataform_repository_release_config.release", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_dataform_repository_release_config.release",
				ImportState:       true,
				ImportStateVerify: true,
				// region and repository are url_param_only and not returned in the API response body
				ImportStateVerifyIgnore: []string{"region", "repository"},
			},
			{
				Config: testAccDataformRepositoryReleaseConfig_update(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataform_repository_release_config.release", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_dataform_repository_release_config.release",
				ImportState:       true,
				ImportStateVerify: true,
				// region and repository are url_param_only and not returned in the API response body
				ImportStateVerifyIgnore: []string{"region", "repository"},
			},
		},
	})
}

func testAccDataformRepositoryReleaseConfig_basic(context map[string]interface{}) string {
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

resource "google_dataform_repository_release_config" "release" {
  project       = google_dataform_repository.repository.project
  region        = google_dataform_repository.repository.region
  repository    = google_dataform_repository.repository.name
  name          = "tf_test_release%{random_suffix}"
  git_commitish = "main"
}
`, context)
}

func testAccDataformRepositoryReleaseConfig_full(context map[string]interface{}) string {
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

resource "google_dataform_repository_release_config" "release" {
  project       = google_dataform_repository.repository.project
  region        = google_dataform_repository.repository.region
  repository    = google_dataform_repository.repository.name
  name          = "tf_test_release%{random_suffix}"
  git_commitish = "main"
  cron_schedule = "0 7 * * *"
  time_zone     = "America/New_York"
  disabled      = true

  code_compilation_config {
    default_database = "gcp-example-project"
    default_schema   = "example-dataset"
    default_location = "us-central1"
    assertion_schema = "example-assertion-dataset"
    database_suffix  = "_db1"
    schema_suffix    = "_schema1"
    table_prefix     = "prefix1_"
    vars = {
      var1 = "value1"
    }
  }
}
`, context)
}

func testAccDataformRepositoryReleaseConfig_update(context map[string]interface{}) string {
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

resource "google_dataform_repository_release_config" "release" {
  project       = google_dataform_repository.repository.project
  region        = google_dataform_repository.repository.region
  repository    = google_dataform_repository.repository.name
  name          = "tf_test_release%{random_suffix}"
  git_commitish = "dev"
  cron_schedule = "0 9 * * *"
  time_zone     = "UTC"
  disabled      = false

  code_compilation_config {
    default_database = "gcp-updated-project"
    default_schema   = "updated-dataset"
    default_location = "us-east1"
    assertion_schema = "updated-assertion-dataset"
    database_suffix  = "_db2"
    schema_suffix    = "_schema2"
    table_prefix     = "prefix2_"
    vars = {
      var1 = "updated"
      var2 = "value2"
    }
  }
}
`, context)
}

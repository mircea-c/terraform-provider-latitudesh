package latitudesh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testProject = "proj_jv6m5J78VaLPe"

func TestServerList_Basic(t *testing.T) {

	recorder, teardown := createTestRecorder(t)
	defer teardown()
	testAccProviders["latitudesh"].ConfigureContextFunc = testProviderConfigure(recorder)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccTokenCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testServerListCheckBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.latitudesh_server_list.test", "project", testProject),
				),
			},
		},
	})
}

func testServerListCheckBasic() string {
	return fmt.Sprintf(`
data "latitudesh_server_llist" "test" {
	project = "%s"
}
`,
		testProject,
	)
}

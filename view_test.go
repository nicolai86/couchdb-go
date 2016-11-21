package couchdb

import "testing"

type employeeResults struct {
	Results
	Employees []struct {
		Value testDoc `json:"value"`
	} `json:"rows"`
}

func TestDatabase_Results(t *testing.T) {
	playground.Put("_design/company", DesignDocument{
		Language: "javascript",
		Views: map[string]View{
			"employees": {
				MapFn: `
function(doc) {
  var type = doc._id.split(':')[0];
  if (type == 'employee') {
    emit(doc._id, doc);
  }
}`,
			},
		},
	})

	t.Run("get", func(t *testing.T) {
		var result = employeeResults{}
		if err := playground.Results("company", "employees", &result); err != nil {
			t.Fatal(err)
		}

		if len(result.Employees) != 2 {
			t.Fatalf("Expected %d results, but got %d", 2, len(result.Employees))
		}

		expectedIDs := []string{"employee:michael", "employee:raphael"}
		for _, id := range expectedIDs {
			known := false
			for _, doc := range result.Employees {
				known = known || doc.Value.ID == id
			}
			if !known {
				t.Fatalf("Expected doc %q to be included, but wasn't", id)
			}
		}
	})
}

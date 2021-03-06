package v1alpha4

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestVaultConnectionURI(t *testing.T) {
	tests := []struct {
		name string
		db   *Database
		want string
	}{
		{
			name: "Postgres",
			db: &Database{
				ObjectMeta: v1.ObjectMeta{
					Name: "testdb",
				},
				Spec: DatabaseSpec{
					Connection: DatabaseConnection{
						Postgres: &PostgresConnection{
							URI: ValueOrValueFrom{
								ValueFrom: &ValueFrom{
									Vault: &Vault{
										AgentInject: true,
										Role:        "test",
										Secret:      "database/creds/test",
									},
								},
							},
						},
					},
				},
			},
			want: `
{{- with secret "database/creds/test" -}}
postgres://{{ .Data.username }}:{{ .Data.password }}@postgres:5432/testdb{{- end }}`,
		},
		{
			name: "MySQL",
			db: &Database{
				ObjectMeta: v1.ObjectMeta{
					Name: "testdb",
				},
				Spec: DatabaseSpec{
					Connection: DatabaseConnection{
						Mysql: &MysqlConnection{
							URI: ValueOrValueFrom{
								ValueFrom: &ValueFrom{
									Vault: &Vault{
										AgentInject: true,
										Role:        "test",
										Secret:      "database/creds/test",
									},
								},
							},
						},
					},
				},
			},
			want: `
{{- with secret "database/creds/test" -}}
{{ .Data.username }}:{{ .Data.password }}@tcp(mysql:3306)/testdb{{- end }}`,
		},
		{
			name: "CockroachDB",
			db: &Database{
				ObjectMeta: v1.ObjectMeta{
					Name: "testdb",
				},
				Spec: DatabaseSpec{
					Connection: DatabaseConnection{
						CockroachDB: &CockroachDBConnection{
							URI: ValueOrValueFrom{
								ValueFrom: &ValueFrom{
									Vault: &Vault{
										AgentInject: true,
										Role:        "test",
										Secret:      "database/creds/test",
									},
								},
							},
						},
					},
				},
			},
			want: `
{{- with secret "database/creds/test" -}}
postgres://{{ .Data.username }}:{{ .Data.password }}@postgres:5432/testdb{{- end }}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			a, err := test.db.GetVaultAnnotations()
			req.NoError(err)

			if got := a["vault.hashicorp.com/agent-inject-template-schemaherouri"]; got != test.want {
				t.Fatalf("Expected:\n%s\ngot:\n%s", test.want, got)
			}
		})
	}
}

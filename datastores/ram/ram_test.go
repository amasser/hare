package ram

import (
	"errors"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestAllRamTests(t *testing.T) {
	var tests = []func(t *testing.T){
		func(t *testing.T) {
			//New...

			ram := newTestRam(t)
			defer ram.Close()

			want := make(map[int][]byte)
			want[1] = []byte(`{"id":1,"first_name":"John","last_name":"Doe","age":37}`)
			want[2] = []byte(`{"id":2,"first_name":"Abe","last_name":"Lincoln","age":52}`)
			want[3] = []byte(`{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":18}`)
			want[4] = []byte(`{"id":4,"first_name":"Helen","last_name":"Keller","age":25}`)

			got := ram.tables["contacts"].records

			if !reflect.DeepEqual(want, got) {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//Close...

			ram := newTestRam(t)
			ram.Close()

			wantErr := ErrNoTable
			_, gotErr := ram.ReadRec("contacts", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}

			got := ram.tables

			if nil != got {
				t.Errorf("want %v; got %v", nil, got)
			}
		},
		func(t *testing.T) {
			//CreateTable...

			ram := newTestRam(t)
			defer ram.Close()

			err := ram.CreateTable("newtable")
			if err != nil {
				t.Fatal(err)
			}

			want := true
			got := ram.TableExists("newtable")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//CreateTable (ErrTableExists)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := ErrTableExists
			gotErr := ram.CreateTable("contacts")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//DeleteRec...

			ram := newTestRam(t)
			defer ram.Close()

			err := ram.DeleteRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := ErrNoRecord
			_, got := ram.ReadRec("contacts", 3)

			if !errors.Is(got, want) {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//DeleteRec (ErrNoTable)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := ErrNoTable
			gotErr := ram.DeleteRec("nonexistent", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//GetLastID...

			ram := newTestRam(t)
			defer ram.Close()

			want := 4
			got, err := ram.GetLastID("contacts")
			if err != nil {
				t.Fatal(err)
			}

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//GetLastID (ErrNoTable)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := ErrNoTable
			_, gotErr := ram.GetLastID("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//IDs...

			ram := newTestRam(t)
			defer ram.Close()

			want := []int{1, 2, 3, 4}
			got, err := ram.IDs("contacts")
			if err != nil {
				t.Fatal(err)
			}

			sort.Ints(got)

			if len(want) != len(got) {
				t.Errorf("want %v; got %v", want, got)
			} else {

				for i := range want {
					if want[i] != got[i] {
						t.Errorf("want %v; got %v", want, got)
					}
				}
			}
		},
		func(t *testing.T) {
			//IDs (ErrNoTable)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := ErrNoTable
			_, gotErr := ram.IDs("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//InsertRec...

			ram := newTestRam(t)
			defer ram.Close()

			err := ram.InsertRec("contacts", 5, []byte(`{"id":5,"first_name":"Rex","last_name":"Stout","age":77}`))
			if err != nil {
				t.Fatal(err)
			}

			rec, err := ram.ReadRec("contacts", 5)
			if err != nil {
				t.Fatal(err)
			}

			want := `{"id":5,"first_name":"Rex","last_name":"Stout","age":77}`
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//InsertRec (ErrNoTable)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := ErrNoTable
			gotErr := ram.InsertRec("nonexistent", 5, []byte(`{"id":5,"first_name":"Rex","last_name":"Stout","age":77}`))

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//InsertRec (ErrIDExists)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := ErrRecordExists
			gotErr := ram.InsertRec("contacts", 3, []byte(`{"id":3,"first_name":"Rex","last_name":"Stout","age":77}`))
			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}

			rec, err := ram.ReadRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := `{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":18}`
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//ReadRec...

			ram := newTestRam(t)
			defer ram.Close()

			rec, err := ram.ReadRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := `{"id":3,"first_name":"Bill","last_name":"Shakespeare","age":18}`
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//ReadRec (ErrNoTable)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := ErrNoTable
			_, gotErr := ram.ReadRec("nonexistent", 3)

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//RemoveTable...

			ram := newTestRam(t)
			defer ram.Close()

			err := ram.RemoveTable("contacts")
			if err != nil {
				t.Fatal(err)
			}

			want := false
			got := ram.TableExists("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//RemoveTable (ErrNoTable)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := ErrNoTable
			gotErr := ram.RemoveTable("nonexistent")

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
		func(t *testing.T) {
			//TableExists...

			ram := newTestRam(t)
			defer ram.Close()

			want := true
			got := ram.TableExists("contacts")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}

			want = false
			got = ram.TableExists("nonexistant")

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//TableNames...

			ram := newTestRam(t)
			defer ram.Close()

			want := []string{"contacts"}
			got := ram.TableNames()

			sort.Strings(got)

			if len(want) != len(got) {
				t.Errorf("want %v; got %v", want, got)
			} else {

				for i := range want {
					if want[i] != got[i] {
						t.Errorf("want %v; got %v", want, got)
					}
				}
			}
		},
		func(t *testing.T) {
			//UpdateRec...

			ram := newTestRam(t)
			defer ram.Close()

			err := ram.UpdateRec("contacts", 3, []byte(`{"id":3,"first_name":"William","last_name":"Shakespeare","age":77}`))
			if err != nil {
				t.Fatal(err)
			}

			rec, err := ram.ReadRec("contacts", 3)
			if err != nil {
				t.Fatal(err)
			}

			want := `{"id":3,"first_name":"William","last_name":"Shakespeare","age":77}`
			got := string(rec)

			if want != got {
				t.Errorf("want %v; got %v", want, got)
			}
		},
		func(t *testing.T) {
			//UpdateRec (ErrNoTable)...

			ram := newTestRam(t)
			defer ram.Close()

			wantErr := ErrNoTable
			gotErr := ram.UpdateRec("nonexistent", 3, []byte(`{"id":3,"first_name":"William","last_name":"Shakespeare","age":77}`))

			if !errors.Is(gotErr, wantErr) {
				t.Errorf("want %v; got %v", wantErr, gotErr)
			}
		},
	}

	for i, fn := range tests {
		t.Run(strconv.Itoa(i), fn)
	}
}
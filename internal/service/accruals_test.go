package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccrualOrderService_CheckLuhn(t *testing.T) {
	type fields struct {
		repo AccrualOrderRepoContract
	}
	type args struct {
		number uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AccrualOrderService{
				repo: tt.fields.repo,
			}
			assert.Equalf(t, tt.want, a.CheckLuhn(tt.args.number), "CheckLuhn(%v)", tt.args.number)
		})
	}
}
